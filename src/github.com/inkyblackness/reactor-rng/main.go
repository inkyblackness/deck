package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/resfile"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/serial"
)

func usage() string {
	return Title + `

This application can be used to randomize the reactor code in save-games
that the "System Shock Enhanced Edition" initializes with a static number.

When multiple files are modified in one go, they all receive the same code.

Usage:
   reactor-rng <savefile>...
   reactor-rng -h | --help
   reactor-rng --version

Options:
   -h --help              Show this screen.
   --version              Show version.
`
}

func main() {
	newCode := newReactorCode()
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)

	savefiles := arguments["<savefile>"].([]string)
	for _, savefile := range savefiles {
		patchSaveFile(savefile, newCode)
		fmt.Println("")
	}
}

func patchSaveFile(filePath string, newCode reactorCode) {
	absolutePath, absErr := filepath.Abs(filePath)
	if absErr != nil {
		fmt.Fprintln(os.Stderr, "Could not resolve <"+filePath+">")
		return
	}

	fmt.Println("Processing <" + absolutePath + ">...")
	fmt.Println("Reading file.")
	archive, openErr := openSaveGameStore(absolutePath)
	if openErr != nil {
		fmt.Fprintln(os.Stderr, "Could not open the file. Is it a proper savegame, such as SAVGAM0x.DAT ?")
		return
	}

	fmt.Println("Modifying game state.")
	patchErr := modifySaveGame(archive, newCode)
	if patchErr != nil {
		fmt.Fprintln(os.Stderr, "Could not patch the file. Please report it to the authors.")
		return
	}
	fmt.Println("Saving file.")
	saveErr := writeSaveGame(absolutePath, archive)
	if saveErr != nil {
		fmt.Fprintln(os.Stderr, "Could not save the file. Is it writable and do you have enough storage space?")
		return
	}
	fmt.Println("Done.")
}

func openSaveGameStore(filePath string) (archive chunk.Store, err error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	reader, err := resfile.ReaderFrom(bytes.NewReader(fileData))
	if err != nil {
		return
	}
	archive = chunk.NewProviderBackedStore(reader)
	gameState, gameStateErr := archive.Chunk(chunk.ID(GameStateResourceID))
	if gameStateErr != nil {
		err = fmt.Errorf("game state chunk not found: %v", gameStateErr)
		return
	}
	if !isActiveGameState(gameState) {
		err = fmt.Errorf("game state does not describe a save-game")
	}
	return
}

func isActiveGameState(blockProvider chunk.BlockProvider) bool {
	blockData, dataErr := extractDataFromBlockProvider(blockProvider)
	if dataErr != nil {
		return false
	}

	// Test that this is probably a game state block from an active save-game.
	// Health should be more than zero at this point. Otherwise it might be a start-game archive.
	return (len(blockData) > HealthOffset) && (blockData[HealthOffset] != 0x00)
}

func extractDataFromBlockProvider(blockProvider chunk.BlockProvider) ([]byte, error) {
	blockReader, readerErr := blockProvider.Block(0)
	if readerErr != nil {
		return nil, readerErr
	}
	blockData, dataErr := ioutil.ReadAll(blockReader)
	if dataErr != nil {
		return nil, dataErr
	}
	return blockData, nil
}

