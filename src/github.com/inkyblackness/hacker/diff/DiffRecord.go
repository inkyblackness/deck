package diff

import "fmt"

type DiffRecord struct {
	Payload byte
	Delta   DeltaType
}

// String returns a string representation of a record. The string is a
// concatenation of the delta type and the payload.
func (record DiffRecord) String() string {
	return fmt.Sprintf("%s %v", record.Delta, record.Payload)
}
