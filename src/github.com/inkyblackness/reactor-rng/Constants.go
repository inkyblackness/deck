package main

const (
	// Version contains the current version number
	Version = "1.0.0"
	// Name is the name of the application
	Name = "InkyBlackness Reactor Randomizer"
	// Title contains a combined string of name and version
	Title = Name + " v." + Version
)

const (
	// GameStateResourceID identifies the game state chunk.
	GameStateResourceID = 0x0FA1
	// HealthOffset is the offset to the health state in the game state.
	HealthOffset = 0x009C
	// ReactorCodeOffset is the offset to the two reactor code variables.
	ReactorCodeOffset = 0x00F6 + (31 * 2)
)

// Constants for level archives of Citadel.
const (
	SceneryListResourceIDOffset           = 7
	PanelsListResourceIDOffset            = 9
	ResourcesPerLevel                     = 100
	ResourceIDLevelBase                   = 4000
	LevelObjectEntryListsResourceIDOffset = 10

	ReactorLevelPanelClassDataResourceID = ResourceIDLevelBase + LevelObjectEntryListsResourceIDOffset + PanelsListResourceIDOffset
	PanelClassDataEntryCount             = 64
	SceneryClassDataEntryCount           = 176

	ReactorCodePanelMasterIndex = 56
	ReactorCodePanelClassIndex  = 9
	ReactorCodePanelTrigger1    = 58
	ReactorCodePanelTrigger2    = 136
	ReactorCodePanelFail        = 55

	RandomDigitPictureSource    = 0x017F
	CodeDigitPictureSourceStart = 0x0134
	CodeDigitPictureSourceEnd   = CodeDigitPictureSourceStart + 10
)

type levelScreenInfo struct {
	ScreenMasterIndex uint16
	ScreenClassIndex  int
}

var levelCodeScreens = [6][]levelScreenInfo{
	{{167, 45}, {238, 93}},
	{{283, 82}},
	{{224, 8}},
	{{64, 41}, {67, 42}},
	{{16, 23}},
	{{382, 127}},
}
