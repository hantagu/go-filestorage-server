package protocol

const (
	UPLOAD_METADATA = "upload_metadata"
	UPLOAD_DATA     = "upload_data"
)

// File metadata
type UploadMetadata struct {
	Name      string `bson:"name"`      // File name
	Encrypted bool   `bson:"encrypted"` // Was the file additionally encrypted on the client side?
	Parts     uint32 `bson:"parts"`     // Total parts
}

// File
type UploadData struct {
	Part    uint32 `bson:"part"`    // Number of this part
	Content []byte `bson:"content"` // Content of this part
}
