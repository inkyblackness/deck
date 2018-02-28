package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/res/chunk/store"
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
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()
	provider, err := dos.NewChunkProvider(file)
	if err != nil {
		return
	}
	archive = store.NewProviderBacked(chunk.NullProvider(), func() {})
	gameStateFound := false
	for _, id := range provider.IDs() {
		blockHolder := provider.Provide(id)
		blockHolder.BlockCount() // this ensures the block data is buffered
		if (id == GameStateResourceID) && isActiveGameState(blockHolder) {
			gameStateFound = true
		}
		archive.Put(id, blockHolder)
	}
	if !gameStateFound {
		err = fmt.Errorf("game state chunk not found. Probably not a save-game")
	}
	return
}

func isActiveGameState(blockHolder chunk.BlockHolder) (isActive bool) {
	if blockHolder.BlockCount() == 1 {
		blockData := blockHolder.BlockData(0)
		// Test that this is probably a game state block from an active save-game.
		// Health should be more than zero at this point. Otherwise it might be a start-game archive.
		isActive = (len(blockData) > HealthOffset) && (blockData[HealthOffset] != 0x00)
	}
	return
}

func modifySaveGame(archive chunk.Store, newCode reactorCode) (err error) {
	fmt.Println("Changing code.")
	{
		gameState := archive.Get(GameStateResourceID)
		gameStateData := serial.NewByteStore()
		_, _ = gameStateData.Write(gameState.BlockData(0))
		_, _ = gameStateData.Seek(ReactorCodeOffset, io.SeekStart)
		_ = binary.Write(gameStateData, binary.LittleEndian, &newCode)
		gameState.SetBlockData(0, gameStateData.Data())
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

		panelsClassInfo := archive.Get(ReactorLevelPanelClassDataResourceID)
		if panelsClassInfo.BlockCount() == 1 {
			panelsClassData := serial.NewByteStore()
			_, _ = panelsClassData.Write(panelsClassInfo.BlockData(0))
			var entries [PanelClassDataEntryCount]hackedPanelClassEntry
			_, _ = panelsClassData.Seek(0, io.SeekStart)
			_ = binary.Read(panelsClassData, binary.LittleEndian, &entries)

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

				_, _ = panelsClassData.Seek(0, io.SeekStart)
				_ = binary.Write(panelsClassData, binary.LittleEndian, &entries)
				panelsClassInfo.SetBlockData(0, panelsClassData.Data())

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

	sceneryClassInfo := archive.Get(res.ResourceID(sceneryResourceID))
	if sceneryClassInfo.BlockCount() == 1 {
		sceneryClassData := serial.NewByteStore()
		_, _ = sceneryClassData.Write(sceneryClassInfo.BlockData(0))
		var entries [SceneryClassDataEntryCount]hackedScreenClassEntry
		_, _ = sceneryClassData.Seek(0, io.SeekStart)
		_ = binary.Read(sceneryClassData, binary.LittleEndian, &entries)

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

		_, _ = sceneryClassData.Seek(0, io.SeekStart)
		_ = binary.Write(sceneryClassData, binary.LittleEndian, &entries)
		sceneryClassInfo.SetBlockData(0, sceneryClassData.Data())
	}
	return
}

func writeSaveGame(filePath string, archive chunk.Store) (err error) {
	buffer := serial.NewByteStore()
	consumer := dos.NewChunkConsumer(buffer)
	for _, id := range archive.IDs() {
		consumer.Consume(id, archive.Get(id))
	}
	consumer.Finish()

	err = ioutil.WriteFile(filePath, buffer.Data(), 0666)

	return
}
