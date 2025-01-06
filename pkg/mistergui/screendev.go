package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mrext"
)

/*
Dev testing screen
*/

type ScreenDev struct {
	name     string
	parent   IScreen
	guiState *GUIState
	list     *List
}

func (screen *ScreenDev) GUIState() *GUIState {
	return screen.guiState
}

func (screen *ScreenDev) Parent() IScreen {
	return screen.parent
}

func (screen *ScreenDev) Name() string {
	return screen.name
}

func (screen *ScreenDev) Setup() {
	var items []IListItem

	fmt.Println("ScreenDev Setup", screen.name)

	screen.list = NewList(screen, screen.guiState, make([]IListItem, 0), 0)

	if screen.parent != nil {
		item := &BasicListItem{
			label:        fmt.Sprintf("Back to %v", screen.parent.Name()),
			list:         screen.list,
			buttonsLabel: "A: Accept",
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				fmt.Println("Back to Parent", screen.parent.Name())
				screen.guiState.PopScreen()
			}
		}
		items = append(items, item)
	}

	{
		// MGL Test (real console)
		item := &BasicListItem{
			label:        fmt.Sprintf("Print MGL Test"),
			list:         screen.list,
			buttonsLabel: "A: Accept",
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				fmt.Println(mrext.GetSampleMgl())
			}
		}
		items = append(items, item)
	}

	{
		// Scan Console Test
		item := &BasicListItem{
			label:        fmt.Sprintf("Console Test"),
			list:         screen.list,
			buttonsLabel: "A: Accept",
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				outBuffer := make(chan string, 1024)
				completedOk := make(chan bool)
				consoleState := ConsoleState{
					outBuffer:   outBuffer,
					completedOk: completedOk,
				}
				newScreen := &ScreenConsole{
					parent:       screen,
					guiState:     screen.guiState,
					name:         item.Label(),
					consoleState: consoleState,
				}
				newScreen.Setup()
				screen.guiState.PushScreen(newScreen)
				screen.guiState.IsChanged = true

				go func() {
					i := 0
					for i < 100000 {
						i++
						consoleState.outBuffer <- fmt.Sprint(i)
					}
					completedOk <- true
				}()
			}
		}
		items = append(items, item)
	}

	for i := 0; i < 100; i++ {
		item := &BasicListItem{
			label:        fmt.Sprintf("Item %v", i),
			list:         screen.list,
			buttonsLabel: fmt.Sprintf("A: Accept, Item %v", i),
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				fmt.Println("OnSelect item", item.Label())
				newScreen := &ScreenDev{}
				newScreen.Setup()
				screen.guiState.PushScreen(newScreen)
				screen.guiState.IsChanged = true
			}

		}
		items = append(items, item)
	}
	screen.list.ReplaceItems(items)
}

func (screen *ScreenDev) OnEnter() {
	fmt.Println("ScreenDev OnEnter", screen.name)
	screen.list.Render()
}

func (screen *ScreenDev) OnExit() {
	fmt.Println("ScreenDev OnExit", screen.name)

}

func (screen *ScreenDev) OnTick(tick TickData) {
	list := screen.list
	if list == nil {
		return
	}

	list.OnTick()

	//fmt.Println("ScreenDev OnTick", tick.FrameCount)
	//fmt.Println(screen.list.CurrentItem().Label())
	//screen.list.NextItem()
	//screen.list.CurrentItem().OnSelect()
	//screen.guiState.IsChanged = true
}

func (screen *ScreenDev) Render() {
	//fmt.Println("ScreenDev Render")
	screen.list.Render()
}

func (screen *ScreenDev) TearDown() {
	fmt.Println("ScreenDev OnTearDown", screen.name)
}
