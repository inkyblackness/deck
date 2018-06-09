package resfile

const (
	headerString                = "LG Res File v2\r\n"
	commentTerminator           = byte(0x1A)
	chunkDirectoryFileOffsetPos = 0x7C
	boundarySize                = 4

	chunkTypeFlagCompressed = byte(0x01)
	chunkTypeFlagFragmented = byte(0x02)
)
