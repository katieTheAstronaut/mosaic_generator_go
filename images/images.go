//#################################
// Verwaltung der Basismotive
//#################################
package images

import (
	"io"
	"net/http"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// truct für Basismotivsammlungen
type ImageSet struct {
	SetName string `bson:"setName"`
	User    string `bson:"user"`
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
	user := "katie"
	// Prüfen ob Nutzer bereits eine gleichnamige Sammlung angelegt hat

	newSet := ImageSet{setName, user}
	imageSetCollection.Insert(newSet)
}

func AddImage(r *http.Request) {

	// var newImg = Image{
	// 	name:   filename,
	// 	user:   "aaaaaa",
	// 	imgSet: "testSet",
	// }

	//Bild(er) aus Form auslesen
	reader, err := r.MultipartReader()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Jeden "part" in die Datenbank speichern
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		// falls part.FileName() leer, überspringen:
		if part.FileName() == "" {
			continue
		}
		// Dateien öffnen
		datei, err := os.Open(part.FileName())

		// GridFile für jedes Bild erstellen
		gridFile, _ := imageCollection.Create(part.FileName())

		// Zusatzinformationen in den Metadaten festhalten
		// Jedes Bild hat eine zugehörige Sammlung
		gridFile.SetMeta(bson.M{"imgSet": "testSet"})

		// Ursprungsbild in die Gridfile kopieren
		_, err = io.Copy(gridFile, datei)

		// Dateien schließen
		err = datei.Close()
		err = gridFile.Close()
	}

}
