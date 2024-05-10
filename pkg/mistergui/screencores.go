package mistergui

import (
	"fmt"
	//"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mister"
)

type ScreenCores struct {
	guiState *GUIState
	list     *List
}

func (screen *ScreenCores) Setup(guiState *GUIState) {
	screen.guiState = guiState

	var items []ListItem
	for i := 0; i < 1000; i++ {
		items = append(items, &BasicListItem{label: fmt.Sprintf("Item%v", i)})
	}
	screen.list = NewList(8, screen.guiState, items)
}

func (screen *ScreenCores) OnEnter() {
	fmt.Println("screenCores OnEnter")

}

func (screen *ScreenCores) OnExit() {

}

func (screen *ScreenCores) OnTick(tick TickData) {
	fmt.Println("screenCores OnTick", tick.FrameCount)
	screen.list.NextPage()
	//screen.guiState.IsChanged = true
}

func (screen *ScreenCores) Render() {
	fmt.Println("screenCores Render")
	screen.list.Render()
}

func (screen *ScreenCores) TearDown() {

}
