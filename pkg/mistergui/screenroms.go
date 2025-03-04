package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
)

type ScreenRoms struct {
	name     string
	parent   IScreen
	guiState *GUIState
	list     *List
	client   *mgdb.MGDBClient
	game     mgdb.Game
}

func (screen *ScreenRoms) GUIState() *GUIState {
	return screen.guiState
}

func (screen *ScreenRoms) Parent() IScreen {
	return screen.parent
}

func (screen *ScreenRoms) Name() string {
	return screen.name
}

func (screen *ScreenRoms) Setup() {
	screen.name = "ROMs"

	var items []IListItem

	fmt.Println("ScreenRoms Setup", screen.name)

	screen.list = NewList(screen, screen.guiState, make([]IListItem, 0), 0)

	{
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

	client := screen.client

	roms, _ := client.GetIndexedRoms(screen.game.GameID)

	for _, rom := range roms {
		item := &BasicListItem{
			label:        fmt.Sprintf("%v%v", rom.FileName, rom.FileExt),
			list:         screen.list,
			buttonsLabel: "A: Choose System to Load ROM",
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				fmt.Println("OnSelect item", item.Label())
				fmt.Println("pressed B1")
				romSystemScreen := &ScreenRomSystems{
					parent:   screen,
					guiState: screen.guiState,
					name:     item.Label(),
					client:   screen.client,
					rom:      rom,
				}
				romSystemScreen.Setup()
				screen.guiState.PushScreen(romSystemScreen)
				screen.guiState.IsChanged = true
			}

		}
		items = append(items, item)
	}

	screen.list.ReplaceItems(items)

}

func (screen *ScreenRoms) OnEnter() {
	fmt.Println("ScreenRoms OnEnter")
	screen.list.Render()
}

func (screen *ScreenRoms) OnExit() {
	fmt.Println("ScreenRoms OnExit", screen.name)

}

func (screen *ScreenRoms) OnTick(tick TickData) {
	list := screen.list
	if list == nil {
		return
	}

	list.OnTick()
}

func (screen *ScreenRoms) Render() {
	//fmt.Println("screenCores Render")
	screen.list.Render()
}

func (screen *ScreenRoms) TearDown() {
	fmt.Println("ScreenRoms OnTearDown", screen.name)
}
