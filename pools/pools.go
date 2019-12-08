//#################################
// Verwaltung der Pools
//#################################
package pools

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// struct für Pool
type Pool struct {
	PoolName string `bson:"poolName"`
	User     string `bson:"user"`
	Size     int    `bson:"size"`
}

// struct für Liste von Pools
type PoolList struct {
	Pools []Pool
}

type Image struct {
	Filename string `bson:"filename"`
	URL      string `bson:"url"`
}

// name beschreibt Namen des Pools oder der Sammlung, zu der es gehört
type ImageList struct {
	Images []Image
	Name   string `bson:"name"`
}

// Collection für Bilder
var imageCollection *mgo.GridFS

// Collection für Pools
var poolCollection *mgo.Collection

// Collection für Pool von main package holen
func GetCollections(gridFS *mgo.GridFS, poolColl *mgo.Collection) {
	imageCollection = gridFS
	poolCollection = poolColl
}

// Neuen Pool für Kachelbilder anlegen
func CreatePool(r *http.Request) {

	// Daten für neue Sammlung auslesen
	poolName := r.PostFormValue("poolName")
	size := r.PostFormValue("poolSize")
	poolSize, _ := strconv.Atoi(size)

	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value

	// Prüfen ob Nutzer bereits einen gleichnamigen Pool angelegt hat
	// Alle Sammlungen des Nutzers auslesen
	poolExists := false
	var allPools []Pool
	poolCollection.Find(bson.M{"user": user}).All(&allPools)
	for i := 0; i < len(allPools); i++ {
		if allPools[i].PoolName == poolName {
			poolExists = true
		}
	}

	// Falls kein Duplikat entstehen würde, Motivsammlung anlegen
	if !poolExists {
		// Neue Sammlung anlegen
		newPool := Pool{poolName, user, poolSize}
		poolCollection.Insert(newPool)
	}

}

// Erstellt eine Liste aller dem Nutzer gehörenden Pools
func GetAllPools(r *http.Request) PoolList {

	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value

	var allPools []Pool
	// Alle Motivsammlungen des angemeldeten Nutzers aus DB abfragen
	poolCollection.Find(bson.M{"user": user}).All(&allPools)

	// ImageSetList mit angeforderten Sammlungen erstellen
	userPools := PoolList{
		Pools: allPools,
	}
	return userPools
}

// Funktion zur Darstellung eines Pools
func DisplayPool(r *http.Request, w http.ResponseWriter) ImageList {

	var result *mgo.GridFile
	newImgList := ImageList{}
	currentPool := ""

	// überprüfen ob cookie für aktuellen Pool bereits gesetzt ist
	cookieExists := CheckCookie(r, "currentPool")

	// ist kein cookie gesetzt, ist es der erste Aufruf der Seite
	if cookieExists == "" {
		// Query auslesen, um aktuelle Sammlung herauszufinden
		currentPool = r.URL.Query().Get("pool")

		// Cookie für aktuell ausgewählte Sammlung anlegen
		cookie := http.Cookie{Name: "currentPool", Value: currentPool}
		http.SetCookie(w, &cookie)
	} else {
		// ist bereits ein cookie gesetzt, ist dies der wiederholte Aufruf der Seite
		cookie, _ := r.Cookie("currentPool")
		currentPool = cookie.Value
	}

	// Alle Bilder passend zur aktuellen Sammlung auslesen
	iter := imageCollection.Find(bson.M{"metadata.pool": currentPool}).Iter()

	for imageCollection.OpenNext(iter, &result) {
		// url zum abrufen jedes bilder in der src im template erstellen
		imgURL := fmt.Sprintf("/showImg?filename=%s", result.Name())
		// neues Struct für jedes Bild erstellen
		newImage := Image{result.Name(), imgURL}
		// Alle Bilder in eine BildListe hinzufügen
		newImgList.Images = append(newImgList.Images, newImage)
	}
	newImgList.Name = currentPool

	return newImgList
}

