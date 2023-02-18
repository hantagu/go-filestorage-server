package protocol

import (
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

// Generic protocol response packet
type Response struct {
	Successful bool     `bson:"successful"` //
	Data       bson.Raw `bson:"data"`       //
}

// Response data with simple string description
type ResponseDescription struct {
	Description string `bson:"description"` // Description
}

func SendDescription(conn net.Conn, successful bool, description string) {

	desc := ResponseDescription{description}
	raw_desc, _ := bson.Marshal(desc)

	response := Response{successful, raw_desc}
	raw_response, _ := bson.Marshal(response)

	conn.Write(raw_response)
}

type GetUsernameResponseData struct {
	Username string `bson:"username"`
}

type SetUsernameResponseData struct {
	// TODO
}
