package mistergui

import (
	"fmt"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mrext"
)

type ScreenRomSystems struct {
	name     string
	parent   IScreen
	guiState *GUIState
	list     *List
	client   *mgdb.MGDBClient
	rom      mgdb.IndexedRom
}

func (screen *ScreenRomSystems) GUIState() *GUIState {
	return screen.guiState
}

func (screen *ScreenRomSystems) Parent() IScreen {
	return screen.parent
}

func (screen *ScreenRomSystems) Name() string {
	return screen.name
}

func (screen *ScreenRomSystems) Setup() {
	screen.name = "ROM Systems"

	var items []IListItem

	fmt.Println("ScreenRomSystems Setup", screen.name)

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

	/*{
		item := &BasicListItem{
			label:        fmt.Sprintf("%v%v", screen.rom.FileName, screen.rom.FileExt),
			list:         screen.list,
			buttonsLabel: "ROM to load below",
		}
		items = append(items, item)
	}*/

	client := screen.client

	// Move to new screen
	info, _ := client.GetMGDBInfo()

	systems := mrext.GetSystemsByIDsString(info.SupportedSystemIds)

	for _, iSystem := range systems {
		system := iSystem
		item := &BasicListItem{
			label:        fmt.Sprintf("Load as %v", system.Name),
			list:         screen.list,
			buttonsLabel: "A: Load ROM via MGL Command and Exit GUI",
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				fmt.Println("OnSelect item", item.Label())
				absPath, found := mrext.GetFirstGamePathFromRelative(screen.rom.Path)
				if !found {
					item.label = fmt.Sprintf("!%v", item.label)
					item.buttonsLabel = "ROM not found as indexed, cannot load."
					return
				}
				mrext.LoadSystemMGLFromPath(system, absPath)
				screen.guiState.QuitChan <- true
			}
		}
		items = append(items, item)
	}

	screen.list.ReplaceItems(items)

}

func (screen *ScreenRomSystems) OnEnter() {
	fmt.Println("ScreenRomSystems OnEnter")
	screen.list.Render()
}

func (screen *ScreenRomSystems) OnExit() {
	fmt.Println("ScreenRomSystems OnExit", screen.name)

}

func (screen *ScreenRomSystems) OnTick(tick TickData) {
	list := screen.list
	if list == nil {
		return
	}

	list.OnTick()
}

func (screen *ScreenRomSystems) Render() {
	//fmt.Println("screenCores Render")
	screen.list.Render()
}

func (screen *ScreenRomSystems) TearDown() {
	fmt.Println("ScreenRomSystems OnTearDown", screen.name)
}
