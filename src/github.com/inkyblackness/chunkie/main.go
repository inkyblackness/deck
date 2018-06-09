package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	goImage "image"
	"image/color"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/resfile"
	"github.com/inkyblackness/res/compress/rle"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/res/movi"
	"github.com/inkyblackness/res/serial"

	"github.com/inkyblackness/chunkie/convert"
	"github.com/inkyblackness/chunkie/convert/wav"
)

const (
	// Version contains the current version number
	Version = "1.1.0"
	// Name is the name of the application
	Name = "InkyBlackness Chunkie"
	// Title contains a combined string of name and version
	Title = Name + " v." + Version
)

func usage() string {
	return Title + `

Usage:
  chunkie export <resource-file> <chunk-id> [--block=<block-id>] [--raw] [--pal=<palette-file>] [--pal-id=<palette-id>] [--fps=<framerate>] [<folder>]
  chunkie import <resource-file> <chunk-id> [--block=<block-id>] [--compressed] [--force-transparency] <source-file>
  chunkie -h | --help
  chunkie --version

Options:
  <resource-file>        The resource file to work on.
  <chunk-id>             The chunk identifier. Defaults to decimal, use "0x" as prefix for hexadecimal. "all" for all.
  --block=<block-id>     The block identifier. Defaults to decimal, use "0x" as prefix for hexadecimal. "all" for all. [default: 0]
  --raw                  With this flag, the chunk will be exported without conversion to a common file format.
  --compressed           With this flag, imported bitmaps will be compressed.
  --force-transparency   With this flag, imported bitmaps will be marked to have transparency. [default: false]
  --pal=<palette-file>   For handling bitmaps & models, use this palette file to write color information
  --pal-id=<palette-id>  Optional palette chunk identifier. If not provided, uses first palette found in palette-file.
  --fps=<framerate>      The frames per second to emulate when exporting movies. 0 names files after timestamp. [default: 0]
  <folder>               The path of the folder to use. [default: .]
  <source-file>          The source file to import.
  -h --help              Show this screen.
  --version              Show version.
`
}

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)
	fmt.Printf("%v\n", arguments)

	if arguments["export"].(bool) {
		resourceFile := arguments["<resource-file>"].(string)
		inFile, inFileErr := os.Open(resourceFile)
		if inFileErr != nil {
			fmt.Printf("Failed to open file\n")
			return
		}
		defer inFile.Close()
		provider, providerErr := resfile.ReaderFrom(inFile)
		if providerErr != nil {
			fmt.Printf("Failed to read resource file: %v\n", providerErr)
			return
		}
		chunkText := arguments["<chunk-id>"].(string)
		chunkSelection := int64(-1)
		if chunkText != "all" {
			chunkSelection, _ = strconv.ParseInt(chunkText, 0, 16)
		}
		blockText := arguments["--block"].(string)
		blockSelection := int64(-1)
		if blockText != "all" {
			blockSelection, _ = strconv.ParseInt(blockText, 0, 16)
		}
		framesPerSecond, _ := strconv.ParseFloat(arguments["--fps"].(string), 32)
		raw := arguments["--raw"].(bool)
		palArgument := arguments["--pal"]
		palIDArgument := arguments["--pal-id"]
		var palette color.Palette
		paletteID := uint64(0)
		folderArgument := arguments["<folder>"]
		folder := "."

		if palIDArgument != nil {
			paletteID, _ = strconv.ParseUint(palIDArgument.(string), 0, 16)
		}
		if palArgument != nil {
			palette = loadPalette(palArgument.(string), chunk.ID(uint16(paletteID)))
		}
		if folderArgument != nil {
			folder = folderArgument.(string)
		}
		os.MkdirAll(folder, os.FileMode(0755))

		processBlock := func(chunkID chunk.Identifier, selectedChunk *chunk.Chunk, blockID int) {
			outFileName := fmt.Sprintf("%04X_%03d", chunkID, blockID)
			exportFile(provider, selectedChunk, blockID, path.Join(folder, outFileName), raw, palette, float32(framesPerSecond))
		}
		processChunk := func(chunkID chunk.Identifier) {
			selectedChunk, chunkErr := provider.Chunk(chunkID)

			if chunkErr != nil {
				fmt.Printf("Failed to read chunk %v: %v\n", chunkID, chunkErr)
				return
			}
			if blockSelection == -1 {
				for blockID := 0; blockID < selectedChunk.BlockCount(); blockID++ {
					processBlock(chunkID, selectedChunk, blockID)
				}
			} else {
				processBlock(chunkID, selectedChunk, int(blockSelection))
			}
		}
		if chunkSelection == -1 {
			for _, chunkID := range provider.IDs() {
				processChunk(chunkID)
			}
		} else {
			processChunk(chunk.ID(uint16(chunkSelection)))
		}

	} else if arguments["import"].(bool) {
		resourceFile := arguments["<resource-file>"].(string)
		chunkID, _ := strconv.ParseUint(arguments["<chunk-id>"].(string), 0, 16)
		blockID, _ := strconv.ParseUint(arguments["--block"].(string), 0, 16)
		sourceFile := arguments["<source-file>"].(string)
		compressed := arguments["--compressed"].(bool)
		forceTransparency := arguments["--force-transparency"].(bool)

		importData(resourceFile, chunk.ID(uint16(chunkID)), int(blockID), sourceFile, compressed, forceTransparency)
	}
}

