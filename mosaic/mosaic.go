//#################################
// Verwaltung der Mosaike
//#################################

package mosaic

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Mosaic struct {
	URL string `bson:"url"`
}

// struct für Pool
type Pool struct {
	PoolName string `bson:"poolName"`
	User     string `bson:"user"`
	Size     int    `bson:"size"`
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

// Funktion zum Generieren eines neuen Mosaics mit dem übergebenem Basismotiv und Pool
func GenerateMosaic(baseImg string, pool string) Mosaic {

	// Mosaic erstellen
	mosaic := CreateNewImg(baseImg, pool)

	// Mosaic in DB speichern
	filename := UploadMosaic(mosaic)

	// URL zum anzeigen des Mosaics erstellen und in ein Struct des Typen Mosaic speichern
	url := fmt.Sprintf("/showMosaic?filename=%s", filename)
	newMosaic := Mosaic{url}

	// Mosaik-Struct zurückliefern
	return newMosaic
}

// Funktion um neues Mosaik zu erstellen
func CreateNewImg(baseImg string, pool string) *image.NRGBA {

	// Bild in gridFS Collection suchen und öffnen
	f, _ := imageCollection.Open(baseImg)

	// poolGröße abfragen
	thisPool := Pool{}
	poolCollection.Find(bson.M{"poolName": pool}).One(&thisPool)
	poolSize := thisPool.Size

	// Größe des Basismotivs abfragen
	decodedImg, _ := imaging.Decode(f)
	width := decodedImg.Bounds().Dx()
	height := decodedImg.Bounds().Dy()

	widthMosaic := width * poolSize
	heightMosaic := height * poolSize

	// Neues, leeres Bild erstellen
	newImage := imaging.New(widthMosaic, heightMosaic, color.NRGBA{237, 250, 250, 1})

	// Bild für den gewünschten Pixel holen
	coverImg, _ := imageCollection.Open("aa_px_20_aa_david-kovalenko-G85VuTpw6jg-unsplash.jpg")
	gridDecoded, _ := imaging.Decode(coverImg)

	tile := imaging.Fill(gridDecoded, poolSize, poolSize, imaging.Center, imaging.Lanczos)

	// Bild an spezifische Stelle über leeres Bild legen
	for i := 1; i < widthMosaic/10; i += poolSize {
		for j := 1; j < heightMosaic/10; j += poolSize {
			// Helligkeit des Pixels ausrechnen
			// brightness := ComputeBrightness(baseImg, i, j)

			// passende Kachel an der Stelle einsetzen
			newImage = imaging.Paste(newImage, tile, image.Pt(i, j))
		}
	}

	// newImage = imaging.Paste(newImage, tile, image.Pt(20, 1))

	f.Close()

	return newImage
}

// Funktion um Mosaic in DB hochzuladen
func UploadMosaic(image *image.NRGBA) string {

	filename := "mosaic28"
	gridFile, _ := imageCollection.Create(filename)

	// 	// Zusatzinformationen in den Metadaten festhalten
	// 	// Jedes Bild hat eine zugehörige Sammlung
	gridFile.SetMeta(bson.M{"mosaic": "true"})

	// in GridFSkopieren:
	png.Encode(gridFile, image)

	_ = gridFile.Close()

	return filename
}

// Funktion zum Darstellen eines Mosaic-Bildes
func ShowMosaic(r *http.Request, w http.ResponseWriter) {

	// Dateinamen des Bildes aus request auslesen
	filename := r.URL.Query().Get("filename")

	// MosaikBild in gridFS Collection suchen und öffnen
	f, err := imageCollection.Open(filename)

	if err != nil {
		log.Printf("Failed to open %s: %v", filename, err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Thumbnail des Bildes erstellen
	// Hierfür muss Bild erst in das passende Image-Typ umgewandelt werden
	img, _ := imaging.Decode(f)
	// thumb := imaging.Thumbnail(img, 100, 100, imaging.CatmullRom)
	png.Encode(w, img)

}

// Funktion um Helligkeit eines Pixels zu berechnen
func ComputeBrightness(filename string, x int, y int) float64 {

	// Bild in gridFS Collection suchen und öffnen
	img, _ := imageCollection.Open(filename)
	decodedImg, _ := imaging.Decode(img)

	// RGB-Werte an der gewünschten Stelle des Bildes auslesen
	r, g, b, _ := decodedImg.At(x, y).RGBA()
	// von uint32 auf float64 umrechnen
	realR := float64(r / 257)
	realG := float64(g / 257)
	realB := float64(b / 257)

	// Helligkeit des Pixels berechnen
	brightness := math.Abs(math.Sqrt(realR*realR) + math.Sqrt(realG*realG) + math.Sqrt(realB*realB))

	return brightness
}
