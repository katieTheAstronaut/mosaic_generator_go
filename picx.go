// // Main Datei
// //
// // Template Management (holt und schickt die nötigen Templates)

package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/globalsign/mgo"
)

// Externe Template-Dateien einlesen
var t = template.Must(template.ParseFiles("templates/picx.html"))

//#################################
// Nutzerverwaltung
//#################################

// deklariert Datenbank
var dataB *mgo.Database

// Collection für Nutzerdaten erstellen
var collectionUser *mgo.Collection

//struct für Nutzer
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// Funktion um neuen Nutzer in der DB zu erstellen
func registerNewUser(user string, pw string) {

	// Neues Nutzer-Dokument in CollectionUser anlegen
	newUser := User{user, pw}
	// Neuen Nutzer in Collection übergeben
	_ = collectionUser.Insert(newUser)

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

	// collection für Nutzerverwaltung erstellen
	collectionUser = dataB.C("collUser")

	registerNewUser("test2", "pw2")
	// // Testnutzer-Dokument in CollectionUser anlegen
	// testuser := User{"test1", "passwort1"}
	// // Testnutzer in Collection übergeben
	// _ = collectionUser.Insert(testuser)

	// ----Handler-------------------------

	// File Server für statische Element (CSS, JS, ..)
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	http.HandleFunc("/home", homeHandler)

	err := http.ListenAndServe(":4242", nil)

	if err != nil {
		fmt.Println(err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "picx.html", nil)
}
