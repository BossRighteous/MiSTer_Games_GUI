package mistergui

import (
	"image"
	"image/color"
	"image/draw"
)

type Surface struct {
	Width     uint16
	Height    uint16
	Interlace bool
	Image     *BGR8
	BgImage   *BGR8
}

func fillImage(img *BGR8, c color.Color) {
	bounds := img.Bounds()
	pixels := bounds.Dx() * bounds.Dy()
	cBGR := ColorToBGR8(c)
	for i := 0; i < pixels*BGR8BytesPerPixel; i += BGR8BytesPerPixel {
		img.Pix[i] = cBGR.B
		img.Pix[i+1] = cBGR.G
		img.Pix[i+2] = cBGR.R
	}
}

func (surface *Surface) Fill(c color.Color) {
	fillImage(surface.Image, c)
}

func (surface *Surface) FillBg(c color.Color) {
	fillImage(surface.BgImage, c)
}

func (surface *Surface) Erase(rect image.Rectangle, sp image.Point) {
	draw.Draw(surface.Image, rect, surface.BgImage, sp, draw.Src)
}

func (surface *Surface) BGRbytes(_ bool) []byte {
	//fmt.Println(surface.Image.Pix[3000:3100])
	return surface.Image.Pix
}

func NewSurface(width, height uint16, interlace bool) *Surface {
	surface := Surface{
		Width:     width,
		Height:    height,
		Interlace: interlace,
	}
	rect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(width), int(height)},
	}
	surface.Image = NewBGR8(rect)
	surface.BgImage = NewBGR8(rect)
	return &surface
}
