package api

import (
	"os"

	controllers "github.com/yoaz/NFTerryAPI/api/Controllers"
)


var Server = controllers.Server{}

//Run API
func APIRun(){
	//Load envoirmental variables
	controllers.LoadEnvVariables()

	Server.InitDBRouter(os.Getenv("MONGO_DB_CONNECTION_STRING"), os.Getenv("MONGO_DB_NAME"))
	Server.RunServer(":" + os.Getenv("PORT"), os.Getenv("PORT"))
}	