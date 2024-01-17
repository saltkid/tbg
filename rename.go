package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func rename() {
	path, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln("error reading path:", err)
	}
	log.Println("[INFO] renaming images in:", path)

	toRename := make([]string, 0)
	alreadyCorrect := make([]string, 0)
	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && !strings.EqualFold(d.Name(), filepath.Base(path)) {
			return filepath.SkipDir
		}

		if IsImageFile(d.Name()) {
			re := regexp.MustCompile(`^(\d+)\.`)
			if re.MatchString(d.Name()) {
				alreadyCorrect = append(alreadyCorrect, d.Name())
			} else {
				toRename = append(toRename, d.Name())
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalln("error when reading dir", filepath.Clean("."), ";", err)
	}

	// fill in gaps in alreadyCorrect files
	// ie: input:	1.png, 3.png, 4.png, 5.png
	//		0.png, 1.png, 2.png, 3.png
	log.Println("sorting:", alreadyCorrect)
	sort.Sort(Natural(alreadyCorrect))
	log.Println("sorted :", alreadyCorrect)
	for i, file := range alreadyCorrect {
		re := regexp.MustCompile(`^(\d+)\.`)
		match := re.FindStringSubmatch(file)
		num, _ := strconv.Atoi(match[1])
		if i != num {
			newName := fmt.Sprintf("%d%s", i, filepath.Ext(file))
			err := os.Rename(file, newName)
			if err != nil {
				log.Println("[ERROR] error renaming", file, "to", newName, ";", err)
			}
		}
	}

	log.Println("[INFO] files to rename:", toRename)
	newIndexStart := len(alreadyCorrect)
	for i, file := range toRename {
		newName := fmt.Sprintf("%d%s", newIndexStart+i, filepath.Ext(file))

		// check if file exists
		if _, err := os.Stat(file); os.IsExist(err) {
			log.Println("[ERROR]:", newName, "already exists:", err)
			continue
		}

		err := os.Rename(file, newName)
		if err != nil {
			log.Println("[ERROR] error renaming", file, "to", newName, ";", err)
			continue
		}
		log.Println("[INFO]: renamed", file, "to", newName)
	}
}
