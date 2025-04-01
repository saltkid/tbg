package main

import (
	"errors"
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
	DefaultInterval  uint16  = 30 * 60
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

// Common config initialization for all commands accepting --config flag.
//
// Reads the config file at the given path and validates it.
// If passed in a nil config path, the default one will be used instead.
//
// returns: config instance, config path, error
func ConfigInit(optConfigPath *string) (*Config, string, error) {
	var configPath string
	if optConfigPath != nil {
		configPath = *optConfigPath
	} else {
		defaultPath, err := ConfigPath()
		if err != nil {
			return nil, "", err
		}
		configPath = defaultPath
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to read config at %s: %s", shrinkHome(configPath), err)
	}
	config := new(Config)
	err = config.Unmarshal(yamlFile)
	if err != nil {
		return nil, "", err
	}
	errs := config.Validate()
	if len(errs) != 0 {
		var errMsg strings.Builder
		fmt.Fprintln(&errMsg, Decorate("CONFIG ERRORS").Bold(), "\n",
			Decorate("Please resolve the ff:").Italic(),
		)
		for _, err := range errs {
			fmt.Fprintln(&errMsg, " ", err)
		}
		return nil, "", fmt.Errorf(errMsg.String())
	}
	return config, configPath, nil
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

// helper struct for AddPath for easier logging changes to a path in the config
type pathEdit struct {
	old ImagesPath
	new ImagesPath
}

// helper function for AddPath to return the edited value if and only if it is
// not equal to the old value
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

func (cfg *Config) AddPath(
	configPath string,
	pathToAdd string,
	cleanPathToAdd string,
	align *string,
	opacity *float32,
	stretch *string,
) error {
	isEditingOptions := align != nil || stretch != nil || opacity != nil
	pathExists := false
	var added *ImagesPath
	var edited *pathEdit
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
			edited = &pathEdit{
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

func (cfg *Config) RemovePath(
	configPath string,
	pathToRemove string,
	cleanPathToRemove string,
	align bool,
	opacity bool,
	stretch bool,
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
				var removedFlags strings.Builder
				if align && path.Alignment != nil {
					fmt.Fprint(&removedFlags, AlignmentFlag.String(), ", ")
					cfg.Paths[i].Alignment = nil
				}
				if opacity && path.Opacity != nil {
					fmt.Fprint(&removedFlags, OpacityFlag.String(), ", ")
					cfg.Paths[i].Opacity = nil
				}
				if stretch && path.Stretch != nil {
					fmt.Fprint(&removedFlags, StretchFlag.String(), ", ")
					cfg.Paths[i].Stretch = nil
				}
				if removedFlags.Len() != 0 {
					tmp := fmt.Sprintf("'%s' from '%s'", strings.TrimSuffix(removedFlags.String(), ", "), path.Path)
					removed = &tmp
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

// helper struct for EditConfig that contains edits info when changing config
// through the `config` command
type configEdits struct {
	title string
	old   string
	new   string
}

func (cfg *Config) EditConfig(
	configPath string,
	interval *uint16,
	port *uint16,
	profile *string,
) error {
	edits := make([]configEdits, 0)
	if interval != nil {
		if cfg.IntervalOrDefault() != *interval {
			edits = append(edits,
				configEdits{
					title: "interval",
					old:   strconv.FormatUint(uint64(cfg.IntervalOrDefault()), 10),
					new:   strconv.FormatUint(uint64(*interval), 10),
				},
			)
			cfg.Interval = interval
		}
	}
	if port != nil {
		if cfg.PortOrDefault() != *port {
			edits = append(edits,
				configEdits{
					title: "port",
					old:   strconv.FormatUint(uint64(cfg.PortOrDefault()), 10),
					new:   strconv.FormatUint(uint64(*port), 10),
				},
			)
			cfg.Port = port
		}
	}
	if profile != nil {
		if cfg.ProfileOrDefault() != *profile {
			edits = append(edits,
				configEdits{
					title: "profile",
					old:   cfg.ProfileOrDefault(),
					new:   *profile,
				},
			)
			cfg.Profile = profile
		}
	}
	template := NewConfigTemplate(configPath)
	template.Content, _ = yaml.Marshal(cfg)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf(
			"error writing to config at %s: %s",
			shrinkHome(configPath),
			err.Error(),
		)
	}
	cfg.Log(configPath).Edited(edits)
	return nil
}

// always initializes the returned error messages so no need to check against
// nil
func (cfg *Config) Validate() []error {
	errs := make([]error, 0)
	if len(cfg.Paths) == 0 {
		errs = append(errs, errors.New("paths: must have at least one path entry"))
	}
	leftPad := "\n           "
	for i, path := range cfg.Paths {
		var errStr strings.Builder
		// validate if path exists
		absPath, err := NormalizePath(path.Path)
		if err != nil {
			fmt.Fprint(&errStr,
				"path ", i+1,
				" (", filepath.Join("..", filepath.Base(path.Path)), ")",
				leftPad, err,
			)
			errs = append(errs, errors.New(errStr.String()))
			errStr.Reset()
		} else if _, err = os.Stat(absPath); os.IsNotExist(err) {
			fmt.Fprint(&errStr,
				"path ", i+1,
				" (", filepath.Join("..", filepath.Base(path.Path)), ")",
				leftPad, path.Path, "does not exist",
			)
			errs = append(errs, errors.New(errStr.String()))
			errStr.Reset()
		}
		// validate path alignment if set
		alignment := path.AlignmentOrDefault()
		if _, err = ValidateAlignment(&alignment); err != nil {
			fmt.Fprint(&errStr,
				"path ", i+1, " alignment",
				" (", filepath.Join("..", filepath.Base(path.Path)), ")",
				leftPad, strings.ReplaceAll(err.Error(), "\n", leftPad),
				"\n",
			)
			errs = append(errs, errors.New(errStr.String()))
			errStr.Reset()
		}
		// validate path opacity if set
		opacity := strconv.FormatFloat(float64(path.OpacityOrDefault()), 'f', -1, 32)
		if _, err = ValidateOpacity(&opacity); err != nil {
			fmt.Fprint(&errStr,
				"path ", i+1, " opacity",
				" (", filepath.Join("..", filepath.Base(path.Path)), ")",
				leftPad, strings.ReplaceAll(err.Error(), "\n", leftPad),
				"\n",
			)
			errs = append(errs, errors.New(errStr.String()))
			errStr.Reset()
		}
		// validate path stretch if set
		stretch := path.StretchOrDefault()
		if _, err = ValidateStretch(&stretch); err != nil {
			fmt.Fprint(&errStr,
				"path ", i+1, " stretch",
				" (", filepath.Join("..", filepath.Base(path.Path)), ")",
				leftPad, strings.ReplaceAll(err.Error(), "\n", leftPad),
				"\n",
			)
			errs = append(errs, errors.New(errStr.String()))
			errStr.Reset()
		}
	}

	// validate config interval if set
	interval := strconv.FormatUint(uint64(cfg.IntervalOrDefault()), 10)
	if _, err := ValidateInterval(&interval); err != nil {
		errs = append(errs, fmt.Errorf("interval: %s", err))
	}
	// validate config port if set
	port := strconv.FormatUint(uint64(cfg.PortOrDefault()), 10)
	if _, err := ValidatePort(&port); err != nil {
		errs = append(errs, fmt.Errorf("port: %s", err))
	}
	// validate config profile if set
	profile := cfg.ProfileOrDefault()
	if _, err := ValidateProfile(&profile); err != nil {
		errs = append(errs, fmt.Errorf("profile: %s", err))
	}
	return errs
}

type ImagesPath struct {
	Path      string   `yaml:"path"`
	Alignment *string  `yaml:"alignment,omitempty"`
	Opacity   *float32 `yaml:"opacity,omitempty"`
	Stretch   *string  `yaml:"stretch,omitempty"`
}

func (path *ImagesPath) String() string {
	return fmt.Sprint(`Path: `, path.Path, `
  Alignment: `, Option(path.Alignment).UnwrapOr("not set"), `
  Stretch: `, Option(path.Stretch).UnwrapOr("not set"), `
  Opacity: `, func() string {
		if path.Opacity != nil {
			return strconv.FormatFloat(float64(*path.Opacity), 'f', -1, 32)
		}
		return "not set"
	}(), `
`)
}

// get alignment if set, otherwise the default value
func (path *ImagesPath) AlignmentOrDefault() string {
	return Option(path.Alignment).UnwrapOr(DefaultAlignment)
}

// get opacity if set, otherwise the default value
func (path *ImagesPath) OpacityOrDefault() float32 {
	return Option(path.Opacity).UnwrapOr(DefaultOpacity)
}

// get stretch if set, otherwise the default value
func (path *ImagesPath) StretchOrDefault() string {
	return Option(path.Stretch).UnwrapOr(DefaultStretch)
}

// get all images under the directory
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
		if IsImageFile(filepath.Join(dir, d.Name())) {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to walk directory %s: %s", dir, err)
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("Found no image files at %s", dir)
	}
	return images, nil
}

type ConfigLogger struct{}

func (cfg *Config) Log(configPath string) ConfigLogger {
	shrunkConfigPath := shrinkHome(configPath)
	fmt.Print("## ", shrunkConfigPath, `
paths:    `, func() string {
		if len(cfg.Paths) == 0 {
			return "[]"
		}
		var ret strings.Builder
		for _, dir := range cfg.Paths {
			fmt.Fprint(&ret, `
    - path: `, dir.Path)
			if dir.Alignment != nil {
				fmt.Fprint(&ret, `
      alignment: `, dir.AlignmentOrDefault())
			}
			if dir.Opacity != nil {
				fmt.Fprint(&ret, `
      opacity: `, dir.OpacityOrDefault())
			}
			if dir.Stretch != nil {
				fmt.Fprint(&ret, `
      stretch:`, dir.StretchOrDefault())
			}
		}
		return ret.String()
	}(), `
profile:  `, cfg.ProfileOrDefault(), `
port:     `, cfg.PortOrDefault(), `
interval: `, cfg.IntervalOrDefault(), `
`, func() string {
		var ret strings.Builder
		if errs := cfg.Validate(); len(errs) > 0 {
			fmt.Fprintln(&ret, "## ERRORS:")
			for _, err := range errs {
				fmt.Fprintln(&ret, "#", strings.ReplaceAll(err.Error(), "\n", "\n#"))
			}
			return ret.String()
		}
		return ""
	}(),
	)
	return ConfigLogger{}
}

func (log ConfigLogger) Added(added *ImagesPath, edited *pathEdit) ConfigLogger {
	if edited == nil && added == nil {
		fmt.Println("# no changes made")
	} else {
		if added != nil {
			fmt.Println("## ADDED: ")
			fmt.Printf("%-5sadded path: %s\n", "#", added.Path)
			if added.Alignment != nil {
				fmt.Printf("%-5s- alignment: %s\n", "#", *added.Alignment)
			}
			if added.Opacity != nil {
				fmt.Printf("%-5s- opacity: %f\n", "#", *added.Opacity)
			}
			if added.Stretch != nil {
				fmt.Printf("%-5s- stretch: %s\n", "#", *added.Stretch)
			}
		}
		if edited != nil {
			fmt.Println("## EDITED: ")
			fmt.Printf("%-5sedited path: %s\n", "#", edited.old.Path)
			if edited.new.Alignment != nil {
				if edited.old.Alignment != nil {
					fmt.Printf("%-5s- old alignment: %s\n", "#", *edited.old.Alignment)
				}
				fmt.Printf("%-5s- new alignment: %s\n", "#", *edited.new.Alignment)
			}
			if edited.new.Opacity != nil {
				if edited.old.Opacity != nil {
					fmt.Printf("%-5s- old opacity: %f\n", "#", *edited.old.Opacity)
				}
				fmt.Printf("%-5s- new opacity: %f\n", "#", *edited.new.Opacity)
			}
			if edited.new.Stretch != nil {
				if edited.old.Stretch != nil {
					fmt.Printf("%-5s- old stretch: %s\n", "#", *edited.old.Stretch)
				}
				fmt.Printf("%-5s- new stretch: %s\n", "#", *edited.new.Stretch)
			}
		}
	}
	return log
}

func (log ConfigLogger) Removed(removed *string) ConfigLogger {
	fmt.Println("## REMOVED: ")
	if removed == nil {
		fmt.Println("# no changes made")
	} else {
		fmt.Println("#", *removed)
	}
	return log
}

func (log ConfigLogger) Edited(edited []configEdits) ConfigLogger {
	fmt.Println("## EDITED")
	if len(edited) == 0 {
		fmt.Println("# no changes made")
	} else {
		for _, edit := range edited {
			fmt.Printf("# %-10s %s --> %s\n", edit.title, edit.old, edit.new)
		}
	}
	return log
}
