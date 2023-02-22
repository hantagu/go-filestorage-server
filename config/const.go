package config

const (
	CONFIG_FILE_PATH                   = "./config.json"                    // Default config path
	PROTOCOL_PREAMBLE                  = "\xB3\xD6\x4C\xA4\xF6\x9B\x71\xF8" // Protocol preamble
	PROTOCOL_PREAMBLE_SIZE             = 8                                  // Preamble lenght (bytes)
	PROTOCOL_CHUNK_SIZE                = 8 * 1024 * 1024                    //
	PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE = 4                                  // Number of bytes indicating the size of the entire BSON document
	MONGODB_CONTEXT_TIMEOUT            = 3                                  // Default timeout for MongoDB contexts
)
