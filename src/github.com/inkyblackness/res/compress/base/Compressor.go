package base

import (
	"io"

	"github.com/inkyblackness/res/serial"
)

type compressor struct {
	writer *wordWriter

	overtime       int
	dictionary     *dictEntry
	dictionarySize int
	curEntry       *dictEntry
}

// NewCompressor creates a new compressor instance over an encoder.
func NewCompressor(coder serial.Coder) io.WriteCloser {
	obj := &compressor{
		writer:         newWordWriter(coder),
		dictionary:     rootDictEntry(),
		dictionarySize: 0,
		overtime:       0}

	obj.resetDictionary()

	return obj
}

func (obj *compressor) resetDictionary() {
	obj.dictionarySize = 0
	for i := 0; i < 0x100; i++ {
		obj.dictionary.Add(byte(i), word(i))
	}
	obj.curEntry = obj.dictionary
}

func (obj *compressor) Close() error {
	obj.writer.write(obj.curEntry.key)
	obj.writer.close()

	return nil
}

func (obj *compressor) Write(p []byte) (n int, err error) {
	n = len(p)

	for _, input := range p {
		obj.addByte(input)
	}

	return
}

func (obj *compressor) addByte(value byte) {
	nextEntry := obj.curEntry.next[int(value)]
	if nextEntry != nil {
		obj.curEntry = nextEntry
	} else {
		obj.writer.write(obj.curEntry.key)

		key := word(int(literalLimit) + obj.dictionarySize)
		if key < reset {
			obj.curEntry.Add(value, key)
			obj.dictionarySize++
		} else {
			obj.onKeySaturation()
		}

		obj.curEntry = obj.dictionary.next[value]
	}
}

func (obj *compressor) onKeySaturation() {
	obj.overtime++
	if obj.overtime > 1000 {
		obj.writer.write(reset)
		obj.resetDictionary()
		obj.overtime = 0
	}
}
