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
	Games *ScreenGames
	Roms  *ScreenRoms
}

/*
Screen Hierarchy
ScreenCores
	fetch defined cores (JSON)
	check for meta pack/index
		if no pack, prefix question mark
			A button does nothing
		has pack
			A button sets GUIState Core goes to ScreenGames
	B button sets GUIState Core, goes to ScreenCoreSettings
ScreenCoreSettings
	if no GUIState core, go to ScreenCores
	Simple CommandList
		Fetch Meta (Game Relative Path)
		Index Local Roms
		Clean Meta (Remove unused)
		Remove All ScreenShots
		Remove All TitleScreens
	C Button goes to ScreenCores
ScreenGames
	if no pack
		Notice missing meta, go to ScreenCoreSettings
	if no index
		Notice missing index, go to ScreenCoreSettings
	GamesList
		label: game title from index
		stored value: DB Game id from index
	Left Right act as page changes
	onChange (input event modifying list)
		set game/title/screen to null-eq
		async load DB Game, Title, Screenshot
	A button, set game ID, go to ScreenGameRoms
	B button cycles view in render branch
	C Button goes to ScreenCores
ScreenGameRoms
	if no gameID, go to ScreenGames
	GameRomsList
		label: rom filename no ext from index
		stored value: DB Game id from index
	A Button loads MGL, exits
	C Button goes to ScreenCores

*/
