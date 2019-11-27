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
var t = template.Must(template.ParseFiles("templates/picx.html", "templates/register.html", "templates/login.html", "templates/login2.html"))

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
	http.HandleFunc("/home", handlerHome)                       // http://localhost:4242/home
	http.HandleFunc("/postUserLoginData", handlerUserLoginData) // http://localhost:4242/postUserLoginData
	http.HandleFunc("/getRegistration", handlerGetRegistration) // http://localhost:4242/getRegistration
	http.HandleFunc("/getLogin", handlerGetLogin)               // http://localhost:4242/getLogin

	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		fmt.Println(err)
	}

}

// Handler für den Aufruf der Startseite (Login)
func handlerHome(w http.ResponseWriter, r *http.Request) {
	// Homeseite darstellen
	t.ExecuteTemplate(w, "picx.html", nil)
}

// Handler für die Verarbeitung der im Client eingegebenen LoginDaten
func handlerUserLoginData(w http.ResponseWriter, r *http.Request) {

	// Logindaten auslesen
	username := r.PostFormValue("usernameInput")
	password := r.PostFormValue("passwordInput")

	// Neuen Nutzer registrieren und falls nötig Fehlermeldung abfangen
	errorMessage := usermanagement.RegisterNewUser(username, password)

	// Fehlermeldung an Client zurückschicken, damit diese dem Nutzer dargestellt werden kann
	fmt.Fprint(w, errorMessage)
}

// Handler für den Aufruf der Registrierungsseite
func handlerGetRegistration(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.ExecuteTemplate(w, "register.html", nil))
}

// Handler für den Aufruf der Registrierungsseite
func handlerGetLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.ExecuteTemplate(w, "login.html", nil))
}
