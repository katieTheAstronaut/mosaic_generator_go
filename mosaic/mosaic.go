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
	"math/rand"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Mosaic struct {
	URL string `bson:"url"`
}

type Image struct {
	Filename string `bson:"filename"`
	URL      string `bson:"url"`
}

type ImageList struct {
	Images []Image
	Name   string `bson:"name"`
}

// struct für Pool
type Pool struct {
	PoolName   string    `bson:"poolName"`
	User       string    `bson:"user"`
	Size       int       `bson:"size"`
	Filenames  []string  `bson:"filenames"`
	Brightness []float64 `bson:"brightness"`
}

type MosaicInfo struct {
	Name   string `bson:"name"`
	Width  int    `bson:"width"`
	Height int    `bson:"height"`
	URL    string `bson:"url"`
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
func GenerateMosaic(baseImg string, pool string, r *http.Request) Mosaic {

	// Mosaic erstellen
	mosaic := CreateNewImg(baseImg, pool)

	// Neuen Dateinamen erstellen
	mosaicName := GetRandomName(r, "mosaic")
	// Mosaic in DB speichern
	UploadMosaic(mosaic, mosaicName)

	// URL zum anzeigen des Mosaics erstellen und in ein Struct des Typen Mosaic speichern
	url := fmt.Sprintf("/showMosaic?filename=%s", mosaicName)
	newMosaic := Mosaic{url}

	// Mosaik-Struct zurückliefern
	return newMosaic
}

// Funktion zum Generieren eines neuen Mosaics mit dem übergebenem Basismotiv und Pool
func GenerateMosaicFast(baseImg string, pool string, r *http.Request) Mosaic {

	// Mosaic erstellen
	mosaic := CreateNewImgByBrightness(baseImg, pool)

	// Neuen Dateinamen erstellen
	mosaicName := GetRandomName(r, "mosaic")
	// Mosaic in DB speichern
	UploadMosaic(mosaic, mosaicName)

	// URL zum anzeigen des Mosaics erstellen und in ein Struct des Typen Mosaic speichern
	url := fmt.Sprintf("/showMosaic?filename=%s", mosaicName)
	newMosaic := Mosaic{url}

	// Mosaik-Struct zurückliefern
	return newMosaic
}

// Funktion um neues Mosaik zu erstellen (Farbabstand-Algo)
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

	// Größe des Mosaiks berechnen
	widthMosaic := width * poolSize
	heightMosaic := height * poolSize

	// Neues, leeres Bild erstellen
	newImage := imaging.New(widthMosaic, heightMosaic, color.NRGBA{0, 0, 0, 1})

	// Zugehörigen Filename-Slice des aktuellen Pools abfragen
	filesSlice := thisPool.Filenames

	// Bild an spezifische Stelle über leeres Bild legen
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {

			// Farbe des Pixels berechnen
			rPx, gPx, bPx := computeColourOfPx(baseImg, i, j)

			// Kachel mit geringstem Farbabstand suchen
			closest := getSmallestColourDifference(rPx, gPx, bPx, filesSlice, poolSize)
			closestTile, _ := imageCollection.Open(closest)
			decodedTile, _ := imaging.Decode(closestTile)
			tile := imaging.Fill(decodedTile, poolSize, poolSize, imaging.Center, imaging.Lanczos)
			// passende Kachel an der Stelle einsetzen
			newImage = imaging.Paste(newImage, tile, image.Pt(i*poolSize, j*poolSize))

			closestTile.Close()
		}
	}
	f.Close()
	return newImage
}

// Funktion um neues Mosaik zu erstellen (Helligkeitsalgorithmus)
func CreateNewImgByBrightness(baseImg string, pool string) *image.NRGBA {

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

	// Größe des Mosaiks berechnen
	widthMosaic := width * poolSize
	heightMosaic := height * poolSize

	// Neues, leeres Bild erstellen
	newImage := imaging.New(widthMosaic, heightMosaic, color.NRGBA{0, 0, 0, 1})

	//Helligkeits-Slice des aktuellen Pools abfragen
	brightnessSlice := thisPool.Brightness
	// Zugehörigen Filename-Slice des aktuellen Pools abfragen
	filesSlice := thisPool.Filenames

	// Bild an spezifische Stelle über leeres Bild legen
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {

			// Helligkeit des Pixels ausrechnen
			brightness := ComputeBrightness(baseImg, i, j)

			// Kachel mit ähnlichster Helligkeit aus Pool suchen
			closest := getClosestBrightness(brightness, brightnessSlice, filesSlice)

			closestTile, _ := imageCollection.Open(closest)
			decodedTile, _ := imaging.Decode(closestTile)
			tile := imaging.Fill(decodedTile, poolSize, poolSize, imaging.Center, imaging.Lanczos)
			// passende Kachel an der Stelle einsetzen
			newImage = imaging.Paste(newImage, tile, image.Pt(i*poolSize, j*poolSize))

			closestTile.Close()
		}
	}
	f.Close()
	return newImage
}

