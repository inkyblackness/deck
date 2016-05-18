package convert

import (
	"encoding/xml"
	"os"

	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/text"
)

// ToTxt extracts all blocks from a given holder and writes
func ToTxt(fileName string, holder chunk.BlockHolder) (result bool) {
	file, _ := os.Create(fileName)

	if file != nil {
		defer file.Close()
		cp := text.DefaultCodepage()
		var decoded Text

		for blockID := uint16(0); blockID < holder.BlockCount(); blockID++ {
			temp := blockID
			blockData := holder.BlockData(blockID)
			decoded.Entries = append(decoded.Entries, TextEntry{Block: &temp, CData: cp.Decode(blockData)})
		}
		enc := xml.NewEncoder(file)
		enc.Indent("", "    ")
		if err := enc.Encode(&decoded); err == nil {
			result = true
		}
	}

	return
}
