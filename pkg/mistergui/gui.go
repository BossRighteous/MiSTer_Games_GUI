package mistergui

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

type GUI struct {
	bgColor   ColorBGR8
	surface   *Surface
	modeline  *groovymister.Modeline
	fontImage *image.NRGBA
	psImage   *image.Image
}

func (gui *GUI) Setup(surface *Surface, modeline *groovymister.Modeline) {
	fmt.Println("setting up GUI")
	gui.surface = surface
	gui.modeline = modeline
	gui.bgColor = ColorBGR8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	fmt.Println(gui.bgColor)
	gui.surface.FillBg(gui.bgColor)
	p0 := image.Point{0, 0}
	surface.Erase(gui.surface.BgImage.Bounds(), p0)
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
	gui.psImage = PowerstoneImg
	draw.Draw(gui.surface.Image, gui.surface.Image.Bounds(), *gui.psImage, p0, draw.Over)
	gui.fontImage = DrawText(text, gui.surface.Image.Bounds(), image.Transparent)
	draw.Draw(gui.surface.Image, gui.surface.Image.Bounds(), gui.fontImage, p0, draw.Over)
	//surface.Erase(gui.surface.BgImage.Bounds(), p0)
}

func (gui *GUI) OnTick(frameCount uint32, delta float64) {
	fpsInt := math.Floor(1 / delta)
	fmt.Printf("%v fps", fpsInt)
	fpsImg := DrawText([]string{fmt.Sprintf("%v", fpsInt)}, image.Rect(0, 0, 40, 30), image.White)
	draw.Draw(gui.surface.Image, fpsImg.Bounds(), fpsImg, image.Point{0, 0}, draw.Src)
}

func (gui *GUI) TearDown() {
	//gui.surface.Fill(gui.bgColor.b, gui.bgColor.g, gui.bgColor.r)
}

func NewGUI() *GUI {
	return &GUI{}
}
