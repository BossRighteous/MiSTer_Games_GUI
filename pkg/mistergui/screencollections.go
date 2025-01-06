package mistergui

import (
	"fmt"
	"path/filepath"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/utils"
)

/*
MDGB Files as entrypoint root screen. May also link to settings
*/

type ScreenCollections struct {
	name     string
	parent   IScreen
	guiState *GUIState
	list     *List
}

func (screen *ScreenCollections) GUIState() *GUIState {
	return screen.guiState
}

func (screen *ScreenCollections) Parent() IScreen {
	return screen.parent
}

func (screen *ScreenCollections) Name() string {
	return screen.name
}

func (screen *ScreenCollections) Setup() {
	screen.name = "Collections"

	var items []IListItem

	fmt.Println("screenCores Setup", screen.name)

	screen.list = NewList(screen, screen.guiState, make([]IListItem, 0), 0)

	{
		item := &BasicListItem{
			label:        "Quit GUI",
			list:         screen.list,
			buttonsLabel: "A: Accept",
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				fmt.Println("Quitting GUI")
				screen.guiState.PopScreen()
				screen.guiState.QuitChan <- true
			}
		}
		items = append(items, item)
	}

	{
		// Fetch list of MGDB Files from main dir
		collectionsPath := filepath.Clean(screen.guiState.Settings.CollectionsPath)
		mgdbs, _ := filepath.Glob(filepath.Join(collectionsPath, "*.mgdb"))
		fmt.Println(collectionsPath, mgdbs)

		if len(mgdbs) == 0 {
			item := &BasicListItem{
				label:        "No Collections found in path",
				list:         screen.list,
				buttonsLabel: fmt.Sprintf("Path: %v", collectionsPath),
			}
			items = append(items, item)
		}

		for _, mgdbPath := range mgdbs {
			name, found := utils.CutSuffix(filepath.Base(mgdbPath), filepath.Ext(mgdbPath))
			if !found {
				name = filepath.Base(mgdbPath)
			}

			item := &BasicListItem{
				label:        name,
				list:         screen.list,
				buttonsLabel: fmt.Sprintf("A: View %v Games", name),
			}
			item.tickCallback = func() {
				if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
					client, err := mgdb.OpenMGDB(mgdbPath)
					if err != nil {
						fmt.Println(fmt.Sprintf("Error loading MGDB: %v", mgdbPath))
						fmt.Println(err)
						item.label = fmt.Sprintf("!%v", item.label)
						item.buttonsLabel = fmt.Sprintf("Error loading MGDB: %v", mgdbPath)
						screen.guiState.IsChanged = true
						return
					}

					gamesScreen := &ScreenGames{
						parent: screen, guiState: screen.guiState, name: item.Label(), client: client,
					}
					gamesScreen.Setup()
					screen.guiState.PushScreen(gamesScreen)
					screen.guiState.IsChanged = true
				}
			}
			items = append(items, item)
		}
		screen.list.ReplaceItems(items)
	}

	screen.list.ReplaceItems(items)
}

func (screen *ScreenCollections) OnEnter() {
	fmt.Println("screenCores OnEnter", screen.name)
	screen.list.Render()
}

func (screen *ScreenCollections) OnExit() {
	fmt.Println("screenCores OnExit", screen.name)
}

func (screen *ScreenCollections) OnTick(tick TickData) {
	list := screen.list
	if list == nil {
		return
	}

	list.OnTick()
}

func (screen *ScreenCollections) Render() {
	//fmt.Println("screenCores Render")
	screen.list.Render()
}

func (screen *ScreenCollections) TearDown() {
	fmt.Println("screenCores OnTearDown", screen.name)
}
