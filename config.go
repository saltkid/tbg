package main

import (
	"log"
)

const defaultConfig = `#------------------------------------------
# this is the default config. you can edit this or use your own config file by doing any of the following:
#  1. run 'tbg config <path/to/config.yaml>'
#  2. edit this file: set use_user_config to true and set user_config to the path to your config file
#------------------------------------------
use_user_config: false
user_config:
image_col_paths:
- ""
interval: 30
target: default
alignment: center
opacity: 0.1
stretch: uniform
#------------------------------------------
# Fields:
#   use_user_config: whether to use the user config set in user_config
#
#   user_config: path to the user config file
#
#   image_col_paths: list of image collection paths
#      note: put directories that contain images, not image filepaths
#
#   interval: time in minutes between each image change
#
#   target: target profile in Windows Terminal (default, list-0, list-1, etc.)
#      see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general for more information
#
#   alignment: image alignment in Windows Terminal (left, center, right)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment for more information
#
#   opacity: image opacity of background images (0-1)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#transparency for more information
#
#   stretch: image stretch in Windows Terminal (uniform, fill)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode for more information
#------------------------------------------
`

type ConfigFile struct {
	UseUserConfig bool     `yaml:"use_user_config"`
	UserConfig    string   `yaml:"user_config"`
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Target        string   `yaml:"target"`
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
