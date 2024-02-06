package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/saltkid/tbg/utils"
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

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	done := make(chan struct{})
	nextDir := make(chan struct{})
	nextImage := make(chan struct{})
	go readUserInput(keysEvents, done, nextDir, nextImage)

	for {
		/* the order of flag importance is:
		 *
		 * 1. flags set by user on execution      (eg. --alignment, --stretch, etc.)
		 * 2. per path options set in config.yaml (eg. C:/Users/username/Pictures | right uniform 0.5)
		 * 3. default fields on config.yaml       (eg. default_alignment, default_stretch, etc.)
		 *
		 */

		// default fields
		overrideAlign, overrideStretch, overrideOpacity := c.Alignment, c.Stretch, strconv.FormatFloat(c.Opacity, 'f', -1, 64)
		for i, dir := range c.ImageColPaths {
			// check if path entry has per path options
			dir, opts, hasOpts := strings.Cut(dir, "|")
			dir = strings.TrimSpace(dir)
			if hasOpts {
				opts = strings.TrimSpace(opts)
				optSlice := strings.Split(opts, " ")
				overrideAlign, overrideStretch, overrideOpacity = strings.TrimSpace(optSlice[0]), strings.TrimSpace(optSlice[1]), strings.TrimSpace(optSlice[2])
			}

			// flags set by user on execution
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
				fmt.Println("Press a key to execute a command ('h' for help): ")

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

func readUserInput(keysEvents <-chan keyboard.KeyEvent, done chan<- struct{}, nextDir chan<- struct{}, nextImage chan<- struct{}) {
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		if keyboard.Key(event.Rune) == keyboard.Key('q') {
			fmt.Println("Exiting...")
			close(done)
			return

		} else if keyboard.Key(event.Rune) == keyboard.Key('c') {
			nextDir <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('n') {
			nextImage <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('h') {
			help()

		} else {
			fmt.Printf("invalid key '%c' ('h' for help)\n", event.Rune)
		}
	}
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

func fetchImages(dir string) ([]string, error) {
	images := make([]string, 0)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != filepath.Base(dir) {
			return filepath.SkipDir
		}
		if utils.IsImageFile(d.Name()) {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to walk directory %s: %s", dir, err)
	}
	return images, nil
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

func help() {
	fmt.Println()
	fmt.Println("q: [q]uit")
	fmt.Println("c: [c]hange dir")
	fmt.Println("n: [n]ext image")
	fmt.Println("h: [h]elp")
}
