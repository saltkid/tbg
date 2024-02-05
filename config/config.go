package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Profile       string   `yaml:"profile"`
	Alignment     string   `yaml:"default_alignment"`
	Stretch       string   `yaml:"default_stretch"`
	Opacity       float64  `yaml:"default_opacity"`
}

func (c *Config) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config: %s", err)
	}
	return nil
}

func (c *Config) AddPath(toAdd string, configPath string, align string, stretch string, opacity string) error {
	// set flags after path only if at least one is set
	if align != "" || opacity != "" || stretch != "" {
		if align == "" {
			align = c.Alignment
		}
		if stretch == "" {
			stretch = c.Stretch
		}
		if opacity == "" {
			opacity = strconv.FormatFloat(c.Opacity, 'f', -1, 64)
		}

		isDefaultAlign := strings.EqualFold(align, c.Alignment)
		isDefaultStretch := strings.EqualFold(stretch, c.Stretch)
		isDefaultOpacity := strings.EqualFold(opacity, strconv.FormatFloat(c.Opacity, 'f', -1, 64))
		// only add flags if all are not the default values
		if !(isDefaultAlign && isDefaultStretch && isDefaultOpacity) {
			toAdd = fmt.Sprintf("%s | %s %s %s", toAdd, align, stretch, opacity)
		}
	}

	for _, path := range c.ImageColPaths {
		pureAbsPath, _, _ := strings.Cut(toAdd, "|")
		pureAbsPath = strings.TrimSpace(pureAbsPath)
		purePath, _, _ := strings.Cut(path, "|")
		purePath = strings.TrimSpace(purePath)

		if strings.EqualFold(pureAbsPath, purePath) {
			return fmt.Errorf("%s already in user config", toAdd)
		}
	}
	c.ImageColPaths = append(c.ImageColPaths, toAdd)

	template := DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(c)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to user config: %s", err.Error())
	}

	c.Log(configPath)
	return nil
}

func (c *Config) RemovePath(absPath string, configPath string) error {
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
func (c *Config) EditPath(arg string, configPath string, profile string, interval string, align string, stretch string, opacity string) error {
	// key:val = old:new
	edited := make(map[string]string, 0)

	// edit these two on a config level
	if profile != "" {
		edited[c.Profile] = profile
		c.Profile = profile
	}
	if interval != "" {
		edited[strconv.Itoa(c.Interval)] = interval
		intervalInt, _ := strconv.Atoi(interval)
		c.Interval = intervalInt
	}

	if arg == "fields" {
		// edit the rest of the fields on a config level too
		if align != "" {
			edited[c.Alignment] = align
			c.Alignment = align
		}
		if stretch != "" {
			edited[c.Stretch] = stretch
			c.Stretch = stretch
		}
		if opacity != "" {
			edited[strconv.FormatFloat(c.Opacity, 'f', -1, 64)] = opacity
			opacityFloat, _ := strconv.ParseFloat(opacity, 64)
			c.Opacity = opacityFloat
		}

	} else {
		// edit the rest of the fields on a per path basis
		for i, path := range c.ImageColPaths {
			purePath, opts, hasOpts := strings.Cut(path, "|")
			purePath, opts = strings.TrimSpace(purePath), strings.TrimSpace(opts)

			// if arg is specific path, skip non equal paths
			if arg != "all" && !strings.EqualFold(arg, purePath) {
				continue
			}

			if hasOpts {
				// use options already present if not set
				optSlice := strings.Split(opts, " ")
				if len(optSlice) != 3 {
					return fmt.Errorf("invalid options for %s: %s", purePath, opts)
				}

				currAlign, currStretch, currOpacity := strings.TrimSpace(optSlice[0]), strings.TrimSpace(optSlice[1]), strings.TrimSpace(optSlice[2])
				if align == "" {
					align = currAlign
				}
				if stretch == "" {
					stretch = currStretch
				}
				if opacity == "" {
					opacity = currOpacity
				}
			} else {
				// use default values if not set
				if align == "" {
					align = c.Alignment
				}
				if stretch == "" {
					stretch = c.Stretch
				}
				if opacity == "" {
					opacity = strconv.FormatFloat(c.Opacity, 'f', -1, 64)
				}

			}

			// if all opts are equal to the defaults set in the config, just remove the options
			isDefaultAlign := strings.EqualFold(align, c.Alignment)
			isDefaultStretch := strings.EqualFold(stretch, c.Stretch)
			isDefaultOpacity := strings.EqualFold(opacity, strconv.FormatFloat(c.Opacity, 'f', -1, 64))
			if isDefaultAlign && isDefaultStretch && isDefaultOpacity {
				c.ImageColPaths[i] = purePath // removed opts after | and just kept the path
			} else {
				c.ImageColPaths[i] = fmt.Sprintf("%s | %s %s %s", purePath, align, stretch, opacity)
			}

			// check if path was edited for logging purposes
			if path != c.ImageColPaths[i] {
				edited[path] = c.ImageColPaths[i]
			}

			// stop editing if only editing one path
			if arg != "all" {
				break
			}
		}
	}
	template := DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(c)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to default config: %s", err.Error())
	}
	if len(edited) == 0 {
		edited["no changes made"] = ""
	}
	c.Log(configPath).LogEdited(edited)
	return nil
}

func (c *Config) Log(configPath string) *Config {
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

func (c *Config) LogEdited(editedPaths map[string]string) *Config {
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

func (c *Config) LogRemoved(path map[string]struct{}) *Config {
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

func (c *Config) LogRunSettings(imagePath string, profile string, interval int, align string, stretch string, opacity float64) *Config {
	fmt.Println("| editing", profile, "profile")
	fmt.Println("| image collection:", filepath.Dir(imagePath))
	fmt.Println("|   image:", filepath.Base(imagePath))
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Printf("%-25s%d%s\n", "| change image every: ", interval, " minutes")
	fmt.Printf("%-25s%s\n", "| alignment:", align)
	fmt.Printf("%-25s%s\n", "| stretch:", stretch)
	fmt.Printf("%-25s%f\n", "| opacity:", opacity)
	fmt.Println("------------------------------------------------------------------------------------")

	return c
}
