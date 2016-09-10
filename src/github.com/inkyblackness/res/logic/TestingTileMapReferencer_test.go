package logic

type TestingTileMapReferencer struct {
	references map[TileLocation]CrossReferenceListIndex
}

func NewTestingTileMapReferencer() *TestingTileMapReferencer {
	return &TestingTileMapReferencer{
		references: make(map[TileLocation]CrossReferenceListIndex)}
}

func (ref *TestingTileMapReferencer) ReferenceIndex(location TileLocation) CrossReferenceListIndex {
	return ref.references[location]
}

func (ref *TestingTileMapReferencer) SetReferenceIndex(location TileLocation, index CrossReferenceListIndex) {
	ref.references[location] = index
}
