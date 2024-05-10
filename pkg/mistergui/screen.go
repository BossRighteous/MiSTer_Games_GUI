package mistergui

type Screen interface {
	Setup(*GUIState)
	OnEnter()
	OnExit()
	OnTick(TickData)
	Render()
	TearDown()
}

type Screens struct {
	Cores *ScreenCores
}
