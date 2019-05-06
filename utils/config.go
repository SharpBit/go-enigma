package utils

import (
	"encoding/json"
	"io/ioutil"
)

// Config : Configuration from config.json
type Config struct {
	Token string `json:"token"`
}

// GetConfig returns config from JSON
func GetConfig() (data Config) {
	file, _ := ioutil.ReadFile("config.json")
	data = Config{}
	_ = json.Unmarshal([]byte(file), &data)
	return
}
