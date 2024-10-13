package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ImageColPaths []string `yaml:"image_col_paths"`
	Interval      int      `yaml:"interval"`
	Profile       string   `yaml:"profile"`
	Alignment     string   `yaml:"default_alignment"`
	Stretch       string   `yaml:"default_stretch"`
	Opacity       float64  `yaml:"default_opacity"`
}

func ConfigPath() (string, error) {
	e, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Failed to get tbg path to get default config: %s", err.Error())
	}

	configPath := filepath.Join(filepath.Dir(e), "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = NewConfigTemplate(configPath).WriteFile()
		if err != nil {
			return "", fmt.Errorf("Failed to create default config: %s", err.Error())
		}
	}

	return configPath, nil
}

func (c *Config) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config: %s", err)
	}
	return nil
}

func (c *Config) AddPath(toAdd string, configPath string, align *string, stretch *string, opacity *string) error {
	// keep values blank if not set
	editFlags := false
	var toAddFlags string
	if align != nil || stretch != nil || opacity != nil {
		editFlags = true
		blank := "_"
		if align == nil {
			align = &blank
		}
		if stretch == nil {
			stretch = &blank
		}
		if opacity == nil {
			opacity = &blank
		}
		toAddFlags = fmt.Sprintf("%s %s %s", *align, *stretch, *opacity)
	}

	added := make(map[string]struct{}, 0)
	for i, path := range c.ImageColPaths {
		purePath, opts, hasOpts := strings.Cut(path, "|")
		purePath, opts = strings.TrimSpace(purePath), strings.TrimSpace(opts)
		purePath = filepath.ToSlash(purePath)

		if strings.EqualFold(toAdd, purePath) {
			if !editFlags {
				added[fmt.Sprintf("'%s' already exists in config as '%s'", toAdd, path)] = struct{}{}
				break
			}
			// add flags to path without flags
			if !hasOpts {
				toAdd = fmt.Sprintf("%s | %s", toAdd, toAddFlags)
				c.ImageColPaths[i] = toAdd
				added[fmt.Sprintf("added flags to '%s': '%s'", purePath, toAddFlags)] = struct{}{}
				break
			}
			// replace flags of path with flags
			optSlice := strings.Split(opts, " ")
			if len(optSlice) < 3 {
				return fmt.Errorf("'%s' are not valid flags for a path", opts)
			}
			currAlign, currStretch, currOpacity := optSlice[0], optSlice[1], optSlice[2]
			// only replace if not equal and not blank
			changes := 0
			if currAlign != *align && *align != "_" {
				currAlign = *align
				changes++
			}
			if currStretch != *stretch && *stretch != "_" {
				currStretch = *stretch
				changes++
			}
			if currOpacity != *opacity && *opacity != "_" {
				currOpacity = *opacity
				changes++
			}
			// no replacements
			if changes == 0 {
				added["no changes made"] = struct{}{}
				break
			}
			// replaced
			toAddFlags = fmt.Sprintf("%s %s %s", currAlign, currStretch, currOpacity)
			c.ImageColPaths[i] = fmt.Sprintf("%s | %s", purePath, toAddFlags)
			added[fmt.Sprintf("'%s' flags: '%s' --> '%s'", purePath, opts, toAddFlags)] = struct{}{}

		}
	}
	// add new path if no collision and not editing flags
	if len(added) == 0 {
		if toAddFlags != "" {
			toAdd = fmt.Sprintf("%s | %s", toAdd, toAddFlags)
		}
		c.ImageColPaths = append(c.ImageColPaths, toAdd)
		added[toAdd] = struct{}{}
	}

	template := NewConfigTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(c)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to config at %s: %s", configPath, err.Error())
	}

	c.Log(configPath).LogAdded(added)
	return nil
}

