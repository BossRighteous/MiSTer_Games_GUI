package mistergui

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand/v2"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

type RGB struct {
	r uint8
	g uint8
	b uint8
}

type GUI struct {
	bgColor   RGB
	surface   *Surface
	modeline  *groovymister.Modeline
	fontImage *image.NRGBA
	psImage   *image.NRGBA
}

func (gui *GUI) Setup(surface *Surface, modeline *groovymister.Modeline) {
	fmt.Println("setting up GUI")
	gui.surface = surface
	gui.modeline = modeline
	gui.bgColor.b = uint8(rand.IntN(255))
	gui.bgColor.g = uint8(rand.IntN(255))
	gui.bgColor.r = uint8(rand.IntN(255))
	text := []string{
		"",
		"",
		"Integer sed est consequat augue scelerisque mollis in at est.",
		"Nam nec augue facilisis, accumsan turpis vitae, elementum quam.",
		"Nullam volutpat maximus ex posuere euismod.",
		"Vivamus nulla nulla, varius ac augue et, vehicula sollicitudin lectus.",
		"Curabitur vel est quis velit mattis sodales.",
		"Donec semper urna eu efficitur facilisis.",
		"Ut rhoncus interdum quam quis malesuada.",
	}
	gui.fontImage = DrawText(text)
	gui.psImage = PowerstoneImg
	gui.OnTick(0, 0.0)
}

func (gui *GUI) OnTick(frameCount uint32, delta float64) {
	gui.surface.Fill(gui.bgColor.b, gui.bgColor.g, gui.bgColor.r)
	gui.surface.DrawImage(gui.psImage, gui.psImage.Bounds(), image.Point{0, 0}, draw.Over)
	gui.surface.DrawImage(gui.psImage, gui.psImage.Bounds(), image.Point{0, 0}, draw.Over)
	gui.surface.DrawImage(gui.psImage, gui.psImage.Bounds(), image.Point{0, 0}, draw.Over)
	gui.surface.DrawImage(gui.fontImage, gui.fontImage.Bounds(), image.Point{0, 0}, draw.Over)

	fpsInt := math.Floor(1 / delta)
	fpsImg := DrawText([]string{fmt.Sprintf("%v", fpsInt)})
	gui.surface.DrawImage(fpsImg, fpsImg.Bounds(), image.Point{0, 0}, draw.Over)

}

func (gui *GUI) TearDown() {
	//gui.surface.Fill(gui.bgColor.b, gui.bgColor.g, gui.bgColor.r)
}

func NewGUI() *GUI {
	return &GUI{}
}
