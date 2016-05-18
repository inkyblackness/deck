package util

// ValueUnifier collects values from several sources. It can be queried whether
// all sources contained the same value - or return a default value.
type ValueUnifier struct {
	defaultValue interface{}
	values       map[interface{}]int
}

// NewValueUnifier returns a new instance of a ValueUnifier for given default value.
func NewValueUnifier(defaultValue interface{}) *ValueUnifier {
	return &ValueUnifier{
		defaultValue: defaultValue,
		values:       make(map[interface{}]int)}
}

// Value returns the unified value - or the default value, if either no value was
// added or the registered values are not equal.
func (unifier *ValueUnifier) Value() (result interface{}) {
	result = unifier.defaultValue
	if len(unifier.values) == 1 {
		for key := range unifier.values {
			result = key
		}
	}

	return
}

// Add registers an additional value.
func (unifier *ValueUnifier) Add(value interface{}) {
	unifier.values[value]++
}
