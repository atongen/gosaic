package util

import (
	"crypto/md5"
	"fmt"
	"gosaic/model"
	"image"
	_ "image/jpeg"
	"io"
	"math"
	"os"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

const (
	// size of partial square for down-sampling
	// before getting slice of lab data
	DATA_SIZE = 10
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

func GetOrientation(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	x, err := exif.Decode(f)
	if err != nil {
		return 0, err
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return 0, err
	}

	val, err := orientation.Int(0)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// FixOrientation modifies image in-place to match exif orientation data
// http://sylvana.net/jpegcrop/exif_orientation.html
func FixOrientation(img *image.Image, orientation int) error {
	switch orientation {
	case 1:
		// do nothing
	case 2:
		// flop!
		*img = imaging.FlipH(*img)
	case 3:
		// rotate!(180)
		*img = imaging.Rotate180(*img)
	case 4:
		// flip!
		*img = imaging.FlipV(*img)
	case 5:
		// transpose!
		*img = imaging.Rotate270(*img)
		*img = imaging.FlipH(*img)
	case 6:
		// rotate!(90)
		*img = imaging.Rotate270(*img)
	case 7:
		// transverse!
		*img = imaging.Rotate90(*img)
		*img = imaging.FlipH(*img)
	case 8:
		// rotate!(270)
		*img = imaging.Rotate90(*img)
	default:
		return fmt.Errorf("Invalid orientation %d", orientation)
	}
	return nil
}

func OpenImage(path string) (*image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func ScaleAspect(image_w, image_h, aspect_w, aspect_h int) (int, int) {
	ratio_image := float64(image_w) / float64(image_h)
	ratio_aspect := float64(aspect_w) / float64(aspect_h)

	var width, height int

	if ratio_image < ratio_aspect {
		width = image_w
		h := float64(aspect_h) * float64(image_w) / float64(aspect_w)
		height = Round(h)
	} else {
		w := float64(aspect_w) * float64(image_h) / float64(aspect_h)
		width = Round(w)
		height = image_h
	}

	return width, height
}

func Round(f float64) int {
	r := math.Floor(f + .5)
	return int(r)
}

func GetAspectLab(gidx *model.Gidx, aspect *model.Aspect) ([]*model.Lab, error) {
	img, err := OpenImage(gidx.Path)
	if err != nil {
		return nil, err
	}

	if gidx.Orientation != 1 {
		err = FixOrientation(img, gidx.Orientation)
		if err != nil {
			return nil, err
		}
	}

	w, h := ScaleAspect(int(gidx.Width), int(gidx.Height), aspect.Columns, aspect.Rows)

	aspectImg := imaging.Fill((*img), w, h, imaging.Center, imaging.Lanczos)
	dataImg := imaging.Resize(aspectImg, DATA_SIZE, DATA_SIZE, imaging.Lanczos)

	labs := make([]*model.Lab, DATA_SIZE*DATA_SIZE)

	for y := 0; y < DATA_SIZE; y++ {
		for x := 0; x < DATA_SIZE; x++ {
			lab := model.RgbaToLab(dataImg.At(x, y))
			labs[y*DATA_SIZE+x] = lab
		}
	}

	return labs, nil
}