package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Paths     []ImagesPath `yaml:"paths"`
	Interval  uint16       `yaml:"interval"`
	Profile   string       `yaml:"profile"`
	Alignment string       `yaml:"alignment"`
	Stretch   string       `yaml:"stretch"`
	Opacity   float32      `yaml:"opacity"`
}

func (cfg *Config) String() string {
	return fmt.Sprint(`
    Paths: `, func() string {
		ret := ""
		for _, path := range cfg.Paths {
			ret += path.String()
		}
		return ret
	}(), `
    Interval: `, cfg.Interval, `
    Profile: `, cfg.Profile, `
    Alignment: `, cfg.Alignment, `
    Stretch: `, cfg.Stretch, `
    Opacity: `, cfg.Opacity,
	)
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

type ImagesPath struct {
	Path      string   `yaml:"path"`
	Alignment *string  `yaml:"alignment,omitempty"`
	Stretch   *string  `yaml:"stretch,omitempty"`
	Opacity   *float32 `yaml:"opacity,omitempty"`
}

func (path *ImagesPath) String() string {
	return fmt.Sprint(`
       Path: `, path.Path, `
       Alignment: `, Option(path.Alignment).UnwrapOr("not set"), `; Stretch: `, Option(path.Stretch).UnwrapOr("not set"), `; Opacity: `, func() string {
		if path.Opacity == nil {
			return "not set"
		}
		return strconv.FormatFloat(float64(*path.Opacity), 'f', -1, 32)
	}(),
	)
}

func (path *ImagesPath) Images() ([]string, error) {
	dir := path.Path
	images := make([]string, 0)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != filepath.Base(dir) {
			return filepath.SkipDir
		}
		if IsImageFile(d.Name()) {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to walk directory %s: %s", dir, err)
	}
	return images, nil
}

type ConfigLogger struct{}

func (cfg *Config) Log(configPath string) ConfigLogger {
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("|", configPath)
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Println("| paths:")
	for _, dir := range cfg.Paths {
		fmt.Printf("%-25spath: %s\n", "|", dir.Path)
		if dir.Alignment != nil {
			fmt.Printf("%-25s- alignment: %s\n", "|", *dir.Alignment)
		}
		if dir.Stretch != nil {
			fmt.Printf("%-25s- stretch: %s\n", "|", *dir.Stretch)
		}
		if dir.Opacity != nil {
			fmt.Printf("%-25s- opacity: %f\n", "|", *dir.Opacity)
		}
	}
	fmt.Printf("|\n%-25s%s\n", "| profile:", cfg.Profile)
	fmt.Printf("%-25s%d\n", "| interval:", cfg.Interval)
	fmt.Printf("%-25s%s\n", "| alignment:", cfg.Alignment)
	fmt.Printf("%-25s%s\n", "| stretch:", cfg.Stretch)
	fmt.Printf("%-25s%f\n", "| opacity:", cfg.Opacity)
	fmt.Println("------------------------------------------------------------------------------------")
	return ConfigLogger{}
}

func (log ConfigLogger) Added(added *ImagesPath, edited map[ImagesPath]ImagesPath) ConfigLogger {
	noChangesMade := len(edited) == 0
	if noChangesMade {
		fmt.Println("| no changes made")
	} else {
		if added != nil {
			fmt.Println("| added: ")
			fmt.Printf("%-25sadded path: %s\n", "|", added.Path)
			if added.Alignment != nil {
				fmt.Printf("%-25s- alignment: %s\n", "|", *added.Alignment)
			}
			if added.Stretch != nil {
				fmt.Printf("%-25s- stretch: %s\n", "|", *added.Stretch)
			}
			if added.Opacity != nil {
				fmt.Printf("%-25s- opacity: %f\n", "|", *added.Opacity)
			}
		}
		for old, new := range edited {
			fmt.Println("| edited: ")
			fmt.Printf("%-25sedited path: %s\n", "|", added.Path)
			if new.Alignment != nil {
				if old.Alignment != nil {
					fmt.Printf("%-25s- old alignment: %s\n", "|", *old.Alignment)
				}
				fmt.Printf("%-25s- new alignment: %s\n", "|", *new.Alignment)
			}
			if new.Stretch != nil {
				if old.Stretch != nil {
					fmt.Printf("%-25s- old stretch: %s\n", "|", *old.Stretch)
				}
				fmt.Printf("%-25s- new stretch: %s\n", "|", *new.Stretch)
			}
			if new.Opacity != nil {
				if old.Opacity != nil {
					fmt.Printf("%-25s- old opacity: %f\n", "|", *old.Opacity)
				}
				fmt.Printf("%-25s- new opacity: %f\n", "|", *new.Opacity)
			}
		}
	}
	fmt.Println("------------------------------------------------------------------------------------")
	return log
}

func (log ConfigLogger) Removed(removed *string) ConfigLogger {
	noChangesMade := removed == nil
	if noChangesMade {
		fmt.Println("| no changes made")
	} else {
		fmt.Println("| removed: ")
		fmt.Println("|", *removed)
	}
	fmt.Println("------------------------------------------------------------------------------------")
	return log
}

func (log ConfigLogger) Edited(edited map[string]string) ConfigLogger {
	fmt.Println("| edited: ")
	fmt.Println("------------------------------------------------------------------------------------")
	noChangesMade := len(edited) == 0
	if noChangesMade {
		fmt.Println("| no changes made")
	} else {
		for old, edited := range edited {
			if old != edited {
				fmt.Printf("%-25s%s\n", "| old:", old)
				fmt.Printf("%-25s%s\n", "| new:", edited)
				fmt.Println("|")
			}
		}
	}
	fmt.Println("------------------------------------------------------------------------------------")
	return log
}

func (log ConfigLogger) RunSettings(
	imagePath string,
	profile string,
	interval uint16,
	alignment string,
	stretch string,
	opacity float32,
) ConfigLogger {
	fmt.Println("| editing", profile, "profile")
	fmt.Println("| image collection:", filepath.Dir(imagePath))
	fmt.Println("|   image:", filepath.Base(imagePath))
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Printf("%-25s%d%s\n", "| change image every: ", interval, " minutes")
	fmt.Printf("%-25s%s\n", "| alignment:", alignment)
	fmt.Printf("%-25s%s\n", "| stretch:", stretch)
	fmt.Printf("%-25s%f\n", "| opacity:", opacity)
	fmt.Println("------------------------------------------------------------------------------------")
	return log
}
