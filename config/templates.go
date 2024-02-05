package config

import (
	"os"
)

type ConfigTemplate struct {
	Path         string
	BeginDesc    []byte
	YamlContents []byte
	EndDesc      []byte
}

func (c *ConfigTemplate) WriteFile() error {
	toWrite := append(append(c.BeginDesc, c.YamlContents...), c.EndDesc...)
	return os.WriteFile(c.Path, toWrite, 0666)
}

func DefaultTemplate(absConfigPath string) *ConfigTemplate {
	return &ConfigTemplate{
		Path: absConfigPath,
		BeginDesc: []byte(`#------------------------------------------
# this is the default config. you can edit this or use your own config file by doing any of the following:
#  1. run 'tbg config <path/to/config.yaml>'
#  2. edit this file: set use_user_config to true and set user_config to the path to your config file
#------------------------------------------
`),
		YamlContents: []byte(`user_config:

image_col_paths: []
interval: 30
profile: default

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
`),
		EndDesc: []byte(`#------------------------------------------
# Fields:
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
		Path: path,
		BeginDesc: []byte(`#------------------------------------------
# this is the user config.
# Notes:
#  1. to update the default config's "user_config" field to point to this file
#     - tbg config ` + path + `
#------------------------------------------
`),
		YamlContents: []byte(`image_col_paths: []
interval: 30
profile: default

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
`),
		EndDesc: []byte(`#------------------------------------------
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
