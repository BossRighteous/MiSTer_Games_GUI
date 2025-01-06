package mistergui

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mrext"
)

type MetaView int

const (
	ListView MetaView = iota
	ScreenshotView
	TitleScreenView
	InfoView
	DescriptionView
)

var LoadingImage image.Image = image.NewNRGBA(image.Rect(0, 0, 0, 0))

type ScreenGames struct {
	name           string
	parent         IScreen
	guiState       *GUIState
	list           *List
	client         *mgdb.MGDBClient
	screenshot     *image.Image
	titleScreen    *image.Image
	infoImg        *image.Image
	descriptionImg *image.Image
	view           MetaView
	didScan        bool
}

func (screen *ScreenGames) GUIState() *GUIState {
	return screen.guiState
}

func (screen *ScreenGames) Parent() IScreen {
	return screen.parent
}

func (screen *ScreenGames) Name() string {
	return screen.name
}

func (screen *ScreenGames) Setup() {
	screen.didScan = false
	screen.list = NewList(screen, screen.guiState, []IListItem{}, 0)
	screen.view = ListView

	var items []IListItem

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

	if screen.client == nil {
		item := &BasicListItem{
			label:        "Unable to load Collection data",
			list:         screen.list,
			buttonsLabel: "",
		}
		items = append(items, item)
		screen.list.ReplaceItems(items)
		return
	}

	{
		// Scan for ROMs
		item := &BasicListItem{
			label:        fmt.Sprintf("Scan for ROMs"),
			list:         screen.list,
			buttonsLabel: "A: Accept",
		}
		item.tickCallback = func() {
			if screen.guiState.Input.IsJustPressed(1, groovymister.InputB1) {
				//
				screen.didScan = true
				outBuffer, completedOk := ScanMGDBGames(screen.client)
				consoleState := ConsoleState{
					outBuffer:   outBuffer,
					completedOk: completedOk,
				}
				consoleScreen := &ScreenConsole{
					parent:       screen,
					guiState:     screen.guiState,
					name:         item.Label(),
					consoleState: consoleState,
				}
				consoleScreen.Setup()
				screen.guiState.PushScreen(consoleScreen)
				screen.guiState.IsChanged = true
			}
		}
		items = append(items, item)
	}

	go func() {
		client := screen.client
		mgdbList, _ := client.GetGameList()

		for _, mgdbGameItem := range mgdbList {
			// Make GameListItem with GameID for additonal use
			item := &GameListItem{
				Game:   mgdbGameItem,
				screen: screen,
				list:   screen.list,
			}
			items = append(items, item)
		}
		fmt.Println("AsyncChan sending")
		screen.guiState.AsyncChan <- func(gui *GUI) {
			screen.list.ReplaceItems(items)
			fmt.Println("AsyncChan callback executed")
		}
	}()

	screen.list.ReplaceItems(items)
}

func (screen *ScreenGames) OnEnter() {
	fmt.Println("screenCores OnEnter")
	if screen.didScan {
		screen.Setup()
	}
	screen.list.Render()
}

func (screen *ScreenGames) OnExit() {
	//screen.list.ReplaceItems([]IListItem{})
	screen.view = ListView
	screen.screenshot = nil
	screen.titleScreen = nil
	screen.infoImg = nil
	screen.descriptionImg = nil
}

func (screen *ScreenGames) OnTick(tick TickData) {
	list := screen.list
	if list == nil {
		// Don't do anything until list is ready
		return
	}
	screen.list.OnTick()

	if screen.view == ListView {
		screen.onTickListView()
	} else if screen.view == ScreenshotView {
		screen.onTickScreenshotView()
	} else if screen.view == TitleScreenView {
		screen.onTickTitleScreenView()
	} else if screen.view == InfoView {
		screen.onTickInfoView()
	} else if screen.view == DescriptionView {
		screen.onTickDescriptionView()
	}
}

func (screen *ScreenGames) CycleMediaView() {
	screen.view = screen.view + 1
	if screen.view > 4 {
		screen.view = 0
	}
	screen.guiState.IsChanged = true
}

func (screen *ScreenGames) ResetGameAssets() {
	screen.screenshot = nil
	screen.titleScreen = nil
	screen.infoImg = nil
	screen.descriptionImg = nil
}

