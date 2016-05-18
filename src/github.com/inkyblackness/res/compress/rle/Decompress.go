package rle

import (
	"fmt"
	"io"
)

// Decompress decompresses from the given reader and writes into the provided output buffer.
func Decompress(reader io.Reader, output []byte) (err error) {
	outIndex := 0
	done := false
	nextByte := func() byte {
		zz := []byte{0x00}
		_, err = reader.Read(zz)
		return zz[0]
	}

	for !done && (err == nil) {
		first := nextByte()

		if first == 0x00 {
			nn := nextByte()
			zz := nextByte()

			outIndex += writeBytesOfValue(output[outIndex:outIndex+int(nn)], func() byte { return zz })
		} else if first < 0x80 {
			outIndex += writeBytesOfValue(output[outIndex:outIndex+int(first)], nextByte)
		} else if first == 0x80 {
			control := uint16(nextByte())
			control += uint16(nextByte()) << 8
			if control == 0x0000 {
				done = true
			} else if control < 0x8000 {
				outIndex += int(control)
			} else if control < 0xC000 {
				outIndex += writeBytesOfValue(output[outIndex:outIndex+int(control&0x3FFF)], nextByte)
			} else if (control & 0xFF00) == 0xC000 {
				err = fmt.Errorf("Undefined case 80 nn C0")
			} else {
				zz := nextByte()

				outIndex += writeBytesOfValue(output[outIndex:outIndex+int(control&0x3FFF)], func() byte { return zz })
			}
		} else {
			outIndex += int(first & 0x7F)
		}
	}

	return
}

func writeBytesOfValue(buffer []byte, producer func() byte) int {
	count := len(buffer)
	for i := 0; i < count; i++ {
		buffer[i] = producer()
	}
	return count
}
