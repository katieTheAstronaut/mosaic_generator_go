//#################################
// Verwaltung der Basismotive
//#################################
package images

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// struct für Basismotivsammlungen
type ImageSet struct {
	SetName string `bson:"setName"`
	User    string `bson:"user"`
}

// struct für Liste von Basismotivsammlungen
type ImageSetList struct {
	ImgSets []ImageSet
}

type Image struct {
	Filename string `bson:"filename"`
	URL      string `bson:"url"`
}

type ImageList struct {
	Images []Image
	Name   string `bson:"name"`
}

type Images struct {
	ImgLists []ImageList
}

type Infos struct {
	Name       string  `bson:"name"`
	Width      int     `bson:"width"`
	Height     int     `bson:"height"`
	Brightness float64 `bson:"brightness"`
	R          float64 `bson:"r"`
	G          float64 `bson:"g"`
	B          float64 `bson:"b"`
}

// Collection für Bilder
var imageCollection *mgo.GridFS

// Collection für Sammlungen
var imageSetCollection *mgo.Collection

// Collection für Basismotive von main package holen
func GetImgCollections(gridFS *mgo.GridFS, imageSetColl *mgo.Collection) {
	imageCollection = gridFS
	imageSetCollection = imageSetColl
}

// Neue Sammlung für Basismotive anlegen
func CreateImageSet(r *http.Request) {

	// Daten für neue Sammlung auslesen
	setName := r.PostFormValue("imgSetName")

	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value

	// Prüfen ob Nutzer bereits eine gleichnamige Sammlung angelegt hat
	// Alle Sammlungen des Nutzers auslesen
	setExists := false
	var allImageSets []ImageSet
	imageSetCollection.Find(bson.M{"user": user}).All(&allImageSets)
	for i := 0; i < len(allImageSets); i++ {
		if allImageSets[i].SetName == setName {
			setExists = true
		}
	}

	// Falls kein Duplikat entstehen würde, Motivsammlung anlegen
	if !setExists {
		// Neue Sammlung anlegen
		newSet := ImageSet{setName, user}
		imageSetCollection.Insert(newSet)
	}

}

// Funktion um Bilder in die Datenbank in die GridFS-Collection hochzuladen
func AddImage(r *http.Request) {

	// multipart-form parsen und lesen:
	err := r.ParseMultipartForm(2000000) // bytes

	// aktuell ausgewählte Sammlung aus Cookie auslesen:
	var currentImgSet, _ = r.Cookie("currentImgSet")
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
				gridFile.SetMeta(bson.M{"imgSet": currentImgSet.Value})

				// Bild verkleinern auf ca. 80 Pixel auf der längsten Seite
				// in GridFSkopieren:
				_, err = io.Copy(gridFile, uplFile)

				err = gridFile.Close()
				// Jedes Bild wird automatisch auf etwa 50Pixel Breihe oder Höhe (Verhältnis bleibt erhalten) runterskaliert
				Resize(newFileName, 50, r, user)

				// Original löschen (da meist viel zu groß)
				imageCollection.Remove(newFileName)

			}
		}
	}

}

// Erstellt eine Liste aller dem Nutzer gehörenden Sammlungen
func GetAllImageSets(r *http.Request) ImageSetList {

	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value

	var allImageSets []ImageSet
	// Alle Motivsammlungen des angemeldeten Nutzers aus DB abfragen
	imageSetCollection.Find(bson.M{"user": user}).All(&allImageSets)

	// ImageSetList mit angeforderten Sammlungen erstellen
	userImageSets := ImageSetList{
		ImgSets: allImageSets,
	}
	return userImageSets
}

// Funktion um ein einzelnes Bild aus der GridFS-Collektion herauszulesen