func (screen *ScreenGames) LoadAsyncGameAssets(gameID int) {
	fmt.Println("LoadAsyncGameAssets", gameID)
	// load all go routines in parallel
	if screen.screenshot != &LoadingImage {
		screen.screenshot = &LoadingImage
		go func() {
			screenshot, err := screen.client.GetGameScreenshot(gameID)
			if err != nil {
				fmt.Println("unable to load screenshot ", gameID)
			}
			screen.guiState.AsyncChan <- func(gui *GUI) {
				screen.screenshot = screenshot
				gui.State.IsChanged = true
				fmt.Println("setting screenshot")
			}
		}()
	}

	if screen.titleScreen != &LoadingImage {
		screen.titleScreen = &LoadingImage
		go func() {
			titleScreen, err := screen.client.GetGameTitleScreen(gameID)
			if err != nil {
				fmt.Println("unable to load titleScreen ", gameID)
			}
			screen.guiState.AsyncChan <- func(gui *GUI) {
				screen.titleScreen = titleScreen
				gui.State.IsChanged = true
				fmt.Println("setting titleScreen")
			}
		}()
	}

	if screen.descriptionImg != &LoadingImage {
		screen.infoImg = &LoadingImage
		screen.descriptionImg = &LoadingImage
		surfaceRect := screen.guiState.Surface.Image.Rect
		go func() {
			game, err := screen.client.GetGame(gameID)
			if err != nil {
				fmt.Println("unable to load Game", gameID)
			}
			var infoImg image.Image
			var descriptionImg image.Image

			infoText := []string{
				fmt.Sprintf("Name: %s", game.Name),
				fmt.Sprintf("Developer: %s", game.Developer),
				fmt.Sprintf("Publisher: %s", game.Publisher),
				fmt.Sprintf("Genre: %s", game.Genre),
				fmt.Sprintf("Rating: %s", game.Rating),
				fmt.Sprintf("Release Date: %s", game.ReleaseDate),
				fmt.Sprintf("Players: %s", game.Players),
			}
			infoImg = DrawText(infoText, surfaceRect, image.Transparent)

			charsPerLine := 45
			descriptionLines := make([]string, 0)
			descriptionLines = append(descriptionLines, "Description:")
			if game.Description != "" {
				desc := game.Description
				offset := 0
				end := len(game.Description) - 1
				for offset < end {
					lookahead := charsPerLine
					if lookahead+offset >= end {
						lookahead = end - offset
					}
					fmt.Println("Game description looping", offset, lookahead, end)
					subslice := desc[offset : offset+lookahead]
					if lookahead == charsPerLine {
						for lookahead > 0 {
							if string(subslice[lookahead-1]) == " " {
								break
							}
							lookahead--
						}
					}
					fmt.Println("Game description looping", offset, lookahead, end)
					descriptionLines = append(descriptionLines, desc[offset:offset+lookahead])
					offset += lookahead
				}
			}
			descriptionImg = DrawText(descriptionLines, surfaceRect, image.Transparent)

			screen.guiState.AsyncChan <- func(gui *GUI) {
				screen.infoImg = &infoImg
				screen.descriptionImg = &descriptionImg
				gui.State.IsChanged = true
				fmt.Println("setting titleScreen")
			}
		}()
	}

}

func (screen *ScreenGames) onTickScreenshotView() {
}

func (screen *ScreenGames) onTickListView() {
}

func (screen *ScreenGames) onTickTitleScreenView() {
}

func (screen *ScreenGames) onTickInfoView() {
}

func (screen *ScreenGames) onTickDescriptionView() {
}

func (screen *ScreenGames) Render() {
	fmt.Println("rendering Screen")
	_, isCurrentItemGame := screen.list.CurrentItem().(*GameListItem)
	if screen.view == ListView {
		screen.renderListView()
	} else if !isCurrentItemGame {
		// Swap back to list view for non game items
		screen.view = ListView
		screen.renderListView()
	} else if screen.view == ScreenshotView {
		screen.renderScreenshotView()
	} else if screen.view == TitleScreenView {
		screen.renderTitleScreenView()
	} else if screen.view == InfoView {
		screen.renderInfoView()
	} else if screen.view == DescriptionView {
		screen.renderDesciptionView()
	}
}

