package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func DefaultConfigPath() string {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/config.yaml", filepath.Dir(e))
}

type DefaultConfig struct {
	UserConfig    string   `yaml:"user_config"`
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Profile       string   `yaml:"profile"`
	Alignment     string   `yaml:"default_alignment"`
	Stretch       string   `yaml:"default_stretch"`
	Opacity       float64  `yaml:"default_opacity"`
}

func (c *DefaultConfig) IsDefaultConfig() bool {
	return true
}

func (c *DefaultConfig) IsUserConfig() bool {
	return false
}

func (c *DefaultConfig) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal default config: %s", err)
	}
	return nil
}

func (c *DefaultConfig) AddPath(absPath string, configPath string) error {
	for _, path := range c.ImageColPaths {
		pureAbsPath, _, _ := strings.Cut(absPath, "|")
		pureAbsPath = strings.TrimSpace(pureAbsPath)
		purePath, _, _ := strings.Cut(path, "|")

		if pureAbsPath == purePath {
			return fmt.Errorf("%s already in user config", absPath)
		}
	}
	c.ImageColPaths = append(c.ImageColPaths, absPath)

	template := DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(c)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to user config: %s", err.Error())
	}

	c.Log(configPath)
	return nil
}

func (c *DefaultConfig) RemovePath(absPath string, configPath string) error {
	removed := make(map[string]struct{})
	for i, path := range c.ImageColPaths {
		purePath, _, _ := strings.Cut(path, "|")
		purePath = strings.TrimSpace(purePath)

		if strings.EqualFold(absPath, purePath) {
			removed[absPath] = struct{}{}
			c.ImageColPaths = append(c.ImageColPaths[:i], c.ImageColPaths[i+1:]...)
			break
		}
	}
	template := DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(c)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to default config: %s", err.Error())
	}
	if len(removed) == 0 {
		removed["no changes made"] = struct{}{}
	}

	c.Log(configPath).LogRemoved(removed)
	return nil
}

func (c *DefaultConfig) Log(configPath string) Config {
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("|", configPath)
	fmt.Println("------------------------------------------------------------------------------------")
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
		for old, edited := range editedPaths {
			if old != edited {
				fmt.Printf("%-25s%s\n", "| old:", old)
				fmt.Printf("%-25s%s\n", "| new:", edited)
				fmt.Println("|")
			}
		}
	}
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}

func (c *DefaultConfig) LogRemoved(path map[string]struct{}) Config {
	fmt.Println("| removed: ")
	if _, ok := path["no changes made"]; ok {
		fmt.Println("| no changes made")
	} else {
		for path := range path {
			fmt.Println("|", path)
		}
	}
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}
