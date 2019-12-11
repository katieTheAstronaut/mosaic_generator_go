// // Main Datei
// //
// // Template Management (holt und schickt die nötigen Templates)

package main

// Importierte Packages
import (
	"fmt"
	"html/template"
	"net/http"

	// Eigene Packages

	"./images"
	"./mosaic"
	"./pools"
	"./usermanagement"

	// Externe Packages
	"github.com/globalsign/mgo"
)

// Externe Template-Dateien einlesen
var t = template.Must(template.ParseFiles("templates/picx.html",
	"templates/register.html", "templates/login.html", "templates/home.html",
	"templates/images.html", "templates/imageSet.html", "templates/pools.html",
	"templates/singlePool.html", "templates/mosaic.html", "templates/mosaicDisplay.html",
	"templates/imageInfo.html", "templates/mosaicSet.html", "templates/mosaicInfo.html"))

// deklariert Datenbank
var dataB *mgo.Database

// Struct für Mosaikinformationen
type MosaicInfo struct {
	ImgList  images.Images
	PoolList pools.PoolList
}

//#################################
// Main
//#################################
func main() {

	// ----- Datenbankzeugs-----------------
	// Verbindung zum Mongo-DBMS herstellen
	dbSession, _ := mgo.Dial("localhost")
	defer dbSession.Close()
	// Datenbank wählen (bzw. neu erstellen, wenn noch nicht vorhanden)
	dataB = dbSession.DB("HA19DB_kathrin_duerkop_630119")
	// Collections an die jeweiligen Packages übergeben
	usermanagement.GetUserCollection(dataB.C("userColl"), dataB.C("poolsColl"), dataB.C("imgSetColl"), dataB.GridFS("imageColl"))
	images.GetImgCollections(dataB.GridFS("imageColl"), dataB.C("imgSetColl"))
	pools.GetCollections(dataB.GridFS("imageColl"), dataB.C("poolsColl"))
	mosaic.GetCollections(dataB.GridFS("imageColl"), dataB.C("poolsColl"))

	// File Server für statische Element (CSS, JS, ..) registrieren
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	// ----Funktions-Handler-------------------------
	http.HandleFunc("/picx", handlerPicx)                             // http://localhost:4242/picx
	http.HandleFunc("/postNewUser", handlerNewUser)                   // http://localhost:4242/postNewUser
	http.HandleFunc("/getRegistration", handlerGetRegistration)       // http://localhost:4242/getRegistration
	http.HandleFunc("/getLogin", handlerGetLogin)                     // http://localhost:4242/getLogin
	http.HandleFunc("/home", handlerHome)                             // http://localhost:4242/home
	http.HandleFunc("/backToHome", handlerBackToHome)                 // http://localhost:4242/backToHome
	http.HandleFunc("/images", handlerImages)                         // http://localhost:4242/images
	http.HandleFunc("/uploadImage", handlerUploadImage)               // http://localhost:4242/uploadImage
	http.HandleFunc("/createSet", handlerCreateSet)                   // http://localhost:4242/createSet
	http.HandleFunc("/showSet", handleShowSet)                        // http://localhost:4242/showSet
	http.HandleFunc("/showImg", handleShowImg)                        // http://localhost:4242/showImg
	http.HandleFunc("/logout", handleLogout)                          // http://localhost:4242/logout
	http.HandleFunc("/pools", handlerPools)                           // http://localhost:4242/pools
	http.HandleFunc("/createPool", handlerCreatePool)                 // http://localhost:4242/createPool
	http.HandleFunc("/showPool", handlerShowPool)                     // http://localhost:4242/showPool
	http.HandleFunc("/uploadImageToPool", handlerUploadImageToPool)   // http://localhost:4242/uploadImageToPool
	http.HandleFunc("/deleteOriginals", handlerDeleteOriginals)       // http://localhost:4242/deleteOriginals
	http.HandleFunc("/deleteUser", handlerDeleteUser)                 // http://localhost:4242/deleteUser
	http.HandleFunc("/mosaic", handlerMosaic)                         // http://localhost:4242/mosaic
	http.HandleFunc("/generateMosaic", handlerGenerateMosaic)         // http://localhost:4242/generateMosaic
	http.HandleFunc("/generateMosaicFast", handlerGenerateMosaicFast) // http://localhost:4242/generateMosaicFast
	http.HandleFunc("/showMosaic", handlerShowMosaic)                 // http://localhost:4242/showMosaic
	http.HandleFunc("/showAllMosaics", handlerShowAllMosaics)         // http://localhost:4242/showAllMosaics
	http.HandleFunc("/getInfo", handlerGetInfo)                       // http://localhost:4242/getInfo
	http.HandleFunc("/getMosaicInfo", handlerGetMosaicInfo)           // http://localhost:4242/getMosaicInfo
	http.HandleFunc("/showMosaicBig", handlerShowMosaicBig)           // http://localhost:4242/showMosaicBig

	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		fmt.Println(err)
	}
}

