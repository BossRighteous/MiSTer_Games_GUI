package mistergui

import (
	"image"
	"image/draw"
	"math"
)

type List struct {
	guiState       *GUIState
	items          []ListItem
	itemIndex      int
	perPage        int
	renderCallback func(*List)
}

func (list *List) ReplaceItems(items []ListItem) {
	list.itemIndex = 0
	list.items = items
}

func (list *List) CurrentItem() ListItem {
	return list.items[list.itemIndex]
}

func (list *List) IsCurrentItem(item ListItem) bool {
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

func (list *List) PageItems() []ListItem {
	return list.items[list.PageInitialIndex() : list.PageFinalIndex()+1]
}

func NewList(perPage int, guiState *GUIState, items []ListItem) *List {
	list := &List{
		perPage:        perPage,
		items:          items,
		guiState:       guiState,
		renderCallback: DefaultListRender,
	}
	return list
}

var ListPadding = 16
var DefaultListRect = image.Rect(0, 0, 320-(ListPadding*2), 240-(ListPadding*2))

var DefaultListRender = func(list *List) {
	items := list.PageItems()
	textStrings := make([]string, list.perPage)
	for i, item := range items {
		textStrings[i] = item.Label()
		if list.IsCurrentItem(item) {
			textStrings[i] = "> " + item.Label()
		}
	}
	surface := list.guiState.Surface
	surface.Erase(DefaultListRect, P0)
	img := DrawText(textStrings, DefaultListRect, image.Transparent)
	draw.Draw(surface.Image, DefaultListRect, img, image.Point{0, 0}, draw.Over)
}

type ListItem interface {
	Label() string
	OnSelect()
	OnEnter()
	OnExit()
	//OnCancel()
	//TearDown()
}

type BasicListItem struct {
	labelPrefix    string
	label          string
	selectCallback func()
}

func (item *BasicListItem) Label() string {
	return item.labelPrefix + item.label
}

func (item *BasicListItem) OnSelect() {
	item.selectCallback()
}

func (item *BasicListItem) OnEnter() {
}

func (item *BasicListItem) OnExit() {
}