func (screen *ScreenGames) renderListView() {
	screen.list.Render()
	fmt.Println("rendering ListView")
}

func (screen *ScreenGames) renderScreenshotView() {
	fmt.Println("rendering ScreenshotView")
	if screen.screenshot == nil {
		fmt.Println("rendering Screenshot NIL")
		return
	}
	surface := screen.guiState.Surface
	surface.Erase(surface.Image.Rect, P0)
	draw.Draw(surface.Image, surface.Image.Rect, *screen.screenshot, P0, draw.Over)
}

func (screen *ScreenGames) renderTitleScreenView() {
	fmt.Println("rendering ScreenshotView")
	if screen.titleScreen == nil {
		fmt.Println("rendering titleScreen NIL")
		return
	}
	surface := screen.guiState.Surface
	surface.Erase(surface.Image.Rect, P0)
	draw.Draw(surface.Image, surface.Image.Rect, *screen.titleScreen, P0, draw.Over)

}

func (screen *ScreenGames) renderInfoView() {
	fmt.Println("rendering InfoView")
	if screen.infoImg == nil {
		fmt.Println("rendering infoImg NIL")
		return
	}
	surface := screen.guiState.Surface
	rect := surface.Image.Rect
	surface.Erase(rect, P0)
	draw.Draw(surface.Image, rect, *screen.infoImg, P0, draw.Over)
}

func (screen *ScreenGames) renderDesciptionView() {
	fmt.Println("rendering DesciptionView")
	if screen.descriptionImg == nil {
		fmt.Println("rendering descriptionImg NIL")
		return
	}
	surface := screen.guiState.Surface
	rect := surface.Image.Rect
	surface.Erase(rect, P0)
	draw.Draw(surface.Image, rect, *screen.descriptionImg, P0, draw.Over)
}

func (screen *ScreenGames) TearDown() {

}

// Utility method for collection scanning
func ScanMGDBGames(client *mgdb.MGDBClient) (chan string, chan bool) {
	outBuffer := make(chan string, 1024)
	completedOk := make(chan bool)

	go func() {
		bPrint := func(msg string) {
			fmt.Println(msg)
			outBuffer <- msg
		}

		info, err := client.GetMGDBInfo()
		if err != nil {
			bPrint(err.Error())
			completedOk <- false
			return
		}

		// TODO: reset IsIndexed on Games table to 0
		if err := client.FlushGamesIndex(); err != nil {
			fmt.Println(err)
			bPrint("Could not flush Games index")
			completedOk <- false
			return
		}

		systems := mrext.GetSystemsByIDsString(info.SupportedSystemIds)
		roms, err := mrext.GetSystemsGamesPaths(systems)
		if err != nil {
			bPrint(err.Error())
			completedOk <- false
			return
		}

		bPrint(fmt.Sprintf("Found %v ROMs", len(roms)))
		indexedCount := 0
		for _, romAbsPath := range roms {
			fmt.Println(romAbsPath)
			romRelPath, ok := mrext.GetRelativeGamePath(romAbsPath)
			if !ok {
				bPrint("Relative Pathing Error on Scanned ROM")
				bPrint(romRelPath)
				continue
			}
			bPrint(fmt.Sprintf("Found %v", romRelPath))

			// Match to DB
			// Update Game row IsIndexed
			// Add IndexedRom record
			indexedRom := mgdb.MakeIndexedRomFromPath(romRelPath, 0)

			gameId, findErr := client.FindGameIdFromFilename(indexedRom.FileName)
			if findErr != nil {
				bPrint("Error attempting filename match in Collection")
				bPrint(romRelPath)
				continue
			}
			indexedRom.GameID = gameId

			// Even 0 IDs should be indexed as unknown
			client.IndexGameRom(indexedRom)

			indexedCount++
			bPrint(fmt.Sprintf("Indexed %v", romRelPath))

		}
		bPrint(fmt.Sprintf("Indexed %v ROMs", len(roms)))

		completedOk <- true
	}()

	return outBuffer, completedOk
}
