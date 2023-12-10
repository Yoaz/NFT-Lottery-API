package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/* -------------------------------- Types ----------------------------------*/

// NFT model, contains all relevant information to store in project db
type NFT struct {
	ID primitive.ObjectID				`json:"_id,omitempty" bson:"_id,omitempty"`
	TokenAddress	string				`json:"tokenaddress,omitempty"`
	TokenID	int64						`json:"tokenid,omitempty"`
	Name string							`json:"name,omitempty"`
	Symol string						`json:"symbol,omitempty"`
	Attributes map[string]interface{} 	`json:"attributes"`
	TokenURI string						`json:"tokenuri,omitempty"`
	OwnerAddress string 				`json:"owneraddress,omitempty"`
	Active bool							`json:"active,omitempty"`
	DateAdded time.Time					`json:"dateadded"`
}

/* -------------------------------- Helpers ----------------------------------*/

//Prepare NFT struct for insertion
func (nft *NFT) Prepare() {
	nft.DateAdded = time.Now()
}



/* -------------------------------- DB Actions ----------------------------------*/

/* Insert 1 record */
func (nft *NFT) InsertOneNFT(db *mongo.Database) (*NFT, error) {
	collection := db.Collection(os.Getenv("MONGO_DB_NFT_COL_NAME"))

	//Prepare NFT struct
	nft.Prepare()

	//Calling mongoDB insert one action
	inserted, err := collection.InsertOne(context.Background(), nft)

	// In case inserted operation failed
	if err != nil {
		return &NFT{}, err
	}

	// In case of success
	fmt.Println("Inserted 1 NFT to treasury DB with ID: ", inserted.InsertedID)

	return nft, nil
}


/* Update one NFT */
func (nft *NFT) UpdateOneNFT(nftID string, db *mongo.Database) (int64, error) {
	var err error

	collection := db.Collection(os.Getenv("MONGO_DB_NFT_COL_NAME"))
	id, err := primitive.ObjectIDFromHex(nftID)

	// In case of ID string conversion to mongoDB id failed
	if err != nil {
		return 0, err
	}

	// In case of sucess, find matching ID and update
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"active": false}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	// In case of update record failed
	if err != nil {
		return result.ModifiedCount, err	
	}

	// In case of update record success
	fmt.Println("Modified one record success, count: ", result.ModifiedCount)
	
	return result.ModifiedCount, nil
}

/* Delete one NFT */
func (nft *NFT) DeleteOneNFT(nftID string, db *mongo.Database) (int64, error) {
	var err error

	collection := db.Collection(os.Getenv("MONGO_DB_NFT_COL_NAME"))
	id, err := primitive.ObjectIDFromHex(nftID)

	// In case of id convertion failed
	if err != nil{
		return 0, err
	}

	// In case of id converation successful
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)

	//In case error deleting from db
	if err != nil{
		return 0, err
	}

	fmt.Printf("%d NFT was deleted, ID is: %s ", deleteCount.DeletedCount, nftID)

	return deleteCount.DeletedCount, nil
}

/* Delete all NFTs */
func (nft *NFT) DeleteAllNFT(db *mongo.Database) (int64, error){
	var err error

	collection := db.Collection(os.Getenv("MONGO_DB_NFT_COL_NAME"))
	deleteCount, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	// In case delete all record failed
	if err != nil {
		log.Fatalf("There was an error deleting all records in db, err: %s", err)
	}
	
	// In case delete all record succesed 
	log.Fatalf("%d Records were deleted, these are all records in DB", deleteCount.DeletedCount)

	return deleteCount.DeletedCount, err
}

/* Get all NFT's */
func (nft *NFT) GetAllNFTs(db *mongo.Database) ([]NFT, error) {
	collection := db.Collection(os.Getenv("MONGO_DB_NFT_COL_NAME"))
	cur, err := collection.Find(context.Background(), bson.D{{}})

	// In case fetching all records failed
	if err != nil {
		return []NFT{}, err
	}

	// In case fetching succeeded
	var nfts []NFT // Will hold all fetched results to return

	// Loop through all records and save in local var
	for cur.Next(context.Background()) {
		var nft NFT// Will hold current fetched nft
		err := cur.Decode(&nft)

		// In case unmarshalling/decoding failed
		if err != nil {
			log.Fatalf("There was an error trying to decode the record, err: %s", err)
		}

		nfts = append(nfts, nft)
	}

	defer cur.Close(context.Background())

	return nfts, nil
}


//Get one NFT by provided ID
func (nft *NFT) GetNFTByID(nftID string, db *mongo.Database) (*NFT, error) {
	var err error
	collection := db.Collection(os.Getenv("MONGO_DB_NFT_COL_NAME"))
	id, err := primitive.ObjectIDFromHex(nftID)

	// In case of ID string conversion to mongoDB id failed
	if err != nil {
		return &NFT{}, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&nft)

	// In case of get record failed
	if err != nil {
		return &NFT{}, err
	}

	return nft, err
}
 