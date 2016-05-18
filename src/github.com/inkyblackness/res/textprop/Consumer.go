package textprop

// Consumer wraps methods to consume texture property data
type Consumer interface {
	// Consume takes the provided data and stores it under the specified ID.
	Consume(id uint32, data []byte)
	// Finish marks the end of consumption. After calling Finish, the consumer can't be used anymore.
	Finish()
}
