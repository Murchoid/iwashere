package models

import (
	"time"
)

type Config struct {
	Project struct {
		Name     string    `json:"name"`
		InitDate time.Time `json:"init_date"`
	} `json:"project"`

	Storage struct {
		Type string `json:"type"` // "sqlite" or "json"
		Path string `json:"path"`
	} `json:"storage"`

	Git struct {
		AutoContext   bool `json:"auto_context"`
		TrackBranches bool `json:"track_branches"`
	} `json:"git"`
	
	Team struct {
		TeamName   string `json:"team_name"`
	} `json:"team"`
}

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.Storage.Type = "json"
	cfg.Storage.Path = "./notes"
	cfg.Project.InitDate = time.Now()
	cfg.Git.AutoContext = true
	cfg.Git.TrackBranches = true
	return cfg
}