func exportFile(provider chunk.Provider, selectedChunk *chunk.Chunk, blockID int,
	outFileName string, raw bool, palette color.Palette, framesPerSecond float32) {
	blockReader, blockErr := selectedChunk.Block(blockID)
	contentType := selectedChunk.ContentType
	exportRaw := raw

	if blockErr != nil {
		fmt.Printf("Failed to access block %d: %v\n", blockID, blockErr)
		return
	}
	blockData, dataErr := ioutil.ReadAll(blockReader)
	if dataErr != nil {
		fmt.Printf("Failed to read block %d: %v\n", blockID, dataErr)
		return
	}
	if !exportRaw {
		if contentType == chunk.Sound {
			soundData, _ := audio.DecodeSoundChunk(blockData)
			wav.ExportToWav(outFileName+".wav", soundData)
		} else if contentType == chunk.Media {
			exportRaw = exportMedia(blockData, outFileName, framesPerSecond)
		} else if contentType == chunk.Bitmap {
			exportRaw = !convert.ToPng(outFileName+".png", blockData, palette)
		} else if contentType == chunk.Geometry {
			exportRaw = !convert.ToWavefrontObj(outFileName, blockData, palette)
		} else if contentType == chunk.VideoClip {
			exportRaw = exportVideoClip(provider, blockData, outFileName, framesPerSecond, palette)
		} else if contentType == chunk.Text {
			// Don't recreate whole XML for each block since convert.ToTxt merge them into one file
			if blockID == 0 {
				exportRaw = !convert.ToTxt(outFileName+".xml", selectedChunk)
			}
		} else {
			exportRaw = true
		}
	}
	if exportRaw {
		ioutil.WriteFile(outFileName+".bin", blockData, os.FileMode(0644))
	}
}

func loadPalette(fileName string, paletteID chunk.Identifier) (pal color.Palette) {
	if len(fileName) > 0 {
		inFile, _ := os.Open(fileName)
		defer inFile.Close()
		reader, readerErr := resfile.ReaderFrom(inFile)

		if readerErr != nil {
			fmt.Printf("Failed to load palette: %v\n", readerErr)
			return
		}
		tryLoad := func(id chunk.Identifier) {
			palChunk, chunkErr := reader.Chunk(id)

			if chunkErr != nil {
				return
			}
			if palChunk != nil && palChunk.ContentType == chunk.Palette && pal == nil {
				palReader, palErr := palChunk.Block(0)
				if palErr != nil {
					return
				}
				pal, _ = image.LoadPalette(palReader)
			}
		}

		tryLoad(paletteID)
		if pal == nil {
			ids := reader.IDs()
			for _, id := range ids {
				tryLoad(id)
			}
		}
	}
	return
}

func exportMedia(blockData []byte, fileBaseName string, framesPerSecond float32) (failed bool) {
	container, err := movi.Read(bytes.NewReader(blockData))

	if err == nil {
		handler := newExportingMediaHandler(fileBaseName, container.MediaDuration(), framesPerSecond, float32(container.AudioSampleRate()))
		dispatcher := movi.NewMediaDispatcher(container, handler)
		more := true

		for more && err == nil {
			more, err = dispatcher.DispatchNext()
		}
		if !more {
			handler.finish()
		}
	}

	if err != nil {
		failed = true
	}
	return
}

