package mistergui

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/groovymister"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
)

type List struct {
	screen         IScreen
	guiState       *GUIState
	items          []IListItem
	itemIndex      int
	perPage        int
	renderCallback func(*List)
}

func (list *List) Screen() IScreen {
	return list.screen
}

func (list *List) OnTick() {

	screen := list.Screen()
	if screen == nil {
		return
	}
	guiState := screen.GUIState()
	input := screen.GUIState().Input
	if list.ItemCount() > 0 {
		if input.IsJustPressed(1, groovymister.InputDown) {
			list.NextItem()
			guiState.IsChanged = true
		} else if input.IsJustPressed(1, groovymister.InputUp) {
			list.PreviousItem()
			guiState.IsChanged = true
		} else if input.IsJustPressed(1, groovymister.InputRight) {
			list.NextPage()
			guiState.IsChanged = true
		} else if input.IsJustPressed(1, groovymister.InputLeft) {
			list.PreviousPage()
			guiState.IsChanged = true
		} else {
			list.CurrentItem().OnTick()
		}
	}
}

func (list *List) ReplaceItems(items []IListItem) {
	list.itemIndex = 0
	list.items = items
}

func (list *List) CurrentItem() IListItem {
	return list.items[list.itemIndex]
}

func (list *List) IsCurrentItem(item IListItem) bool {
	return list.CurrentItem() == item
}

func (list *List) NextItem() {
	if list.itemIndex < len(list.items)-1 {
		list.guiState.IsChanged = true
		list.CurrentItem().OnExit()
		list.itemIndex++
		list.CurrentItem().OnEnter()
	}
}

func (list *List) PreviousItem() {
	if list.itemIndex > 0 {
		list.guiState.IsChanged = true
		list.CurrentItem().OnExit()
		list.itemIndex--
		list.CurrentItem().OnEnter()
	}
}

func (list *List) NextPage() {
	itemLen := len(list.items)
	if list.itemIndex+list.perPage < itemLen-1 {
		list.guiState.IsChanged = true
		list.CurrentItem().OnExit()
		list.itemIndex += list.perPage
		list.CurrentItem().OnEnter()
	} else {
		list.guiState.IsChanged = true
		list.CurrentItem().OnExit()
		list.itemIndex = itemLen - 1
		list.CurrentItem().OnEnter()
	}
}

func (list *List) PreviousPage() {
	// uints overflow!
	if list.itemIndex-list.perPage >= 0 {
		list.guiState.IsChanged = true
		list.CurrentItem().OnExit()
		list.itemIndex -= list.perPage
		list.CurrentItem().OnEnter()
	} else {
		list.guiState.IsChanged = true
		list.CurrentItem().OnExit()
		list.itemIndex = 0
		list.CurrentItem().OnEnter()
	}

}

func (list *List) Render() {
	list.renderCallback(list)
}

func (list *List) ItemCount() int {
	return len(list.items)
}

func (list *List) PageCount() int {
	return int(math.Floor(float64(len(list.items)+1) / float64(list.perPage)))
}

func (list *List) CurrentPage() int {
	return int(math.Floor(float64(list.itemIndex) / float64(list.perPage)))
}

func (list *List) PageInitialIndex() int {
	return list.CurrentPage() * list.perPage
}

func (list *List) PageItemLocalIndex() int {
	//fmt.Println(list.itemIndex, list.PageInitialIndex(), list.currentPage, list.perPage)
	return list.itemIndex - list.PageInitialIndex()
}

func (list *List) PageFinalIndex() int {
	itemsLen := len(list.items)
	if list.PageInitialIndex()+list.perPage >= itemsLen-1 {
		return itemsLen - 1
	}
	return list.PageInitialIndex() + list.perPage - 1
}

func (list *List) PageItems() []IListItem {
	return list.items[list.PageInitialIndex() : list.PageFinalIndex()+1]
}

func NewList(screen IScreen, guiState *GUIState, items []IListItem, perPage int) *List {
	if perPage == 0 {
		perPage = 10
	}
	list := &List{
		screen:         screen,
		perPage:        perPage,
		items:          items,
		guiState:       guiState,
		renderCallback: DefaultListRender,
	}
	return list
}

var DefaultListRender = func(list *List) {
	items := list.PageItems()
	textStrings := make([]string, 0)
	for _, item := range items {
		line := item.Label()
		if list.IsCurrentItem(item) {
			line = "> " + item.Label()
		}
		textStrings = append(textStrings, line)
		//fmt.Println(textStrings[i])
	}
	surface := list.guiState.Surface
	surface.Erase(surface.Image.Rect, P0)
	img := DrawText(textStrings, surface.Image.Rect, image.Transparent)
	draw.Draw(surface.Image, surface.Image.Rect, img, P0, draw.Over)

	// Draw button label
	btnLabel := []string{list.CurrentItem().ButtonsLabel()}
	surfRect := surface.Image.Rect
	labelRect := image.Rect(0, 0, surfRect.Max.X, 30)
	btnLabelImg := DrawText(btnLabel, labelRect, image.Transparent)
	draw.Draw(surface.Image, image.Rect(0, 202, surfRect.Max.X, 232), btnLabelImg, P0, draw.Over)

}

type IListItem interface {
	List() *List
	Label() string
	OnEnter()
	OnExit()
	//OnCancel()
	//TearDown()
	// Rework this to be OnButton() instead of onSelect
	OnTick()
	ButtonsLabel() string
}

type BasicListItem struct {
	list          *List
	labelPrefix   string
	label         string
	tickCallback  func()
	enterCallback func()
	exitCallback  func()
	buttonsLabel  string
}

func (item *BasicListItem) List() *List {
	return item.list
}

func (item *BasicListItem) Label() string {
	return item.labelPrefix + item.label
}

func (item *BasicListItem) OnEnter() {
	if item.enterCallback == nil {
		return
	}
	item.enterCallback()
}

func (item *BasicListItem) OnExit() {
	if item.exitCallback == nil {
		return
	}
	item.exitCallback()
}

func (item *BasicListItem) OnTick() {
	if item.tickCallback == nil {
		return
	}
	item.tickCallback()
}

func (item *BasicListItem) ButtonsLabel() string {
	return item.buttonsLabel
}

/*
 * GameListItem
 */

type GameListItem struct {
	list   *List
	Game   mgdb.GameListItem
	screen *ScreenGames
}

func (item *GameListItem) List() *List {
	return item.list
}

func (item *GameListItem) Label() string {
	return item.Game.Name
}

func (item *GameListItem) OnTick() {
	fmt.Println("OnSelect game item", item.Label())
}

func (item *GameListItem) OnEnter() {
}

func (item *GameListItem) OnExit() {
}

func (item *GameListItem) ButtonsLabel() string {
	return ""
}
