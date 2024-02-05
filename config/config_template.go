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

func NewConfigTemplate(path string) *ConfigTemplate {
	return &ConfigTemplate{
		Path: path,
		BeginDesc: []byte(`#------------------------------------------
# this is a tbg config. Whenver tbg is ran, it will load this config file if it's the config in tbg's profile
# and use the fields below to control the behavior of tbg when changing background images of Windows Terminal
#------------------------------------------
`),
		YamlContents: []byte(`profile: default
interval: 30

image_col_paths: []

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
`),
		EndDesc: []byte(`#------------------------------------------
# Fields:
#   profile: profile profile in Windows Terminal (default, list-0, list-1, etc...)
#      see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general for more information
#
#   interval: time in minutes between each image change
#
#   image_col_paths: list of image collection paths
#      notes:
#        - put directories that contain images, not image filepaths
#        - can override default options for a specific path by putting a | after the path
#          and putting "alignment", "stretch", and "opacity" after the |
#          eg. abs/path/to/images/dir | right uniform 0.1
#
#------------------------------------------
#   Below are default options which can be overriden on a per-path basis by putting a pipe (|)
#   after the path and putting "alignment", "stretch", and "opacity" values after the | in order
#
#   example: abs/path/to/images/dir | right uniform 0.1
#
#   whatever values the values below have, the options after the | will override
#   the values in the default values for that specific path
#------------------------------------------
#
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
