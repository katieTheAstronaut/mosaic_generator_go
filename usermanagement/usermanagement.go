//#################################
// Nutzerverwaltung
//#################################

package usermanagement

import (
	"regexp"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//struct für Nutzer
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// Collection für Nutzerdaten erstellen
var collectionUser *mgo.Collection

// Funktion um neue Collection in der Datenbank zu erstellen
func GetUserCollection(collection *mgo.Collection) {
	collectionUser = collection
}

// Funktion um neuen Nutzer in der DB zu erstellen
func RegisterNewUser(user string, pw string) (errorMessage string) {

	// Prüfen ob Nutzername und Passwort angegeben wurden
	if user == "" || pw == "" {
		errorMessage = "Bitte Nutzername und Passwort angeben!"
		return errorMessage
	}

	// Prüfen ob Passwort lang genug ist
	if len(pw) <= 5 {
		errorMessage = "Das Passwort sollte aus mindesten 6 Zeichen bestehen!"
		return errorMessage
	}

	// Prüfen ob nur zulässige Zeichen verwendet wurden
	if !(isValidString(user) && isValidString(pw)) {
		errorMessage = "Bitte nur gültige Zeichen verwenden!"
		return errorMessage
	}

	// Prüfen ob Nutzer bereits existiert
	userExists := User{}
	collectionUser.Find(bson.M{"username": user}).One(&userExists)
	if userExists.Username != "" {
		errorMessage = "Der Nutzer existiert bereits!"
		return errorMessage
	}

	// Neues Nutzer-Dokument in CollectionUser anlegen
	newUser := User{user, pw}

	// Neuen Nutzer in Collection übergeben
	collectionUser.Insert(newUser)

	return ""
}

// // Testnutzer-Dokument in CollectionUser anlegen
// testuser := User{"test1", "passwort1"}
// // Testnutzer in Collection übergeben
// _ = collectionUser.Insert(testuser)

// Funktion zum Prüfen, ob übergebener String nur aus Buchstaben besteht
func isValidString(value string) (isValid bool) {
	//nutzt regexp, ein package von golang.org um auf regulären Ausdruck zu testen
	//MustCompile panicked, wenn der Ausdruck nicht geparst werden kann
	checkString := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString

	if !checkString(value) {
		return false
	}

	return true
}
