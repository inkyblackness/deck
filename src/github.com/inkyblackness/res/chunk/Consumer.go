package chunk

import "github.com/inkyblackness/res"

// Consumer is filled with chunks until it is finished.
type Consumer interface {
	// Consume adds the given chunk to the consumer.
	Consume(id res.ResourceID, chunk BlockHolder)
	// Finish marks the end of consumption. After calling Finish, the consumer can't be used anymore.
	Finish()
}
