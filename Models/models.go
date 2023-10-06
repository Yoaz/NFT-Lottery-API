package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NFT struct {
	ID primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	TokenAddress	string	`json:"tokenaddress,omitempty"`
	TokenID	int64	`json:"tokenid,omitempty"`
	Name string		`json:"name,omitempty"`
	Symol string	`json:"symbol,omitempty"`
	Attributes map[string]interface{} 	`json:"attributes"`
	TokenURI string		`json:"tokenuri,omitempty"`
	OwnerAddress string `json:"owneraddress,omitempty"`
	Active bool		`json:"active,omitempty"`
}