// Funktion zum Prüfen, ob Cookie gesetzt ist
func CheckCookie(r *http.Request, name string) string {
	// leerer Cookie-wert
	value := ""

	// holt Cookie mit gewünschten Namen
	cookie, _ := r.Cookie(name)

	// wenn Cookie existiert wird der Wert gespeichert
	if cookie != nil {
		value = cookie.Value
	}

	// CookieWert zurückgeben
	return value
}

func AddImage(r *http.Request) {

	// multipart-form parsen und lesen:
	err := r.ParseMultipartForm(2000000) // bytes

	// aktuell ausgewählte Sammlung aus Cookie auslesen:
	var currentPool, _ = r.Cookie("currentPool")

	// zugehörige Größe des Pools auslesen
	curPool := Pool{}
	poolCollection.Find(bson.M{"poolName": currentPool.Value}).One(&curPool)
	currentSize := curPool.Size

	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value

	if err == nil { // => alles ok
		formdataZeiger := r.MultipartForm

		if formdataZeiger != nil { // beim ersten request ist die Form leer!
			files := formdataZeiger.File["newImg"]

			for i := range files {
				// upload-files öffnen:
				uplFile, _ := files[i].Open()
				defer uplFile.Close()

				// bild neu benennen mit nutzername am Anfang
				newFileName := user + "_" + files[i].Filename
				// grid-file mit diesem Namen erzeugen:
				gridFile, _ := imageCollection.Create(newFileName)

				// 	// Zusatzinformationen in den Metadaten festhalten
				// 	// Jedes Bild hat eine zugehörige Sammlung
				gridFile.SetMeta(bson.M{"pool": currentPool.Value})

				// in GridFSkopieren:
				_, err = io.Copy(gridFile, uplFile)

				err = gridFile.Close()
				CropAndScale(newFileName, currentSize, r, user)

			}
		}
	}

}

// Funktion zum Ausschneiden und skalieren der Originalbilder zu Kacheln der richtigen Größe
func CropAndScale(filename string, size int, r *http.Request, user string) {
	var resizedImg image.Image

	// aktuellen pool abrufen
	var currentPool, _ = r.Cookie("currentPool")

	// Bild in gridFS Collection suchen und öffnen
	f, err := imageCollection.Open(filename)

	if err != nil {
		log.Printf("Failed to open %s: %v", filename, err)
		return
	}
	defer f.Close()

	// Bild aus GridFS zu imaging.Image umwandeln
	newImg, _ := imaging.Decode(f)
	// image duplizieren
	clonedImg := imaging.Clone(newImg)

	// wenn Bild breiter als höher ist:
	if clonedImg.Bounds().Dx() > clonedImg.Bounds().Dy() {
		// Bild skalieren: Höhe ist übergebene Größe
		resizedImg = imaging.Resize(clonedImg, 0, size, imaging.Box)
	} else {
		// Bild skalieren: Breite ist übergebene Größe
		resizedImg = imaging.Resize(clonedImg, size, 0, imaging.Box)
	}
	// skaliertes Bild quadratisch zuschneiden
	croppedImg := imaging.CropCenter(resizedImg, size, size)

	// grid-file mit diesem Namen erzeugen:
	gridFile, _ := imageCollection.Create(user + "_" + "px_" + strconv.Itoa(size) + "_" + filename)

	// 	// Zusatzinformationen in den Metadaten festhalten
	// 	// Jedes Bild hat eine zugehörige Sammlung
	gridFile.SetMeta(bson.M{"pool": currentPool.Value})

	// in GridFSkopieren:
	png.Encode(gridFile, croppedImg)

	err = gridFile.Close()

	err = f.Close()
}

// Funktion zum Löschen der Originale von Kacheln
func DeleteOriginals(r *http.Request) {
	var result *mgo.GridFile

	// aktuellen pool abrufen
	var cookie, _ = r.Cookie("currentPool")
	currentPool := cookie.Value

	// Alle Poolkacheln auslesen
	iter := imageCollection.Find(bson.M{"metadata.pool": currentPool}).Iter()

	for imageCollection.OpenNext(iter, &result) {

		// wenn Dateiname mit "px" beginnt, ist es die verkleinerte Version, alternativ ist es das Original
		if !strings.HasPrefix(result.Name(), "px") {
			// Original löschen
			imageCollection.Remove(result.Name())
		}
	}

}
