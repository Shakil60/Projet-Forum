package config

// Lecture des variables d'environnement et du fichier .env.

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Charge le fichier .env s'il existe.
func LoadEnv() {
	errLoad := godotenv.Load("./.env")
	if errLoad != nil {
		log.Println("Aucun fichier .env trouve, utilisation des variables d'environnement systeme")
	}
}

// Renvoie la variable d'environnement ou une valeur par defaut si absente.
func GetEnvWithDefault(key, defaultValue string) string {
	envVar, envExist := os.LookupEnv(key)
	if !envExist {
		return defaultValue
	}
	return envVar
}

// Renvoie une variable d'environnement obligatoire, ou arrete le programme si elle manque.
func GetRequiredEnv(key string) string {
	envVar, envExist := os.LookupEnv(key)
	if !envExist || envVar == "" {
		log.Fatalf("Erreur configuration - Variable d'environnement manquante : %s", key)
	}
	return envVar
}
