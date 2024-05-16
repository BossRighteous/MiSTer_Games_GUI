package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

type ScreenRoms struct {
	guiState *GUIState
	list     *List
}

func (screen *ScreenRoms) Setup(guiState *GUIState) {
	screen.guiState = guiState

	var items []ListItem
	for i := 0; i < 1000; i++ {
		item := &BasicListItem{
			label: fmt.Sprintf("Item%v", i),
		}
		item.selectCallback = func() {
			fmt.Println("OnSelect item", item.Label())
		}
		items = append(items, item)
	}
	screen.list = NewList(8, screen.guiState, items)
}

func (screen *ScreenRoms) OnEnter() {
	fmt.Println("screenCores OnEnter")
	screen.list.Render()
}

func (screen *ScreenRoms) OnExit() {

}

func (screen *ScreenRoms) OnTick(tick TickData) {

	input := screen.guiState.Input
	if input.IsJustPressed(1, groovymister.InputDown) {
		screen.list.NextItem()
	} else if input.IsJustPressed(1, groovymister.InputUp) {
		screen.list.PreviousItem()
	} else if input.IsJustPressed(1, groovymister.InputRight) {
		screen.list.NextPage()
	} else if input.IsJustPressed(1, groovymister.InputLeft) {
		screen.list.PreviousPage()
	}
	//fmt.Println("screenCores OnTick", tick.FrameCount)
	//fmt.Println(screen.list.CurrentItem().Label())
	//screen.list.NextItem()
	//screen.list.CurrentItem().OnSelect()
	//screen.guiState.IsChanged = true
}

func (screen *ScreenRoms) Render() {
	//fmt.Println("screenCores Render")
	screen.list.Render()
}

func (screen *ScreenRoms) TearDown() {

}
