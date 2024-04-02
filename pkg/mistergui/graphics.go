package mistergui

import (
	"bytes"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

func DecodeImageBytes(imageBytes *[]byte) *image.NRGBA {
	reader := bytes.NewReader(*imageBytes)
	img, err := imaging.Decode(reader)
	if err != nil {
		panic("cant decode image from bytes")
	}
	img = imaging.Fit(img, 320, 240, imaging.Lanczos)
	rgba, _ := img.(*image.NRGBA)
	bounds := img.Bounds()
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			r, g, b, _ := rgba.At(x, y).RGBA()
			rgba.SetNRGBA(x, y, color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(20)})
		}
	}
	return rgba
}

var PowerstoneImg = DecodeImageBytes(Embeds.Powerstone)
