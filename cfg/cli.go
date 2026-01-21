package cfg

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CliConfig struct {
	ApiBaseUrl string `json:"apiBaseUrl"`
	Token      string `json:"token"`
}

func createDefaultConfig() *CliConfig {
	return &CliConfig{
		ApiBaseUrl: "https://perfin.zivlak.rs/api",
		Token:      "",
	}
}

func GetCliConfigDir() string {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	perfinDir := filepath.Join(cfgDir, "perfin-cli")
	err = os.MkdirAll(perfinDir, 0o755)
	if err != nil {
		panic(err)
	}

	return perfinDir
}

func LoadCliConfig() *CliConfig {
	cfgFile := filepath.Join(GetCliConfigDir(), "config.json")
	return LoadCustomConfig(cfgFile)
}

func LoadCustomConfig(cfgFile string) *CliConfig {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		cfg := createDefaultConfig()
		SaveCustomCliConfig(cfgFile, cfg)
	}

	var cfg *CliConfig
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func SaveCliConfig(cfg *CliConfig) {
	cfgFile := filepath.Join(GetCliConfigDir(), "config.json")
	SaveCustomCliConfig(cfgFile, cfg)
}

func SaveCustomCliConfig(cfgFile string, cfg *CliConfig) {
	data, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(cfgFile, data, 0o644)
	if err != nil {
		panic(err)
	}
}
