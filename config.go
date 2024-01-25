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
		yamlContents: []byte(`use_user_config: false
user_config:

image_col_paths: []
interval: 30
profile: default

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
`),
		endDesc: []byte(`#------------------------------------------
# Fields:
#   use_user_config: whether to use the user config set in user_config
#
#   user_config: path to the user config file
#
#   image_col_paths: list of image collection paths
#      notes:
#        - put directories that contain images, not image filepaths
#        - can override default options for a specific path by putting a | after the path
#          and putting "alignment", "stretch", and "opacity" after the |
#          eg. abs/path/to/images/dir | right uniform 0.1
#
#   interval: time in minutes between each image change
#
#   profile: profile profile in Windows Terminal (default, list-0, list-1, etc...)
#      see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general for more information
#
#   ---
#   Below are default options which can be overriden on a per-path basis by putting a | after the path
#   and putting "alignment", "stretch", and "opacity" values after the |
#   ---
#   default_alignment: image alignment in Windows Terminal (left, center, right, etc...)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment for more information
#
#   default_opacity: image opacity of background images in Windows Terminal (0.0 - 1.0)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity for more information
#
#   default_stretch: image stretch in Windows Terminal (uniform, fill, etc...)
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
		yamlContents: []byte(`image_col_paths: []
interval: 30
profile: default

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
`),
		endDesc: []byte(`#------------------------------------------
# Fields:
#   use_user_config: whether to use the user config set in user_config
#
#   user_config: path to the user config file
#
#   image_col_paths: list of image collection paths
#      notes:
#        - put directories that contain images, not image filepaths
#        - can override default options for a specific path by putting a | after the path
#          and putting "alignment", "stretch", and "opacity" after the |
#          eg. abs/path/to/images/dir | right uniform 0.1
#
#   interval: time in minutes between each image change
#
#   profile: profile profile in Windows Terminal (default, list-0, list-1, etc...)
#      see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general for more information
#
#   ---
#   Below are default options which can be overriden on a per-path basis by putting a | after the path
#   and putting "alignment", "stretch", and "opacity" values after the |
#   ---
#   default_alignment: image alignment in Windows Terminal (left, center, right, etc...)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment for more information
#
#   default_opacity: image opacity of background images in Windows Terminal (0.0 - 1.0)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity for more information
#
#   default_stretch: image stretch in Windows Terminal (uniform, fill, etc...)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode for more information
#------------------------------------------
`),
	}
}

type Config interface {
	Log(string) Config
	LogRemoved(string) Config
	LogEdited(map[string]string) Config
}

type DefaultConfig struct {
	UseUserConfig bool     `yaml:"use_user_config"`
	UserConfig    string   `yaml:"user_config"`
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Profile       string   `yaml:"profile"`
	Alignment     string   `yaml:"default_alignment"`
	Stretch       string   `yaml:"default_stretch"`
	Opacity       float64  `yaml:"default_opacity"`
}

func (c *DefaultConfig) Log(configPath string) Config {
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("|", configPath)
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Printf("%-25s%s\n", "| use_user_config:", strconv.FormatBool(c.UseUserConfig))
	fmt.Printf("%-25s%s\n", "| user_config:", c.UserConfig)
	fmt.Println("| image_col_paths:")
	for _, path := range c.ImageColPaths {
		fmt.Printf("%-25s%s\n", "|", path)
	}
	fmt.Printf("|\n%-25s%s\n", "| interval:", strconv.Itoa(c.Interval))
	fmt.Printf("%-25s%s\n", "| profile:", c.Profile)
	fmt.Printf("%-25s%s\n", "| default_alignment:", c.Alignment)
	fmt.Printf("%-25s%s\n", "| default_stretch:", c.Stretch)
	fmt.Printf("%-25s%s\n", "| default_opacity:", strconv.FormatFloat(c.Opacity, 'f', -1, 64))
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}

func (c *DefaultConfig) LogEdited(editedPaths map[string]string) Config {
	fmt.Println("| edited: ")
	fmt.Println("------------------------------------------------------------------------------------")
	if _, ok := editedPaths["no changes made"]; ok {
		fmt.Println("| no changes made")
	} else {
		for old, new := range editedPaths {
			if old != new {
				fmt.Printf("%-25s%s\n", "| old:", old)
				fmt.Printf("%-25s%s\n", "| new:", new)
				fmt.Println("|")
			}
		}
	}
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}

func (c *DefaultConfig) LogRemoved(path string) Config {
	fmt.Printf("%-25s%s\n", "| removed:", path)
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}

type UserConfig struct {
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Profile       string   `yaml:"profile"`
	Alignment     string   `yaml:"default_alignment"`
	Stretch       string   `yaml:"default_stretch"`
	Opacity       float64  `yaml:"default_opacity"`
}

func (c *UserConfig) Log(configPath string) Config {
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("|", configPath)
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("| image_col_paths:")
	for _, path := range c.ImageColPaths {
		fmt.Printf("%-25s%s\n", "|", path)
	}
	fmt.Printf("|\n%-25s%s\n", "| interval:", strconv.Itoa(c.Interval))
	fmt.Printf("%-25s%s\n", "| profile:", c.Profile)
	fmt.Printf("%-25s%s\n", "| default_alignment:", c.Alignment)
	fmt.Printf("%-25s%s\n", "| default_stretch:", c.Stretch)
	fmt.Printf("%-25s%s\n", "| default_opacity:", strconv.FormatFloat(c.Opacity, 'f', -1, 64))
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}

func (c *UserConfig) LogEdited(editedPaths map[string]string) Config {
	fmt.Println("| edited: ")
	fmt.Println("------------------------------------------------------------------------------------")
	if _, ok := editedPaths["no changes made"]; ok {
		fmt.Println("| no changes made")
	} else {
		for old, new := range editedPaths {
			if old != new {
				fmt.Printf("%-25s%s\n", "| old:", old)
				fmt.Printf("%-25s%s\n", "| new:", new)
				fmt.Println("|")
			}
		}
	}
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}

func (c *UserConfig) LogRemoved(path string) Config {
	fmt.Printf("%-25s%s\n", "| removed:", path)
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}
