package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultAlignment string  = "center"
	DefaultInterval  uint16  = 30
	DefaultOpacity   float32 = 1.0
	DefaultPort      uint16  = 9545
	DefaultProfile   string  = "default"
	DefaultStretch   string  = "uniformToFill"
)

type Config struct {
	Interval *uint16      `yaml:"interval,omitempty"`
	Port     *uint16      `yaml:"port,omitempty"`
	Profile  *string      `yaml:"profile,omitempty"`
	Paths    []ImagesPath `yaml:"paths"`
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
    Port: `, cfg.Port, `
    Profile: `, cfg.Profile,
	)
}

// returns the interval if it is set. otherwise, it returns the default
// interval (30)
func (cfg *Config) IntervalOrDefault() uint16 {
	return Option(cfg.Interval).UnwrapOr(DefaultInterval)
}

// returns the port if it is set. otherwise, it returns the default
// port (9545)
func (cfg *Config) PortOrDefault() uint16 {
	return Option(cfg.Port).UnwrapOr(DefaultPort)
}

// returns the profile if it is set. otherwise, it returns the default
// profile ("default")
func (cfg *Config) ProfileOrDefault() string {
	return Option(cfg.Profile).UnwrapOr(DefaultProfile)
}

func ConfigPath() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "", fmt.Errorf("Failed to get config path: LOCALAPPDATA environment variable is not set/does not exist")
	}
	configPath := filepath.Join(localAppData, "tbg", "config.yml")
	configPathDir := filepath.Dir(configPath)
	if _, err := os.Stat(configPathDir); os.IsNotExist(err) {
		os.MkdirAll(configPathDir, os.ModePerm)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = NewConfigTemplate(configPath).WriteFile()
		if err != nil {
			return "", fmt.Errorf("Failed to create default config: %s", err.Error())
		}
	}

	return configPath, nil
}

func (cfg *Config) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config: %s", err)
	}
	return nil
}

func (cfg *Config) AddPath(
	configPath string,
	pathToAdd string,
	cleanPathToAdd string,
	align *string,
	stretch *string,
	opacity *float32,
) error {
	isEditingOptions := align != nil || stretch != nil || opacity != nil
	pathExists := false
	var added *ImagesPath
	var edited *pathEditForLogs
	errors := make([]string, 1)
	for i, path := range cfg.Paths {
		cleanPath, err := NormalizePath(path.Path)
		if err != nil {
			return fmt.Errorf("Failed to normalize path in config %s: %s", path.Path, err)
		}
		if strings.EqualFold(cleanPathToAdd, cleanPath) {
			pathExists = true
			if !isEditingOptions {
				errors = append(errors, fmt.Sprintf("'%s' already exists in config as '%s'", pathToAdd, path.Path))
				break
			}
			finalAlign := getEditedIfChanged(path.Alignment, align)
			finalOpacity := getEditedIfChanged(path.Opacity, opacity)
			finalStretch := getEditedIfChanged(path.Stretch, stretch)
			if finalAlign == nil && finalOpacity == nil && finalStretch == nil {
				break
			}
			edited = &pathEditForLogs{
				old: path,
				new: ImagesPath{
					Path:      pathToAdd,
					Alignment: finalAlign,
					Opacity:   finalOpacity,
					Stretch:   finalStretch,
				},
			}
			cfg.Paths[i] = ImagesPath{
				Path:      pathToAdd,
				Alignment: Option(align).Or(path.Alignment).val,
				Opacity:   Option(opacity).Or(path.Opacity).val,
				Stretch:   Option(stretch).Or(path.Stretch).val,
			}
			break
		}
	}
	if edited == nil && !pathExists {
		addedPath := ImagesPath{
			Path:      pathToAdd,
			Alignment: align,
			Stretch:   stretch,
			Opacity:   opacity,
		}
		cfg.Paths = append(cfg.Paths, addedPath)
		added = &addedPath
	}
	template := NewConfigTemplate(configPath)
	template.Content, _ = yaml.Marshal(cfg)
	if err := template.WriteFile(); err != nil {
		return err
	}
	cfg.Log(configPath).Added(added, edited)
	return nil
}

// helper struct for easier logging changes to a path in the config
type pathEditForLogs struct {
	old ImagesPath
	new ImagesPath
}

// helper function to return the edited value if and only if it is not equal to
// the old value
func getEditedIfChanged[T comparable](old, edited *T) *T {
	if edited != nil {
		if old != nil {
			if *edited != *old {
				return edited
			}
		} else {
			return edited
		}
	}
	return nil
}

func (cfg *Config) RemovePath(
	configPath string,
	pathToRemove string,
	cleanPathToRemove string,
	align bool,
	stretch bool,
	opacity bool,
) error {
	var removed *string
	for i, path := range cfg.Paths {
		cleanPath, err := NormalizePath(path.Path)
		if err != nil {
			return fmt.Errorf("Failed to normalize path in config %s: %s", path.Path, err)
		}
		if strings.EqualFold(cleanPathToRemove, cleanPath) {
			removePath := !align && !stretch && !opacity
			if removePath {
				removed = &pathToRemove
				cfg.Paths = slices.Delete(cfg.Paths, i, i+1)
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
				tmp := fmt.Sprintf("'%s' from '%s'", strings.TrimSuffix(removedFlags, ", "), path.Path)
				removed = &tmp
				cfg.Paths[i] = ImagesPath{
					Path: path.Path,
				}
			}
			break
		}
	}
	template := NewConfigTemplate(configPath)
	var err error
	template.Content, err = yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Failed to marshal yaml contents: %s", err.Error())
	}
	err = template.WriteFile()
	if err != nil {
		return fmt.Errorf("Error writing to config at %s: %s", shrinkHome(configPath), err.Error())
	}
	cfg.Log(configPath).Removed(removed)
	return nil
}

func (cfg *Config) EditConfig(
	configPath string,
	interval *uint16,
	port *uint16,
	profile *string,
) error {
	// key:val = old:new
	edited := make(map[string]string, 0)
	if interval != nil {
		edited[strconv.Itoa(int(cfg.IntervalOrDefault()))] = strconv.Itoa(int(*interval))
		cfg.Interval = interval
	}
	if port != nil {
		edited[strconv.Itoa(int(cfg.PortOrDefault()))] = strconv.Itoa(int(*port))
		cfg.Port = port
	}
	if profile != nil {
		edited[cfg.ProfileOrDefault()] = *profile
		cfg.Profile = profile
	}
	template := NewConfigTemplate(configPath)
	template.Content, _ = yaml.Marshal(cfg)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to config at %s: %s", shrinkHome(configPath), err.Error())
	}
	cfg.Log(configPath).Edited(edited)
	return nil
}

type ImagesPath struct {
	Path      string   `yaml:"path"`
	Alignment *string  `yaml:"alignment,omitempty"`
	Opacity   *float32 `yaml:"opacity,omitempty"`
	Stretch   *string  `yaml:"stretch,omitempty"`
}

func (path *ImagesPath) String() string {
	return fmt.Sprint(`
       Path: `, path.Path, `
       Alignment: `, Option(path.Alignment).UnwrapOr("not set"),
		`; Stretch: `, Option(path.Stretch).UnwrapOr("not set"),
		`; Opacity: `, func() string {
			if path.Opacity == nil {
				return "not set"
			}
			return strconv.FormatFloat(float64(*path.Opacity), 'f', -1, 32)
		}(),
	)
}

func (path *ImagesPath) AlignmentOrDefault() string {
	return Option(path.Alignment).UnwrapOr(DefaultAlignment)
}

func (path *ImagesPath) OpacityOrDefault() float32 {
	return Option(path.Opacity).UnwrapOr(DefaultOpacity)
}

func (path *ImagesPath) StretchOrDefault() string {
	return Option(path.Stretch).UnwrapOr(DefaultStretch)
}

func (path *ImagesPath) Images() ([]string, error) {
	dir, err := NormalizePath(path.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize path %s: %s", path.Path, err)
	}
	images := make([]string, 0)
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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
	fmt.Println(`
------------------------------------------------------------------------------------
| ` + shrinkHome(configPath) + `
------------------------------------------------------------------------------------
| paths:` + func() string {
		var ret strings.Builder
		for _, dir := range cfg.Paths {
			ret.WriteString(`
|              path: ` + dir.Path +
				func() string {
					if dir.Alignment != nil {
						return `
|              - alignment: ` + dir.AlignmentOrDefault()
					}
					return ""
				}() +
				func() string {
					if dir.Opacity != nil {
						return `
|              - opacity: ` + strconv.FormatFloat(float64(dir.OpacityOrDefault()), 'f', -1, 32)
					}
					return ""
				}() +
				func() string {
					if dir.Stretch != nil {
						return `
|              - stretch: ` + dir.StretchOrDefault()
					}
					return ""
				}() + `
|`)
		}
		return ret.String()
	}() + `
| profile:     ` + cfg.ProfileOrDefault() + `
| port:        ` + strconv.FormatUint(uint64(cfg.PortOrDefault()), 10) + `
| interval:    ` + strconv.FormatUint(uint64(cfg.IntervalOrDefault()), 10) + `
------------------------------------------------------------------------------------`)
	return ConfigLogger{}
}

func (log ConfigLogger) Added(added *ImagesPath, edited *pathEditForLogs) ConfigLogger {
	if edited == nil && added == nil {
		fmt.Println("| no changes made")
	} else {
		if added != nil {
			fmt.Println("| added: ")
			fmt.Printf("%-25sadded path: %s\n", "|", added.Path)
			if added.Alignment != nil {
				fmt.Printf("%-25s- alignment: %s\n", "|", *added.Alignment)
			}
			if added.Opacity != nil {
				fmt.Printf("%-25s- opacity: %f\n", "|", *added.Opacity)
			}
			if added.Stretch != nil {
				fmt.Printf("%-25s- stretch: %s\n", "|", *added.Stretch)
			}
		}
		if edited != nil {
			fmt.Println("| edited: ")
			fmt.Printf("%-25sedited path: %s\n", "|", edited.old.Path)
			if edited.new.Alignment != nil {
				if edited.old.Alignment != nil {
					fmt.Printf("%-25s- old alignment: %s\n", "|", *edited.old.Alignment)
				}
				fmt.Printf("%-25s- new alignment: %s\n", "|", *edited.new.Alignment)
			}
			if edited.new.Opacity != nil {
				if edited.old.Opacity != nil {
					fmt.Printf("%-25s- old opacity: %f\n", "|", *edited.old.Opacity)
				}
				fmt.Printf("%-25s- new opacity: %f\n", "|", *edited.new.Opacity)
			}
			if edited.new.Stretch != nil {
				if edited.old.Stretch != nil {
					fmt.Printf("%-25s- old stretch: %s\n", "|", *edited.old.Stretch)
				}
				fmt.Printf("%-25s- new stretch: %s\n", "|", *edited.new.Stretch)
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
	opacity float32,
	stretch string,
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
