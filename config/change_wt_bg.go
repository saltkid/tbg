package config

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/saltkid/tbg/utils"
)

func (c *Config) ChangeBgImage(configPath string, profile *string, interval *string, align *string, stretch *string, opacity *string) error {
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
	prevDir := make(chan struct{})
	prevImage := make(chan struct{})
	randomDir := make(chan struct{})
	randomImage := make(chan struct{})
	go readUserInput(keysEvents, done, nextDir, nextImage, prevDir, prevImage, randomDir, randomImage)

	for {
		// indices to allow going to previous dir/image
		dirIndex, imgIndex := 0, 0

		/* bool to determine direction
		 *  only applicable in nextImage and prevImage
		 *  when going nextDir/prevDir, it will always start at the first image (always true)
		 *
		 *  if true : when going to next dir, start at first image of next dir
		 *          : when dirs are exhausted, start at first image of first dir
		 *  if false: when going to prev  dir, start at last image of prev dir
		 *          : when dirs are exhausted, start at last image of last dir
		 */
		startAtFirstImage := true

		for dirIndex >= 0 && dirIndex < len(c.ImageColPaths) {
			/* the order of flag importance is:
			 *
			 * 1. flags set by user on execution      (eg. --alignment, --stretch, etc.)
			 * 2. per path options set in config.yaml (eg. C:/Users/username/Pictures | right uniform 0.5)
			 * 3. default fields on config.yaml       (eg. default_alignment, default_stretch, etc.)
			 */

			// default fields
			overrideAlign, overrideStretch, overrideOpacity := c.Alignment, c.Stretch, strconv.FormatFloat(c.Opacity, 'f', -1, 64)

			// check if path entry has per path options
			dir := c.ImageColPaths[dirIndex]
			dir, opts, hasOpts := strings.Cut(dir, "|")
			dir = strings.TrimSpace(dir)
			images, err := fetchImages(dir)
			if err != nil {
				return err
			}
			if startAtFirstImage {
				imgIndex = 0
			} else {
				imgIndex = len(images) - 1
			}

			if hasOpts {
				opts = strings.TrimSpace(opts)
				optSlice := strings.Split(opts, " ")
				overrideAlign, overrideStretch, overrideOpacity = strings.TrimSpace(optSlice[0]), strings.TrimSpace(optSlice[1]), strings.TrimSpace(optSlice[2])
			}

			// flags set by user on execution
			var intervalInt int
			if profile == nil {
				profile = &c.Profile
			}
			if interval == nil {
				intervalInt = c.Interval
			} else {
				intervalInt, _ = strconv.Atoi(*interval)
			}
			if align != nil {
				overrideAlign = *align
			}
			if stretch != nil {
				overrideStretch = *stretch
			}
			if opacity != nil {
				overrideOpacity = *opacity
			}

		imageLoop:
			for imgIndex >= 0 && imgIndex < len(images) {
				image := images[imgIndex]

				ticker := time.Tick(time.Duration(intervalInt) * time.Minute)
				// ticker := time.Tick(time.Second * 10) // for debug purposes

				fmt.Println()
				opacityF, _ := strconv.ParseFloat(overrideOpacity, 64)
				c.Log(configPath).LogRunSettings(image, *profile, intervalInt, overrideAlign, overrideStretch, opacityF)

				err = updateWtJsonFields(allData, settingsPath, *profile, image, overrideAlign, overrideStretch, overrideOpacity)
				if err != nil {
					return err
				}
				// prompt
				fmt.Println("Press a key to execute a command ('c' for list of commands): ")

				select {
				case <-ticker:
					imgIndex++
					if imgIndex == len(images) {
						fmt.Print("no more images. going to next dir: ")
						dirIndex++
						startAtFirstImage = true
					}
				case <-done:
					fmt.Println("Goodbye!")
					return nil
				case <-randomImage:
					fmt.Println("randomizing from current image up to last image...")
					fmt.Println("(previous images will not be randomized so you can go back)")
					for range images[imgIndex:] {
						i := rand.Intn(len(images) - imgIndex)
						i += imgIndex
						images[i], images[imgIndex] = images[imgIndex], images[i]
					}
				case <-randomDir:
					fmt.Println("randomizing from current dir up to last dir...")
					fmt.Println("(previous dirs will not be randomized so you can go back)")
					for range c.ImageColPaths[dirIndex:] {
						i := rand.Intn(len(c.ImageColPaths) - dirIndex)
						i += dirIndex
						c.ImageColPaths[i], c.ImageColPaths[dirIndex] = c.ImageColPaths[dirIndex], c.ImageColPaths[i]
					}
					break imageLoop
				case <-nextImage:
					fmt.Println("using next image...")
					imgIndex++
					if imgIndex == len(images) {
						fmt.Print("no more images. going to next dir: ")
						dirIndex++
						startAtFirstImage = true
					}
				case <-prevImage:
					fmt.Println("using previous image...")
					imgIndex--
					if imgIndex < 0 {
						fmt.Print("no more images. going to previous dir: ")
						dirIndex--
						startAtFirstImage = false
					}
				case <-nextDir:
					fmt.Println("using next dir...")
					dirIndex++
					startAtFirstImage = true
					break imageLoop
				case <-prevDir:
					fmt.Println("using previous dir...")
					dirIndex--
					startAtFirstImage = true
					break imageLoop
				}

			}

			if dirIndex >= len(c.ImageColPaths) {
				fmt.Println("no more next dirs. going to first dir again: ", c.ImageColPaths[0])
				dirIndex = 0
			} else if dirIndex < 0 {
				fmt.Println("no more previous dirs. going to last dir again: ", c.ImageColPaths[len(c.ImageColPaths)-1])
				dirIndex = len(c.ImageColPaths) - 1
			}

		}
	}
}

func readUserInput(keysEvents <-chan keyboard.KeyEvent, done chan<- struct{},
	nextDir chan<- struct{}, nextImage chan<- struct{},
	prevDir chan<- struct{}, prevImage chan<- struct{},
	randomDir chan<- struct{}, randomImage chan<- struct{}) {
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		if keyboard.Key(event.Rune) == keyboard.Key('q') {
			fmt.Println("Exiting...")
			close(done)
			return

		} else if keyboard.Key(event.Rune) == keyboard.Key('R') {
			randomDir <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('r') {
			randomImage <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('N') {
			nextDir <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('P') {
			prevDir <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('n') {
			nextImage <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('p') {
			prevImage <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('c') {
			commandList()

		} else {
			fmt.Printf("invalid key '%c' ('h' for help)\n", event.Rune)
		}
	}
}

func updateWtJsonFields(allData map[string]json.RawMessage, settingsPath string, profile string, image string, align string, stretch string, opacity string) error {
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

func commandList() {
	fmt.Println()
	fmt.Println("q: [q]uit")
	fmt.Println("n: [n]ext image")
	fmt.Println("p: [p]revious image")
	fmt.Println("f: [N]ext dir")
	fmt.Println("b: [P]revious dir")
	fmt.Println("c: [c]ommand list")
}
