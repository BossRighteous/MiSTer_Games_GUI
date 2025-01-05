package mistergui

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

/*
Screen Console is for streaming status lines during async work, much like a stdout console stream
*/

type ConsoleState struct {
	outBuffer   chan string
	completedOk chan bool
	canExit     bool
}

type ScreenConsole struct {
	name         string
	parent       IScreen
	guiState     *GUIState
	textLines    []string
	consoleState ConsoleState
}

func (screen *ScreenConsole) GUIState() *GUIState {
	return screen.guiState
}

func (screen *ScreenConsole) Parent() IScreen {
	return screen.parent
}

func (screen *ScreenConsole) Name() string {
	return screen.name
}

func (screen *ScreenConsole) Setup() {
	fmt.Println("ScreenConsole Setup", screen.name)
	screen.textLines = make([]string, 10)
}

func (screen *ScreenConsole) OnEnter() {
	fmt.Println("ScreenConsole OnEnter", screen.name)
}

func (screen *ScreenConsole) OnExit() {
	fmt.Println("ScreenConsole OnExit", screen.name)
}

func (screen *ScreenConsole) OnTick(tick TickData) {
	if !screen.consoleState.canExit {
		for len(screen.consoleState.outBuffer) > 0 {
			// Receive all
			line := <-screen.consoleState.outBuffer
			if len(screen.consoleState.outBuffer) <= 10 {
				// Only push tailend from tick
				screen.PushText(line)
			}
		}

		select {
		case ok := <-screen.consoleState.completedOk:
			fmt.Println("got completedOk", ok)
			screen.consoleState.canExit = true
			if ok {
				screen.PushText("Process completed successfully")
			} else {
				screen.PushText("Process completed with errors")
			}
			screen.PushText("Press A to return to menu...")
		default:
		}
	} else {
		if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
			screen.guiState.PopScreen()
		}
	}
}

func (screen *ScreenConsole) PushText(newLine string) {
	for i := range screen.textLines {
		if i < 9 {
			screen.textLines[i] = screen.textLines[i+1]
		} else {
			screen.textLines[9] = newLine
		}
	}
	screen.guiState.IsChanged = true
}

func (screen *ScreenConsole) Render() {
	//fmt.Println("ScreenConsole Render")
	surface := screen.guiState.Surface
	surface.Erase(surface.Image.Rect, P0)
	img := DrawText(screen.textLines, surface.Image.Rect, image.Transparent)
	draw.Draw(surface.Image, surface.Image.Rect, img, P0, draw.Over)
}

func (screen *ScreenConsole) TearDown() {
	fmt.Println("ScreenConsole OnTearDown", screen.name)
}
