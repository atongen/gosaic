package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	exiftoolName = "exiftool"
)

func ExiftoolPath(toolPath string) (myPath string, err error) {
	if toolPath == "" {
		// provided path is empty, search PATH, it is not an error
		// if exiftool isn't found
		myPath, _ = ExiftoolPathFind()
	} else {
		myPath, err = ExiftoolPathValidate(toolPath)
	}
	return
}

// ExiftoolPathValidate validates the provided path to exiftool exists,
// and if so, returns the absolute path to it.
func ExiftoolPathValidate(toolPath string) (myPath string, err error) {
	_, err = os.Stat(toolPath)
	if err != nil {
		return
	}
	myPath, err = filepath.Abs(toolPath)
	return
}

// ExiftoolPathFind searches PATH for exiftool and returns the absolute
// path to it if found.
func ExiftoolPathFind() (myPath string, err error) {
	myPath, err = exec.LookPath(exiftoolName)
	if err != nil {
		return
	}
	myPath, err = filepath.Abs(myPath)
	return
}

// YYYY:MM:DD HH:MM:SS
func ExifTags() map[string]string {
	now := time.Now()
	return map[string]string{
		"ResolutionUnit": "inches",
		"XResolution":    "300",
		"YResolution":    "300",
		"Software":       "https://github.com/atongen/gosaic",
		"DateTime":       now.Format("2006:01:02 15:04:05"),
	}
}

// ExifCp copies exif data from src image to dst image.
// It excludes Orientation tag because the orientation has been
// normalized in the processed images.
func ExifCp(toolPath string, src, dst string) (string, error) {
	args := []string{"-overwrite_original_in_place", "-tagsFromFile", src, "-x", "Orientation"}
	for key, value := range ExifTags() {
		args = append(args, fmt.Sprintf("-%s=\"%s\"", key, value))
	}
	args = append(args, dst)
	cmd := exec.Command(toolPath, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
