package video

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// UnpackControlWords reads from an encoded string of bytes a series of packed control words.
// If all is OK, the control words are returned. If the function returns an error, something
// could not be read/decoded properly.
func UnpackControlWords(data []byte) (words []ControlWord, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if len(data) < 4 {
		err = FormatError
	} else {
		var controlBytes uint32
		reader := bytes.NewReader(data)

		binary.Read(reader, binary.LittleEndian, &controlBytes)
		if controlBytes%bytesPerControlWord != 0 {
			err = FormatError
		} else {
			wordCount := int(controlBytes / bytesPerControlWord)
			unpacked := 0

			words = make([]ControlWord, wordCount)
			for unpacked < wordCount {
				var packed PackedControlWord
				binary.Read(reader, binary.LittleEndian, &packed)
				times := packed.Times()
				for i := 0; i < times; i++ {
					words[unpacked+i] = packed.Value()
				}
				unpacked += times
			}
		}
	}

	return
}