func modifySaveGame(archive chunk.Store, newCode reactorCode) (err error) {
	fmt.Println("Changing code.")
	{
		gameStateChunk, _ := archive.Chunk(chunk.ID(GameStateResourceID))
		gameStateData, _ := extractDataFromBlockProvider(gameStateChunk)
		gameStateStore := serial.NewByteStore()
		_, _ = gameStateStore.Write(gameStateData)
		_, _ = gameStateStore.Seek(ReactorCodeOffset, io.SeekStart)
		_ = binary.Write(gameStateStore, binary.LittleEndian, &newCode)
		gameStateChunk.SetBlock(0, gameStateStore.Data())
	}

	// The following code assumes all panels to be code input panels.
	// It checks whether the content matches up to avoid messing up other archives.
	{
		type hackedPanelClassEntry struct {
			data.LevelObjectPrefix

			Unused    [2]byte
			Condition [4]byte

			Code1          uint16
			TriggerObject1 uint16
			Code2          uint16
			TriggerObject2 uint16
			Code3          uint16
			TriggerObject3 uint16
			FailObject     uint16
			Unknown        [4]byte
		}

		panelsClassChunk, _ := archive.Chunk(chunk.ID(ReactorLevelPanelClassDataResourceID))
		if (panelsClassChunk != nil) && (panelsClassChunk.BlockCount() == 1) {
			panelsClassData, _ := extractDataFromBlockProvider(panelsClassChunk)
			panelsClassStore := serial.NewByteStore()
			_, _ = panelsClassStore.Write(panelsClassData)
			var entries [PanelClassDataEntryCount]hackedPanelClassEntry
			_, _ = panelsClassStore.Seek(0, io.SeekStart)
			_ = binary.Read(panelsClassStore, binary.LittleEndian, &entries)

			reactorPanel := &entries[ReactorCodePanelClassIndex]

			hasProperMasterIndex := reactorPanel.LevelObjectTableIndex == ReactorCodePanelMasterIndex
			hasProperTriggerObjects :=
				reactorPanel.TriggerObject1 == ReactorCodePanelTrigger1 &&
					reactorPanel.TriggerObject2 == ReactorCodePanelTrigger2 &&
					reactorPanel.FailObject == ReactorCodePanelFail
			if hasProperMasterIndex && hasProperTriggerObjects {
				fmt.Println("Applying code where it is needed on Citadel.")
				reactorPanel.Code1 = newCode.one
				reactorPanel.Code2 = newCode.two

				_, _ = panelsClassStore.Seek(0, io.SeekStart)
				_ = binary.Write(panelsClassStore, binary.LittleEndian, &entries)
				panelsClassChunk.SetBlock(0, panelsClassStore.Data())

				for codeDigitIndex := 0; codeDigitIndex < 6; codeDigitIndex++ {
					err = modifyLevelScreen(archive, newCode, codeDigitIndex)
					if err != nil {
						return
					}
				}
			} else {
				fmt.Println("Could not find game object to modify. This may be no issue if already used or you are running a fan-mission.")
			}
		}
	}
	return
}

func modifyLevelScreen(archive chunk.Store, newCode reactorCode, codeDigitIndex int) (err error) {
	// As with the panel, the following code assumes all scenery to be display screens.
	// Again, list entries are cross-checked to see whether they are most likely the right objects.
	type hackedScreenClassEntry struct {
		data.LevelObjectPrefix

		FrameCount             uint16
		Mixed                  uint16 // LoopType for screens, or trigger object for pedestals
		AlternationType        uint16
		PictureSource          uint16
		AlternatePictureSource uint16
	}
	levelNumber := codeDigitIndex + 1
	sceneryResourceID := ResourceIDLevelBase + (ResourcesPerLevel * levelNumber) +
		LevelObjectEntryListsResourceIDOffset + SceneryListResourceIDOffset

	sceneryClassChunk, _ := archive.Chunk(chunk.ID(uint16(sceneryResourceID)))
	if (sceneryClassChunk != nil) && (sceneryClassChunk.BlockCount() == 1) {
		sceneryClassData, _ := extractDataFromBlockProvider(sceneryClassChunk)
		sceneryClassStore := serial.NewByteStore()
		_, _ = sceneryClassStore.Write(sceneryClassData)
		var entries [SceneryClassDataEntryCount]hackedScreenClassEntry
		_, _ = sceneryClassStore.Seek(0, io.SeekStart)
		_ = binary.Read(sceneryClassStore, binary.LittleEndian, &entries)

		codeScreenInfoList := levelCodeScreens[codeDigitIndex]
		for _, codeScreenInfo := range codeScreenInfoList {
			codeScreen := &entries[codeScreenInfo.ScreenClassIndex]

			isCodeScreen := codeScreen.LevelObjectTableIndex == codeScreenInfo.ScreenMasterIndex
			showsRandomDigits := codeScreen.PictureSource == RandomDigitPictureSource
			showsDigit := (codeScreen.PictureSource >= CodeDigitPictureSourceStart) &&
				(codeScreen.PictureSource < CodeDigitPictureSourceEnd)

			if isCodeScreen && showsRandomDigits {
				// nothing to do here
			} else if isCodeScreen && showsDigit {
				fullCode := uint32(newCode.one)<<12 | uint32(newCode.two)
				codeScreen.PictureSource = CodeDigitPictureSourceStart + uint16((fullCode>>(4*(5-uint32(codeDigitIndex))))&0x0F)
			} else {
				err = fmt.Errorf("can not patch level %v", levelNumber)
				return
			}
		}

		_, _ = sceneryClassStore.Seek(0, io.SeekStart)
		_ = binary.Write(sceneryClassStore, binary.LittleEndian, &entries)
		sceneryClassChunk.SetBlock(0, sceneryClassStore.Data())
	}
	return
}

func writeSaveGame(filePath string, archive chunk.Provider) (err error) {
	buffer := serial.NewByteStore()
	bufferErr := resfile.Write(buffer, archive)
	if bufferErr != nil {
		return bufferErr
	}
	err = ioutil.WriteFile(filePath, buffer.Data(), 0666)

	return
}
