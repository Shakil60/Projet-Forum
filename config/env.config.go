package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	errLoad := godotenv.Load("./.env")
	if errLoad != nil {
		log.Println("Aucun fichier .env trouve, utilisation des variables d'environnement systeme")
	}
}

func GetEnvWithDefault(key, defaultValue string) string {
	envVar, envExist := os.LookupEnv(key)
	if !envExist {
		return defaultValue
	}
	return envVar
}

func GetRequiredEnv(key string) string {
	envVar, envExist := os.LookupEnv(key)
	if !envExist || envVar == "" {
		log.Fatalf("Erreur configuration - Variable d'environnement manquante : %s", key)
	}
	return envVar
}
