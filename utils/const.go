package utils

const (
	PROTO_NET_BUFFER_SIZE           = 4096                               // Buffer size
	PROTO_NET_PREAMBLE              = "\xB3\xD6\x4C\xA4\xF6\x9B\x71\xF8" //	Protocol preamble
	PROTO_NET_PREAMBLE_SIZE         = 8                                  // Preamble lenght (bytes)
	PROTO_BSON_DOCUMENT_LENGTH_SIZE = 4                                  //	BSON Document first N bytes for total document lenght
	CONFIG_FILE_PATH                = "./config.json"
)
