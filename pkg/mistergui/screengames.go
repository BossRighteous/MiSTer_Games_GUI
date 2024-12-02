package mistergui

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
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
	screen.list = NewList(screen, screen.guiState, []IListItem{}, 0)
	screen.view = ListView
}

func (screen *ScreenGames) OnEnter() {
	fmt.Println("screenCores OnEnter")
	client, err := mgdb.OpenMGDB("/mnt/c/Users/bossr/Code/MiSTer_Games_GUI/games/N64/_N64.mgdb")
	if err != nil {
		fmt.Println(err)
	}
	screen.client = client

	go func() {
		mgdbList, _ := client.GetGameList()
		//fmt.Printf("%+v\n", list)

		var items []IListItem
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
	fmt.Println("Async called")
}

func (screen *ScreenGames) OnExit() {
	screen.list.ReplaceItems([]IListItem{})
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

func (screen *ScreenGames) onTickListView() {
	list := screen.list
	input := screen.guiState.Input
	if list.ItemCount() > 0 {
		itemChanged := false
		if input.IsJustPressed(1, groovymister.InputDown) {
			list.NextItem()
			itemChanged = true
		} else if input.IsJustPressed(1, groovymister.InputUp) {
			list.PreviousItem()
			itemChanged = true
		} else if input.IsJustPressed(1, groovymister.InputRight) {
			list.NextPage()
			itemChanged = true
		} else if input.IsJustPressed(1, groovymister.InputLeft) {
			list.PreviousPage()
			itemChanged = true
		} else {
			list.CurrentItem().OnTick()
		}
		/*else if input.IsJustPressed(1, groovymister.InputB1) {
			list.CurrentItem().OnSelect()
		} else if input.IsJustPressed(1, groovymister.InputB3) {
			item := list.CurrentItem()
			gameItem, ok := item.(*GameListItem)
			if ok {
				screen.loadAsyncGameAssets(gameItem.Game.GameID)
			}
			fmt.Println("changing view to ScreenshotView")
			screen.view = ScreenshotView
			screen.guiState.IsChanged = true
		}
		*/
		if itemChanged {
			screen.guiState.IsChanged = true
		}
	}
	/*
		if input.IsJustPressed(1, groovymister.InputB2) {
			fmt.Println("back button pressed, return to cores")
		}
	*/
}

func (screen *ScreenGames) loadAsyncGameAssets(gameID int) {
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
	input := screen.guiState.Input
	if input.IsJustPressed(1, groovymister.InputDown) ||
		input.IsJustPressed(1, groovymister.InputUp) ||
		input.IsJustPressed(1, groovymister.InputRight) ||
		input.IsJustPressed(1, groovymister.InputLeft) ||
		input.IsJustPressed(1, groovymister.InputB2) {
		screen.view = ListView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to ListView")
	} else if input.IsJustPressed(1, groovymister.InputB3) {
		screen.view = TitleScreenView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to TitleScreenView")
	}
}

func (screen *ScreenGames) onTickTitleScreenView() {
	input := screen.guiState.Input
	if input.IsJustPressed(1, groovymister.InputDown) ||
		input.IsJustPressed(1, groovymister.InputUp) ||
		input.IsJustPressed(1, groovymister.InputRight) ||
		input.IsJustPressed(1, groovymister.InputLeft) ||
		input.IsJustPressed(1, groovymister.InputB2) {
		screen.view = ListView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to ListView")
	} else if input.IsJustPressed(1, groovymister.InputB3) {
		screen.view = InfoView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to InfoView")
	}
}

func (screen *ScreenGames) onTickInfoView() {
	input := screen.guiState.Input
	if input.IsJustPressed(1, groovymister.InputDown) ||
		input.IsJustPressed(1, groovymister.InputUp) ||
		input.IsJustPressed(1, groovymister.InputRight) ||
		input.IsJustPressed(1, groovymister.InputLeft) ||
		input.IsJustPressed(1, groovymister.InputB2) {
		screen.view = ListView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to ListView")
	} else if input.IsJustPressed(1, groovymister.InputB3) {
		screen.view = DescriptionView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to DescriptionView")
	}
}

func (screen *ScreenGames) onTickDescriptionView() {
	input := screen.guiState.Input
	if input.IsJustPressed(1, groovymister.InputDown) ||
		input.IsJustPressed(1, groovymister.InputUp) ||
		input.IsJustPressed(1, groovymister.InputRight) ||
		input.IsJustPressed(1, groovymister.InputLeft) ||
		input.IsJustPressed(1, groovymister.InputB2) {
		screen.view = ListView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to ListView")
	} else if input.IsJustPressed(1, groovymister.InputB3) {
		screen.view = ScreenshotView
		screen.guiState.IsChanged = true
		fmt.Println("changing view to ScreenshotView")
	}
}

func (screen *ScreenGames) Render() {
	fmt.Println("rendering Screen")
	if screen.view == ListView {
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
	draw.Draw(surface.Image, surface.Image.Rect, *screen.screenshot, P0, draw.Over)
}

func (screen *ScreenGames) renderTitleScreenView() {
	fmt.Println("rendering ScreenshotView")
	if screen.titleScreen == nil {
		fmt.Println("rendering titleScreen NIL")
		return
	}
	surface := screen.guiState.Surface
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
