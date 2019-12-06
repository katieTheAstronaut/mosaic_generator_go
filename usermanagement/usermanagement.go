//#################################
// Nutzerverwaltung
//#################################

package usermanagement

import (
	"net/http"
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

// Funktion um NutzerCollection aus Main Package zu holen
func GetUserCollection(collection *mgo.Collection) {
	collectionUser = collection
}

// Funktion um neuen Nutzer in der DB zu erstellen
func RegisterNewUser(user string, pw string) (errorMessage string) {

	// Prüfen ob Nutzername und Passwort angegeben wurden
	if user == "" || pw == "" {
		errorMessage = "Fehler: Bitte Nutzername und Passwort angeben!"
		return errorMessage
	}

	// Prüfen ob Passwort lang genug ist
	if len(pw) <= 1 {
		errorMessage = "Fehler: Das Passwort sollte aus mindesten 6 Zeichen bestehen!"
		return errorMessage
	}

	// Prüfen ob nur zulässige Zeichen verwendet wurden
	if !(isValidString(user) && isValidString(pw)) {
		errorMessage = "Fehler: Bitte nur gültige Zeichen verwenden!"
		return errorMessage
	}

	// Prüfen ob Nutzer bereits existiert
	userExists := User{}
	collectionUser.Find(bson.M{"username": user}).One(&userExists)
	if userExists.Username != "" {
		errorMessage = "Fehler: Der Nutzer existiert bereits!"
		return errorMessage
	}

	// Neues Nutzer-Dokument in CollectionUser anlegen
	newUser := User{user, pw}

	// Neuen Nutzer in Collection übergeben
	collectionUser.Insert(newUser)

	return ""
}

// Funktion zum Prüfen, ob Nutzer angemeldet werden kann
func LoginUser(user string, pw string) (errorMessage string) {

	// Prüfen ob Nutzername und Passwort angegeben wurden
	if user == "" || pw == "" {
		errorMessage = "Fehler: Bitte Nutzername und Passwort angeben!"
		return errorMessage
	}

	// Prüfen ob Nutzer existiert
	userExists := User{}
	collectionUser.Find(bson.M{"username": user}).One(&userExists)
	if userExists.Username == "" {
		errorMessage = "Fehler: Der Nutzer existiert nicht!"
		return errorMessage
	}

	// Prüfen ob Passwort richtig ist
	if userExists.Password != pw {
		errorMessage = "Fehler: Falsches Passwort angegeben!"
		return errorMessage
	}

	// wenn alles richtig ist, kann Nutzer angemeldet werden

	return ""
}

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

// Funktion um cookie zu erstellen
func CreateCookie(name string, value string, w http.ResponseWriter) {
	cookie := http.Cookie{Name: name, Value: value}
	http.SetCookie(w, &cookie)
}

// Funktion um cookie zu löschen
func DeleteCookie(name string, w http.ResponseWriter) {
	cookie := http.Cookie{Name: name, MaxAge: -1}
	http.SetCookie(w, &cookie)
}