func (c *Config) RemovePath(absPath string, configPath string, align *string, stretch *string, opacity *string) error {
	removePath := true
	removeAllFlags := false
	replaceFlags := false

	// if any flags assigned, do not remove path
	if align != nil || stretch != nil || opacity != nil {
		removePath = false
		// specifically, if assigned all, remove all flags
		if align != nil && stretch != nil && opacity != nil {
			removeAllFlags = true
		} else {
			replaceFlags = true
		}
	}
	removed := make(map[string]struct{})
	for i, path := range c.ImageColPaths {
		purePath, opts, hasOpts := strings.Cut(path, "|")
		purePath, opts = strings.TrimSpace(purePath), strings.TrimSpace(opts)
		purePath = filepath.ToSlash(purePath)

		if strings.EqualFold(absPath, purePath) {
			if removePath {
				removed[absPath] = struct{}{}
				c.ImageColPaths = append(c.ImageColPaths[:i], c.ImageColPaths[i+1:]...)

			} else if removeAllFlags {
				if hasOpts {
					removed[fmt.Sprintf("'%s' flags from '%s'", opts, purePath)] = struct{}{}
					c.ImageColPaths[i] = purePath
				}
			} else if replaceFlags {
				if !hasOpts {
					removed[fmt.Sprintf("'%s' does not have any flags set", purePath)] = struct{}{}
					break
				}
				optSlice := strings.Split(opts, " ")
				if len(optSlice) != 3 {
					return fmt.Errorf("invalid options for '%s': '%s'", purePath, opts)
				}
				alignOpt, stretchOpt, opacityOpt := optSlice[0], optSlice[1], optSlice[2]
				alignOpt, stretchOpt, opacityOpt = strings.TrimSpace(alignOpt), strings.TrimSpace(stretchOpt), strings.TrimSpace(opacityOpt)

				// if set to remove, keep blank
				// if not set to remove, keep the current flag
				blank := "_"
				if align == nil {
					align = &alignOpt
				} else {
					align = &blank
				}
				if stretch == nil {
					stretch = &stretchOpt
				} else {
					stretch = &blank
				}
				if opacity == nil {
					opacity = &opacityOpt
				} else {
					opacity = &blank
				}
				// final check if all flags are all blank
				if *align == blank && *stretch == blank && *opacity == blank {
					// remove flags
					c.ImageColPaths[i] = purePath
					removed[fmt.Sprintf("'%s' flags from '%s'", opts, purePath)] = struct{}{}
					break
				}
				// replace flags only if they're different
				newFlags := fmt.Sprintf("%s %s %s", *align, *stretch, *opacity)
				if newFlags != opts {
					c.ImageColPaths[i] = fmt.Sprintf("%s | %s", purePath, newFlags)
					if hasOpts {
						removed[fmt.Sprintf("'%s' flags: '%s' --> '%s'", purePath, opts, newFlags)] = struct{}{}
					}
				}
			}
			break
		}
	}
	template := NewConfigTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(c)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to config at %s: %s", configPath, err.Error())
	}
	if len(removed) == 0 {
		removed["no changes made"] = struct{}{}
	}

	c.Log(configPath).LogRemoved(removed)
	return nil
}

func (c *Config) EditConfig(configPath string, profile *string, interval *string, align *string, stretch *string, opacity *string) error {
	// key:val = old:new
	edited := make(map[string]string, 0)
	if profile != nil {
		edited[c.Profile] = *profile
		c.Profile = *profile
	}
	if interval != nil {
		intervalInt, _ := strconv.Atoi(*interval)
		edited[strconv.Itoa(c.Interval)] = *interval
		c.Interval = intervalInt
	}
	if align != nil {
		edited[c.Alignment] = *align
		c.Alignment = *align
	}
	if stretch != nil {
		edited[c.Stretch] = *stretch
		c.Stretch = *stretch
	}
	if opacity != nil {
		edited[strconv.FormatFloat(c.Opacity, 'f', -1, 64)] = *opacity
		opacityFloat, _ := strconv.ParseFloat(*opacity, 64)
		c.Opacity = opacityFloat
	}

	template := NewConfigTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(c)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to config at %s: %s", configPath, err.Error())
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
	fmt.Printf("|\n%-25s%s\n", "| profile:", c.Profile)
	fmt.Printf("%-25s%s\n", "| interval:", strconv.Itoa(c.Interval))
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

func (c *Config) LogAdded(path map[string]struct{}) *Config {
	fmt.Println("| added: ")
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
