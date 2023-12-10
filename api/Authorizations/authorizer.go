package authorizations

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

// Authorziation model, will be used to handle JWT authorization proccess to protect private endpoints
type AuthHandler struct{
	key []byte 
}

// Custom claim to issue with eveery request, embeds jwt standart claim
type CustomClaims struct{
	Username string		`json:"username"`
	jwt.StandardClaims
}


// Attached method for the AuthHandler to generate token based on client request
func (h *AuthHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	// Check user credentials 
	// Issue a token using HS256 symetric crypotgraphic method
	req := struct {
		User, Pwd string
	}{}
	
	// In case decoding body request fails
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatalf("There is an error decoding body request, err: %s", err)
	}

	// Check client credentials
}


//Extract token from request
func (h *AuthHandler) ExtractToken(r *http.Request) string {
	params := mux.Vars(r)
	token := params["token"]
	
	//Return extracted token if not empty
	if token != " " {
		return token
	}

	

}
