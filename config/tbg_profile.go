package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TbgProfile contains the path to the used config. That's it for now
type TbgProfile struct {
	UsedConfig string `yaml:"used_config"`
}

func UsedConfig() (string, error) {
	e, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Failed to get tbg path to get tbg profile: %s", err.Error())
	}

	profilePath := filepath.Join(filepath.Dir(e), "tbg_profile.yaml")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		err = TbgProfileTemplate(profilePath).WriteFile()
		if err != nil {
			return "", fmt.Errorf("Failed to create default tbg profile: %s", err.Error())
		}
	}

	yamlFile, err := os.ReadFile(profilePath)
	if err != nil {
		return "", fmt.Errorf("Failed to read tbg profile: %s", err.Error())
	}

	var profile TbgProfile
	err = yaml.Unmarshal(yamlFile, &profile)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal tbg profile: %s", err.Error())
	}

	if profile.UsedConfig == "" {
		return "", fmt.Errorf("Used config is empty in tbg profile")
	}

	return profile.UsedConfig, nil
}

func TbgProfileTemplate(path string) *ConfigTemplate {
	return &ConfigTemplate{
		Path: path,
		BeginDesc: []byte(`#---------------------------------------------
# this is a tbg profile. Whenver tbg is ran, it will load this profile to get the currently used config
# currently, it only has one field: used_config
# I'll add more if the need arises
#---------------------------------------------
`),
		YamlContents: []byte(`used_config: ""`),
		EndDesc: []byte(`#---------------------------------------------
# Fields:
#   used_config: path to the used config
#------------------------------------------
#`),
	}
}
