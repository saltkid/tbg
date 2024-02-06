// common utils

package utils

import (
	"path/filepath"
	"strings"
)

func IsImageFile(f string) bool {
	f = strings.ToLower(filepath.Ext(f))
	return f == ".png" ||
		f == ".jpg" ||
		f == ".jpeg"
}
