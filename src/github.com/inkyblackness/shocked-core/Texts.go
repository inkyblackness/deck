package core

import (
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

type textInfo struct {
	multiblock bool
	limit      uint16
}

var knownTexts = map[model.ResourceType]textInfo{
	model.ResourceTypeTrapMessages:     {false, model.MaxTrapMessages},
	model.ResourceTypeWords:            {false, model.MaxWords},
	model.ResourceTypeLogCategories:    {false, model.MaxLogCategories},
	model.ResourceTypeScreenMessages:   {false, model.MaxScreenMessages},
	model.ResourceTypeInfoNodeMessages: {false, model.MaxInfoNodeMessages},
	model.ResourceTypeAccessCardNames:  {false, model.MaxAccessCardNames},
	model.ResourceTypeDataletMessages:  {false, model.MaxDataletMessages},
	model.ResourceTypePaperTexts:       {true, model.MaxPaperTexts}}

// Texts is the adapter for general texts.
type Texts struct {
	cybstrng [model.LanguageCount]chunk.Store
	cp       text.Codepage
}

// NewTexts returns a new Texts instance, if possible.
func NewTexts(library io.StoreLibrary) (texts *Texts, err error) {
	var cybstrng [model.LanguageCount]chunk.Store

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = library.ChunkStore(localized[i].cybstrng)
	}
	if err == nil {
		texts = &Texts{
			cybstrng: cybstrng,
			cp:       text.DefaultCodepage()}
	}

	return
}

// Text returns the string of identified text.
func (texts *Texts) Text(key model.ResourceKey) (result string, err error) {
	info, known := knownTexts[key.Type]
	if known && (key.Index < info.limit) && key.HasValidLanguage() {
		if info.multiblock {
			holder := texts.cybstrng[key.Language.ToIndex()].Get(res.ResourceID(int(key.Type) + int(key.Index)))

			if holder != nil {
				for blockIndex := uint16(0); blockIndex < holder.BlockCount(); blockIndex++ {
					blockData := holder.BlockData(blockIndex)
					result += texts.cp.Decode(blockData)
				}
			}
		} else {
			holder := texts.cybstrng[key.Language.ToIndex()].Get(res.ResourceID(key.Type))

			if key.Index < holder.BlockCount() {
				blockData := holder.BlockData(key.Index)
				result = texts.cp.Decode(blockData)
			}
		}
	} else {
		err = fmt.Errorf("Unsupported resource key: %v", key)
	}

	return
}

// SetText requests to set the string of a text resource.
func (texts *Texts) SetText(key model.ResourceKey, value string) (resultKey model.ResourceKey, err error) {
	info, known := knownTexts[key.Type]
	emptyString := texts.cp.Encode("")

	if known && (key.Index < info.limit) && key.HasValidLanguage() {
		if info.multiblock {
			store := texts.cybstrng[key.Language.ToIndex()]
			chunkID := res.ResourceID(int(key.Type) + int(key.Index))
			blockData := [][]byte{texts.cp.Encode(value)}

			store.Put(chunkID, chunk.NewBlockHolder(chunk.BasicChunkType.WithDirectory(), res.Text, blockData))
		} else {
			holder := texts.cybstrng[key.Language.ToIndex()].Get(res.ResourceID(key.Type))

			for holder.BlockCount() < key.Index {
				holder.SetBlockData(holder.BlockCount(), emptyString)
			}
			holder.SetBlockData(key.Index, texts.cp.Encode(value))
		}
		resultKey = key
	} else {
		err = fmt.Errorf("Unsupported resource key: %v", key)
	}

	return
}
