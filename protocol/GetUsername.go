package protocol

import "crypto/ed25519"

const GET_USERNAME = "get_username"

type GetUsername struct {
	Publickey ed25519.PublicKey `bson:"publickey"`
}

type GetUsernameResponse struct {
	Username string `bson:"username"`
}