//#################################
// Handler für Login/Registrierung
//#################################

// Handler für die Verarbeitung der im Client eingegebenen Registrierungsdaten
func handlerNewUser(w http.ResponseWriter, r *http.Request) {

	// Logindaten für neuen Nutzer auslesen
	username := r.PostFormValue("usernameInput")
	password := r.PostFormValue("passwordInput")

	// Neuen Nutzer registrieren und falls nötig Fehlermeldung abfangen
	errorMessage := usermanagement.RegisterNewUser(username, password)

	// Falls kein Fehler geworfen wurde, wurde der Nutzer in der DB registriert und wir können zur Anmeldeseite wechseln
	if errorMessage == "" {
		// Login darstellen
		fmt.Fprint(w, t.ExecuteTemplate(w, "login.html", nil))
	} else {
		// Fehlermeldung an Client zurückschicken, damit diese dem Nutzer dargestellt werden kann
		fmt.Fprint(w, errorMessage)
	}
}

// Handler für den Aufruf der Übersichtsseite (Home) von der Anmeldung aus
func handlerHome(w http.ResponseWriter, r *http.Request) {

	// Logindaten für neuen Nutzer auslesen
	user := r.PostFormValue("usernameInput")
	pw := r.PostFormValue("passwordInput")

	// Prüfen ob der Nutzer angemeldet werden kann
	errorMessage := usermanagement.LoginUser(user, pw)

	// Falls kein Fehler geworfen wurde, kann der Nutzer angemeldet werden
	if errorMessage == "" {

		// Cookie erstellen für Nutzer
		usermanagement.CreateCookie("currentUser", user, w)

		// Homeseite darstellen
		fmt.Fprint(w, t.ExecuteTemplate(w, "home.html", nil))

	} else {
		fmt.Fprint(w, errorMessage)
	}

}

// Handler für den Aufruf der Übersichtsseite (Home) von einer anderen Seite aus
func handlerBackToHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.ExecuteTemplate(w, "home.html", nil))
}

//#################################
// Handler für Template- bzw. Seitenaufrufe
//#################################

// Handler für den Initialen Aufruf der Pixc-Seite
func handlerPicx(w http.ResponseWriter, r *http.Request) {

	// Bei Neuaufruf der Seite alle Cookies löschen
	usermanagement.DeleteCookie("currentUser", w)
	usermanagement.DeleteCookie("currentImgSet", w)
	usermanagement.DeleteCookie("currentPool", w)

	// Base-Template mit integriertem Login-Bereich aufrufen
	t.ExecuteTemplate(w, "picx.html", nil)
}

// Handler für den Aufruf der Registrierungsseite
func handlerGetRegistration(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.ExecuteTemplate(w, "register.html", nil))
}

// Handler für den Aufruf der Registrierungsseite
func handlerGetLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.ExecuteTemplate(w, "login.html", nil))
}

// Handler für Aufruf der Motiv-Übersichtsseite
func handlerImages(w http.ResponseWriter, r *http.Request) {

	// Falls Cookie einer Sammlung existiert, diesen löschen
	cookieExists := images.CheckCookie(r, "currentImgSet")
	if cookieExists != "" {
		usermanagement.DeleteCookie("currentImgSet", w)
	}
	// Alle Sammlungen des Nutzers holen, damit diese im Template dargestellt werden können
	userImageSets := images.GetAllImageSets(r)

	fmt.Fprint(w, t.ExecuteTemplate(w, "images.html", userImageSets))
}

