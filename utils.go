package main

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

func IsImageFile(f string) bool {
	f = strings.ToLower(filepath.Ext(f))
	return f == ".png" ||
		f == ".jpg" ||
		f == ".jpeg"
}

// natural sorting of filenames containing only numbers
//
// usage: sort.Sort(FilenameSort(slice of paths))
type Natural []string

// implement sort.Interface (Len, Less, Swap)

func (f Natural) Len() int {
	return len(f)
}

// natural sorting comparing of filenames consisting of only strings
func (f Natural) Less(i, j int) bool {
	n1, _ := strconv.Atoi(f[i])
	n2, _ := strconv.Atoi(f[j])
	if n1 != n2 {
		return n1 < n2
	}
	// values are the same, compare length (ie: 0001 and 1)
	return len(f[i]) < len(f[j])
}
func (f Natural) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func LogParsedArgs(command *Command) {
	log.Println("command: ", command.name)
	log.Println("flags:")
	for _, flag := range command.flags {
		_, isFlag := flag.(*Flag)
		_, isCommand := flag.(*Command)
		if isFlag {
			log.Println(flag.(*Flag).name, "=", flag.(*Flag).value)
		} else if isCommand {
			log.Println(flag.(*Command).name, "=", flag.(*Command).value)
		}
	}
}
