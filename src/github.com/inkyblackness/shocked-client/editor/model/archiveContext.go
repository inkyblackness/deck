package model

type archiveContext interface {
	projectContext

	ActiveArchiveID() string
}
