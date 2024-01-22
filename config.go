package main

import (
// "gopkg.in/yaml.v3"
)

type ConfigFile struct {
	useUserConfig bool     `yaml:"use_user_config"`
	userConfig    string   `yaml:"user_config"`
	imageColPaths []string `yaml:"image_col_paths"`
	target        string   `yaml:"target"`
	interval      int      `yaml:"interval"`
	alignment     string   `yaml:"alignment"`
	opacity       float64  `yaml:"opacity"`
	stretch       string   `yaml:"stretch"`
}
