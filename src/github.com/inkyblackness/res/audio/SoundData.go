package audio

// SoundData wraps the basic information about a collection of sound samples.
// The format is mono with usigned 8-bit PCM coding.
type SoundData interface {
	// SampleRate returns the amount of samples per second.
	SampleRate() float32
	// SampleCount returns the number of samples available from this data.
	SampleCount() int
	// Samples returns the samples in the given range
	Samples(from, to int) []byte
}
