package types

import "crypto/ed25519"

type User struct {
	PublicKey ed25519.PublicKey `bson:"public_key"`
	Username  string            `bson:"username"`
}
