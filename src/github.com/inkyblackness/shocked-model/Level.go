package model

type Level struct {
	Identifiable

	Properties LevelProperties `json:"properties"`

	Links []Link `json:"links"`
}
