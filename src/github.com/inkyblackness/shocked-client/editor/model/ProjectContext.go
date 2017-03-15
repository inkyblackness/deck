package model

import (
	"github.com/inkyblackness/shocked-model"
)

type projectContext interface {
	simpleStoreFailure(info string) model.FailureFunc
	ActiveProjectID() string
}
