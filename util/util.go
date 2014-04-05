package util

import (
	"crypto/md5"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gosexy/exif"
)

func Md5sum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	hash := md5.New()
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		if _, err := io.WriteString(hash, string(buf[:n])); err != nil {
			panic(err)
		}
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// http://sylvana.net/jpegcrop/exif_orientation.html
func OpenImage(path string) (*image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}

	reader := exif.New()
	err = reader.Open(path)
	if err == nil {
		if orientation, ok := reader.Tags["Orientation"]; ok {
			switch orientation {
			case "Top-left":
				// 1
			case "Top-right":
				// 2 flop!
				img = imaging.FlipH(img)
			case "Bottom-right":
				// 3 rotate!(180)
				img = imaging.Rotate180(img)
			case "Bottom-left":
				// 4 flip!
				img = imaging.FlipV(img)
			case "Left-top":
				// 5 transpose!
				img = imaging.Rotate270(img)
				img = imaging.FlipH(img)
			case "Right-top":
				// 6 rotate!(90)
				img = imaging.Rotate270(img)
			case "Right-bottom":
				// 7 transverse!
				img = imaging.Rotate90(img)
				img = imaging.FlipH(img)
			case "Left-bottom":
				// 8 rotate!(270)
				img = imaging.Rotate90(img)
			}
		}
	}

	return &img, nil
}
