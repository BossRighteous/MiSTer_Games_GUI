package mistergui

import (
	"bytes"
	"image"

	"github.com/disintegration/imaging"
)

func DecodeImageBytes(imageBytes *[]byte) (*image.Image, error) {
	reader := bytes.NewReader(*imageBytes)
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func DecodeImageBytesFit(imageBytes *[]byte) (*image.Image, error) {
	reader := bytes.NewReader(*imageBytes)
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}
	img = imaging.Fit(img, 320, 240, imaging.Lanczos)
	return &img, nil
}
