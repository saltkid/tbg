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
	Interval uint16       `yaml:"interval"`
	Port     uint16       `yaml:"port"`
	Profile  string       `yaml:"profile"`
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
	pathExists := false
	var added *ImagesPath
	var edited *pathEditForLogs
	errors := make([]string, 1)
	for i, path := range cfg.Paths {
		cleanPath := filepath.ToSlash(path.Path)
		if strings.EqualFold(pathToAdd, cleanPath) {
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
	template.YamlContents, _ = yaml.Marshal(cfg)
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

func (cfg *Config) EditConfig(
	configPath string,
	interval *uint16,
	port *uint16,
	profile *string,
) error {
	// key:val = old:new
	edited := make(map[string]string, 0)
	if interval != nil {
		edited[strconv.Itoa(int(cfg.Interval))] = strconv.Itoa(int(*interval))
		cfg.Interval = *interval
	}
	if port != nil {
		edited[strconv.Itoa(int(cfg.Port))] = strconv.Itoa(int(*port))
		cfg.Port = *port
	}
	if profile != nil {
		edited[cfg.Profile] = *profile
		cfg.Profile = *profile
	}
	template := NewConfigTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(cfg)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to config at %s: %s", configPath, err.Error())
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
	fmt.Printf("|\n%-25s%d\n", "| port:", cfg.Port)
	fmt.Printf("%-25s%d\n", "| interval:", cfg.Interval)
	fmt.Println("------------------------------------------------------------------------------------")
	return ConfigLogger{}
}

func (log ConfigLogger) Added(added *ImagesPath, edited *pathEditForLogs) ConfigLogger {
	if edited == nil {
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
				fmt.Printf("%-25s- old alignment: %f\n", "|", *edited.old.Opacity)
			}
			fmt.Printf("%-25s- new alignment: %f\n", "|", *edited.new.Opacity)
		}
		if edited.new.Stretch != nil {
			if edited.old.Stretch != nil {
				fmt.Printf("%-25s- old alignment: %s\n", "|", *edited.old.Stretch)
			}
			fmt.Printf("%-25s- new alignment: %s\n", "|", *edited.new.Stretch)
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
