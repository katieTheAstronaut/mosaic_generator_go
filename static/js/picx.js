window.addEventListener("load", function () {

    //#################################
    //LOGIN-Template
    //#################################

    // Ausgelöst, wenn auf den Registrier-Button geklickt wird
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

    })


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
            $("errorOutput").innerText = xhrPostUserInfo.responseText;
        });

        // Anfrage definieren und mit FormValues absenden
        xhrPostUserInfo.open('POST', 'http://localhost:4242/postUserLoginData');
        xhrPostUserInfo.send(new FormData($('loginForm')));

    });


    // EventListener für den "Zurück Zum Login"-Button, der wieder zur Anmeldeseite wechselt
    $("changeToLogin").addEventListener("click", function () {

        //------------------------------------------------
        // XMLHttpRequest um Nutzerdaten zu schicken 
        //------------------------------------------------
        var xhrGetLogin = new XMLHttpRequest();

        xhrGetLogin.addEventListener('load', function () {
            $("template").innerHTML = xhrGetLogin.responseText;
        });

        // Anfrage definieren und mit FormValues absenden
        xhrGetLogin.open('GET', 'http://localhost:4242/getLogin');
        xhrGetLogin.send();

    })



};

/* ----------------------- Helferfunktionen -----------------------*/
// Funktionen, die das Programmieren lediglich einfacher und effizienter machen




function $(id) {
    return document.getElementById(id);
}