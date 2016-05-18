package diff

// DeltaType describes the relationship of elements in two
// sequences. The following table provides a summary:
//
// Constant Code Meaning
// ---------- ------ ---------------------------------------
// Common " " The element occurs in both sequences.
// LeftOnly "-" The element is unique to sequence 1.
// RightOnly "+" The element is unique to sequence 2.
type DeltaType int

const (
	Common DeltaType = iota
	LeftOnly
	RightOnly
)

// String returns a string representation for DeltaType.
func (t DeltaType) String() string {
	switch t {
	case Common:
		return " "
	case LeftOnly:
		return "-"
	case RightOnly:
		return "+"
	}
	return "?"
}
