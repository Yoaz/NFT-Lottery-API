package models

import (
	"context"
	"errors"
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

/* -------------------------------- Types ----------------------------------*/

//User models to use the API
type User struct {
	ID 			uint32		`json:"id"`
	User		string		`json:"user"`
	Password	string		`json:"password"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
}



/* -------------------------------- Helpers ----------------------------------*/

//Hash password using buit in golang package bcrypr
func Hash(pwd string)([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

//Verify password hashed againt pure password string
func VerifyPassword(hashedPwd string, pwd string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
}

//Before adding user to DB
func (u *User) BeforeSave() error {
	
	hashedPassword, err := Hash(u.Password)
	
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}	


//Prepare user struct to save in DB
func (u *User) Prepare(){
	u.ID = 0
	u.User = html.EscapeString(strings.TrimSpace(u.User))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}


//Validate the required user's fields
func (u *User) Validate(action string) error{
	switch strings.ToLower(action) {
		//Updating case
		case "update":
			if u.User == "" {
				return errors.New("Required User Name")
			}
			if u.Password == "" {
				return errors.New("Required Password")
			}
			return nil

		//Login Case	
		case "login":
			if u.User == "" {
				return errors.New("Required User Name")
			}
			if u.Password == "" {
				return errors.New("Required Password")
			}
			return nil
			
		//Defualt case
		default:
		if u.User == "" {
			return errors.New("Required User Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		} 
		return nil
	}
}



/* -------------------------------- DB Actions ----------------------------------*/


//Save user to DB
func (u *User) InsertOneUser(db *mongo.Database) (*User, error){
	var err error

	collection := db.Collection(os.Getenv("MONGO_DB_USERS_COL_NAME"))
	inserted, err := collection.InsertOne(context.Background(), u)

	// In case inserted operation failed
	if err != nil {
		return &User{}, err
	}

	// In case of success
	fmt.Println("Inserted 1 User to users DB with ID: ", inserted.InsertedID)

	return u, nil
}

//Update a user
func (u *User) UpdateOneUser(db *mongo.Database, userID string) (int64, error){
	var err error
	collection := db.Collection(os.Getenv("MONGO_DB_USERS_COL_NAME"))

	id, err := primitive.ObjectIDFromHex(userID)

	// In case of ID string conversion to mongoDB id failed
	if err != nil {
		return 0, err
	}

	// In case of sucess, find matching ID and update
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"password": u.Password,
		"updated_at": time.Now(),
	}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	// In case of update record failed
	if err != nil {
		return result.ModifiedCount, err
	}

	return result.ModifiedCount, nil
}	

//Get user by ID
func (u *User) GetUserByID(db *mongo.Database, userID string) (*User, error) {
	var err error
	collection := db.Collection(os.Getenv("MONGO_DB_USERS_COL_NAME"))

	id, err := primitive.ObjectIDFromHex(userID)

	// In case of ID string conversion to mongoDB id failed
	if err != nil {
		return &User{}, err
	}

	//In case of success, find matching ID 
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&u)

	//In case no such userID in DB
	if err == mongo.ErrNoDocuments {
		return &User{}, err
	}

	return u, err
}	


//Get all user
func (u *User) GetAllUsers(db *mongo.Database) (*[]User, error){
	var err error
	collection := db.Collection(os.Getenv("MONGO_DB_USERS_COL_NAME"))
	cur, err := collection.Find(context.Background(), bson.M{})

	// In case fetching all records failed
	if err != nil {
		return &[]User{}, err
	}

	var users []User
	for cur.Next(context.Background()) {
		var user User
		err = cur.Decode(&user)

		// In case unmarshalling/decoding failed
		if err != nil {
			return &[]User{}, err
		}
		users = append(users, user)
	}

	defer cur.Close(context.Background())

	return &users, nil
}

//Delete a user
func (u *User) DeleteUser(db *mongo.Database, userID string) (int64, error){
	var err error
	collection := db.Collection(os.Getenv("MONGO_DB_USERS_COL_NAME"))

	id, err := primitive.ObjectIDFromHex(userID)

	// In case of ID string conversion to mongoDB id failed
	if err != nil {
		return 0, err
	}

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	
	//In case error deleting from db
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

