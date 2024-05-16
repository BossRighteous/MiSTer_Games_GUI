package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
)

type ScreenGames struct {
	guiState *GUIState
	list     *List
}

func (screen *ScreenGames) Setup(guiState *GUIState) {
	screen.guiState = guiState

	screen.list = NewList(8, screen.guiState, []ListItem{})
}

func (screen *ScreenGames) OnEnter() {
	fmt.Println("screenCores OnEnter")

	client, _ := mgdb.OpenMGDB("/mnt/c/Users/bossr/Code/MiSTer_Games_GUI/games/N64/_N64.mgdb")
	fmt.Printf("%+v\n", client)
	info, _ := client.GetMGDBInfo()
	fmt.Printf("%+v\n", info)

	list, _ := client.GetGameList()
	fmt.Printf("%+v\n", list)

	var items []ListItem
	for _, gameItem := range list {
		item := &BasicListItem{
			label: gameItem.Name,
		}
		item.selectCallback = func() {
			fmt.Println("OnSelect item", item.Label())
		}
		items = append(items, item)
	}
	screen.list.ReplaceItems(items)

	screen.list.Render()
}

func (screen *ScreenGames) OnExit() {

}

func (screen *ScreenGames) OnTick(tick TickData) {

	input := screen.guiState.Input
	if input.IsJustPressed(1, groovymister.InputDown) {
		screen.list.NextItem()
	} else if input.IsJustPressed(1, groovymister.InputUp) {
		screen.list.PreviousItem()
	} else if input.IsJustPressed(1, groovymister.InputRight) {
		screen.list.NextPage()
	} else if input.IsJustPressed(1, groovymister.InputLeft) {
		screen.list.PreviousPage()
	} else if input.IsJustPressed(1, groovymister.InputB1) {
		screen.list.CurrentItem().OnSelect()
	}
	//fmt.Println("screenCores OnTick", tick.FrameCount)
	//fmt.Println(screen.list.CurrentItem().Label())
	//screen.list.NextItem()
	//screen.list.CurrentItem().OnSelect()
	//screen.guiState.IsChanged = true
}

func (screen *ScreenGames) Render() {
	//fmt.Println("screenCores Render")
	screen.list.Render()
}

func (screen *ScreenGames) TearDown() {

}
