package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func InitConfig() (*Config, error) {
	body, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return nil, err
	}
	newConfig := &Config{}
	if err := json.Unmarshal(body, newConfig); err != nil {
		return nil, err
	}
	return newConfig, nil
}

const (
	ArtistURL           string = "https://groupietrackers.herokuapp.com/api/artists"
	RelationURL         string = "https://groupietrackers.herokuapp.com/api/relation"
	IndexTmplPath       string = "ui/html/index.html"
	ArtistTmplPath      string = "ui/html/artist.html"
	ErrorTmplPath       string = "ui/html/error.html"
	SearchErrorTmplPath string = "ui/html/searchError.html"
	ConfigFilePath      string = "./config/config.json"
)
