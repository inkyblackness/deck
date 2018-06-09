package core

import (
	"bytes"
	"fmt"
	goIo "io"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/audio"
	memAudio "github.com/inkyblackness/res/audio/mem"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/movi"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

// ElectronicMessages handles all data related to electronic messages.
type ElectronicMessages struct {
	cybstrng [model.LanguageCount]*io.DynamicChunkStore
	cp       text.Codepage
	citalog  [model.LanguageCount]*io.DynamicChunkStore
}

type messageRange struct {
	start int
	end   int
}

func (msgRange messageRange) isRelativeIDValid(id int) bool {
	return ((msgRange.start + id) < msgRange.end) && ((msgRange.end - id) >= msgRange.start)
}

var electronicMessageBases = map[model.ElectronicMessageType]messageRange{
	model.ElectronicMessageTypeMail:     {0x0989, 0x09B8},
	model.ElectronicMessageTypeLog:      {0x09B8, 0x0A98},
	model.ElectronicMessageTypeFragment: {0x0A98, 0x0AA8}}

// NewElectronicMessages returns a new instance of ElectronicMessages.
func NewElectronicMessages(library io.StoreLibrary) (messages *ElectronicMessages, err error) {
	var cybstrng [model.LanguageCount]*io.DynamicChunkStore
	var citalog [model.LanguageCount]*io.DynamicChunkStore

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = library.ChunkStore(localized[i].cybstrng)
	}
	for i := 0; i < model.LanguageCount && err == nil; i++ {
		citalog[i], err = library.ChunkStore(localized[i].citalog)
	}
	if err == nil {
		messages = &ElectronicMessages{
			cybstrng: cybstrng,
			cp:       text.DefaultCodepage(),
			citalog:  citalog}
	}

	return
}

// Remove tries to remove the message.
func (messages *ElectronicMessages) Remove(messageType model.ElectronicMessageType, id int) (err error) {
	msgRange, properType := electronicMessageBases[messageType]
	if properType && msgRange.isRelativeIDValid(id) {
		textChunkID := res.ResourceID(msgRange.start + id)
		audioChunkID := res.ResourceID(msgRange.start + id + 300)

		for languageIndex := 0; languageIndex < model.LanguageCount; languageIndex++ {
			messages.cybstrng[languageIndex].Del(textChunkID)
			messages.citalog[languageIndex].Del(audioChunkID)
		}
	} else {
		err = fmt.Errorf("Wrong message type/range: %v", messageType)
	}

	return
}

// Message tries to retrieve the message data for given identification.
func (messages *ElectronicMessages) Message(messageType model.ElectronicMessageType, id int) (message model.ElectronicMessage, err error) {
	msgRange, properType := electronicMessageBases[messageType]
	setMessageText := func(language int, dataMessage *data.ElectronicMessage) {
		message.Title[language] = stringAsPointer(dataMessage.Title())
		message.Sender[language] = stringAsPointer(dataMessage.Sender())
		message.Subject[language] = stringAsPointer(dataMessage.Subject())
		message.VerboseText[language] = stringAsPointer(dataMessage.VerboseText())
		message.TerseText[language] = stringAsPointer(dataMessage.TerseText())
	}

	if properType && msgRange.isRelativeIDValid(id) {
		chunkID := res.ResourceID(msgRange.start + id)
		holder := messages.cybstrng[0].Get(chunkID)
		emptyString := stringAsPointer("")
		emptyText := [model.LanguageCount]*string{emptyString, emptyString, emptyString}

		message.Title = emptyText
		message.Sender = emptyText
		message.Subject = emptyText
		message.VerboseText = emptyText
		message.TerseText = emptyText

		if holder != nil {
			var dataMessage *data.ElectronicMessage
			dataMessage, err = messages.decodeMessage(holder)

			if err == nil {
				message.NextMessage = intAsPointer(dataMessage.NextMessage())
				message.IsInterrupt = boolAsPointer(dataMessage.IsInterrupt())
				message.ColorIndex = intAsPointer(dataMessage.ColorIndex())
				message.LeftDisplay = intAsPointer(dataMessage.LeftDisplay())
				message.RightDisplay = intAsPointer(dataMessage.RightDisplay())

				setMessageText(0, dataMessage)
			}

			for language := 1; (err == nil) && (language < len(messages.cybstrng)); language++ {
				holder = messages.cybstrng[language].Get(chunkID)
				dataMessage, err = messages.decodeMessage(holder)

				if err == nil {
					setMessageText(language, dataMessage)
				}
			}
		}
	} else {
		err = fmt.Errorf("Wrong message type/range: %v", messageType)
	}

	return
}

