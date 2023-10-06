package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/yoaz/NFTerryAPI/models"
	response "github.com/yoaz/NFTerryAPI/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo DataBase Connection Vars
//TODO: Check why .env not getting loaded
var connectionString = os.Getenv("MONGO_DB_CONNECTION_STRING")
var dbName = os.Getenv("MONGO_DB_NAME")
var colName = os.Getenv("MONGO_DB_COL_NAME")

// Establish instance of mongo database collection based on documnetation
var Collection *mongo.Collection

// Connect with mongoDB
func init(){
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//Connect
	client, err := mongo.Connect(context.TODO(), clientOption)

	// In case connection failed
	if err != nil {
		log.Fatalf("There was an error with the MongoDB connection %s", err)
	}

	// Conncetion is seccessful
	fmt.Println("MongoDB connection success!!")

	// Collection will hold the collection reference holder based on given db details
	Collection = client.Database(dbName).Collection(colName)

}


/* -------------------- MongoDB Internal Helpers --------------------*/

/* Insert 1 record */
func insertOneNFT (nft models.NFT) {
	inserted, err := Collection.InsertOne(context.Background(), nft)

	// In case inserted operation failed
	if err != nil {
		log.Fatalf("There was an error inserted the data to the DB, %s", err)
	}

	// In case of success
	fmt.Println("Inserted 1 NFT to treasury DB with ID: ", inserted.InsertedID )
}

/* Update one NFT */
func updateOneNFT(nftID string) {
	id, err := primitive.ObjectIDFromHex(nftID)

	// In case of ID string conversion to mongoDB id failed
	if err != nil {
		log.Fatalf("There was an error conveting string ID to mongoDB _id!, %s", err)
	}

	// In case of sucess, find matching ID and update
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"active": false}}

	result, err := Collection.UpdateOne(context.Background(), filter, update)

	// In case of update record failed
	if err != nil {
		log.Fatalf("There was an error updating record id %s, the error: %s", nftID, err)
	}

	// In case of update record success
	fmt.Println("Modified one record success, count: ", result.ModifiedCount)
}

/* Delete one NFT */
func deleteOneNFT(nftID string) {
	id, err := primitive.ObjectIDFromHex(nftID)

	// In case of id convertion failed
	if err != nil{
		log.Fatalf("The id conversation failed %s", nftID)
	}

	// In case of id converation successful
	filter := bson.M{"_id": id}
	deleteCount, err := Collection.DeleteOne(context.Background(), filter)

	fmt.Printf("%d NFT was deleted, ID is: %s ", deleteCount.DeletedCount, nftID )
}

/* Delete all NFTs */
func deleteAllNFT() {
	
	deleteCount, err := Collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	// In case delete all record failed
	if err != nil {
		log.Fatalf("There was an error deleting all records in db, err: %s", err)
	}
	
	// In case delete all record succesed 
	log.Fatalf("%d Records were deleted, these are all records in DB", deleteCount.DeletedCount)
}

/* Get all NFT's */
func getAllNFTs() []primitive.M {
	cur, err := Collection.Find(context.Background(), bson.D{{}})

	// In case fetching all records failed
	if err != nil {
		log.Fatalf("Fetching all records has been failed, this is the err: %s", err)
	}

	// In case fetching succeeded
	var nfts []primitive.M // Will hold all fetched results to return

	// Loop through all records and save in local var
	for cur.Next(context.Background()) {
		var nft bson.M // Will hold current fetched nft
		err := cur.Decode(&nft)

		// In case unmarshalling/decoding failed
		if err != nil {
			log.Fatalf("There was an error trying to decode the record, err:%s", err)
		}

		nfts = append(nfts, nft)
	}

	defer cur.Close(context.Background())

	return nfts
}

func getOneNFT(nftID string) map[string]interface{} {
	id, err := primitive.ObjectIDFromHex(nftID)

	// In case of ID string conversion to mongoDB id failed
	if err != nil {
		log.Fatalf("There was an error conveting string ID to mongoDB _id!, %s", err)
	}

	var nft map[string]interface{}
	e := Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&nft)

	// In case of get record failed
	if e != nil {
		log.Fatalf("There was an error getting record id %s, the error: %s", nftID, e)
	}

	return nft
}
 

/* -------------------- API Externals --------------------*/

/* TODO: Add validate for Request BODY in API calls using Validator github repo */

func Home (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	resp := "<h1>Welcome To NFTerry API<h1>"
	w.Write([]byte(resp))
	return
}

// Get all NFT's
func GetAllNFTs(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	nfts := getAllNFTs()

	// Craft a layout response
	resp := response.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": nfts}}
	json.NewEncoder(w).Encode(resp)
	return
}

// Get one NFT
func GetOneNFT(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	
	params := mux.Vars(r)
	nft := getOneNFT(params["id"])

	// Craft a layout response
	resp := response.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": nft}}
	json.NewEncoder(w).Encode(resp)
	return
}

// Create one NFT
func CreateOneNFT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var nft models.NFT

	//TODO: Validate the request body using validator
	
	err := json.NewDecoder(r.Body).Decode(&nft)

	// In case decoding failed
	if err != nil {
		log.Fatalf("There was an error decoding body request! %s", err)
	}

	// Call mongoDB associated helper
	insertOneNFT(nft)

	// Craft a layout response
	resp := response.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": nft}}
	json.NewEncoder(w).Encode(resp)
	return
}


// Delete one NFT
func DeleteOneNFT(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)
	deleteOneNFT(params["id"])

	// Craft a layout response
	resp := response.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "NFT Seccessfully Deleted!" }}
	json.NewEncoder(w).Encode(resp)
	return
}


// Delete all NFT's
func DeleteAllNFT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	deleteAllNFT()

	// Craft a layout response
	resp := response.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "All NFT's Seccessfully Deleted!"}}
	json.NewEncoder(w).Encode(resp)
	return
}

// Update one NFT
func UpdateOneNFT(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	// Expecting NFT "id" to update just the "Active" bool field
	params := mux.Vars(r)

	// Setting "Active" field False 
	updateOneNFT(params["id"])

	// Craft a layout response
	resp := response.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "NFT active set to false seccessfully!"}}
	json.NewEncoder(w).Encode(resp)
	return
}