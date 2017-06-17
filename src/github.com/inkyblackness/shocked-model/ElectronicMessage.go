package model

// ElectronicMessage describes the base properties of an electronic message.
type ElectronicMessage struct {
	// NextMessage describes the message that will interrupt this one. -1 for no interrupt.
	NextMessage *int
	// IsInterrupt is set for interrupting messages.
	IsInterrupt *bool
	// ColorIndex for special colored headers. -1 for default.
	ColorIndex *int
	// LeftDisplay identifies the image for the left display. -1 for none.
	LeftDisplay *int
	// RightDisplay identifies the image for the right display. -1 for none.
	RightDisplay *int

	Title       [LanguageCount]*string
	Sender      [LanguageCount]*string
	Subject     [LanguageCount]*string
	VerboseText [LanguageCount]*string
	TerseText   [LanguageCount]*string
}
