//#################################
// Verwaltung der Pools
//#################################
package pools

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// struct für Pool
type Pool struct {
	PoolName   string    `bson:"poolName"`
	User       string    `bson:"user"`
	Size       int       `bson:"size"`
	Filenames  []string  `bson:"filenames"`
	Brightness []float64 `bson:"brightness"`
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
		newPool := Pool{poolName, user, poolSize, []string{}, []float64{}}
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
				// Bild auf geringere Pixelgröße verkleinern
				croppedFileName := CropAndScale(newFileName, currentSize, r, user)

				// Kachelhelligkeit auslesen
				brightness := ComputeBrightnessOfImg(croppedFileName, currentSize)
				// Bild und Kachel-Helligkeit in Slice speichern

				curPool.Filenames = append(curPool.Filenames, croppedFileName)
				curPool.Brightness = append(curPool.Brightness, brightness)
				poolCollection.Update(bson.M{"poolName": currentPool.Value}, curPool)
				poolCollection.Update(bson.M{"poolName": currentPool.Value}, curPool)

			}
		}
	}

}

// Funktion zum Ausschneiden und skalieren der Originalbilder zu Kacheln der richtigen Größe
func CropAndScale(filename string, size int, r *http.Request, user string) string {
	var resizedImg image.Image

	// aktuellen pool abrufen
	var currentPool, _ = r.Cookie("currentPool")

	// Bild in gridFS Collection suchen und öffnen
	f, _ := imageCollection.Open(filename)

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
	croppedFilename := user + "_" + "px_" + strconv.Itoa(size) + "_" + filename
	// grid-file mit diesem Namen erzeugen:
	gridFile, _ := imageCollection.Create(croppedFilename)

	// 	// Zusatzinformationen in den Metadaten festhalten
	// 	// Jedes Bild hat eine zugehörige Sammlung
	gridFile.SetMeta(bson.M{"pool": currentPool.Value})

	// in GridFSkopieren:
	png.Encode(gridFile, croppedImg)

	gridFile.Close()

	f.Close()

	return croppedFilename
}

// Funktion zum Löschen der Originale von Kacheln
func DeleteOriginals(r *http.Request) {
	var result *mgo.GridFile
	user := getCurrentUser(r)

	// aktuellen pool abrufen
	var cookie, _ = r.Cookie("currentPool")
	currentPool := cookie.Value
	prefix := user + "_" + "px"

	// Alle Poolkacheln auslesen
	iter := imageCollection.Find(bson.M{"metadata.pool": currentPool}).Iter()

	for imageCollection.OpenNext(iter, &result) {

		// wenn Dateiname mit "px" beginnt, ist es die verkleinerte Version, alternativ ist es das Original
		if !strings.HasPrefix(result.Name(), prefix) {
			// Original löschen
			imageCollection.Remove(result.Name())
		}
	}
}

// Funktion um aktuellen Nutzer auszulesen
func getCurrentUser(r *http.Request) string {
	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value
	return user
}

// Funktion um Helligkeit eines ganzen Bildes zu berechnen
func ComputeBrightnessOfImg(filename string, poolSize int) float64 {

	// Mittlere Farbwerte der Kachel abrufen
	rMid, gMid, bMid := computeColour(filename, poolSize)

	// Helligkeit auf Basis des Mittelwertes auslesen
	brightness := math.Sqrt((rMid * rMid) + (gMid * gMid) + (bMid * bMid))

	return brightness
}

// Funktion um die mittlere Farbe einer Kachel zu berechnen
func computeColour(filename string, poolSize int) (float64, float64, float64) {
	// Gewünschtes Bild öffnen
	img, _ := imageCollection.Open(filename)
	decodedImg, _ := imaging.Decode(img)

	poolsize := float64(poolSize)
	// R-,G-,B-Mittelwerte auf 0 setzen
	rMid := float64(0)
	gMid := float64(0)
	bMid := float64(0)

	// Über gesamte Kachel iterieren und R-,G-,B-Werte addieren
	for i := 1; i < poolSize; i++ {
		for j := 1; j < poolSize; j++ {
			r, g, b, _ := decodedImg.At(i, j).RGBA()

			realR := float64(r / 257)
			realG := float64(g / 257)
			realB := float64(b / 257)

			rMid = rMid + realR
			gMid = gMid + realG
			bMid = bMid + realB

		}
	}
	// RGB-Werte durch Pixel teilen um Mittelwerte zu erhalten
	rMid = rMid / (poolsize * poolsize)
	gMid = gMid / (poolsize * poolsize)
	bMid = bMid / (poolsize * poolsize)

	img.Close()

	return rMid, gMid, bMid

}
