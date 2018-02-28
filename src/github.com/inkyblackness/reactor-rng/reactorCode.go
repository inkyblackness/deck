package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type reactorCode struct {
	one uint16
	two uint16
}

func newReactorCode() (code reactorCode) {
	limit := big.NewInt(10)
	newDigit := func() uint16 {
		digit, err := rand.Int(rand.Reader, limit)
		if err != nil {
			panic(err)
		}
		return uint16(digit.Int64())
	}
	saveCode := func() (value uint16) {
		for value == 0 {
			value = (newDigit() << 8) | (newDigit() << 4) | (newDigit() << 0)
		}
		return
	}
	code.one = saveCode()
	code.two = saveCode()
	return
}

func (code reactorCode) String() string {
	return fmt.Sprintf("%03X%03X", code.one, code.two)
}
