package config

import (
	"encoding/json"
	"io/ioutil"
)

// GlobalConfig contains config information
type GlobalConfig struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	Download   DownloadConfig   `json:"download"`
	Pagination PaginationConfig `json:"pagination"`
	Cron       CronConfig       `json:"cron"`
}

// ServerConfig contains server information
type ServerConfig struct {
	Address string `json:"address"`
	Mode    string `json:"mode"`
	Secret  string `json:"secret"`
}

// DatabaseConfig contains database information
type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Name     string `json:"name"`
}

// DownloadConfig contains rcp connect information
type DownloadConfig struct {
	Address string `json:"address"`
	Path    string `json:"path"`
}

// PaginationConfig contains rcp connect information
type PaginationConfig struct {
	Limit string `json:"limit"`
	Page  string `json:"page"`
}

// CronConfig contains info about cron jobs
type CronConfig struct {
	Update string `json:"update"`
	Clean  string `json:"clean"`
}

// Config is a point to GlobalConfig
var Config = new(GlobalConfig)

// Init function read and parse config.json
func Init() error {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, Config)
	if err != nil {
		return err
	}

	return nil
}
