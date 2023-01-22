package types

import "net"

const UPLOAD = "upload"

type UploadHead struct {
	Name      string `bson:"name,omitempty"`      // File name
	Encrypted bool   `bson:"encrypted,omitempty"` // Was the file additionally encrypted on the client side?
	Parts     uint32 `bson:"parts,omitempty"`     // Total Parts

	Part uint32 `bson:"part,omitempty"` // Number of this part
}

type UploadBody []byte

func HandleUpload(connection net.Conn, packet *Packet) {
	// TODO
}
