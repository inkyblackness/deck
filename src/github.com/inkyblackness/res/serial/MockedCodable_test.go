package serial

type MockedCodable struct {
	calledCoder Coder
}

func (codable *MockedCodable) Code(coder Coder) {
	codable.calledCoder = coder
}
