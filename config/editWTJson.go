package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	// "strings"
	"strconv"
	"time"
)

type Profile struct {
	Profiles WTProfiles `json:"profiles"`
}

type WTProfiles struct {
	Default WTProfile   `json:"defaults"`
	List    []WTProfile `json:"list"`
}

type WTProfile struct {
	BGImage   string  `json:"backgroundImage"`
	BGAlign   string  `json:"backgroundImageAlignment"`
	BGStretch string  `json:"backgroundImageStretchMode"`
	BGOpacity float64 `json:"backgroundImageOpacity"`
}

func (c *DefaultConfig) EditWTJson(configPath string, profile string, interval string, align string, stretch string, opacity string) error {
	// set values only if specified by user
	intervalInt, _ := strconv.Atoi(interval)
	opacityF, _ := strconv.ParseFloat(opacity, 64)
	if profile == "" {
		profile = c.Profile
	}
	if interval == "" {
		intervalInt = c.Interval
	}
	if align == "" {
		align = c.Alignment
	}
	if stretch == "" {
		stretch = c.Stretch
	}
	if opacity == "" {
		opacityF = c.Opacity
	}
	println(profile, intervalInt, align, stretch, opacityF)

	// read settings.json
	settingsPath, err := settingsJsonPath()
	if err != nil {
		return err
	}
	settingsData, err := os.ReadFile(settingsPath)
	var p Profile
	err = json.Unmarshal(settingsData, &p)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal settings.json at %s: %s", settingsPath, err)
	}

	for {
		for i, dir := range c.ImageColPaths {
			images, err := fetchImages(dir)
			if err != nil {
				return err
			}

			done := make(chan struct{})
			nextDir := make(chan struct{})
			nextImage := make(chan struct{})

			go readUserInput(done, nextDir, nextImage)
			for j, image := range images {

				// ticker := time.Tick(time.Duration(intervalInt) * time.Minute)
				ticker := time.Tick(time.Second * 10)

				fmt.Println()
				c.Log(configPath).LogRunSettings(image, profile, intervalInt, align, stretch, opacityF)
				fmt.Println("Enter a command ('h' for help): ")
				fmt.Print("> ")

				select {
				case <-ticker:
				case <-done:
					fmt.Println("Goodbye!")
					return nil
				case <-nextDir:
					fmt.Println("using next dir...")
					break
				case <-nextImage:
					fmt.Println("using next image...")
					if j == len(images)-1 {
						fmt.Print("no more images. going to next dir: ")
						if i+1 < len(c.ImageColPaths) {
							fmt.Println(c.ImageColPaths[i+1])
						}
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

func (c *UserConfig) EditWTJson(configPath string, profile string, interval string, align string, stretch string, opacity string) error {
	// set values only if specified by user
	intervalInt, _ := strconv.Atoi(interval)
	opacityF, _ := strconv.ParseFloat(opacity, 64)
	if profile == "" {
		profile = c.Profile
	}
	if interval == "" {
		intervalInt = c.Interval
	}
	if align == "" {
		align = c.Alignment
	}
	if stretch == "" {
		stretch = c.Stretch
	}
	if opacity == "" {
		opacityF = c.Opacity
	}
	println(profile, intervalInt, align, stretch, opacityF)

	// read settings.json
	settingsPath, err := settingsJsonPath()
	if err != nil {
		return err
	}
	settingsData, err := os.ReadFile(settingsPath)
	var p Profile
	err = json.Unmarshal(settingsData, &p)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal settings.json at %s: %s", settingsPath, err)
	}

	for {
		for i, dir := range c.ImageColPaths {
			images, err := fetchImages(dir)
			if err != nil {
				return err
			}

			done := make(chan struct{})
			nextDir := make(chan struct{})
			nextImage := make(chan struct{})

			go readUserInput(done, nextDir, nextImage)
			for j, image := range images {

				// ticker := time.Tick(time.Duration(intervalInt) * time.Minute)
				ticker := time.Tick(time.Second * 10)

				fmt.Println()
				c.Log(configPath).LogRunSettings(image, profile, intervalInt, align, stretch, opacityF)
				fmt.Println("Enter a command ('h' for help): ")
				fmt.Print("> ")

				select {
				case <-ticker:
				case <-done:
					fmt.Println("Goodbye!")
					return nil
				case <-nextDir:
					fmt.Println("using next dir...")
					break
				case <-nextImage:
					fmt.Println("using next image...")
					if j == len(images)-1 {
						fmt.Print("no more images. going to next dir: ")
						if i+1 < len(c.ImageColPaths) {
							fmt.Println(c.ImageColPaths[i+1])
						}
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
				close(nextDir)
				return

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

func printInfo(p *Profile) {
	fmt.Println("Defaults:")
	fmt.Println("  image:", p.Profiles.Default.BGImage)
	fmt.Println("  align:", p.Profiles.Default.BGAlign)
	fmt.Println("  stretch:", p.Profiles.Default.BGStretch)
	fmt.Println("  opacity:", p.Profiles.Default.BGOpacity)

	fmt.Println("List:")
	for i, profile := range p.Profiles.List {
		fmt.Println(i)
		fmt.Println("  image:", profile.BGImage)
		fmt.Println("  align:", profile.BGAlign)
		fmt.Println("  stretch:", profile.BGStretch)
		fmt.Println("  opacity:", profile.BGOpacity)
		fmt.Println()
	}
}
