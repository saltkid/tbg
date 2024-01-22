package main

import (
	"log"
)

type ConfigFile struct {
	UseUserConfig bool     `yaml:"use_user_config"`
	UserConfig    string   `yaml:"user_config"`
	ImageColPaths []string `yaml:"image_col_paths"`
	Target        string   `yaml:"target"`
	Interval      int      `yaml:"interval"`
	Alignment     string   `yaml:"alignment"`
	Opacity       float64  `yaml:"opacity"`
	Stretch       string   `yaml:"stretch"`
}

func LogConfig(config ConfigFile, configPath string) {
	log.Println("using config at", configPath)
	log.Println("use user config:", config.UseUserConfig)
	log.Println("user config:", config.UserConfig)
	log.Println("image collection paths:", config.ImageColPaths)
	log.Println("target:", config.Target)
	log.Println("interval:", config.Interval)
	log.Println("alignment:", config.Alignment)
	log.Println("opacity:", config.Opacity)
	log.Println("stretch:", config.Stretch)
}
