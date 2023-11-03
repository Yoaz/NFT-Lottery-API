package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/yoaz/NFTerryAPI/api/Models"
)

/* ------------------------------------ API External --------------------------------------*/

/* TODO: Add validate for Request BODY in API calls using Validator github repo */

func Home(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	resp := "<h1>Welcome To NFTerry API<h1>"
	w.Write([]byte(resp))
	return
}

// Get all NFT's
func (server *Server) GetAllNFTs(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	
	var resp models.Response
	var nft models.NFT

	nfts, err := nft.GetAllNFTs(server.DB)
	if err != nil {
		resp.BadResponse(w, http.StatusInternalServerError, "failed to fetch users from db", err)
		return
	}

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": nfts})
	return
}

// Get one NFT
func (server *Server) GetOneNFT(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	
	var resp models.Response
	var nft models.NFT

	params := mux.Vars(r)
	nftID := params["id"]

	fetchedNFT, err := nft.GetNFTByID(nftID, server.DB)
	if err != nil {
		resp.BadResponse(w, http.StatusInternalServerError, "couldn't fetch NFT from DB", err)
		return
	}

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": fetchedNFT})
	return
}

// Create one NFT
func (server *Server) CreateOneNFT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var resp models.Response
	var nft models.NFT

	//TODO: Validate the request body using validator
	
	err := json.NewDecoder(r.Body).Decode(&nft)

	// In case decoding failed
	if err != nil {
		log.Fatalf("There was an error decoding body request! %s", err)
	}

	// Call mongoDB associated helper
	insertedNFT, err := nft.InsertOneNFT(server.DB)
	if err != nil{
		resp.BadResponse(w, http.StatusInternalServerError, "failed to insert NFT to db", err)
		return
	}

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": insertedNFT})
	return
}


// Delete one NFT
func (server *Server) DeleteOneNFT(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var resp models.Response
	var nft models.NFT

	params := mux.Vars(r)
	nftID := params["id"]

	delCount, err := nft.DeleteOneNFT(nftID, server.DB)
	if err != nil {
		resp.BadResponse(w, http.StatusInternalServerError, "failed to delete NFT", err)
		return 
	}

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": delCount})
	return
}


// Delete all NFT's
func (server *Server) DeleteAllNFT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var resp models.Response
	var nft models.NFT

	delCount, err := nft.DeleteAllNFT(server.DB)
	if err != nil {
		resp.BadResponse(w, http.StatusInternalServerError, "failed to delete all NFTs", err)
		return 
	}

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": delCount})
	return
}

// Update one NFT
func (server *Server) UpdateOneNFT(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	var resp models.Response
	var nft models.NFT

	// Expecting NFT "id" to update just the "Active" bool field
	params := mux.Vars(r)
	userID := params["id"]
	
	updatedCount, err := nft.UpdateOneNFT(userID, server.DB)
	if err != nil {
		resp.BadResponse(w, http.StatusInternalServerError, "failed to update NFT in db", err)
		return
	}

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": updatedCount})
	return
}