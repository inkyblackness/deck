package objprop

import (
	"github.com/inkyblackness/res"
)

// Consumer wraps methods to consume object data
type Consumer interface {
	// Consume takes the provided data and associates it with the given ID.
	Consume(id res.ObjectID, data ObjectData)
	// Finish marks the end of consumption. After calling Finish, the consumer can't be used anymore.
	Finish()
}
