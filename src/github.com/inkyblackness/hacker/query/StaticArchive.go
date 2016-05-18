package query

import (
	"bytes"
	"fmt"
)

// StaticArchive queries the data source about a list of chunks in the current archive and verifies whether
// they contain a statically known value.
func StaticArchive(dataSource DataSource) (result string) {
	levelIDs := dataSource.LevelIDs()

	for _, level := range levelIDs {
		result += verifyStaticLevel(dataSource, level)
	}

	return
}

func verifyStaticLevel(dataSource DataSource, level int) (result string) {
	expected := map[int][]byte{
		2:  []byte{0x0B, 0x00, 0x00, 0x00},
		3:  []byte{0x1B, 0x00, 0x00, 0x00},
		40: []byte{0x0D, 0x00, 0x00, 0x00},
		41: []byte{0x00},
		49: make([]byte, 0x1C0),
		50: make([]byte, 2),
		53: make([]byte, 0x40)}

	for levelChunk, expectedData := range expected {
		data := dataSource.LevelChunkData(level, levelChunk)

		if bytes.Compare(data, expectedData) != 0 {
			result += fmt.Sprintf("%02d L%02d (%04X) static mismatch", level, levelChunk, 4000+100*level+levelChunk)
			if data == nil {
				result += " (chunk is empty / not present)"
			}
			result += "\n"
		}
	}

	return
}
