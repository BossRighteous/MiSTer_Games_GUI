package mistergui

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/settings"
)

type TickData struct {
	FrameCount  uint32
	InputPacket groovymister.GroovyInputPacket
}

type GUIState struct {
	Screen    IScreen
	IsChanged bool
	IsLoading bool
	Surface   *Surface
	Input     *groovymister.GroovyInput
	Modeline  *groovymister.Modeline
	QuitChan  chan bool
	Settings  *settings.Settings
}

func (state *GUIState) PushScreen(newScreen IScreen) {
	if state.Screen != nil {
		state.Screen.OnExit()
	}
	state.Screen = newScreen
	state.Screen.OnEnter()
}

func (state *GUIState) PopScreen() {
	if state.Screen != nil {
		state.Screen.OnExit()
	}
	parent := state.Screen.Parent()
	if parent != nil {
		state.Screen = parent
		state.Screen.OnEnter()
	}
	state.IsChanged = true
}

type GUI struct {
	State     *GUIState
	QuitChan  chan bool
	UdpClient *groovymister.UdpClient
}

type AsyncCallback func(gui *GUI)

var P0 image.Point = image.Point{0, 0}

func (gui *GUI) Setup(
	modeline *groovymister.Modeline,
	settings *settings.Settings,
	udpClient *groovymister.UdpClient,
	quitChan chan bool,
) {
	fmt.Println("setting up GUI")
	gui.UdpClient = udpClient
	gui.QuitChan = quitChan
	p0 := image.Point{0, 0}
	surface := NewSurface(modeline.HActive, modeline.VActive, modeline.Interlace)
	bgColor := ColorBGR8{uint8(104), uint8(66), uint8(13)}
	surface.FillBg(bgColor)
	listingBgImg, err := DecodeImageBytes(Embeds.ListingBg)
	if err == nil {
		draw.Draw(surface.BgImage, surface.Image.Bounds(), *listingBgImg, p0, draw.Over)
	}

	surface.Erase(surface.BgImage.Bounds(), p0)

	gui.State = &GUIState{
		Screen:    nil,
		IsChanged: false,
		Surface:   surface,
		Input:     &groovymister.GroovyInput{},
		Modeline:  modeline,
		QuitChan:  gui.QuitChan,
		Settings:  settings,
	}

	rootScreen := &ScreenCollections{parent: nil, guiState: gui.State, name: "Root"}
	rootScreen.Setup()
	gui.State.PushScreen(rootScreen)

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

	gui.UdpClient.CmdBlit(gui.State.Surface.BGRbytes(true))
}

func (gui *GUI) OnTick(tick TickData) {
	//fpsInt := math.Floor(1 / tick.Delta)
	//fmt.Printf("%v fps", fpsInt)
	//fpsImg := DrawText([]string{fmt.Sprintf("%v", fpsInt)}, image.Rect(0, 0, 40, 30), image.White)
	//draw.Draw(gui.surface.Image, fpsImg.Bounds(), fpsImg, P0, draw.Src)
	// Observe inputs
	gui.State.Input.AddInputPacket(tick.InputPacket)

	gui.State.Screen.OnTick(tick)

	if gui.State.IsChanged {
		gui.State.Screen.Render()
		gui.UdpClient.CmdBlit(gui.State.Surface.BGRbytes(true))
		gui.State.IsChanged = false
	}
}

func (gui *GUI) TearDown() {
	if gui.State.Screen != nil {
		gui.State.Screen.TearDown()
	}
}
