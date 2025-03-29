package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
  # alignment: right # optional (default: center)
  # stretch: fill    # optional (default: uniformToFill)
  # opacity: 0.25    # optional (default: 1.0)
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
port:` + strconv.FormatUint(uint64(DefaultPort), 10) + `
profile:` + DefaultProfile + `
interval:` + strconv.FormatUint(uint64(DefaultInterval), 10) + `

alignment: ` + DefaultAlignment + `
stretch: ` + DefaultStretch + `
opacity: ` + strconv.FormatFloat(float64(DefaultOpacity), 'f', -1, 32) + `
`),
		EndDesc: []byte(`
#------------------------------------------
# Fields:
#   paths: list of image collection paths
#     - path: directory that contain images
#       alignment: (optional) image alignment in Windows Terminal
#                  valid values: topLeft, top, topRight, left, center, right, bottomLeft, bottom, bottomRight
#                  https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment
#
#       opacity:   (optional) image opacity of background images in Windows Terminal
#                  valid values: 0.0 - 1.0 (inclusive)
#                  https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity
#
#       stretch:   (optional) image stretch in Windows Terminal
#                  valid values: fill, none, uniform, uniformToFill
#                  https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode 
#
#   port: port used by tbg server to send POST requests to, to trigger tbg actions such as change image and quit server
#
#   profile: profile profile in Windows Terminal
#      valid values: default, 0, 1, ..., n
#      https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general
#
#   interval: time in minutes between each image change
#------------------------------------------
`),
	}
}
