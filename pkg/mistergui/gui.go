package mistergui

import (
	"fmt"
	"math/rand/v2"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

type RGB struct {
	r uint8
	g uint8
	b uint8
}

type GUI struct {
	bgColor  RGB
	surface  *Surface
	modeline *groovymister.Modeline
}

func (gui *GUI) Setup(surface *Surface, modeline *groovymister.Modeline) {
	fmt.Println("setting up GUI")
	gui.surface = surface
	gui.modeline = modeline
	gui.bgColor.b = uint8(rand.IntN(255))
	gui.bgColor.g = uint8(rand.IntN(255))
	gui.bgColor.r = uint8(rand.IntN(255))
	gui.surface.Fill(gui.bgColor.b, gui.bgColor.g, gui.bgColor.r)
}

func (gui *GUI) OnTick(frameCount uint32, delta float64) {
	gui.surface.Fill(gui.bgColor.b, gui.bgColor.g, gui.bgColor.r)
}

func (gui *GUI) TearDown() {
	gui.surface.Fill(gui.bgColor.b, gui.bgColor.g, gui.bgColor.r)
}

func NewGUI() *GUI {
	return &GUI{}
}
