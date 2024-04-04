package mistergui

import (
	"image"
	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func parseFont(fontBytes []byte) *truetype.Font {
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	return font
}

var Roboto *truetype.Font = parseFont(*Embeds.RobotoBold)

func DrawText(text []string, rect image.Rectangle, bg *image.Uniform) *image.NRGBA {
	var (
		dpi     float64        = 72
		hinting string         = "full"
		size    float64        = 12
		spacing float64        = 1.5
		ttfont  *truetype.Font = Roboto
	)

	fg := image.Black
	rgba := image.NewNRGBA(rect)
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{0, 0}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(ttfont)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, _ = c.DrawString(string(s), pt)
		pt.Y += c.PointToFixed(size * spacing)
	}

	return rgba
}
