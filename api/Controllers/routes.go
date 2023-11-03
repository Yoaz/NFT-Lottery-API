package controllers

//Initializing server routes
func (server *Server) InitRoutes() { 

	//NFT Routes
	server.Router.HandleFunc("/", Home)
	server.Router.HandleFunc("/api/nfts", server.GetAllNFTs).Methods("GET")
	server.Router.HandleFunc("/api/nft", server.CreateOneNFT).Methods("POST")
	server.Router.HandleFunc("/api/nft/{id}", server.GetOneNFT).Methods("GET")
	server.Router.HandleFunc("/api/nft/{id}", server.UpdateOneNFT).Methods("PUT")
	server.Router.HandleFunc("/api/nft/{id}", server.DeleteOneNFT).Methods("DELETE")
	server.Router.HandleFunc("/api/nft/{id}", server.DeleteOneNFT).Methods("DELETE")
	server.Router.HandleFunc("/api/deleteallnfts", server.DeleteAllNFT).Methods("DELETE") 
	
	//User Routes
	server.Router.HandleFunc("/api/user", server.CreateUser).Methods("POST")
	server.Router.HandleFunc("/api/user/{id}", server.DeleteOneUser).Methods("DELETE")
	server.Router.HandleFunc("/api/user/{id}", server.GetOneUser).Methods("GET")
	server.Router.HandleFunc("/api/users", server.GetAllUsers).Methods("GET")
}