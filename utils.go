package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// shuffles the slice from the current index up to the end
//
// it does not affect the elements before the current index
func ShuffleFrom[T any](currentIndex int, slice []T) {
	for range slice[currentIndex:] {
		i := rand.IntN(len(slice)-currentIndex) + currentIndex
		slice[i], slice[currentIndex] = slice[currentIndex], slice[i]
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
	return os.ExpandEnv(expandWin)
}

// Expands prefixed tilde to the user's home directory
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
