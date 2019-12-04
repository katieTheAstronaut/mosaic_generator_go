// // Main Datei
// //
// // Template Management (holt und schickt die nötigen Templates)

package main

// Importierte Packages
import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	// Eigene Packages

	"./images"
	"./usermanagement"

	// Externe Packages
	"github.com/globalsign/mgo"
)

// Externe Template-Dateien einlesen
var t = template.Must(template.ParseFiles("templates/picx.html", "templates/register.html", "templates/login.html", "templates/home.html", "templates/images.html", "templates/imageSet.html"))

// deklariert Datenbank
var dataB *mgo.Database

// Structs für Basismotive und -Sammlungen
type ImageSet struct {
	SetName string `bson:"setName"`
	User    string `bson:"user"`
}

type ImageSetList struct {
	ImgSets []ImageSet
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

	// Collection für Nutzerverwaltung an usermanagement package übergeben
	usermanagement.GetUserCollection(dataB.C("userColl"))

	// Collections für Basismotive und Motivsammlungen an package images übergeben
	images.GetImgCollections(dataB.GridFS("imageColl"), dataB.C("imgSetColl"))

	// ----Handler-------------------------

	// File Server für statische Element (CSS, JS, ..) registrieren
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	// Handler
	http.HandleFunc("/picx", handlerPicx)                       // http://localhost:4242/picx
	http.HandleFunc("/postNewUser", handlerNewUser)             // http://localhost:4242/postNewUser
	http.HandleFunc("/getRegistration", handlerGetRegistration) // http://localhost:4242/getRegistration
	http.HandleFunc("/getLogin", handlerGetLogin)               // http://localhost:4242/getLogin
	http.HandleFunc("/home", handlerHome)                       // http://localhost:4242/home
	http.HandleFunc("/images", handlerImages)                   // http://localhost:4242/images
	http.HandleFunc("/uploadImage", handlerUploadImage)         // http://localhost:4242/uploadImage
	http.HandleFunc("/createSet", handlerCreateSet)             // http://localhost:4242/createSet
	http.HandleFunc("/showSet", handleShowSet)                  // http://localhost:4242/showSet
	http.HandleFunc("/showImg", handleShowImg)                  // http://localhost:4242/showImg
	http.HandleFunc("/logout", handleLogout)                    // http://localhost:4242/logout

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
		cookie := http.Cookie{Name: "currentUser", Value: user}
		http.SetCookie(w, &cookie)
		// Homeseite darstellen
		fmt.Fprint(w, t.ExecuteTemplate(w, "home.html", nil))

	} else {
		fmt.Fprint(w, errorMessage)
	}

}

//#################################
// Handler für Template- bzw. Seitenaufrufe
//#################################

// Handler für den Initialen Aufruf der Pixc-Seite
func handlerPicx(w http.ResponseWriter, r *http.Request) {

	// Bei Neuaufruf der Seite alle Cookies löschen
	cookie := http.Cookie{Name: "currentUser", MaxAge: -1}
	http.SetCookie(w, &cookie)
	cookieSet := http.Cookie{Name: "currentImgSet", MaxAge: -1}
	http.SetCookie(w, &cookieSet)

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

	// Alle Sammlungen des Nutzers holen, damit diese im Template dargestellt werden können
	userImageSets := images.GetAllImageSets(r)

	fmt.Fprint(w, t.ExecuteTemplate(w, "images.html", userImageSets))
}

// Handler für den Upload von Bildern
func handlerUploadImage(w http.ResponseWriter, r *http.Request) {
	// Funktion aufrufen (package images) um Bilder hochzuladen
	// Hier müssen Datenbank, Collection für die Bilder übergeben werden
	images.AddImage(r)
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
	// Query auslesen, um aktuelle Sammlung herauszufinden
	currentImgSet := r.URL.Query().Get("imgSet")

	// Cookie für aktuell ausgewählte Sammlung anlegen
	cookie := http.Cookie{Name: "currentImgSet", Value: currentImgSet}
	http.SetCookie(w, &cookie)

	fmt.Fprint(w, t.ExecuteTemplate(w, "imageSet.html", currentImgSet))

}

func handleShowImg(w http.ResponseWriter, r *http.Request) {

	// angemeldeten Nutzer von Cookie auslesen
	cookie, _ := r.Cookie("currentUser")
	user := cookie.Value

	// aktuelle Motivsammlung aus Cookie auslesen
	cookieSet, _ := r.Cookie("currentImgSet")
	currentImageSet := cookieSet.Value

	// Alle Bilder passend zum Nutzer und zur Sammlung auslesen

	gridfs := dataB.GridFS("imageColl")
	name := "gridFileName"
	f, err := gridfs.Open(name)
	if err != nil {
		log.Printf("Failed to open %s: %v", name, err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Add("Content-Type", "image/jpeg")
	_, err = io.Copy(w, f)

	err = f.Close()

	fmt.Println(err)

}

// Handler für den Logout des Nutzers
func handleLogout(w http.ResponseWriter, r *http.Request) {

	// Cookie des Nutzers löschen
	cookie := http.Cookie{Name: "currentUser", MaxAge: -1}
	http.SetCookie(w, &cookie)
	// Cookie für ausgewählte Sammlung löschen
	cookieSet := http.Cookie{Name: "currentImgSet", MaxAge: -1}
	http.SetCookie(w, &cookieSet)

	// Login-Seite schicken, damit Nutzer sich wieder anmelden kann
	fmt.Fprint(w, t.ExecuteTemplate(w, "login.html", nil))
}
