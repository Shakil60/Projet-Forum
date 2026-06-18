package config

// Initialisation de la connexion a la base de donnees MySQL.

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Ouvre la connexion MySQL a partir des variables d'environnement et verifie qu'elle repond.
func InitDB() *sql.DB {
	user := GetRequiredEnv("DB_USER")
	pwd := GetEnvWithDefault("DB_PWD", "")
	host := GetRequiredEnv("DB_HOST")
	port := GetRequiredEnv("DB_PORT")
	name := GetRequiredEnv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci", user, pwd, host, port, name)

	dbContext, dbContextErr := sql.Open("mysql", connectionString)
	if dbContextErr != nil {
		log.Fatalf("Erreur connexion base de donnees - Erreur : \n\t %s", dbContextErr.Error())
	}

	pingErr := dbContext.Ping()
	if pingErr != nil {
		dbContext.Close()
		log.Fatalf("Erreur ping base de donnees - Erreur : \n\t %s", pingErr.Error())
	}

	log.Printf("BDD - Connexion reussie")
	return dbContext
}
