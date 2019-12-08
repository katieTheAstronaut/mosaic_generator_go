// EventListener für den initialen Aufruf der picx-Seite lädt lediglich die spezifischen 
// JS-Funktionen für das Login-Template und die EventListener für den Logout-Button im Header
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
                // LogoutButton darstellen
                $("logout").innerText = "Abmelden";
                addJSforLogoutButton();
                // JavaScript Funktionen für das Home Template initialisieren.
                addJSforHome();
            }


        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetHome.open('POST', 'http://localhost:4242/home');
        xhrGetHome.send(new FormData($('loginForm')));

    })


};



function addJSforLogoutButton() {
    // EventListener für den Logout-Button
    $("logout").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Nutzer auszuloggen und ihn zur Login-Seite zurückzuschicken
        //------------------------------------------------
        // neues XMLHttpRequest anlegen
        var xhrLogout = new XMLHttpRequest();

        // callback, um Fehlermeldungen als Antwort zu erhalten und diese im html einzusetzen
        xhrLogout.addEventListener('load', function () {
            $("logout").innerText = "";
            $("template").innerHTML = xhrLogout.responseText;
            addJSforLogin();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrLogout.open('GET', 'http://localhost:4242/logout');
        xhrLogout.send();

    });
}


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


    // Eventlistener für den Pools-Button
    $("poolWrapper").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Pool-Seite anzufordern
        //------------------------------------------------
        var xhrGetPoolTemplate = new XMLHttpRequest();

        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        xhrGetPoolTemplate.addEventListener('load', function () {
            $("template").innerHTML = xhrGetPoolTemplate.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforPools();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetPoolTemplate.open('GET', 'http://localhost:4242/pools');
        xhrGetPoolTemplate.send();
    });


    // Eventlistener für den Mosaik-Button
    $("mosaicWrapper").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Pool-Seite anzufordern
        //------------------------------------------------
        var xhrGetMosaicTemplate = new XMLHttpRequest();

        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        xhrGetMosaicTemplate.addEventListener('load', function () {
            $("template").innerHTML = xhrGetMosaicTemplate.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforMosaic();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetMosaicTemplate.open('GET', 'http://localhost:4242/mosaic');
        xhrGetMosaicTemplate.send();
    });



    // Eventlistener für den Nutzer löschen
    $("deleteUser").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Pool-Seite anzufordern
        //------------------------------------------------
        var xhrDeleteUser = new XMLHttpRequest();

        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        xhrDeleteUser.addEventListener('load', function () {
            $("template").innerHTML = xhrDeleteUser.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforLogin();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrDeleteUser.open('POST', 'http://localhost:4242/deleteUser');
        xhrDeleteUser.send();
    });

};



//#################################
//Images-Template
//#################################

function addJSforImages() {

    // Eventlistener für den Zurück-Button setzen
    backToMain();

    // Eventlistener für den "Sammlung Erstellen"-Button
    $("newImgSetSubmit").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um neue Sammlung zu erstellen
        //------------------------------------------------
        var xhrPostImageSet = new XMLHttpRequest();

        // Nachdem neue Sammlung erstellt wurde, soll Template neu aufgerufen werden (mit aktualisierten Informationen)
        xhrPostImageSet.addEventListener('load', function () {
            $("template").innerHTML = xhrPostImageSet.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforImages();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrPostImageSet.open('POST', 'http://localhost:4242/createSet');
        xhrPostImageSet.send(new FormData($('imageSetForm')));

    });


    // Eventlistener für alle Sammlungs-Buttons, die zur jeweiligen Sammlungsseite führen
    var buttons = document.getElementsByClassName("setButtons");

    for (i = 0; i < buttons.length; i++) {
        buttons[i].addEventListener("click", function () {

            //------------------------------------------------
            // XMLHttpRequest um Sammlung darzustellen
            //------------------------------------------------
            var xhrShowImageSet = new XMLHttpRequest();

            // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
            xhrShowImageSet.addEventListener('load', function () {
                $("template").innerHTML = xhrShowImageSet.responseText;
                // JavaScript Funktionen für das nächste Template initialisieren.
                addJSforImageSet();
            });
            // Anfrage definieren und Sammlungsname als Query mit absenden
            var setName = this.value; // Wert des aktuell gedrückten Buttons auslesen = Sammlungsname
            xhrShowImageSet.open('POST', `http://localhost:4242/showSet?imgSet=${setName}`);
            xhrShowImageSet.send();

        });
    }


}


//#################################
//ImageSet-Template
//#################################
function addJSforImageSet() {

    // Eventlistener für den Zurück-Button
    $("backToImages").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um neue Sammlung zu erstellen
        //------------------------------------------------
        var xhrBackToImages = new XMLHttpRequest();

        // Nachdem neue Sammlung erstellt wurde, soll Template neu aufgerufen werden (mit aktualisierten Informationen)
        xhrBackToImages.addEventListener('load', function () {
            $("template").innerHTML = xhrBackToImages.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforImages();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrBackToImages.open('GET', 'http://localhost:4242/images');
        xhrBackToImages.send();

    });

    // Eventlistener für den Bild hochladen - Button
    $("newImgSubmit").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Bild hochzuladen
        //------------------------------------------------
        var xhrPostImage = new XMLHttpRequest();


        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        xhrPostImage.addEventListener('load', function () {


            $("template").innerHTML = xhrPostImage.responseText;

            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforImageSet();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrPostImage.open('POST', 'http://localhost:4242/uploadImage');
        xhrPostImage.send(new FormData($('imageUploadForm')));
    });

}



