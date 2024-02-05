package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (c *Config) EditWTJson(configPath string, profile string, interval string, align string, stretch string, opacity string) error {
	// read settings.json
	settingsPath, err := settingsJsonPath()
	if err != nil {
		return err
	}
	settingsData, err := os.ReadFile(settingsPath)

	// get all data first to keep unused fields
	var allData map[string]json.RawMessage
	err = json.Unmarshal(settingsData, &allData)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal settings.json at %s: %s", settingsPath, err)
	}

	done := make(chan struct{})
	nextDir := make(chan struct{})
	nextImage := make(chan struct{})
	go readUserInput(done, nextDir, nextImage)

	for {
		overrideAlign, overrideStretch, overrideOpacity := c.Alignment, c.Stretch, strconv.FormatFloat(c.Opacity, 'f', -1, 64)
		for i, dir := range c.ImageColPaths {
			// per path options to override defaults
			dir, opts, hasOpts := strings.Cut(dir, "|")
			dir = strings.TrimSpace(dir)
			// use defaults of no options
			if hasOpts {
				opts = strings.TrimSpace(opts)
				optSlice := strings.Split(opts, " ")
				overrideAlign, overrideStretch, overrideOpacity = strings.TrimSpace(optSlice[0]), strings.TrimSpace(optSlice[1]), strings.TrimSpace(optSlice[2])
			}

			// set values only if specified by user
			intervalInt, _ := strconv.Atoi(interval)
			if profile == "" {
				profile = c.Profile
			}
			if interval == "" {
				intervalInt = c.Interval
			}
			if align != "" {
				overrideAlign = align
			}
			if stretch != "" {
				overrideStretch = stretch
			}
			if opacity != "" {
				overrideOpacity = opacity
			}

			images, err := fetchImages(dir)
			if err != nil {
				return err
			}

		imageLoop:
			for j, image := range images {
				ticker := time.Tick(time.Duration(intervalInt) * time.Minute)
				// ticker := time.Tick(time.Second * 10) // for debug purposes

				fmt.Println()
				opacityF, _ := strconv.ParseFloat(overrideOpacity, 64)
				c.Log(configPath).LogRunSettings(image, profile, intervalInt, overrideAlign, overrideStretch, opacityF)

				err = changeBackgroundImage(allData, settingsPath, profile, image, overrideAlign, overrideStretch, overrideOpacity)
				if err != nil {
					return err
				}
				// prompt
				fmt.Println("Enter a command ('h' for help): ")
				fmt.Print("> ")

				select {
				case <-ticker:
				case <-done:
					fmt.Println("Goodbye!")
					return nil
				case <-nextDir:
					fmt.Println("using next dir...")
					break imageLoop
				case <-nextImage:
					fmt.Println("using next image...")
					if j == len(images)-1 {
						fmt.Print("no more images. going to next dir: ")
					}
					continue
				}

			}
			if i == len(c.ImageColPaths)-1 {
				fmt.Println("no more dirs. going to first dir again: ", c.ImageColPaths[0])
			}
		}
	}
}

func fetchImages(dir string) ([]string, error) {
	images := make([]string, 0)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != filepath.Base(dir) {
			return filepath.SkipDir
		}
		if isImageFile(d.Name()) {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to walk directory %s: %s", dir, err)
	}
	return images, nil
}

func readUserInput(done chan<- struct{}, nextDir chan<- struct{}, nextImage chan<- struct{}) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if scanner.Text() == "q" {
				fmt.Println("Exiting...")
				close(done)
				return

			} else if scanner.Text() == "c" {
				nextDir <- struct{}{}

			} else if scanner.Text() == "n" {
				nextImage <- struct{}{}

			} else if scanner.Text() == "h" {
				help()

			} else {
				fmt.Printf("invalid input '%s' ('h' for help)\n", scanner.Text())
				fmt.Print("> ")
			}
		} else {
			return
		}
	}
}

