package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	lev "github.com/agnivade/levenshtein"
)

type WTSettings struct {
	Data map[string]json.RawMessage
	Path string
}

func NewWTSettings() (*WTSettings, error) {
	ret := new(WTSettings)
	var err error
	ret.Path, err = settingsJsonPath()
	if err != nil {
		return nil, err
	}
	if err = ret.readSettings(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (wt *WTSettings) Write(
	image string,
	profile string,
	alignment string,
	opacity float32,
	stretch string,
) error {
	if err := wt.readSettings(); err != nil {
		return err
	}
	// normalize fields to be json friendly
	image = fmt.Sprintf(`"%s"`, strings.ReplaceAll(image, `\`, `\\`))
	alignment = fmt.Sprintf(`"%s"`, alignment)
	stretch = fmt.Sprintf(`"%s"`, stretch)
	// read profiles
	var profiles map[string]json.RawMessage
	err := json.Unmarshal(wt.Data["profiles"], &profiles)
	if err != nil {
		return fmt.Errorf(`Failed to unmarshal field "profiles" from settings.json: %s`, err)
	}
	// edit profiles
	switch profile {
	case "default":
		err = writeToDefaultProfile(profiles, image, alignment, stretch, opacity)
	default:
		err = writeToListProfile(profiles, profile, image, alignment, stretch, opacity)
	}
	if err != nil {
		return err
	}
	// write profiles
	updatedProfiles, err := json.Marshal(profiles)
	if err != nil {
		return fmt.Errorf(`Failed to marshal field "profiles" in settings.json: %s`, err)
	}
	wt.Data["profiles"] = updatedProfiles
	updatedJson, err := json.Marshal(wt.Data)
	err = os.WriteFile(wt.Path, updatedJson, 0666)
	return nil
}

func (wt *WTSettings) readSettings() error {
	settingsData, err := os.ReadFile(wt.Path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(settingsData, &wt.Data)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal settings.json at %s: %s", wt.Path, err)
	}
	return nil
}

func writeToDefaultProfile(
	profiles map[string]json.RawMessage,
	image string,
	alignment string,
	stretch string,
	opacity float32,
) error {
	var defaultProfile map[string]json.RawMessage
	err := json.Unmarshal(profiles["defaults"], &defaultProfile)
	if err != nil {
		return fmt.Errorf(`Failed to unmarshal field "defaults" from field "profiles" in settings.json: %s`, err)
	}
	defaultProfile["backgroundImage"] = json.RawMessage([]byte(image))
	defaultProfile["backgroundImageAlignment"] = json.RawMessage([]byte(alignment))
	defaultProfile["backgroundImageStretchMode"] = json.RawMessage([]byte(stretch))
	defaultProfile["backgroundImageOpacity"] = json.RawMessage([]byte(strconv.FormatFloat(float64(opacity), 'f', -1, 32)))
	profiles["defaults"], err = json.Marshal(defaultProfile)
	if err != nil {
		return fmt.Errorf(`Failed to marshal field "defaults" to field "profiles" in settings.json: %s`, err)
	}
	return nil
}

func writeToListProfile(
	profiles map[string]json.RawMessage,
	profile string,
	image string,
	alignment string,
	stretch string,
	opacity float32,
) error {
	profileNum, err := strconv.Atoi(profile)
	if err != nil {
		return writeToListProfileByName(profiles, profile, image, alignment, stretch, opacity)
	} else {
		return writeToListProfileByNumber(profiles, uint16(profileNum), image, alignment, stretch, opacity)
	}
}

func writeToListProfileByName(
	profiles map[string]json.RawMessage,
	profileName string,
	image string,
	alignment string,
	stretch string,
	opacity float32,
) error {
	var profileList []map[string]json.RawMessage
	err := json.Unmarshal(profiles["list"], &profileList)
	if err != nil {
		return fmt.Errorf(`Failed to unmarshal field "list" from field "profiles" in settings.json: %s`, err)
	}
	// TODO: handle profiles with the same name; e.g. multiple "Arch"
	closestNameIndices := make([]uint8, 0)
	for i, profile := range profileList {
		wtProfileName := strings.ToLower(string(profile["name"]))
		// json strings are enclosed in `"` so I enclosed the user arg
		distance := lev.ComputeDistance(fmt.Sprintf(`"%s"`, strings.ToLower(profileName)), wtProfileName)
		if distance == 0 {
			profileList[i]["backgroundImage"] = json.RawMessage([]byte(image))
			profileList[i]["backgroundImageAlignment"] = json.RawMessage([]byte(alignment))
			profileList[i]["backgroundImageStretchMode"] = json.RawMessage([]byte(stretch))
			profileList[i]["backgroundImageOpacity"] = json.RawMessage([]byte(strconv.FormatFloat(float64(opacity), 'f', -1, 32)))
			profiles["list"], err = json.Marshal(profileList)
			if err != nil {
				return fmt.Errorf(`Failed to marshal field "list" to field "profiles" in settings.json: %s`, err)
			}
			return nil
		} else if distance < 3 {
			closestNameIndices = append(closestNameIndices, uint8(i))
		}
	}
	switch len(closestNameIndices) {
	case 0:
		fmt.Printf("Failed to find profile with name equal or similar to \"%s\".\n", profileName)
		fmt.Println("Here is the list of available WT profiles with their respective")
		fmt.Println("profile number incase there are multiple with the same name:")
		var ret strings.Builder
		for i, profile := range profileList {
			ret.WriteString(fmt.Sprintf("%d) %s\n", i+1, profile["name"]))
		}
		return errors.New(ret.String())
	case 1:
		similarNameIndex := closestNameIndices[0]
		similarProfileName := string(profileList[similarNameIndex]["name"])
		fmt.Printf("There is no profile \"%s\"\n", profileName)
		return fmt.Errorf("Did you mean %s (profile #%d)?", similarProfileName, similarNameIndex)
	default:
		fmt.Printf("There is no profile \"%s\"\n", profileName)
		fmt.Println("Did you mean any of the following?")
		var ret strings.Builder
		for _, index := range closestNameIndices {
			ret.WriteString(fmt.Sprintf("%d) %s ", index+1, profileList[index]["name"]))
		}
		return errors.New(ret.String())
	}
}

func writeToListProfileByNumber(
	profiles map[string]json.RawMessage,
	profileNum uint16,
	image string,
	alignment string,
	stretch string,
	opacity float32,
) error {
	var profileList []map[string]json.RawMessage
	err := json.Unmarshal(profiles["list"], &profileList)
	if err != nil {
		return fmt.Errorf(`Failed to unmarshal field "list" from field "profiles" in settings.json: %s`, err)
	}
	if int(profileNum) > len(profileList) {
		fmt.Printf("Profile number \"%d\" does not exist. Only %d profiles exist:\n", profileNum, len(profileList))
		var ret strings.Builder
		for i, profile := range profileList {
			ret.WriteString(fmt.Sprintf("%d) %s\n", i+1, profile["name"]))
		}
		return errors.New(ret.String())
	}
	profileList[profileNum-1]["backgroundImage"] = json.RawMessage([]byte(image))
	profileList[profileNum-1]["backgroundImageAlignment"] = json.RawMessage([]byte(alignment))
	profileList[profileNum-1]["backgroundImageStretchMode"] = json.RawMessage([]byte(stretch))
	profileList[profileNum-1]["backgroundImageOpacity"] = json.RawMessage([]byte(strconv.FormatFloat(float64(opacity), 'f', -1, 32)))
	profiles["list"], err = json.Marshal(profileList)
	if err != nil {
		return fmt.Errorf(`Failed to marshal field "list" to field "profiles" in settings.json: %s`, err)
	}
	return nil
}

func settingsJsonPath() (string, error) {
	localAppDataPath, exists := os.LookupEnv("LOCALAPPDATA")
	if !exists {
		return "", fmt.Errorf("LOCALAPPDATA environment variable is not set")
	} else if exists && localAppDataPath == "" {
		return "", fmt.Errorf("LOCALAPPDATA environment variable is empty")
	}
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
	return "", fmt.Errorf("Windows Terminal's settings.json not found")
}
