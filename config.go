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
		return "", fmt.Errorf("Failed to get tbg executable path to get default config: %s", err.Error())
	}

	configPath := filepath.Join(filepath.Dir(e), ".tbg.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = NewConfigTemplate(configPath).WriteFile()
		if err != nil {
			return "", fmt.Errorf("Failed to create default config: %s", err.Error())
		}
	}

	return configPath, nil
}

func (cfg *Config) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config: %s", err)
	}
	return nil
}

func (cfg *Config) AddPath(
	configPath string,
	pathToAdd string,
	align *string,
	stretch *string,
	opacity *float32,
) error {
	isEditingOptions := align != nil || stretch != nil || opacity != nil
	var added *ImagesPath
	// where key == old, val == new
	// didn't use make so i can check against nil
	var edited map[ImagesPath]ImagesPath
	errors := make([]string, 1)
	for i, path := range cfg.Paths {
		cleanPath := filepath.ToSlash(path.Path)
		if strings.EqualFold(pathToAdd, cleanPath) {
			if !isEditingOptions {
				errors = append(errors, fmt.Sprintf("'%s' already exists in config as '%s'", pathToAdd, path.Path))
				break
			}
			dirToAdd := ImagesPath{
				Path:      pathToAdd,
				Alignment: Option(align).Or(path.Alignment).Pointer(),
				Stretch:   Option(stretch).Or(path.Stretch).Pointer(),
				Opacity:   Option(opacity).Or(path.Opacity).Pointer(),
			}
			cfg.Paths[i] = dirToAdd
			edited = map[ImagesPath]ImagesPath{
				path: dirToAdd,
			}
			break
		}
	}
	noChangesMade := len(edited) == 0
	if noChangesMade {
		dirToAdd := ImagesPath{
			Path:      pathToAdd,
			Alignment: align,
			Stretch:   stretch,
			Opacity:   opacity,
		}
		cfg.Paths = append(cfg.Paths, dirToAdd)
		added = &dirToAdd
	}
	template := NewConfigTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(cfg)
	if err := template.WriteFile(); err != nil {
		return err
	}
	cfg.Log(configPath).Added(added, edited)
	return nil
}

func (cfg *Config) RemovePath(
	configPath string,
	pathToRemove string,
	align bool,
	stretch bool,
	opacity bool,
) error {
	var removed *string
	for i, path := range cfg.Paths {
		cleanPath := filepath.ToSlash(path.Path)
		if strings.EqualFold(pathToRemove, cleanPath) {
			removePath := !align && !stretch && !opacity
			if removePath {
				removed = &pathToRemove
				cfg.Paths = append(cfg.Paths[:i], cfg.Paths[i+1:]...)
			} else {
				removedFlags := ""
				if align {
					removedFlags += AlignmentFlag.String() + ", "
				}
				if stretch {
					removedFlags += StretchFlag.String() + ", "
				}
				if opacity {
					removedFlags += OpacityFlag.String() + ", "
				}
				tmp := fmt.Sprintf("'%s' from '%s'", strings.TrimSuffix(removedFlags, ", "), cleanPath)
				removed = &tmp
				cfg.Paths[i] = ImagesPath{
					Path: cleanPath,
				}
			}
			break
		}
	}
	template := NewConfigTemplate(configPath)
	var err error
	template.YamlContents, err = yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Failed to marshal yaml contents: %s", err.Error())
	}
	err = template.WriteFile()
	if err != nil {
		return fmt.Errorf("Error writing to config at %s: %s", configPath, err.Error())
	}
	cfg.Log(configPath).Removed(removed)
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
