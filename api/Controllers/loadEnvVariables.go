package controllers

import (
	"log"

	"github.com/joho/godotenv"
)

// Initialize env variables
func LoadEnvVariables(){
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("There was an error to load envoirmental variables, %s", err)
	}
}