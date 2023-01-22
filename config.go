package main

const (
	LISTEN_NETWORK        = "tcp4"                 // Listen on TCP in IPv4
	LISTEN_ADDRESS        = "0.0.0.0:18123"        // Listen on all interfaces
	TLS_CRT_PATH          = "./tls/server-crt.pem" // Path to server's TLS certificate
	TLS_KEY_PATH          = "./tls/server-key.pem" // Path to server's TLS key
	USERDATA_PATH         = "./user_data"          // Path to User Data directory
	USERDATA_DEFAULT_PERM = 0o700                  // Default User Data directory permissions
)

const (
	MONGO_URI = "mongodb://127.0.0.1:27017/"
)

const (
	PROTO_NET_BUFFER_SIZE           = 4096                               // Buffer size
	PROTO_NET_PREAMBLE              = "\xB3\xD6\x4C\xA4\xF6\x9B\x71\xF8" //	Protocol preamble
	PROTO_NET_PREAMBLE_SIZE         = 8                                  // Preamble lenght (bytes)
	PROTO_BSON_DOCUMENT_LENGTH_SIZE = 4                                  //	BSON Document first N bytes for total document lenght
)