// Handler für den Upload von Bildern
func handlerUploadImage(w http.ResponseWriter, r *http.Request) {
	// Funktion aufrufen (package images) um Bilder zur DB hinzuzufügen
	images.AddImage(r)

	//Aktuelle Liste an Basismotiven abrufen
	newImgList := images.DisplaySet(r, w)
	// Template für Motivsammlung aufrufen und bildliste übergeben
	fmt.Fprint(w, t.ExecuteTemplate(w, "imageSet.html", newImgList))
}

// Handler für die Erstellung einer neuen Basismotiv-Sammlung
func handlerCreateSet(w http.ResponseWriter, r *http.Request) {
	// Ruft Funktion im package images auf, um neue Basismotiv-Sammlung
	// in der DB anzulegen
	images.CreateImageSet(r)
	//aktualisierte Liste an Sammlungen aufrufen
	userImageSets := images.GetAllImageSets(r)

	fmt.Fprint(w, t.ExecuteTemplate(w, "images.html", userImageSets))
}

// Handler für die Darstellung einer Sammlung
func handleShowSet(w http.ResponseWriter, r *http.Request) {
	// Liste aller Bilder der aktuell ausgewählten Sammlung auslesen
	newImgList := images.DisplaySet(r, w)
	// Template für Motivsammlung aufrufen und bildliste übergeben
	fmt.Fprint(w, t.ExecuteTemplate(w, "imageSet.html", newImgList))
}

func handleShowImg(w http.ResponseWriter, r *http.Request) {
	// Funktion zum auslesen und anzeigen eines einzelnen Bildes im Paket images aufrufen
	images.ShowImg(r, w)
}

// Handler für den Logout des Nutzers
func handleLogout(w http.ResponseWriter, r *http.Request) {
	// Cookie des Nutzers löschen
	usermanagement.DeleteCookie("currentUser", w)

	// Cookie für ausgewählte Sammlung löschen
	usermanagement.DeleteCookie("currentImgSet", w)

	// Cookie für ausgewählten Pool löschen
	usermanagement.DeleteCookie("currentPool", w)

	// Login-Seite schicken, damit Nutzer sich wieder anmelden kann
	fmt.Fprint(w, t.ExecuteTemplate(w, "login.html", nil))
}

// Handler für Aufruf der Pool-Übersichtsseite
func handlerPools(w http.ResponseWriter, r *http.Request) {

	// Falls Cookie eines Pools existiert, diesen löschen
	cookieExists := pools.CheckCookie(r, "currentPool")
	if cookieExists != "" {
		usermanagement.DeleteCookie("currentPool", w)
	}

	// Alle Sammlungen des Nutzers holen, damit diese im Template dargestellt werden können
	userPools := pools.GetAllPools(r)

	fmt.Fprint(w, t.ExecuteTemplate(w, "pools.html", userPools))
}

// Handler für die Erstellung eines neuen Pools
func handlerCreatePool(w http.ResponseWriter, r *http.Request) {
	// // Ruft Funktion im package pools auf, um neuen Pool in der DB anzulegen
	pools.CreatePool(r)
	//aktualisierte Liste an Pools aufrufen
	userPools := pools.GetAllPools(r)

	fmt.Fprint(w, t.ExecuteTemplate(w, "pools.html", userPools))
}

// Handler für die Darstellung einer Sammlung
func handlerShowPool(w http.ResponseWriter, r *http.Request) {
	// Liste aller Bilder des aktuell ausgewählten Pools auslesen
	newImgList := pools.DisplayPool(r, w)
	// Template für Motivsammlung aufrufen und bildliste übergeben
	fmt.Fprint(w, t.ExecuteTemplate(w, "singlePool.html", newImgList))
}

// Handler für den Upload von Kacheln
func handlerUploadImageToPool(w http.ResponseWriter, r *http.Request) {
	// Funktion aufrufen (package pools) um Kacheln zur DB hinzuzufügen
	pools.AddImage(r)
	//Aktuelle Liste an Basismotiven abrufen
	newImgList := pools.DisplayPool(r, w)
	// Template für Motivsammlung aufrufen und bildliste übergeben
	fmt.Fprint(w, t.ExecuteTemplate(w, "singlePool.html", newImgList))
}

