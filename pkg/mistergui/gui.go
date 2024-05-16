package mistergui

import (
	"fmt"
	"image"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mister"
)

type TickData struct {
	FrameCount  uint32
	Delta       float64
	InputPacket groovymister.GroovyInputPacket
}

// mgdb version later
type Game struct {
	Name string
	Path string
}

type GUIState struct {
	Screen    Screen
	Screens   *Screens
	IsChanged bool
	IsLoading bool
	Core      *mister.Core
	Cores     *[]mister.Core
	Game      *Game
	Modal     *Modal
	AsyncChan chan AsyncCallback
	Surface   *Surface
	Input     *groovymister.GroovyInput
	Modeline  *groovymister.Modeline
}

func (state *GUIState) ChangeScreen(newScreen Screen) {
	if state.Screen != nil {
		state.Screen.OnExit()
	}
	state.Screen = newScreen
	state.Screen.OnEnter()
}

type GUI struct {
	State *GUIState

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

	surface := NewSurface(modeline.HActive, modeline.VActive, modeline.Interlace)
	bgColor := ColorBGR8{uint8(104), uint8(66), uint8(13)}
	surface.FillBg(bgColor)
	p0 := image.Point{0, 0}
	surface.Erase(surface.BgImage.Bounds(), p0)

	cores := mister.GetCoresFromJSON()

	gui.State = &GUIState{
		Screen: nil,
		Screens: &Screens{
			Cores: &ScreenCores{},
			Games: &ScreenGames{},
			Roms:  &ScreenRoms{},
		},
		IsChanged: false,
		Core:      &cores[0],
		Cores:     &cores,
		Game:      nil,
		Modal:     nil,
		AsyncChan: gui.AsyncCallbackChan,
		Surface:   surface,
		Input:     &groovymister.GroovyInput{},
		Modeline:  modeline,
	}

	gui.State.Screens.Cores.Setup(gui.State)
	gui.State.Screens.Games.Setup(gui.State)
	gui.State.Screens.Roms.Setup(gui.State)

	gui.State.ChangeScreen(gui.State.Screens.Games)

	//gui.psImage = ListingBg
	//draw.Draw(gui.surface.Image, gui.surface.Image.Bounds(), *gui.psImage, p0, draw.Over)
	//gui.fontImage = DrawText(text, gui.surface.Image.Bounds(), image.Transparent)
	//draw.Draw(gui.surface.Image, gui.surface.Image.Bounds(), gui.fontImage, p0, draw.Over)
	//surface.Erase(gui.surface.BgImage.Bounds(), p0)
	/*go func() {
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
	}()*/

	gui.FrameBufferChan <- gui.State.Surface.BGRbytes(true)
}

func (gui *GUI) OnTick(tick TickData) {
	//gui.redraw = true
	//fpsInt := math.Floor(1 / tick.Delta)
	//fmt.Printf("%v fps", fpsInt)
	//fpsImg := DrawText([]string{fmt.Sprintf("%v", fpsInt)}, image.Rect(0, 0, 40, 30), image.White)
	//draw.Draw(gui.surface.Image, fpsImg.Bounds(), fpsImg, P0, draw.Src)

	// Observe inputs
	gui.State.Input.AddInputPacket(tick.InputPacket)

	if gui.State.Modal != nil {
		modal := *gui.State.Modal
		if gui.State.IsChanged {
			modal.Render()
		}

		if gui.State.Input.IsJustPressed(1, groovymister.InputB1) {
			modal.OnAccept()
			gui.State.Modal = nil
			return
		} else if gui.State.Input.IsJustPressed(1, groovymister.InputB2) {
			modal.OnReject()
			gui.State.Modal = nil
			return
		}
		return
	}

	gui.State.Screen.OnTick(tick)

	if gui.State.IsChanged {
		gui.State.Screen.Render()
		gui.FrameBufferChan <- gui.State.Surface.BGRbytes(true)
		gui.State.IsChanged = false
	}
}

func (gui *GUI) TearDown() {
	if gui.State.Screen != nil {
		gui.State.Screen.TearDown()
	}
}

func listen(gui *GUI) {
	for {
		select {
		case <-gui.QuitChan:
			fmt.Println("gui.quitChan recv, closing goroutine")
			return
		case tickData := <-gui.TickChan:
			gui.OnTick(tickData)
		case promiseFn := <-gui.AsyncCallbackChan:
			gui.State.IsChanged = true
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
