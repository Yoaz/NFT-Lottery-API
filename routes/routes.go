package routes

import (
	"github.com/gorilla/mux"
	"github.com/yoaz/NFTerryAPI/controllers"
)

func Routes() *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/", controllers.Home)
	router.HandleFunc("/api/nfts", controllers.GetAllNFTs).Methods("GET")
	router.HandleFunc("/api/nft", controllers.CreateOneNFT).Methods("POST")
	router.HandleFunc("/api/nft/{id}", controllers.GetOneNFT).Methods("GET")
	router.HandleFunc("/api/nft/{id}", controllers.UpdateOneNFT).Methods("PUT")
	router.HandleFunc("/api/nft/{id}", controllers.DeleteOneNFT).Methods("DELETE")
	router.HandleFunc("/api/nft/{id}", controllers.DeleteOneNFT).Methods("DELETE")
	router.HandleFunc("/api/deleteallnfts", controllers.DeleteAllNFT).Methods("DELETE")


	return router
}