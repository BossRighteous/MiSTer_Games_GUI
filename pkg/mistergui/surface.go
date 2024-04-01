package mistergui

import (
	"image"
	"image/color"
	"image/draw"
)

const (
	bytesPerPixel byte = 3
)

type Surface struct {
	Width     uint16
	Height    uint16
	Interlace byte
	bgr8      []byte
	image     *image.RGBA
}

func (surface *Surface) UpdateFromImage(img image.Image, rect image.Rectangle) {
	// void
	draw.Draw(surface.image, rect, img, image.Point{0, 0}, draw.Over)
	surface.redrawBGR8()
}

func (surface *Surface) Fill(r, g, b uint8) {
	for y := range surface.image.Bounds().Dy() {
		for x := range surface.image.Bounds().Dx() {
			surface.image.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	surface.redrawBGR8()
}

func (surface *Surface) redrawBGR8() {
	var offset int = 0
	for y := range surface.image.Bounds().Dy() {
		for x := range surface.image.Bounds().Dx() {
			r, g, b, _ := surface.image.At(x, y).RGBA()
			surface.bgr8[offset] = uint8(b)
			surface.bgr8[offset+1] = uint8(g)
			surface.bgr8[offset+2] = uint8(r)
			offset += 3
		}
	}
}

func (surface *Surface) BGRbytes() []byte {
	return surface.bgr8[:]
}

func NewSurface(width, height uint16, interlace byte) *Surface {
	surface := Surface{
		Width:     width,
		Height:    height,
		Interlace: interlace,
	}
	var targetSize uint32 = uint32(width) * uint32(height) * uint32(bytesPerPixel)
	surface.bgr8 = make([]byte, targetSize)
	rect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(width), int(height)},
	}
	surface.image = image.NewRGBA(rect)
	surface.Fill(0, 0, 0)
	return &surface
}
