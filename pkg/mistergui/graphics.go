package mistergui

import (
	"bytes"
	"image"

	"github.com/disintegration/imaging"
)

func DecodeImageBytes(imageBytes *[]byte) *image.Image {
	reader := bytes.NewReader(*imageBytes)
	img, err := imaging.Decode(reader)
	if err != nil {
		panic("cant decode image from bytes")
	}
	img = imaging.Fit(img, 320, 240, imaging.Lanczos)
	return &img
}

//var ListingBg = DecodeImageBytes(Embeds.ListingBg)
