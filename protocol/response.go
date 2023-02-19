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

type GetUsernameResponseData struct {
	Username string `bson:"username"`
}

type SetUsernameResponseData struct {
	// TODO
}

type UploadFileMetadataResponseData struct {
	// TODO
}

type UploadFileContentResponseData struct {
	// TODO
}

func SendResponse(conn net.Conn, successful bool, response_data interface{}) {

	raw_response_data, _ := bson.Marshal(response_data)

	raw_response, _ := bson.Marshal(&Response{
		Successful: successful,
		Data:       raw_response_data,
	})

	conn.Write(raw_response)
}
