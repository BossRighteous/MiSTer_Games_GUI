package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
)

type ScreenRoms struct {
	name     string
	parent   IScreen
	guiState *GUIState
	list     *List
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
	var items []IListItem
	for i := 0; i < 1000; i++ {
		item := &BasicListItem{
			label: fmt.Sprintf("Item%v", i),
			list:  screen.list,
		}
		items = append(items, item)
	}
	screen.list = NewList(screen, screen.guiState, items, 0)
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
