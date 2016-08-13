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
			return "", err
		}
		if n == 0 {
			break
		}
		if _, err := io.WriteString(hash, string(buf[:n])); err != nil {
			return "", err
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
	// no exif data
	if err != nil {
		return 1, nil
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return 1, nil
	}

	val, err := orientation.Int(0)
	if err != nil {
		return 1, nil
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

func GetAspectLab(i model.Image, aspect *model.Aspect) ([]*model.Lab, error) {
	img, err := OpenImage(i.GetPath())
	if err != nil {
		return nil, err
	}

	if i.GetOrientation() != 1 {
		err = FixOrientation(img, i.GetOrientation())
		if err != nil {
			return nil, err
		}
	}

	return GetImgAspectLab(img, i, aspect)
}

func GetImgAspectLab(img *image.Image, i model.Image, aspect *model.Aspect) ([]*model.Lab, error) {
	w, h := ScaleAspect(int(i.GetWidth()), int(i.GetHeight()), aspect.Columns, aspect.Rows)

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

func GetPartialLab(i model.Image, coverPartial *model.CoverPartial) ([]*model.Lab, error) {
	img, err := OpenImage(i.GetPath())
	if err != nil {
		return nil, err
	}

	if i.GetOrientation() != 1 {
		err = FixOrientation(img, i.GetOrientation())
		if err != nil {
			return nil, err
		}
	}

	return GetImgPartialLab(img, i, coverPartial)
}

func GetImgPartialLab(img *image.Image, i model.Image, coverPartial *model.CoverPartial) ([]*model.Lab, error) {
	cropImg := imaging.Crop((*img), coverPartial.Rectangle())
	dataImg := imaging.Resize(cropImg, DATA_SIZE, DATA_SIZE, imaging.Lanczos)

	labs := make([]*model.Lab, DATA_SIZE*DATA_SIZE)

	for y := 0; y < DATA_SIZE; y++ {
		for x := 0; x < DATA_SIZE; x++ {
			lab := model.RgbaToLab(dataImg.At(x, y))
			labs[y*DATA_SIZE+x] = lab
		}
	}

	return labs, nil
}

func FillAspect(img *image.Image, aspectWidth, aspectHeight int) *image.Image {
	bounds := (*img).Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	w, h := ScaleAspect(width, height, aspectWidth, aspectHeight)

	var myImg image.Image = imaging.Fill((*img), w, h, imaging.Center, imaging.Lanczos)
	return &myImg
}

func SliceContainsInt64(s []int64, a int64) bool {
	for _, b := range s {
		if a == b {
			return true
		}
	}
	return false
}
