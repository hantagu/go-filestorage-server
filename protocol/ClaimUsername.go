package protocol

const CLAIM_USERNAME = "claim_username"

type ClaimUsername struct {
	Username string `bson:"username"` // Username to be taken
}
