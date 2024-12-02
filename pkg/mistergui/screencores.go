package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mrext"
)

type ScreenCores struct {
	name     string
	parent   IScreen
	guiState *GUIState
	list     *List
}

func (screen *ScreenCores) GUIState() *GUIState {
	return screen.guiState
}

func (screen *ScreenCores) Parent() IScreen {
	return screen.parent
}

func (screen *ScreenCores) Name() string {
	return screen.name
}

func (screen *ScreenCores) Setup() {
	var items []IListItem

	fmt.Println("screenCores Setup", screen.name)

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

	for i := 0; i < 100; i++ {
		item := &BasicListItem{
			label:        fmt.Sprintf("Item %v", i),
			list:         screen.list,
			buttonsLabel: fmt.Sprintf("A: Accept, Item %v", i),
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				fmt.Println("OnSelect item", item.Label())
				newScreen := &ScreenCores{parent: screen, guiState: screen.guiState, name: item.Label()}
				newScreen.Setup()
				screen.guiState.PushScreen(newScreen)
				screen.guiState.IsChanged = true
			}

		}
		items = append(items, item)
	}
	screen.list.ReplaceItems(items)
}

func (screen *ScreenCores) OnEnter() {
	fmt.Println("screenCores OnEnter", screen.name)
	screen.list.Render()
}

func (screen *ScreenCores) OnExit() {
	fmt.Println("screenCores OnExit", screen.name)

}

func (screen *ScreenCores) OnTick(tick TickData) {
	list := screen.list
	if list == nil {
		return
	}

	list.OnTick()

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
	fmt.Println("screenCores OnTearDown", screen.name)
}
