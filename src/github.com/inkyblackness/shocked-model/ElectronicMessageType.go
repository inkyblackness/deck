package model

// ElectronicMessageType classifies the messages.
type ElectronicMessageType string

const (
	// ElectronicMessageTypeMail is for wireless mails.
	ElectronicMessageTypeMail ElectronicMessageType = "mail"
	// ElectronicMessageTypeLog is for collected logs.
	ElectronicMessageTypeLog = "log"
	// ElectronicMessageTypeFragment is for downloaded fragments.
	ElectronicMessageTypeFragment = "fragment"
)

// ElectronicMessageTypes returns all known message types.
func ElectronicMessageTypes() []ElectronicMessageType {
	return []ElectronicMessageType{ElectronicMessageTypeMail, ElectronicMessageTypeLog, ElectronicMessageTypeFragment}
}
