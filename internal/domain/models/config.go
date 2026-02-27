package models

import (
	"time"
)

type Config struct {
	Project struct {
		Name     string    `yaml:"name"`
		InitDate time.Time `yaml:"init_date"`
	} `yaml:"project"`

	Storage struct {
		Type string `yaml:"type"` // "sqlite" or "json"
		Path string `yaml:"path"`
	} `yaml:"storage"`

	Git struct {
		AutoContext   bool `yaml:"auto_context"`
		TrackBranches bool `yaml:"track_branches"`
	} `yaml:"git"`
}

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.Storage.Type = "json"
	cfg.Storage.Path = "./notes"
	cfg.Git.AutoContext = true
	cfg.Git.TrackBranches = true
	return cfg
}

// func LoadConfig(filepath string) *Config {

// 	if data, err := os.Open(filepath); err != nil {

// 	}
// }