// Handler für das Löschen der Originale der Kacheln eines Pools
func handlerDeleteOriginals(w http.ResponseWriter, r *http.Request) {
	// Originale in dem Pool löschen
	pools.DeleteOriginals(r)
	// Liste aller Bilder des aktuell ausgewählten Pools auslesen
	newImgList := pools.DisplayPool(r, w)
	// Template für Motivsammlung aufrufen und bildliste übergeben
	fmt.Fprint(w, t.ExecuteTemplate(w, "singlePool.html", newImgList))
}

// Handler für das Löschen eines Nutzers
func handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	usermanagement.DeleteUser(r, w)

	// Login-Seite schicken
	fmt.Fprint(w, t.ExecuteTemplate(w, "login.html", nil))

}

// Handler für Aufruf der Mosaik-Übersichtsseite
func handlerMosaic(w http.ResponseWriter, r *http.Request) {

	// Alle Sammlungen des Nutzers und die darin enthaltenen Bilder holen
	mosaicBaseImages := images.GetAllImagesAndSets(r, w)

	pools := pools.GetAllPools(r)

	// Neues Struct für alle nötigen Infos für die Mosaikseite anlegen
	mosaicInfo := MosaicInfo{mosaicBaseImages, pools}

	fmt.Fprint(w, t.ExecuteTemplate(w, "mosaic.html", mosaicInfo))
}

// Handler für das Erstellen des Mosaiks mit Farbabstands-Algorithmus
func handlerGenerateMosaic(w http.ResponseWriter, r *http.Request) {
	// Motivbild und Pool aus Form auslesen
	baseImage := r.PostFormValue("image")
	pool := r.PostFormValue("pool")

	newMosaic := mosaic.GenerateMosaic(baseImage, pool, r)
	fmt.Fprint(w, t.ExecuteTemplate(w, "mosaicDisplay.html", newMosaic))
}

// Handler für das Erstellen des Mosaiks mit Helligkeits-Algorithmus
func handlerGenerateMosaicFast(w http.ResponseWriter, r *http.Request) {

	// Motivbild und Pool aus Form auslesen
	baseImage := r.PostFormValue("image")
	pool := r.PostFormValue("pool")

	newMosaic := mosaic.GenerateMosaicFast(baseImage, pool, r)
	fmt.Fprint(w, t.ExecuteTemplate(w, "mosaicDisplay.html", newMosaic))
}

// Handler für die Darstellung eines einzelnen Mosaiks
func handlerShowMosaic(w http.ResponseWriter, r *http.Request) {
	// Funktion zum auslesen und anzeigen eines einzelnen Bildes im Paket images aufrufen
	mosaic.ShowMosaic(r, w)
}

// Handler für die Darstellung aller bisherigen Mosaike
func handlerShowAllMosaics(w http.ResponseWriter, r *http.Request) {
	// Funktion zum auslesen und anzeigen eines einzelnen Bildes im Paket images aufrufen
	mosaics := mosaic.GetAllMosaics(r, w)
	fmt.Fprint(w, t.ExecuteTemplate(w, "mosaicSet.html", mosaics))
}

// Handler um Bildinformationen darzustellen
func handlerGetInfo(w http.ResponseWriter, r *http.Request) {
	// Funktion zum Abfragen der Bildinformationen eines übergebenen Bildes
	image := r.URL.Query().Get("img")

	imageInfo := images.GetImageInfo(image)

	fmt.Fprint(w, t.ExecuteTemplate(w, "imageInfo.html", imageInfo))
}

// Handler um Bildinformationen der Mosaike darzustellen
func handlerGetMosaicInfo(w http.ResponseWriter, r *http.Request) {
	// Funktion zum Abfragen der Bildinformationen eines übergebenen Mosaikbildes
	mosaicName := r.URL.Query().Get("mosaicName")

	fmt.Print(mosaicName)
	mosaicInfo := mosaic.GetMosaicInfo(mosaicName)
	fmt.Print(mosaicInfo)
	fmt.Fprint(w, t.ExecuteTemplate(w, "mosaicInfo.html", mosaicInfo))
}

// Handler um Mosaik in Originalgröße darzustellen
func handlerShowMosaicBig(w http.ResponseWriter, r *http.Request) {
	mosaic.ShowMosaicBig(r, w)
}