func exportVideoClip(provider chunk.Provider, blockData []byte, fileBaseName string, framesPerSecond float32, pal color.Palette) (failed bool) {
	reader := bytes.NewReader(blockData)
	sequence := data.DefaultVideoClipSequence((len(blockData) - data.VideoClipSequenceBaseSize) / data.VideoClipSequenceEntrySize)
	var err error
	clipPalette := make([]color.Color, len(pal))

	clipPalette[0] = color.NRGBA{R: 0, G: 0, B: 0, A: 0xFF}
	copy(clipPalette[1:], pal[1:])

	sequence.Code(serial.NewDecoder(reader))
	{
		var times []float32
		mediaDuration := float32(0.0)
		for _, entry := range sequence.Entries {
			frameTime := float32(entry.FrameTime) / 1000.0
			for i := 0; i < int(entry.LastFrame-entry.FirstFrame)+1; i++ {
				times = append(times, mediaDuration)
				mediaDuration += frameTime
			}
		}

		framesChunk, framesErr := provider.Chunk(chunk.ID(sequence.FramesID))

		if framesErr != nil {
			fmt.Printf("Failed to access chunk for frames: %v", framesErr)
			return
		}
		imageRect := goImage.Rect(0, 0, int(sequence.Width), int(sequence.Height))
		img := goImage.NewPaletted(imageRect, clipPalette)
		handler := newExportingMediaHandler(fileBaseName, mediaDuration, framesPerSecond, 0.0)
		for frameID := 0; frameID < framesChunk.BlockCount() && err == nil; frameID++ {
			frameReader, frameErr := framesChunk.Block(frameID)
			var header image.BitmapHeader

			if frameErr != nil {
				fmt.Printf("Failed to load frame %v.%d: %v\n", chunk.ID(sequence.FramesID), frameID, frameErr)
			} else {
				binary.Read(frameReader, binary.LittleEndian, &header)
				err = rle.Decompress(frameReader, img.Pix)
				handler.OnVideo(times[int(frameID)], img)
			}
		}
		handler.finish()
	}

	if err != nil {
		fmt.Printf("error exporting video clip: %v\n", err)
		failed = true
	}
	return
}

func importData(resourceFile string, chunkID chunk.Identifier, blockID int, sourceFile string,
	compressed, forceTransparency bool) {
	inFile, inFileErr := os.Open(resourceFile)
	if inFileErr != nil {
		fmt.Printf("Failed to open input file: %v\n", inFileErr)
		return
	}
	defer inFile.Close()
	reader, readerErr := resfile.ReaderFrom(inFile)
	if readerErr != nil {
		fmt.Printf("Failed to read resources from input file: %v\n", readerErr)
		return
	}
	store := chunk.NewProviderBackedStore(reader)

	modChunk, chunkErr := store.Chunk(chunkID)
	if chunkErr != nil {
		fmt.Printf("Failed to access chunk to modify: %v\n", chunkErr)
		return
	}
	modChunk.SetBlock(blockID, importFile(sourceFile, modChunk.ContentType, compressed, forceTransparency))

	buffer := serial.NewByteStore()
	writeErr := resfile.Write(buffer, store)
	if writeErr != nil {
		fmt.Printf("Failed to re-encode data: %v\n", writeErr)
		return
	}
	err := ioutil.WriteFile(resourceFile, buffer.Data(), os.FileMode(0644))
	if err != nil {
		fmt.Printf("Failed to save file: %v\n", err)
		return
	}
}

func importFile(sourceFile string, contentType chunk.ContentType, compressed, forceTransparency bool) (data []byte) {
	extension := path.Ext(sourceFile)
	switch extension {
	case ".wav":
		{
			soundData := wav.ImportFromWav(sourceFile)
			if contentType == chunk.Sound {
				data = audio.EncodeSoundChunk(soundData)
			} else if contentType == chunk.Media {
				data = movi.ContainSoundData(soundData)
			}
		}
	case ".png":
		{
			if contentType == chunk.Bitmap {
				data = convert.FromPng(sourceFile, false, compressed, forceTransparency)
			}
		}
	default:
		{
			var dataErr error
			data, dataErr = ioutil.ReadFile(sourceFile)
			if dataErr != nil {
				fmt.Printf("Failed to read from source file: %v\n", dataErr)
			}
		}
	}
	if data == nil {
		fmt.Printf("No data produced from source file - Is the target chunk compatible with input file?\n")
	}

	return
}
