package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/yoaz/NFTerryAPI/api/Models"
)

// Insert One Jackpot Details to DB
func (server *Server) CreateOneJackpot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Method", "POST")

	var resp models.Response
	var jackpot models.Jackpot

	//TODO: validate the request body using validator
	
	//Prepare for insertion
	err := jackpot.Prepare(r); 
	// err := json.NewDecoder(r.Body).Decode(&jackpot)

	// In case decoding failed
	if err != nil {
		log.Fatalf("There was an error decoding body request! %s", err)
		return
	}

	// Call mongoDB associated helper
	insertedJackpot, err := jackpot.InsertOneJackpot(server.DB)
	if err != nil{
		resp.BadResponse(w, http.StatusInternalServerError, "failed to insert Jackpot to db", err)
		return
	} 

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": insertedJackpot})
	return
}


// Get one jackpot
func (server *Server) GetOneJackpotByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Method", "GET")

	var resp models.Response
	var jackpot models.Jackpot

	//TODO: validate the request body using validator

	params := mux.Vars(r)
	jackpotDate := params["jackpotdate"]

	// Call mongoDB associated helper
	fetchedJackpot, err := jackpot.GetJackpotByDate(server.DB, jackpotDate)
	if err != nil{
		resp.BadResponse(w, http.StatusInternalServerError, "failed to fetch Jackpot from db", err)
		return
	} 

	// Craft a layout response
	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": fetchedJackpot})
	return
}