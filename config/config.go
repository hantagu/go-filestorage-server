package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

var Config *config = &config{
	"0.0.0.0:18123",
	"./tls/server-crt.pem",
	"./tls/server-key.pem",
	"./userdata",
	"mongodb://127.0.0.1:27017",
	"filestorage",
	"files",
	"users",
}

type config struct {
	ListenAddress            string `json:"listen_address"`
	TLSCertificatePath       string `json:"tls_cert_path"`
	TLSKeyPath               string `json:"tls_key_path"`
	UserdataPath             string `json:"userdata_path"`
	MongoDB_URI              string `json:"mongodb_uri"`
	MongoDB_DB               string `json:"mongodb_db"`
	MongoDB_Files_Collection string `json:"mongodb_files_collection"`
	MongoDB_Users_Collection string `json:"mongodb_users_collection"`
}

func Init() {

	if raw_cfg, err := os.ReadFile(CONFIG_FILE_PATH); errors.Is(err, os.ErrNotExist) {

		raw_cfg, _ = json.MarshalIndent(Config, "", "    ")

		if err := os.WriteFile(CONFIG_FILE_PATH, raw_cfg, 0o666); err != nil {
			log.Default().Fatalln(err)
		}

	} else if err != nil {
		log.Default().Fatalln(err)
	} else {
		json.Unmarshal(raw_cfg, Config)
	}
}
