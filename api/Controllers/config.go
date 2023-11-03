package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Server type to pass around contains mongoDB & mux Router instances
type Server struct {
	DB *mongo.Database
	Router *mux.Router
}

//Initialize Server struct's MongoDB database & mux Router instances
func (server *Server) InitDBRouter(connectionString, dbName string){
	//Context with timout declaration
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//Connect
	client, err := mongo.Connect(ctx, clientOption)

	// In case connection failed
	if err != nil {
		log.Fatalf("There was an error with the MongoDB connection %s", err)
	}

	//Based on Mongo documentation recommedation
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
	
	// Call Ping to verify that the deployment is up and the Client was
	// configured successfully. As mentioned in the Ping documentation, this
	// reduces application resiliency as the server may be temporarily
	// unavailable when Ping is called.
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	//Assign values to type fields
	server.Router = mux.NewRouter()
	server.DB = client.Database(dbName)
	server.InitRoutes()

	// Conncetion is seccessful
	fmt.Println("MongoDB connection success!!")
}


//Run server
func (server *Server) RunServer(addr , port string) {
	fmt.Printf("Listening to port %s\n", port)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