// SetMessage updates the properties of a message.
func (messages *ElectronicMessages) SetMessage(messageType model.ElectronicMessageType, id int, message model.ElectronicMessage) (err error) {
	msgRange, properType := electronicMessageBases[messageType]
	setMessageData := func(language int, dataMessage *data.ElectronicMessage) {
		if message.NextMessage != nil {
			dataMessage.SetNextMessage(*message.NextMessage)
		}
		if message.IsInterrupt != nil {
			dataMessage.SetInterrupt(*message.IsInterrupt)
		}
		if message.ColorIndex != nil {
			dataMessage.SetColorIndex(*message.ColorIndex)
		}
		if message.LeftDisplay != nil {
			dataMessage.SetLeftDisplay(*message.LeftDisplay)
		}
		if message.RightDisplay != nil {
			dataMessage.SetRightDisplay(*message.RightDisplay)
		}

		if message.Title[language] != nil {
			dataMessage.SetTitle(*message.Title[language])
		}
		if message.Sender[language] != nil {
			dataMessage.SetSender(*message.Sender[language])
		}
		if message.Subject[language] != nil {
			dataMessage.SetSubject(*message.Subject[language])
		}
		if message.VerboseText[language] != nil {
			dataMessage.SetVerboseText(*message.VerboseText[language])
		}
		if message.TerseText[language] != nil {
			dataMessage.SetTerseText(*message.TerseText[language])
		}
	}

	if properType && msgRange.isRelativeIDValid(id) {
		chunkID := res.ResourceID(msgRange.start + id)

		for language := 0; language < model.LanguageCount; language++ {
			holder := messages.cybstrng[language].Get(chunkID)
			var dataMessage *data.ElectronicMessage
			var langErr error

			if holder != nil {
				dataMessage, langErr = messages.decodeMessage(holder)
			}
			if (dataMessage == nil) || (langErr != nil) {
				dataMessage = data.NewElectronicMessage()
			}
			setMessageData(language, dataMessage)

			messages.cybstrng[language].Put(chunkID, dataMessage.Encode(messages.cp))
		}
	} else {
		err = fmt.Errorf("Wrong message type/range: %v", messageType)
	}

	return
}

type blockProviderStore struct {
	provider *io.DynamicBlockStore
}

func (store blockProviderStore) BlockCount() int {
	return int(store.provider.BlockCount())
}

func (store blockProviderStore) Block(index int) (goIo.Reader, error) {
	return bytes.NewReader(store.provider.BlockData(uint16(index))), nil
}

func (messages *ElectronicMessages) decodeMessage(blockStore *io.DynamicBlockStore) (message *data.ElectronicMessage, err error) {
	wrapper := blockProviderStore{blockStore}
	return data.DecodeElectronicMessage(messages.cp, wrapper)
}

// MessageAudio tries to retrieve the audio data for given key.
func (messages *ElectronicMessages) MessageAudio(messageType model.ElectronicMessageType, id int, language model.ResourceLanguage) (data audio.SoundData, err error) {
	msgRange := electronicMessageBases[messageType]
	if ((messageType == model.ElectronicMessageTypeLog) || (messageType == model.ElectronicMessageTypeMail)) && msgRange.isRelativeIDValid(id) {
		holder := messages.citalog[language.ToIndex()].Get(res.ResourceID(msgRange.start + id + 300))
		if holder != nil {
			blockData := holder.BlockData(0)
			var container movi.Container
			container, err = movi.Read(bytes.NewReader(blockData))

			if err == nil {
				samples := []byte{}
				for index := 0; index < container.EntryCount(); index++ {
					entry := container.Entry(index)
					if entry.Type() == movi.Audio {
						samples = append(samples, entry.Data()...)
					}
				}
				data = memAudio.NewL8SoundData(float32(container.AudioSampleRate()), samples)
			}
		}
	} else {
		err = fmt.Errorf("Wrong message type/range: %v", messageType)
	}
	return
}

// SetMessageAudio tries to set the audio data for given key.
func (messages *ElectronicMessages) SetMessageAudio(messageType model.ElectronicMessageType, id int, language model.ResourceLanguage,
	soundData audio.SoundData) (err error) {
	msgRange := electronicMessageBases[messageType]
	if ((messageType == model.ElectronicMessageTypeLog) || (messageType == model.ElectronicMessageTypeMail)) && msgRange.isRelativeIDValid(id) {
		store := messages.citalog[language.ToIndex()]
		resourceID := res.ResourceID(msgRange.start + id + 300)
		if soundData != nil {
			encodedData := movi.ContainSoundData(soundData)
			store.Put(resourceID,
				&chunk.Chunk{
					ContentType:   chunk.Media,
					BlockProvider: chunk.MemoryBlockProvider([][]byte{encodedData})})
		} else {
			store.Del(resourceID)
		}
	} else {
		err = fmt.Errorf("Wrong message type/range: %v", messageType)
	}
	return
}
