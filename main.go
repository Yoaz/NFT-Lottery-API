/*                                                                                                           *\
***************************************************************************************************************
AUTHOR: Yoaz Sh.
ASSIGNMENT: NFT Lottery Application Programming Interface (API)
VERSION: 1.0
LICENSE: Use For Educational ONLY
PUBLISHED: 04.10.2023
CONTACT ME AT: https://yoaz.info
***************************************************************************************************************
*/

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/yoaz/NFTerryAPI/routes"
)

func Init(){
	if err := godotenv.Load(".env"); err != nil{
		log.Fatalf("There was an error loading the envoirmental variables, error: %s", err)
	}
}

func main(){
	fmt.Println("NFTerry API")
	Init()

	// Getting the router
	r := routes.Routes()
	
	// Listening to server
	fmt.Println("Server is getting started...")
	http.ListenAndServe(":4000", r)
	
}