// Funktion um Mosaic in DB hochzuladen
func UploadMosaic(image *image.NRGBA, filename string) {

	gridFile, _ := imageCollection.Create(filename)

	// 	// Zusatzinformationen in den Metadaten festhalten
	// 	// Jedes Bild hat eine zugehörige Sammlung
	gridFile.SetMeta(bson.M{"mosaic": "true"})

	// in GridFSkopieren:
	png.Encode(gridFile, image)

	_ = gridFile.Close()
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
	thumb := imaging.Thumbnail(img, 100, 100, imaging.CatmullRom)
	png.Encode(w, thumb)

}

// Funktion um Helligkeit eines Pixels zu berechnen
func ComputeBrightness(filename string, x int, y int) float64 {

	// RGB-Werte an der gewünschten Stelle des Bildes auslesen
	r, g, b := computeColourOfPx(filename, x, y)
	// Helligkeit des Pixels berechnen
	brightness := math.Sqrt((r * r) + (g * g) + (b * b))

	return brightness
}

// Funktion um Farbe eines Pixels zu berechnen
func computeColourOfPx(filename string, x int, y int) (float64, float64, float64) {
	// Bild in gridFS Collection suchen und öffnen
	img, _ := imageCollection.Open(filename)
	decodedImg, _ := imaging.Decode(img)

	// RGB-Werte an der gewünschten Stelle des Bildes auslesen
	r, g, b, _ := decodedImg.At(x, y).RGBA()
	// von uint32 auf float64 umrechnen
	realR := float64(r / 257)
	realG := float64(g / 257)
	realB := float64(b / 257)

	return realR, realG, realB
}

func ComputeColourDifference(filename string, rPx float64, gPx float64, bPx float64, poolSize int) float64 {

	// Mittlere Farbwerte der Kachel abrufen
	rMid, gMid, bMid := computeColour(filename, poolSize)

	// Farbunterschiede berechnen
	diffR := rMid - rPx
	diffG := gMid - gPx
	diffB := bMid - bPx

	// Farbabstand insgesamt berechnen:
	diff := math.Sqrt((diffR * diffR) + (diffG * diffG) + (diffB * diffB))

	return diff
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

// Funktion um ähnlichste Helligkeit zu einem Pixel in einem Pool zu finden
func getClosestBrightness(brightness float64, brightnessPool []float64, filePool []string) string {
	minimum := float64(1000)
	var closest = ""
	// Alle Helligkeiten im Pool durchgehen und jeweils die ähnlichste Helligkeit speichern
	for i := range brightnessPool {
		difference := math.Abs(brightness - brightnessPool[i])
		if difference < minimum {
			minimum = difference
			closest = filePool[i]
		}
	}
	return closest
}

// Funktion um ähnlichste Helligkeit zu einem Pixel in einem Pool zu finden
func getSmallestColourDifference(rPx float64, gPx float64, bPx float64, filePool []string, poolSize int) string {

	minimum := float64(1000)
	var closest = ""

	// Jede Kachel im Pool durchgehen und Farbabstand berechnen, kleinsten Farbabstand speichern
	for i, filename := range filePool {

		difference := ComputeColourDifference(filename, rPx, gPx, bPx, poolSize)
		if difference < minimum {
			minimum = difference
			closest = filePool[i]
		}
	}
	return closest
}

// Funktion um zufälligen Dateinamen zu erstellen
func GetRandomName(r *http.Request, filetype string) string {
	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value

	// zufälligen String erstellen
	randomString := getRandomString(5)

	// Dateiname erstellen
	filename := user + "_" + filetype + "_" + randomString
	return filename
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// Funktion um zufällige Zeichenkette zu erstellen
func getRandomString(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[seededRand.Intn(len(characters))]
	}
	return string(b)
}

// Funktion zum Erstellen der Liste aller Mosaike des Nutzers
func GetAllMosaics(r *http.Request, w http.ResponseWriter) ImageList {
	var result *mgo.GridFile
	newImgList := ImageList{}

	// Alle Bilder passend zur aktuellen Sammlung auslesen
	iter := imageCollection.Find(bson.M{"metadata.mosaic": "true"}).Iter()

	for imageCollection.OpenNext(iter, &result) {
		// url zum abrufen jedes bilder in der src im template erstellen
		imgURL := fmt.Sprintf("/showMosaic?filename=%s", result.Name())
		// neues Struct für jedes Bild erstellen
		newImage := Image{result.Name(), imgURL}
		// Alle Bilder in eine BildListe hinzufügen
		newImgList.Images = append(newImgList.Images, newImage)
	}
	newImgList.Name = "Mosaike"

	return newImgList
}

// Funktion zum Erstellen der Infos für ein jeweiliges Mosaik
func GetMosaicInfo(filename string) MosaicInfo {
	width, height := GetSizeOfImg(filename)
	url := fmt.Sprintf("/showMosaicBig?filename=%s", filename)

	newInfo := MosaicInfo{
		filename, width, height, url}

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

// Funktion um Mosaik in Originalgröße darzustellen
func ShowMosaicBig(r *http.Request, w http.ResponseWriter) {
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

	img, _ := imaging.Decode(f)
	png.Encode(w, img)
}
