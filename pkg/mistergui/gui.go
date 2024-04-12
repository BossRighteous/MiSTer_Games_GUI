package mistergui

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

type State uint8

const (
	StateInit State = iota
	StateWaiting
	StateWorking
	StateGameLoading
)

type Screen uint8

const (
	ScreenCores Screen = iota
	ScreenGameListing
	ScreenGameScreenshot
	ScreenGameMeta
)

type TickData struct {
	FrameCount uint32
	Delta      float64
}

type GUI struct {
	bgColor   ColorBGR8
	surface   *Surface
	modeline  *groovymister.Modeline
	fontImage *image.NRGBA
	psImage   *image.Image
	state     State
	screen    Screen
	redraw    bool

	TickChan chan TickData
	QuitChan chan bool
	// Allow channel based mutation via callbacks to avoid races
	AsyncCallbackChan chan AsyncCallback
	FrameBufferChan   chan []uint8
}

type AsyncCallback func(gui *GUI)

var P0 image.Point = image.Point{0, 0}

func (gui *GUI) Setup(modeline *groovymister.Modeline) {
	fmt.Println("setting up GUI")

	gui.surface = NewSurface(modeline.HActive, modeline.VActive, modeline.Interlace)
	gui.modeline = modeline
	gui.bgColor = ColorBGR8{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
	fmt.Println(gui.bgColor)
	gui.surface.FillBg(gui.bgColor)
	p0 := image.Point{0, 0}
	gui.surface.Erase(gui.surface.BgImage.Bounds(), p0)
	/*text := []string{
		"",
		"",
		"Integer sed est consequat augue scelerisque mollis in at est.",
		"Nam nec augue facilisis, accumsan turpis vitae, elementum quam.",
		"Nullam volutpat maximus ex posuere euismod.",
		"Vivamus nulla nulla, varius ac augue et, vehicula sollicitudin lectus.",
		"Curabitur vel est quis velit mattis sodales.",
		"Donec semper urna eu efficitur facilisis.",
		"Ut rhoncus interdum quam quis malesuada.",
	}*/
	gui.psImage = PowerstoneImg
	draw.Draw(gui.surface.Image, gui.surface.Image.Bounds(), *gui.psImage, p0, draw.Over)
	//gui.fontImage = DrawText(text, gui.surface.Image.Bounds(), image.Transparent)
	//draw.Draw(gui.surface.Image, gui.surface.Image.Bounds(), gui.fontImage, p0, draw.Over)
	//surface.Erase(gui.surface.BgImage.Bounds(), p0)
	gui.FrameBufferChan <- gui.surface.BGRbytes(true)

	go func() {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		metaImages := LoadMetaImages(
			"games/.mistergamesgui/N64/citytourgpzennihongtsenshukenjapan.json",
			gui.surface.BgImage.Rect,
		)
		gui.AsyncCallbackChan <- func(gui *GUI) {
			if len(metaImages) > 0 {
				fmt.Println("can draw")
				draw.Draw(gui.surface.Image, gui.surface.Image.Bounds(), &metaImages[0], p0, draw.Over)
			}
		}
	}()
}

func (gui *GUI) OnTick(frameCount uint32, delta float64) {
	gui.redraw = true
	fpsInt := math.Floor(1 / delta)
	fmt.Printf("%v fps", fpsInt)
	fpsImg := DrawText([]string{fmt.Sprintf("%v", fpsInt)}, image.Rect(0, 0, 40, 30), image.White)
	draw.Draw(gui.surface.Image, fpsImg.Bounds(), fpsImg, P0, draw.Src)

	if gui.redraw {
		gui.redraw = false
		gui.FrameBufferChan <- gui.surface.BGRbytes(true)
	}
}

func (gui *GUI) TearDown() {
	//gui.surface.Fill(gui.bgColor.b, gui.bgColor.g, gui.bgColor.r)
}

func listen(gui *GUI) {
	for {
		select {
		case <-gui.QuitChan:
			fmt.Println("gui.quitChan recv, closing goroutine")
			return
		case tickData := <-gui.TickChan:
			gui.OnTick(tickData.FrameCount, tickData.Delta)
		case promiseFn := <-gui.AsyncCallbackChan:
			gui.redraw = true
			promiseFn(gui)
		}
	}
}

func NewGUI() *GUI {
	gui := &GUI{}
	// Receieve
	gui.TickChan = make(chan TickData, 1)
	gui.QuitChan = make(chan bool, 1)
	// Local loopback
	gui.AsyncCallbackChan = make(chan AsyncCallback, 10)
	// Send
	gui.FrameBufferChan = make(chan []uint8, 1)
	go listen(gui)
	return gui
}
