package modes

import (
	"github.com/inkyblackness/shocked-client/editor/model"
)

type objectTypeItem struct {
	id          model.ObjectID
	displayName string
}

func (item *objectTypeItem) String() string {
	return item.displayName + " (" + item.id.String() + ")"
}
