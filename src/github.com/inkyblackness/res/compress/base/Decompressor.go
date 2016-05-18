package base

import (
	"io"

	"github.com/inkyblackness/res/serial"
)

type decompressor struct {
	reader *wordReader

	isEndOfStream  bool
	dictionary     *dictEntry
	dictionarySize int
	lastEntry      *dictEntry
	lookup         map[word]*dictEntry

	leftover []byte
}

// NewDecompressor creates a new decompressor instance over a decoder.
func NewDecompressor(coder serial.Coder) io.Reader {
	obj := &decompressor{
		reader: newWordReader(coder), dictionary: rootDictEntry()}
	obj.resetDictionary()

	return obj
}

func (obj *decompressor) resetDictionary() {
	obj.dictionarySize = 0
	obj.lookup = make(map[word]*dictEntry)
	obj.dictionary = rootDictEntry()
	for i := 0; i < 0x100; i++ {
		entry := obj.dictionary.Add(byte(i), word(i))
		obj.lookup[word(i)] = entry
	}
	obj.lastEntry = obj.dictionary
}

func (obj *decompressor) Read(p []byte) (n int, err error) {
	requested := len(p)

	for n < requested && !obj.isEndOfStream {
		n += obj.takeFromLeftover(p, n)
		if n < requested {
			obj.readNextWord()
			n += obj.takeFromLeftover(p, n)
		}
	}

	return
}

func (obj *decompressor) takeFromLeftover(dest []byte, destOffset int) (provided int) {
	requested := len(dest) - destOffset
	available := len(obj.leftover)

	if available > 0 && requested > 0 {
		provided = available
		if provided > requested {
			provided = requested
		}
		copy(dest[destOffset:destOffset+provided], obj.leftover)
		obj.leftover = obj.leftover[provided:]
	}

	return
}

func (obj *decompressor) readNextWord() {
	nextWord := obj.reader.read()

	obj.leftover = obj.lastEntry.Data()
	if nextWord == endOfStream {
		obj.isEndOfStream = true
	} else if nextWord == reset {
		obj.resetDictionary()
	} else {
		nextEntry, nextExisting := obj.lookup[nextWord]

		if nextExisting {
			if obj.lastEntry.depth > 0 {
				obj.addToDictionary(nextEntry.FirstByte())
			}
			obj.lastEntry = nextEntry
		} else if nextWord >= literalLimit {
			nextValue := obj.lastEntry.FirstByte()
			obj.addToDictionary(nextValue)
			obj.lastEntry = obj.lastEntry.next[nextValue]
		} else {
			nextValue := byte(nextWord)
			obj.addToDictionary(nextValue)
			obj.lastEntry = obj.dictionary.next[nextValue]
		}
	}
}

func (obj *decompressor) addToDictionary(value byte) {
	key := word(int(literalLimit) + obj.dictionarySize)
	nextEntry := obj.lastEntry.Add(value, key)
	obj.lookup[key] = nextEntry
	obj.dictionarySize++
}
