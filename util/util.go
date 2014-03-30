package util

import (
	"crypto/md5"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"
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

func ImageBounds(path string) (image.Rectangle, error) {
	result := image.Rectangle{}
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return result, err
	}
	return img.Bounds(), nil
}
