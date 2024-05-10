package mistergui

import (
	"fmt"
	"math"
)

type List struct {
	guiState    *GUIState
	items       []ListItem
	itemIndex   uint
	perPage     uint
	currentPage uint
}

func (list *List) CurrentItem() ListItem {
	return list.items[list.itemIndex]
}

func (list *List) NextItem() {
	if list.itemIndex < uint(len(list.items))-1 {
		list.guiState.IsChanged = true
		list.itemIndex++
	}
}

func (list *List) PreviousItem() {
	if list.itemIndex > 0 {
		list.guiState.IsChanged = true
		list.itemIndex--
	}
}

func (list *List) NextPage() {
	if list.itemIndex+list.perPage < uint(len(list.items))-1 {
		list.guiState.IsChanged = true
		list.itemIndex += list.perPage
	}
}

func (list *List) PreviousPage() {
	if list.itemIndex-list.perPage > 0 {
		list.guiState.IsChanged = true
		list.itemIndex -= list.perPage
	}
}

func (list *List) Render() {
	fmt.Println(list.CurrentItem().Label())
	// Maybe only need to render page changes based on underlay
	pageCalc := uint(math.Floor(float64(list.itemIndex) / float64(list.perPage)))
	if list.currentPage == pageCalc {
		return
	}
	fmt.Println("New Page", pageCalc)
	list.currentPage = pageCalc
}

func (list *List) ItemCount() uint {
	return uint(len(list.items))
}

func (list *List) PageCount() uint {
	return uint(math.Floor(float64(len(list.items)) / float64(list.perPage)))
}

func (list *List) CurrentPage() uint {
	return list.currentPage
}

func NewList(perPage uint, guiState *GUIState, items []ListItem) *List {
	list := &List{
		perPage:  perPage,
		items:    items,
		guiState: guiState,
	}
	return list
}

type ListItem interface {
	Label() string
	//Setup(GUIState)
	//OnEnter()
	//OnExit()
	//OnSelect()
	//OnCancel()
	//TearDown()
}

type BasicListItem struct {
	label string
	//value uint
}

func (item *BasicListItem) Label() string {
	return item.label
}
