//#################################
// Verwaltung der Basismotive
//#################################
package images

import (
	"io"
	"net/http"

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

// Collection für Basismotive
var imageCollection *mgo.GridFS
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

	if err == nil { // => alles ok
		formdataZeiger := r.MultipartForm

		if formdataZeiger != nil { // beim ersten request ist die Form leer!
			files := formdataZeiger.File["newImg"]

			for i := range files {
				// upload-files öffnen:
				uplFile, _ := files[i].Open()
				defer uplFile.Close()

				// grid-file mit diesem Namen erzeugen:
				gridFile, _ := imageCollection.Create(files[i].Filename)

				newID := bson.NewObjectId()
				gridFile.SetId(newID)
				// 	// Zusatzinformationen in den Metadaten festhalten
				// 	// Jedes Bild hat eine zugehörige Sammlung
				gridFile.SetMeta(bson.M{"imgSet": currentImgSet.Value})

				// in GridFSkopieren:
				_, err = io.Copy(gridFile, uplFile)

				err = gridFile.Close()
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

// func GetAllImages(r *http.Request) {

// 	// // angemeldeten Nutzer von Cookie auslesen
// 	// cookie, _ := r.Cookie("currentUser")
// 	// user := cookie.Value

// 	// // aktuelle Motivsammlung aus Cookie auslesen
// 	// cookieSet, _ := r.Cookie("currentImgSet")
// 	// currentImageSet := cookieSet.Value

// }
