package models

/* -------------------------------- Types ----------------------------------*/

// API client access model
type Client struct {
	Name string `json:"name"`
	Password string `json:"pwd"`
}

/* -------------------------------- Helpers ----------------------------------*/

//TODO: Finish JWT auth procedure

// Client method to check credentials
func (c *Client) CheckCredentials(userName, userPassword string) bool{
	return userName == c.Name && userPassword == c.Password
}

