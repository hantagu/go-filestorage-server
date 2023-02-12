package db_types

import "crypto/ed25519"

type User struct {
	PublicKey ed25519.PublicKey `bson:"public_key"` // User's public key
	Username  string            `bson:"username"`   // Unique username associated with public key
}
