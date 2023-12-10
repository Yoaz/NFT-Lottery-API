package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* -------------------------------- Types ----------------------------------*/

//Will hold data in regards to jackpot winners wallets
type Jackpot struct {
	WinnerAddress string `json:"winneraddress"`
	LotteryDate time.Time	`json:"lotterydate"`
	TreasuryPrize NFT	`json:"treasuryprize"`
}


/* -------------------------------- Helpers ----------------------------------*/

func (j *Jackpot) Prepare(r *http.Request) error {
	// Custom jackpot type to unmarshell time correctly
	type parseType struct {
		WinnerAddress string `json:"winneraddress"`
		LotteryDate string	`json:"lotterydate"`
		TreasuryPrize NFT	`json:"treasuryprize"`
	}
	var res parseType

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return err
	}

	// Parse string date to RFC3339
	println(res.LotteryDate)
	parsedDate, err := time.Parse(time.RFC3339, res.LotteryDate)
	if err != nil {
		return err
	}

	// Parsing success
	j.WinnerAddress = res.WinnerAddress
	j.LotteryDate = parsedDate
	j.TreasuryPrize = res.TreasuryPrize

	println(j)

	return nil
}


/* -------------------------------- DB Actions ----------------------------------*/

// Insert one jackpot to database
func (j *Jackpot) InsertOneJackpot(db *mongo.Database) (*Jackpot, error) {
	var err error

	collection := db.Collection(os.Getenv("MONGO_DB_JACKPOT_COL_NAME"))
	inserted, err := collection.InsertOne(context.Background(), j)

	// In case inserted operation failed
	if err != nil {
		return &Jackpot{}, err
	}

	//In case of success
	fmt.Println("Inserted 1 Jackpot details to Jackpot DB with ID: ", inserted.InsertedID)
	
	return j, nil 
}



// Get one jackpot from DB based on date
func (j *Jackpot) GetJackpotByDate(db *mongo.Database, jackpotDate string) (*Jackpot, error) {
	var err error

	collection := db.Collection(os.Getenv("MONGO_DB_JACKPOT_COL_NAME"))
	// Filter by lottery date
	// There shouold be only ONE jackpot entry for each date (daily lottery)
	// err = j.Prepare(jackpotDate)
	// if err != nil {
	// 	return &Jackpot{}, err
	// }

	filter := bson.M{"lotterydate": j.LotteryDate} 

	err = collection.FindOne(context.Background(), filter).Decode(&j)

	// In case of no such jackpot in the provided date in DB
	if err == mongo.ErrNoDocuments {
		return &Jackpot{}, err
	}

	return j, nil
}


// Get all jackpot history from DB
func (j *Jackpot) GetAllJackpotHistory(db *mongo.Database) (*[]Jackpot, error) {
	var err error

	collection := db.Collection(os.Getenv("MONGO_DB_JACKPOT_COL_NAME"))
	cur, err := collection.Find(context.Background(), bson.M{})

	//In case fetching fails
	if err != nil {
		return &[]Jackpot{}, err
	}

	var jackpots []Jackpot
	for cur.Next(context.Background()) {
		var jackpot Jackpot
		err = cur.Decode(&jackpot)

		// In case unmarshalling/decoding failed
		if err != nil {
			return &[]Jackpot{}, err
		}
		jackpots = append(jackpots, jackpot)
	}

	defer cur.Close(context.Background())

	return &jackpots, nil
}