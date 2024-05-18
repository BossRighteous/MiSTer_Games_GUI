package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

type ScreenCores struct {
	guiState *GUIState
	list     *List
}

func (screen *ScreenCores) Setup(guiState *GUIState) {
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
	screen.list = NewList(screen.guiState, items, 0)
}

func (screen *ScreenCores) OnEnter() {
	fmt.Println("screenCores OnEnter")
	screen.list.Render()
}

func (screen *ScreenCores) OnExit() {

}

func (screen *ScreenCores) OnTick(tick TickData) {

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

func (screen *ScreenCores) Render() {
	//fmt.Println("screenCores Render")
	screen.list.Render()
}

func (screen *ScreenCores) TearDown() {

}
