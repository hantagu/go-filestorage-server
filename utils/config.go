package utils

import (
	"encoding/json"
	"os"
)

var Config *config = DefaultConfig()

type config struct {
	ListenAddress           string `json:"listen_address"`
	TLSCertificatePath      string `json:"tls_cert_path"`
	TLSKeyPath              string `json:"tls_key_path"`
	UserdataPath            string `json:"userdata_path"`
	MongoDB_URI             string `json:"mongodb_uri"`
	MongoDB_DB              string `json:"mongodb_db"`
	MongoDB_FilesCollection string `json:"mongodb_filescollection"`
	MongoDB_UsersCollection string `json:"mongodb_userscollection"`
}

func DefaultConfig() *config {
	return &config{
		"0.0.0.0:18123",
		"./tls/server-crt.pem",
		"./tls/server-key.pem",
		"./userdata",
		"mongodb://127.0.0.1:27017",
		"filestorage",
		"files",
		"users",
	}
}

func InitConfig() {

	if stat, err := os.Stat(CONFIG_FILE_PATH); err != nil {
		raw_cfg, _ := json.MarshalIndent(Config, "", "    ")
		os.WriteFile(CONFIG_FILE_PATH, raw_cfg, 0o600)
	} else if stat.IsDir() {
		Logger.Fatalf("`%s` is not a file\n", stat.Name())
	}

	if raw_cfg, err := os.ReadFile(CONFIG_FILE_PATH); err != nil {
		Logger.Fatalln(err)
	} else {
		json.Unmarshal(raw_cfg, Config)
	}
}
