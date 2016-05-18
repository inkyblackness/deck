package data

type String struct {
	Value string
}

func NewString(value string) *String {
	return &String{Value: value}
}

func (str *String) String() string {
	return "\"" + str.Value + "\""
}
