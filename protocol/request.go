package protocol

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	REQ_GET_USERNAME = "get_username"
	REQ_SET_USERNAME = "set_username"

	REQ_LIST_FILES    = "list_files"
	REQ_UPLOAD_FILE   = "upload_file"
	REQ_DOWNLOAD_FILE = "download_file"
)

// Generic request packet
type Request struct {
	Type      string            `bson:"type"`       // Type of packet
	PublicKey ed25519.PublicKey `bson:"public_key"` // Sender's public key
	Data      bson.Raw          `bson:"data"`       // Data of packet
}

// Request data to get username by public key
type GetUsernameRequestData struct {
	PublicKey ed25519.PublicKey `bson:"public_key"`
}

// Request data to set username
type SetUsernameRequestData struct {
	Username string `bson:"username"`
}

// Request data with file metadata
type UploadFileMetadataRequestData struct {
	Name      string `bson:"name"`      // File name
	Encrypted bool   `bson:"encrypted"` // Was the file additionally encrypted on the client side?
	Parts     uint32 `bson:"parts"`     // Total parts
}

// Request data with file content
type UploadFileContentRequestData struct {
	Part    uint32 `bson:"part"`    // Number of this part
	Content []byte `bson:"content"` // Content of this part
}