//#################################
//Pools-Template
//#################################
function addJSforPools() {

    // Eventlistener für den Zurück-Button setzen
    backToMain();

    // Eventlistener für den "Pool Erstellen"-Button
    $("newPoolSubmit").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um neue Sammlung zu erstellen
        //------------------------------------------------
        var xhrPostPool = new XMLHttpRequest();

        // Nachdem neue Sammlung erstellt wurde, soll Template neu aufgerufen werden (mit aktualisierten Informationen)
        xhrPostPool.addEventListener('load', function () {
            $("template").innerHTML = xhrPostPool.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforPools();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrPostPool.open('POST', 'http://localhost:4242/createPool');
        xhrPostPool.send(new FormData($('poolForm')));
    });


    // Eventlistener für alle Sammlungs-Buttons, die zur jeweiligen Sammlungsseite führen
    var buttons = document.getElementsByClassName("poolButtons");

    for (i = 0; i < buttons.length; i++) {
        buttons[i].addEventListener("click", function () {

            //------------------------------------------------
            // XMLHttpRequest um Sammlung darzustellen
            //------------------------------------------------
            var xhrShowPool = new XMLHttpRequest();

            // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
            xhrShowPool.addEventListener('load', function () {
                $("template").innerHTML = xhrShowPool.responseText;
                // JavaScript Funktionen für das nächste Template initialisieren.
                addJSforSinglePool();
            });
            // Anfrage definieren und Sammlungsname als Query mit absenden
            var pool = this.value; // Wert des aktuell gedrückten Buttons auslesen = Sammlungsname
            xhrShowPool.open('POST', `http://localhost:4242/showPool?pool=${pool}`);
            xhrShowPool.send();

        });
    }


}



//#################################
//Single Pool-Template
//#################################
function addJSforSinglePool() {

    // Eventlistener für den Zurück-Button
    $("backToPools").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um neue Sammlung zu erstellen
        //------------------------------------------------
        var xhrBackToPools = new XMLHttpRequest();

        // Nachdem neue Sammlung erstellt wurde, soll Template neu aufgerufen werden (mit aktualisierten Informationen)
        xhrBackToPools.addEventListener('load', function () {
            $("template").innerHTML = xhrBackToPools.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforPools();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrBackToPools.open('GET', 'http://localhost:4242/pools');
        xhrBackToPools.send();

    });



    // Eventlistener für den Bild hochladen - Button
    $("newImageSubmit").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Bild hochzuladen
        //------------------------------------------------
        var xhrPostImageToPool = new XMLHttpRequest();


        // callback, um Template als Antwort zu erhalten und diese im html einzusetzen
        xhrPostImageToPool.addEventListener('load', function () {


            $("template").innerHTML = xhrPostImageToPool.responseText;

            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforSinglePool();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrPostImageToPool.open('POST', 'http://localhost:4242/uploadImageToPool');
        xhrPostImageToPool.send(new FormData($('imageUplForm')));

    });


    // EventListener für den Originale-Löschen-Button
    $("deleteOriginals").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um neue Sammlung zu erstellen
        //------------------------------------------------
        var xhrDeleteOriginals = new XMLHttpRequest();

        // Nachdem neue Sammlung erstellt wurde, soll Template neu aufgerufen werden (mit aktualisierten Informationen)
        xhrDeleteOriginals.addEventListener('load', function () {
            $("template").innerHTML = xhrDeleteOriginals.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforSinglePool();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrDeleteOriginals.open('POST', 'http://localhost:4242/deleteOriginals');
        xhrDeleteOriginals.send();

    });

}


//#################################
//Mosaik-Template
//#################################
function addJSforMosaic() {

    // Eventlistener für den Mosaik-Erstellen-Button
    $("mosaicSubmit").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um neue Sammlung zu erstellen
        //------------------------------------------------
        var xhrGetMosaic = new XMLHttpRequest();

        // Nachdem neue Sammlung erstellt wurde, soll Template neu aufgerufen werden (mit aktualisierten Informationen)
        xhrGetMosaic.addEventListener('load', function () {
            $("mosaicPlaceholder").innerHTML = xhrGetMosaic.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforMosaic();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetMosaic.open('POST', 'http://localhost:4242/generateMosaic');
        xhrGetMosaic.send(new FormData($('mosaicForm')));

    });



}




//#################################
//Funktionen für mehrere Templates
//#################################
function backToMain() {
    // Eventlistener für den Zurück-Button, um zur Übersichtsseite zurück zu kommen
    $("backToMain").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um neue Sammlung zu erstellen
        //------------------------------------------------
        var xhrBackToMain = new XMLHttpRequest();

        // Nachdem neue Sammlung erstellt wurde, soll Template neu aufgerufen werden (mit aktualisierten Informationen)
        xhrBackToMain.addEventListener('load', function () {
            $("template").innerHTML = xhrBackToMain.responseText;
            // JavaScript Funktionen für das nächste Template initialisieren.
            addJSforHome();
        });

        // Anfrage definieren und mit FormValues absenden
        xhrBackToMain.open('GET', 'http://localhost:4242/backToHome');
        xhrBackToMain.send();

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

