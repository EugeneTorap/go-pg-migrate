package migrate

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultTimeFormat = "20060102150405"
	defaultTimezone   = "UTC"
	defaultExt        = "sql"
)

// CreateCmd (meant to be called via a CLI command) creates a new migration
func CreateCmd(dir string, name string, print bool) error {
	dir = filepath.Clean(dir)
	ext := "." + strings.TrimPrefix(defaultExt, ".")

	currentTime := time.Now()
	timezone, err := time.LoadLocation(defaultTimezone)
	if err != nil {
		log.Fatal(err)
	}
	startTime := currentTime.In(timezone)

	version := startTime.Format(defaultTimeFormat)
	versionGlob := filepath.Join(dir, version+"_*"+ext)
	matches, err := filepath.Glob(versionGlob)
	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(dir, basename)

		if err = createFile(filename); err != nil {
			return err
		}

		if print {
			absPath, _ := filepath.Abs(filename)
			log.Println(absPath)
		}
	}

	return nil
}

func createFile(filename string) error {
	// create exclusive (fails if file already exists)
	// os.Create() specifies 0666 as the FileMode, so we're doing the same
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	return f.Close()
}
