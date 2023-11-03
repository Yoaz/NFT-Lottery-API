package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/yoaz/NFTerryAPI/api/Models"
)

//Create one user in database
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user models.User
	var res models.Response
	
	err := json.NewDecoder(r.Body).Decode(&user)
	// In case decoding failed
	if err != nil {
		log.Fatalf("There was an error decoding body request! %s", err)
	}

	//Call mongoDB associated helper
	user.Prepare()
	err = user.Validate("")
	//In case of error in validation
	if err != nil {
		log.Fatalf("There was an error validating user data! %s", err)
	}
	//To hash the password
	err = user.BeforeSave()
	//In case of error in validation
	if err != nil {
		log.Fatalf("There was an error encrypting user data! %s", err)
	}

	insertedUser, err := user.InsertOneUser(server.DB)
	//In case of error in inserting user to db
	if err != nil {
		res.BadResponse(w, http.StatusInternalServerError, "failed to insert user to db", err)
		return
	}

	// Craft a layout response
	res.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": insertedUser})
	return
	}	
	
	
//Delete one user from database
func (server *Server) DeleteOneUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user models.User
	var resp models.Response

	params := mux.Vars(r)
	userID := params["id"]

	// //TODO: Take care of validating valid userID
	// uid, err := strconv.ParseUint(userID, 10, 32)
	// if err != nil {
	// 	resp.BadResponse(w, http.StatusBadRequest, "icorrect user id", err)
	// 	return 
	// }
	
	delCount, err := user.DeleteUser(server.DB, userID)
	if err != nil {
		resp.BadResponse(w, http.StatusInternalServerError, "failed to delete user", err)
		return 
	}

	resp.OkResponse(w, http.StatusOK, "user deleted", map[string]interface{}{"data": delCount})
	return
}

//Get one user
func (server *Server) GetOneUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	var user models.User
	var resp models.Response

	params := mux.Vars(r)
	userID := params["id"]

	// //TODO: Take care of validating valid userID
	// uid, err := strconv.ParseUint(userID, 10, 32)
	// if err != nil {
	// 	resp.BadResponse(w, http.StatusBadRequest, "icorrect user id", err)
	// 	return 
	// }

	fetchedUser, err := user.GetUserByID(server.DB, userID)
	if err != nil {
		resp.BadResponse(w, http.StatusInternalServerError, "couldn't fetch user from DB", err)
		return
	}

	resp.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": fetchedUser})
	return
}

//Get All Users
func (server *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	var user models.User
	var res models.Response

	users, err := user.GetAllUsers(server.DB)
	if err != nil {
		res.BadResponse(w, http.StatusInternalServerError, "failed to fetch users from db", err)
		return
	}

	res.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": users})
	return
} 



//Update a user
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-from-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userID := params["id"]
	
	// //TODO: Take care of validating valid userID
	// uid, err := strconv.ParseUint(userID, 10, 32)
	// if err != nil {
	// 	resp.BadResponse(w, http.StatusBadRequest, "icorrect user id", err)
	// 	return 
	// }

	var user models.User
	var res models.Response
		
	updatedCount, err := user.UpdateOneUser(server.DB, userID)
	if err != nil {
		res.BadResponse(w, http.StatusInternalServerError, "failed to update user in db", err)
		return
	}

	res.OkResponse(w, http.StatusOK, "success", map[string]interface{}{"data": updatedCount})
	return
}