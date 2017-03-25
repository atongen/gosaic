package util

import (
	"crypto/md5"
	"fmt"
	"github.com/atongen/gosaic/model"
	"image"
	_ "image/jpeg"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
		return 1, err
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

	if val == 0 {
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

func OpenImg(i model.Image) (*image.Image, error) {
	return OpenImage(i.GetPath())
}

func OpenImage(path string) (*image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func GetAspectLab(i model.Image, aspect *model.Aspect) ([]*model.Lab, error) {
	img, err := OpenImg(i)
	if err != nil {
		return nil, err
	}

	if i.GetOrientation() != 1 {
		err = FixOrientation(img, i.GetOrientation())
		if err != nil {
			return nil, err
		}
	}

	return GetImgAspectLab(img, i, aspect), nil
}

func GetImgAspectLab(img *image.Image, i model.Image, aspect *model.Aspect) []*model.Lab {
	w, h := aspect.ScaleRound(int(i.GetWidth()), int(i.GetHeight()))

	aspectImg := imaging.Fill((*img), w, h, imaging.Center, imaging.Lanczos)
	dataImg := imaging.Resize(aspectImg, DATA_SIZE, DATA_SIZE, imaging.Lanczos)

	labs := make([]*model.Lab, DATA_SIZE*DATA_SIZE)

	for y := 0; y < DATA_SIZE; y++ {
		for x := 0; x < DATA_SIZE; x++ {
			lab := model.RgbaToLab(dataImg.At(x, y))
			labs[y*DATA_SIZE+x] = lab
		}
	}

	return labs
}

func GetImageCoverPartial(i model.Image, coverPartial *model.CoverPartial) (*image.Image, error) {
	img, err := OpenImg(i)
	if err != nil {
		return nil, err
	}

	if i.GetOrientation() != 1 {
		err = FixOrientation(img, i.GetOrientation())
		if err != nil {
			return nil, err
		}
	}

	return GetImgCoverPartial(img, coverPartial), nil
}

func GetImgCoverPartial(img *image.Image, coverPartial *model.CoverPartial) *image.Image {
	var myImg image.Image = imaging.Fill((*img), coverPartial.Width(), coverPartial.Height(), imaging.Center, imaging.Lanczos)
	return &myImg
}

func GetPartialLab(i model.Image, coverPartial *model.CoverPartial) ([]*model.Lab, error) {
	img, err := OpenImg(i)
	if err != nil {
		return nil, err
	}

	if i.GetOrientation() != 1 {
		err = FixOrientation(img, i.GetOrientation())
		if err != nil {
			return nil, err
		}
	}

	return GetImgPartialLab(img, coverPartial), nil
}

func GetImgPartialLab(img *image.Image, coverPartial *model.CoverPartial) []*model.Lab {
	cropImg := imaging.Crop((*img), coverPartial.Rectangle())
	dataImg := imaging.Resize(cropImg, DATA_SIZE, DATA_SIZE, imaging.Lanczos)

	labs := make([]*model.Lab, DATA_SIZE*DATA_SIZE)

	for y := 0; y < DATA_SIZE; y++ {
		for x := 0; x < DATA_SIZE; x++ {
			lab := model.RgbaToLab(dataImg.At(x, y))
			labs[y*DATA_SIZE+x] = lab
		}
	}

	return labs
}

func GetImgAvgDist(img *image.Image, coverPartial *model.CoverPartial) float64 {
	labs := GetImgPartialLab(img, coverPartial)
	avgLab := LabAvg(labs)
	dist := float64(0.0)
	for _, lab := range labs {
		dist += lab.Dist(avgLab)
	}
	return dist
}

func LabAvg(labs []*model.Lab) *model.Lab {
	if len(labs) == 0 {
		return &model.Lab{}
	}

	sL := float64(0.0)
	sA := float64(0.0)
	sB := float64(0.0)
	sAlpha := float64(0.0)

	for _, lab := range labs {
		sL += lab.L
		sA += lab.A
		sB += lab.B
		sAlpha += lab.Alpha
	}

	l := float64(len(labs))
	return &model.Lab{
		L:     sL / l,
		A:     sA / l,
		B:     sB / l,
		Alpha: sAlpha / l,
	}
}

func FillAspect(img *image.Image, aspect *model.Aspect) *image.Image {
	bounds := (*img).Bounds()

	w, h := aspect.Scale(bounds.Max.X, bounds.Max.Y)

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

func SliceContainsString(s []string, a string) bool {
	for _, b := range s {
		if a == b {
			return true
		}
	}
	return false
}

func Round(f float64) int {
	var r float64
	if f >= float64(0.0) {
		r = math.Floor(f + 0.5)
	} else {
		r = math.Ceil(f - 0.5)
	}
	return int(r)
}

var (
	cleanStr1Re = regexp.MustCompile("[^0-9a-z-]+")
	cleanStr2Re = regexp.MustCompile("(?:^_|_$)")
)

func CleanStr(s string) string {
	s = strings.ToLower(s)
	s = cleanStr1Re.ReplaceAllString(s, "_")
	s = strings.Replace(s, "_-_", "-", -1)
	s = strings.Replace(s, "_-", "-", -1)
	s = strings.Replace(s, "-_", "-", -1)
	return cleanStr2Re.ReplaceAllString(s, "")
}

func NextAvailableFilename(myPath string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(myPath))
	if err != nil {
		return "", err
	}

	myName := filepath.Base(myPath)
	ext := filepath.Ext(myName)
	baseName := myName[:len(myName)-len(ext)]

	var iName string
	for i := 0; ; i++ {
		if i == 0 {
			iName = myPath
		} else {
			iName = filepath.Join(dir, fmt.Sprintf("%s-%d%s", baseName, i, ext))
		}
		if _, err := os.Stat(iName); os.IsNotExist(err) {
			break
		}
	}

	return iName, nil
}
