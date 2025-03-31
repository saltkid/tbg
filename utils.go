package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Specifically, checks if the given path points to an image file that is
// supported by Windows Terminal (png, jpeg, gif)
func IsImageFile(path string) bool {
	handle, err := os.Open(path)
	if err != nil {
		return false
	}
	defer handle.Close()
	buf := make([]byte, 512) // 512 according to docs
	if _, err := handle.Read(buf); err != nil {
		return false
	}
	mime := strings.Split(http.DetectContentType(buf), "/")
	if len(mime) != 2 {
		return false
	}
	if mime[0] != "image" {
		return false
	}
	switch mime[1] {
	case "png", "jpeg", "gif":
		return true
	default:
		return false
	}
}

// normalized path to:
//
// 1. be absolute
//
// 2. use "/" as separator
//
// 3. expand prefixed tilde (~) to user home dir
//
// 4. expand environment variables
func NormalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(expandEnv(path))
	if err != nil {
		return "", fmt.Errorf("Failed to get absolute path of %s: %s", path, err)
	}
	return filepath.ToSlash(absPath), nil
}

// Expands both "%windows_style%" and "$unix_style  ${env_vars}" in the string,
// along with tilde (~) expansion
//
// Can expand a string that contains both styles: "$URL:%PORT%/${ENDPOINT}"
func expandEnv(s string) string {
	expanded := expandTilde(s)
	re, err := regexp.Compile(`%([^%]+)%`)
	if err != nil {
		return expanded
	}
	expandWin := re.ReplaceAllStringFunc(expanded, func(str string) string {
		envVar := str[1 : len(str)-1]
		if val, exists := os.LookupEnv(envVar); exists {
			return val
		}
		return str
	})
	expandUnix := os.ExpandEnv(expandWin)
	return expandUnix
}

// Expands a path containing prefixed tilde to the user's home directory.
// Opposite of shrinkHome
//
// ~/Pictures --> C:/Users/username/Pictures
func expandTilde(path string) string {
	if strings.HasPrefix(filepath.ToSlash(path), "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// Shrinks the user home directory path to a "~"
// Opposite of expandTilde
//
// C:/Users/username/Pictures --> ~/Pictures
func shrinkHome(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(filepath.ToSlash(path), filepath.ToSlash(homeDir)) {
		return filepath.Join("~", strings.TrimPrefix(path, homeDir))
	}
	return path
}