func ShowImg(r *http.Request, w http.ResponseWriter) {

	// Dateinamen des Bildes aus request auslesen
	filename := r.URL.Query().Get("filename")
	// Bild in gridFS Collection suchen und öffnen
	f, err := imageCollection.Open(filename)

	if err != nil {
		log.Printf("Failed to open %s: %v", filename, err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Thumbnail des Bildes erstellen
	// Hierfür muss Bild erst in das passende Image-Typ umgewandelt werden
	thumbnailSrc, _ := imaging.Decode(f)
	thumb := imaging.Thumbnail(thumbnailSrc, 100, 100, imaging.CatmullRom)
	png.Encode(w, thumb)

}

// Funktion zur Darstellung einer Motivsammlung
func DisplaySet(r *http.Request, w http.ResponseWriter) ImageList {
	newImgList := ImageList{}
	currentImgSet := ""

	// überprüfen ob cookie für aktuelle Sammlung bereits gesetzt ist
	cookieExists := CheckCookie(r, "currentImgSet")

	// ist kein cookie gesetzt, ist es der erste Aufruf der Seite
	if cookieExists == "" {
		// Query auslesen, um aktuelle Sammlung herauszufinden
		currentImgSet = r.URL.Query().Get("imgSet")

		// Cookie für aktuell ausgewählte Sammlung anlegen
		cookie := http.Cookie{Name: "currentImgSet", Value: currentImgSet}
		http.SetCookie(w, &cookie)
	} else {
		// ist bereits ein cookie gesetzt, ist dies der wiederholte Aufruf der Seite
		cookie, _ := r.Cookie("currentImgSet")
		currentImgSet = cookie.Value
	}
	// Funktion aufrufen zum erstellen der Liste aller Bilder der Sammlung
	newImgList = createImgList(r, w, currentImgSet)

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

// Funktionen zum Abrufen aller Bilder des Nutzers
func GetAllImagesAndSets(r *http.Request, w http.ResponseWriter) Images {
	newImages := Images{}

	imgSetList := GetAllImageSets(r)
	for _, imageSet := range imgSetList.ImgSets {
		imagesOfSingleSet := createImgList(r, w, imageSet.SetName)
		newImages.ImgLists = append(newImages.ImgLists, imagesOfSingleSet)
	}

	return newImages
}

// Funktion zum Erstellen der Liste aller Bilder
func createImgList(r *http.Request, w http.ResponseWriter, currentImgSet string) ImageList {
	var result *mgo.GridFile
	newImgList := ImageList{}

	// Alle Bilder passend zur aktuellen Sammlung auslesen
	iter := imageCollection.Find(bson.M{"metadata.imgSet": currentImgSet}).Iter()

	for imageCollection.OpenNext(iter, &result) {
		// url zum abrufen jedes bilder in der src im template erstellen
		imgURL := fmt.Sprintf("/showImg?filename=%s", result.Name())
		// neues Struct für jedes Bild erstellen
		newImage := Image{result.Name(), imgURL}
		// Alle Bilder in eine BildListe hinzufügen
		newImgList.Images = append(newImgList.Images, newImage)
	}
	newImgList.Name = currentImgSet

	return newImgList
}

// Funktion zum Ausschneiden und skalieren der Originalbilder zu Kacheln der richtigen Größe
func Resize(filename string, size int, r *http.Request, user string) {
	var resizedImg image.Image

	// aktuelle Sammlung abrufen
	var currentImgSet, _ = r.Cookie("currentImgSet")

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

	// grid-file mit diesem Namen erzeugen:
	gridFile, _ := imageCollection.Create(user + "_" + "px_" + strconv.Itoa(size) + "_" + filename)

	// 	// Zusatzinformationen in den Metadaten festhalten
	// 	// Jedes Bild hat eine zugehörige Sammlung
	gridFile.SetMeta(bson.M{"imgSet": currentImgSet.Value})

	// in GridFSkopieren:
	png.Encode(gridFile, resizedImg)

	err = gridFile.Close()

	err = f.Close()
}

func GetImageInfo(filename string) Infos {

	width, height := GetSizeOfImg(filename)
	brightness, r, g, b := ComputeBrightnessOfImg(filename, width, height)

	newInfo := Infos{
		filename, width, height, brightness, r, g, b}

	return newInfo
}

func GetSizeOfImg(filename string) (int, int) {

	f, _ := imageCollection.Open(filename)
	// Bild aus GridFS zu imaging.Image umwandeln
	newImg, _ := imaging.Decode(f)
	width := newImg.Bounds().Dx()
	height := newImg.Bounds().Dy()

	return width, height
}

func ComputeBrightnessOfImg(filename string, width int, height int) (float64, float64, float64, float64) {

	// Gewünschtes Bild öffnen
	img, _ := imageCollection.Open(filename)
	decodedImg, _ := imaging.Decode(img)

	pixels := float64(width * height)

	// R-,G-,B-Mittelwerte auf 0 setzen
	rMid := float64(0)
	gMid := float64(0)
	bMid := float64(0)

	// Über gesamte Kachel iterieren und R-,G-,B-Werte addieren
	for i := 1; i < width; i++ {
		for j := 1; j < height; j++ {
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
	rMid = rMid / pixels
	gMid = gMid / pixels
	bMid = bMid / pixels

	// Helligkeit auf Basis des Mittelwertes auslesen
	brightness := math.Sqrt(rMid*rMid) + math.Sqrt(gMid*gMid) + math.Sqrt(bMid*bMid)

	img.Close()

	return brightness, rMid, gMid, bMid
}
