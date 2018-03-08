package resfile

type chunkDirectoryHeader struct {
	ChunkCount       uint16
	FirstChunkOffset uint32
}
