// // Main Datei
// //
// // Template Management (holt und schickt die nötigen Templates)

package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./usermanagement"

	"github.com/globalsign/mgo"
)

// Externe Template-Dateien einlesen
var t = template.Must(template.ParseFiles("templates/picx.html", "templates/register.html", "templates/login.html", "templates/home.html"))

// deklariert Datenbank
var dataB *mgo.Database

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
	usermanagement.GetUserCollection(dataB.C("collUser"))

	// ----Handler-------------------------

	// File Server für statische Element (CSS, JS, ..) registrieren
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	// Handler
	http.HandleFunc("/picx", handlerPicx)                       // http://localhost:4242/picx
	http.HandleFunc("/postNewUser", handlerNewUser)             // http://localhost:4242/postNewUser
	http.HandleFunc("/getRegistration", handlerGetRegistration) // http://localhost:4242/getRegistration
	http.HandleFunc("/getLogin", handlerGetLogin)               // http://localhost:4242/getLogin
	http.HandleFunc("/home", handlerHome)                       // http://localhost:4242/home

	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		fmt.Println(err)
	}

}

// Handler für den Initialen Aufruf der Pixc-Seite
func handlerPicx(w http.ResponseWriter, r *http.Request) {
	// Base-Template mit integriertem Login-Bereich aufrufen
	t.ExecuteTemplate(w, "picx.html", nil)
}

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

// Handler für den Aufruf der Registrierungsseite
func handlerGetRegistration(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.ExecuteTemplate(w, "register.html", nil))
}

// Handler für den Aufruf der Registrierungsseite
func handlerGetLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.ExecuteTemplate(w, "login.html", nil))
}

// Handler für den Aufruf der Übersichtsseite (Home) von der Anmeldung aus
func handlerHome(w http.ResponseWriter, r *http.Request) {

	// Logindaten für neuen Nutzer auslesen
	user := r.PostFormValue("usernameInput")
	pw := r.PostFormValue("passwordInput")

	fmt.Println(user)
	fmt.Println(pw)

	errorMessage := usermanagement.LoginUser(user, pw)

	fmt.Println(errorMessage)

	// Falls kein Fehler geworfen wurde, kann der Nutzer angemeldet werden
	if errorMessage == "" {

		// Homeseite darstellen
		fmt.Fprint(w, t.ExecuteTemplate(w, "home.html", nil))

	} else {
		fmt.Fprint(w, errorMessage)
	}

}
