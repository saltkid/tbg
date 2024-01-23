package main

import (
	"fmt"
	"os"
	"strconv"
)

type ConfigTemplate struct {
	path         string
	beginDesc    []byte
	yamlContents []byte
	endDesc      []byte
}

func (c *ConfigTemplate) WriteFile() error {
	toWrite := append(append(c.beginDesc, c.yamlContents...), c.endDesc...)
	return os.WriteFile(c.path, toWrite, 0666)
}

func DefaultTemplate(absConfigPath string) *ConfigTemplate {
	return &ConfigTemplate{
		path: absConfigPath,
		beginDesc: []byte(`#------------------------------------------
# this is the default config. you can edit this or use your own config file by doing any of the following:
#  1. run 'tbg config <path/to/config.yaml>'
#  2. edit this file: set use_user_config to true and set user_config to the path to your config file
#------------------------------------------
`),
		yamlContents: []byte(`image_col_paths: []
interval: 30
target: default
alignment: center
opacity: 0.1
stretch: uniform
`),
		endDesc: []byte(`#------------------------------------------
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
#   opacity: image opacity of background images in Windows Terminal (0.0 - 1.0)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity for more information
#
#   stretch: image stretch in Windows Terminal (uniform, fill)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode for more information
#------------------------------------------
`),
	}

}

func UserTemplate(path string) *ConfigTemplate {
	return &ConfigTemplate{
		path: path,
		beginDesc: []byte(`#------------------------------------------
# this is the user config.
# Notes:
#  1. to update the default config's "user_config" field to point to this file
#     - tbg config ` + path + `
#------------------------------------------
`),
		yamlContents: []byte(`image_col_paths:
- ""
interval: 30
target: default
alignment: center
opacity: 0.1
stretch: uniform
`),
		endDesc: []byte(`#------------------------------------------
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
#   opacity: image opacity of background images in Windows Terminal (0.0 - 1.0)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity for more information
#
#   stretch: image stretch in Windows Terminal (uniform, fill)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode for more information
#------------------------------------------
`),
	}
}

type Config interface {
	Log(string)
}

type DefaultConfig struct {
	UseUserConfig bool     `yaml:"use_user_config"`
	UserConfig    string   `yaml:"user_config"`
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Target        string   `yaml:"target"`
	Alignment     string   `yaml:"alignment"`
	Opacity       float64  `yaml:"opacity"`
	Stretch       string   `yaml:"stretch"`
}

func (c *DefaultConfig) Log(configPath string) {
	fmt.Println("------------------------------------------")
	fmt.Println("|", configPath)
	fmt.Println("------------------------------------------")
	fmt.Printf("%-20s%s\n", "| use_user_config:", strconv.FormatBool(c.UseUserConfig))
	fmt.Printf("%-20s%s\n", "| user_config:", c.UserConfig)
	fmt.Println("| image_col_paths:")
	for _, path := range c.ImageColPaths {
		fmt.Printf("%-20s%s\n", "|", path)
	}
	fmt.Printf("|\n%-20s%s\n", "| interval:", strconv.Itoa(c.Interval))
	fmt.Printf("%-20s%s\n", "| target:", c.Target)
	fmt.Printf("%-20s%s\n", "| alignment:", c.Alignment)
	fmt.Printf("%-20s%s\n", "| opacity:", strconv.FormatFloat(c.Opacity, 'f', -1, 64))
	fmt.Printf("%-20s%s\n", "| stretch:", c.Stretch)
	fmt.Println("------------------------------------------")
}

type UserConfig struct {
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Target        string   `yaml:"target"`
	Alignment     string   `yaml:"alignment"`
	Opacity       float64  `yaml:"opacity"`
	Stretch       string   `yaml:"stretch"`
}

func (c *UserConfig) Log(configPath string) {
	fmt.Println("------------------------------------------")
	fmt.Println("|", configPath)
	fmt.Println("------------------------------------------")
	fmt.Println("| image_col_paths:")
	for _, path := range c.ImageColPaths {
		fmt.Printf("%-20s%s\n", "|", path)
	}
	fmt.Printf("|\n%-20s%s\n", "| interval:", strconv.Itoa(c.Interval))
	fmt.Printf("%-20s%s\n", "| target:", c.Target)
	fmt.Printf("%-20s%s\n", "| alignment:", c.Alignment)
	fmt.Printf("%-20s%s\n", "| opacity:", strconv.FormatFloat(c.Opacity, 'f', -1, 64))
	fmt.Printf("%-20s%s\n", "| stretch:", c.Stretch)
	fmt.Println("------------------------------------------")
}
