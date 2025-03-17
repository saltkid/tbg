package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type ConfigTemplate struct {
	Path         string
	BeginDesc    []byte
	YamlContents []byte
	EndDesc      []byte
}

func (cfg *ConfigTemplate) WriteFile() error {
	toWrite := append(append(cfg.BeginDesc, cfg.YamlContents...), cfg.EndDesc...)
	if err := os.WriteFile(cfg.Path, toWrite, 0666); err != nil {
		return fmt.Errorf("Error writing to config at %s: %s", cfg.Path, err.Error())
	}
	return nil
}

func NewConfigTemplate(path string) *ConfigTemplate {
	// put C:\Users\username\Pictures as an initial value
	userProfile, err := os.UserHomeDir()
	paths := `image_col_paths : []`
	if err == nil {
		picturesDir := filepath.Join(userProfile, "Pictures")
		picturesDir = filepath.ToSlash(picturesDir)
		paths = fmt.Sprintf(`
paths:
- path: %s
  # alignment: right # uncomment this if you want to override the alignment field below
  # stretch: fill    # uncomment this if you want to override the stretch field below
  # opacity: 0.25    # uncomment this if you want to override the opacity field below
`,
			picturesDir)
	}

	return &ConfigTemplate{
		Path: path,
		BeginDesc: []byte(`#------------------------------------------
# this is a tbg config. Whenver tbg is ran, it will load this config file
# and use the fields below to control the behavior of tbg when changing
# background images of Windows Terminal
#------------------------------------------
`),
		YamlContents: []byte(paths + `
profile: default
interval: 30

alignment: center
stretch: uniformToFill
opacity: 1.0
`),
		EndDesc: []byte(`
#------------------------------------------
# Fields:
#   paths: list of image collection paths
#     - path: directory that contain images
#       alignment: optional alignment value applied only to this path.
#                  uses default alignment (see "alignment" field) if not specified
#       stretch:   optional stretch value applied only to this path
#                  uses default stretch (see "stretch" field) if not specified
#       opacity:   poptional opacity value applied only to this path
#                  uses default opacity (see "opacity" field) if not specified
#
#   profile: profile profile in Windows Terminal
#      valid values: default, 0, 1, ..., n
#      https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general
#
#   interval: time in minutes between each image change
#
#   alignment: image alignment in Windows Terminal
#     valid values: topLeft, top, topRight, left, center, right, bottomLeft, bottom, bottomRight
#     https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment
#
#   opacity: image opacity of background images in Windows Terminal
#     valid values: 0.0 - 1.0 (inclusive)
#     https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity
#
#   stretch: image stretch in Windows Terminal
#     valid values: fill, none, uniform, uniformToFill
#     https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode 
#------------------------------------------
`),
	}
}
