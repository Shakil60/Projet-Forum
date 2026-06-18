package main

// Point d'entree du programme : initialise l'application et lance le serveur HTTP.

import (
	"forum/app"
	"log"
	"net/http"
)

func main() {
	application := app.InitApp()
	defer application.Close()

	address := ":" + application.Port
	log.Printf("Serveur lance : http://localhost%s", address)

	serveErr := http.ListenAndServe(address, application.Router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur - %s", serveErr.Error())
	}
}
