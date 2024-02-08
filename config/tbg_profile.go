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

func (c *TbgProfile) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config: %s", err)
	}
	return nil
}

// creates tbg_profile.yaml if it does not exist
func TbgProfilePath() (string, error) {
	e, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Failed to get executable path to get tbg profile: %s", err.Error())
	}

	profilePath := filepath.Join(filepath.Dir(e), "tbg_profile.yaml")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		err = TbgProfileTemplate(profilePath).WriteFile()
		if err != nil {
			return "", fmt.Errorf("Failed to create default tbg profile: %s", err.Error())
		}
	}
	return profilePath, nil
}

func UsedConfig() (string, error) {
	profilePath, err := TbgProfilePath()
	if err != nil {
		return "", err
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
	// set default used_config to the default config.yaml
	defaultConfig, err := DefaultConfigPath()
	usedConfig := `used_config: ""`
	if err == nil {
		usedConfig = fmt.Sprintf(`used_config: "%s"`, defaultConfig)
	}

	return &ConfigTemplate{
		Path: path,
		BeginDesc: []byte(`#---------------------------------------------
# this is a tbg profile. Whenver tbg is ran, it will
# load this profile to get the currently used config
#
# currently, it only has one field: used_config
# I'll add more if the need arises
#---------------------------------------------
`),
		YamlContents: []byte(usedConfig),
		EndDesc: []byte(`
#---------------------------------------------
# Fields:
#   used_config: path to the config used by tbg
#------------------------------------------
`),
	}
}
