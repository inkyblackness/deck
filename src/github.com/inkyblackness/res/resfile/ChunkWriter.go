package resfile

type chunkWriter interface {
	// finish signals the writer that the chunk is to be finished.
	// Any pending data needs to be sent to the output and the uncompressed
	// size shall be returned.
	finish() uint32
}