func help() {
	fmt.Println("q: [q]uit")
	fmt.Println("c: [c]hange dir")
	fmt.Println("n: [n]ext image")
	fmt.Println("h: [h]elp")
	fmt.Print("> ")
}

func settingsJsonPath() (string, error) {
	localAppDataPath := os.Getenv("LOCALAPPDATA")

	// stable release
	settingsJson := filepath.Join(localAppDataPath, "Packages", "Microsoft.WindowsTerminal_8wekyb3d8bbwe", "LocalState", "settings.json")
	if _, err := os.Stat(settingsJson); !os.IsNotExist(err) {
		return settingsJson, nil
	}

	// preview release
	settingsJson = filepath.Join(localAppDataPath, "Packages", "Microsoft.WindowsTerminalPreview_8wekyb3d8bbwe", "LocalState", "settings.json")
	if _, err := os.Stat(settingsJson); !os.IsNotExist(err) {
		return settingsJson, nil
	}

	// through package managers (chocolatey, scoop, etc)
	settingsJson = filepath.Join(localAppDataPath, "Microsoft", "Windows Terminal", "settings.json")
	if _, err := os.Stat(settingsJson); !os.IsNotExist(err) {
		return settingsJson, nil
	}

	return "", fmt.Errorf("settings.json not found")

}

func changeBackgroundImage(allData map[string]json.RawMessage, settingsPath string, profile string, image string, align string, stretch string, opacity string) error {
	// normalize fields to be json friendly
	image = strings.ReplaceAll(image, `\`, `\\`)
	image = fmt.Sprintf(`"%s"`, image)
	align = fmt.Sprintf(`"%s"`, align)
	stretch = fmt.Sprintf(`"%s"`, stretch)

	// read profiles
	var profiles map[string]json.RawMessage
	err := json.Unmarshal(allData["profiles"], &profiles)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal field \"profiles\" from settings.json: %s", err)
	}

	// edit profiles
	if profile == "default" {
		var defaultProfile map[string]json.RawMessage
		err = json.Unmarshal(profiles["defaults"], &defaultProfile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal field \"defaults\" from field \"profiles\" in settings.json: %s", err)
		}

		defaultProfile["backgroundImage"] = json.RawMessage([]byte(image))
		defaultProfile["backgroundImageAlignment"] = json.RawMessage([]byte(align))
		defaultProfile["backgroundImageStretchMode"] = json.RawMessage([]byte(stretch))
		defaultProfile["backgroundImageOpacity"] = json.RawMessage([]byte(opacity))

		profiles["defaults"], err = json.Marshal(defaultProfile)
		if err != nil {
			return fmt.Errorf("Failed to marshal field \"defaults\" to field \"profiles\" in settings.json: %s", err)
		}
	} else {
		var profileList []map[string]json.RawMessage
		err = json.Unmarshal(profiles["list"], &profileList)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal field \"list\" from field \"profiles\" in settings.json: %s", err)
		}
		_, num, _ := strings.Cut(profile, "-")
		profileNum, _ := strconv.Atoi(num)
		profileList[profileNum-1]["backgroundImage"] = json.RawMessage([]byte(image))
		profileList[profileNum-1]["backgroundImageAlignment"] = json.RawMessage([]byte(align))
		profileList[profileNum-1]["backgroundImageStretchMode"] = json.RawMessage([]byte(stretch))
		profileList[profileNum-1]["backgroundImageOpacity"] = json.RawMessage([]byte(opacity))

		profiles["list"], err = json.Marshal(profileList)
		if err != nil {
			return fmt.Errorf("Failed to marshal field \"list\" to field \"profiles\" in settings.json: %s", err)
		}
	}

	// write profiles
	updatedProfiles, err := json.Marshal(profiles)
	if err != nil {
		return fmt.Errorf("Failed to marshal field \"profiles\" in settings.json: %s", err)
	}
	allData["profiles"] = updatedProfiles
	updatedJson, err := json.Marshal(allData)
	err = os.WriteFile(settingsPath, updatedJson, 0666)

	return nil
}
