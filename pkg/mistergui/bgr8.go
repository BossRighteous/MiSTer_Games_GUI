package mistergui

import (
	"fmt"
	"image"
	"image/color"
	"math/bits"
)

/*
	Custom Image interface for GrovyMister's BGR byte ordering
*/

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
//
// This panics instead of returning an error because of backwards
// compatibility. The NewXxx functions do not return an error.
func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

// Define color interface

type ColorBGR8 struct {
	B, G, R uint8
}

func (c ColorBGR8) RGBA() (r, g, b, a uint32) {
	a = 0xff
	r = uint32(c.R)
	r |= r << 8
	r *= uint32(a)
	r /= 0xff
	g = uint32(c.G)
	g |= g << 8
	g *= uint32(a)
	g /= 0xff
	b = uint32(c.B)
	b |= b << 8
	b *= uint32(a)
	b /= 0xff
	a = uint32(a)
	a |= a << 8
	return
}

var BGR8Model color.Model = color.ModelFunc(bgr8Model)

func bgr8Model(c color.Color) color.Color {
	return ColorToBGR8(c)
}

func ColorToBGR8(c color.Color) ColorBGR8 {
	if co, ok := c.(ColorBGR8); ok {
		return co
	}
	r, g, b, a := c.RGBA()
	fmt.Println(r, g, b, a)
	if a == 0xffff {
		fmt.Println("full alpha")
		return ColorBGR8{uint8(b >> 8), uint8(g >> 8), uint8(r >> 8)}
	}
	if a == 0 {
		fmt.Println("no alpha")
		return ColorBGR8{0, 0, 0}
	}
	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	r = (r * 0xffff) / a
	g = (g * 0xffff) / a
	b = (b * 0xffff) / a

	fmt.Println(ColorBGR8{uint8(b >> 8), uint8(g >> 8), uint8(r >> 8)})
	return ColorBGR8{uint8(b >> 8), uint8(g >> 8), uint8(r >> 8)}
}

var BGR8BytesPerPixel int = 3

// BGR8 is an in-memory image whose At method returns color.BGR8 values.
type BGR8 struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *BGR8) ColorModel() color.Model { return BGR8Model }

func (p *BGR8) Bounds() image.Rectangle { return p.Rect }

func (p *BGR8) At(x, y int) color.Color {
	return p.BGR8At(x, y)
}

func (p *BGR8) RGBA64At(x, y int) color.RGBA64 {
	r, g, b, a := p.BGR8At(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (p *BGR8) BGR8At(x, y int) ColorBGR8 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return ColorBGR8{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+BGR8BytesPerPixel : i+BGR8BytesPerPixel] // Small cap improves performance, see https://golang.org/issue/27857
	return ColorBGR8{s[0], s[1], s[2]}
}

func (p *BGR8) NRGBAAt(x, y int) color.NRGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.NRGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+BGR8BytesPerPixel : i+BGR8BytesPerPixel] // Small cap improves performance, see https://golang.org/issue/27857
	return color.NRGBA{s[2], s[1], s[0], 0xff}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *BGR8) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*BGR8BytesPerPixel
}

func (p *BGR8) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := BGR8Model.Convert(c).(ColorBGR8)
	s := p.Pix[i : i+BGR8BytesPerPixel : i+BGR8BytesPerPixel] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.B
	s[1] = c1.G
	s[2] = c1.R
}

func (p *BGR8) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	r, g, b, a := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	if (a != 0) && (a != 0xffff) {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+BGR8BytesPerPixel : i+BGR8BytesPerPixel] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(b >> 8)
	s[1] = uint8(g >> 8)
	s[2] = uint8(r >> 8)
}

func (p *BGR8) SetNRGBA(x, y int, c color.NRGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+BGR8BytesPerPixel : i+BGR8BytesPerPixel] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.B
	s[1] = c.G
	s[2] = c.R
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *BGR8) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &BGR8{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &BGR8{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *BGR8) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	return true
}

// NewBGR8 returns a new BGR8 image with the given bounds.
func NewBGR8(r image.Rectangle) *BGR8 {
	return &BGR8{
		Pix:    make([]uint8, pixelBufferLength(BGR8BytesPerPixel, r, "BGR8")),
		Stride: BGR8BytesPerPixel * r.Dx(),
		Rect:   r,
	}
}
