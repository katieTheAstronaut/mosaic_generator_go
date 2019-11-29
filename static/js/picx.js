// EventListener für den initialen Aufruf der picx-Seite lädt lediglich die spezifischen 
// JS-Funktionen für das Login-Template
window.addEventListener("load", function () {
    addJSforLogin();
});

// //#################################
// //Register-Template
// //#################################
function addJSforRegister() {

    // EventListener für den Registrier-Button 
    $("registerButton").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Nutzerdaten zu schicken 
        //------------------------------------------------
        // neues XMLHttpRequest anlegen
        var xhrPostUserInfo = new XMLHttpRequest();

        // callback, um Fehlermeldungen als Antwort zu erhalten und diese im html einzusetzen
        xhrPostUserInfo.addEventListener('load', function () {
            // Prüfen, ob Antwort Fehlercode beinhaltet oder nicht
            if (xhrPostUserInfo.responseText.indexOf("Fehler") == 0) {
                // Falls Fehler
                $("errorOutput").innerText = xhrPostUserInfo.responseText;
            } else {
                // Falls kein Fehler, wurde nur Template geliefert, also zurück zur Login-Seite wechseln
                $("template").innerHTML = xhrPostUserInfo.responseText;
                addJSforLogin();
            }
        });

        // Anfrage definieren und mit FormValues absenden
        xhrPostUserInfo.open('POST', 'http://localhost:4242/postNewUser');
        xhrPostUserInfo.send(new FormData($('registerForm')));

    });


    // EventListener für den "Zurück Zum Login"-Button, der wieder zur Anmeldeseite wechselt
    $("changeToLogin").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Nutzerdaten zu schicken 
        //------------------------------------------------
        var xhrGetLogin = new XMLHttpRequest();

        xhrGetLogin.addEventListener('load', function () {

            $("template").innerHTML = xhrGetLogin.responseText;
            addJSforLogin();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetLogin.open('GET', 'http://localhost:4242/getLogin');
        xhrGetLogin.send();

    });



};


//#################################
//LOGIN-Template
//#################################
function addJSforLogin() {

    // EventListener für den "Neues Konto Erstellen" Button, der die Registrier-Seite wechselt. 
    $("changeToRegister").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Nutzerdaten zu schicken 
        //------------------------------------------------
        // neues XMLHttpRequest anlegen
        var xhrGetRegistration = new XMLHttpRequest();

        // callback, um Fehlermeldungen als Antwort zu erhalten und diese im html einzusetzen
        xhrGetRegistration.addEventListener('load', function () {

            $("template").innerHTML = xhrGetRegistration.responseText;

            // JavaScript Funktionen für das Register Template initialisieren.
            addJSforRegister();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetRegistration.open('GET', 'http://localhost:4242/getRegistration');
        xhrGetRegistration.send();

    });

    // EventListener für den Anmelde-Button
    $("loginButton").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Nutzerdaten zu schicken und Home Template anzufordern
        //------------------------------------------------
        // neues XMLHttpRequest anlegen
        var xhrGetHome = new XMLHttpRequest();

        // callback, um Fehlermeldungen als Antwort zu erhalten und diese im html einzusetzen
        xhrGetHome.addEventListener('load', function () {

            // Prüfen, ob Antwort Fehlercode beinhaltet oder nicht
            if (xhrGetHome.responseText.indexOf("Fehler") == 0) {
                // Falls Fehler
                $("errorOutput").innerText = xhrGetHome.responseText;
            } else {
                // Falls kein Fehler, wurde nur Template geliefert
                $("template").innerHTML = xhrGetHome.responseText;
                // JavaScript Funktionen für das Home Template initialisieren.
                addJSforHome();
            }


        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetHome.open('POST', 'http://localhost:4242/home');
        xhrGetHome.send(new FormData($('loginForm')));

    })


};

//#################################
//Home-Template
//#################################

function addJSforHome() {

    // Eventlistener für den Motive-Button
    $("motiveWrapper").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Motive-Seite anzufordern
        //------------------------------------------------
        var xhrGetImageTemplate = new XMLHttpRequest();

        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        xhrGetImageTemplate.addEventListener('load', function () {

            $("template").innerHTML = xhrGetImageTemplate.responseText;
            addJSforImages();
            // JavaScript Funktionen für das nächste Template initialisieren.

        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetImageTemplate.open('GET', 'http://localhost:4242/images');
        xhrGetImageTemplate.send();

    });




};



//#################################
//Images-Template
//#################################

function addJSforImages() {
    // Eventlistener für den Motive-Button
    $("newImgSubmit").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Bild hochzuladen
        //------------------------------------------------
        var xhrPostImage = new XMLHttpRequest();


        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        // xhrPostImage.addEventListener('load', function () {

        //    console.log(responseText);
        //     // $("template").innerHTML = xhrPostImage.responseText;

        //     // JavaScript Funktionen für das nächste Template initialisieren.

        // });

        // Anfrage definieren und mit FormValues absenden
        xhrPostImage.open('POST', 'http://localhost:4242/uploadImage');
        xhrPostImage.send(new FormData($('imageUploadForm')));

    });



    // Eventlistener für den  Sammlung Erstellen-Button
    $("newImgSetSubmit").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Bild hochzuladen
        //------------------------------------------------
        var xhrPostImageSet = new XMLHttpRequest();


        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        // xhrPostImage.addEventListener('load', function () {

        //    console.log(responseText);
        //     // $("template").innerHTML = xhrPostImage.responseText;

        //     // JavaScript Funktionen für das nächste Template initialisieren.

        // });

        // Anfrage definieren und mit FormValues absenden
        xhrPostImageSet.open('POST', 'http://localhost:4242/createSet');
        xhrPostImageSet.send(new FormData($('imageSetForm')));




    });

}


//#################################
//Helferfunktionen
//#################################
// Funktionen, die das Programmieren lediglich einfacher und effizienter machen



// Funktion um ein DOM-Element zu holen
function $(id) {
    return document.getElementById(id);
}

