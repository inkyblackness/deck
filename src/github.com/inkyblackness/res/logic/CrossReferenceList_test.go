package logic

import (
	"github.com/inkyblackness/res/data"

	check "gopkg.in/check.v1"
)

type CrossReferenceListSuite struct {
	referencer *TestingTileMapReferencer
}

var _ = check.Suite(&CrossReferenceListSuite{})

func (suite *CrossReferenceListSuite) SetUpTest(c *check.C) {
	suite.referencer = NewTestingTileMapReferencer()
}

func (suite *CrossReferenceListSuite) aListOfSize(size int) *CrossReferenceList {
	references := make([]data.LevelObjectCrossReference, size)
	list := &CrossReferenceList{references: references}

	return list
}

func (suite *CrossReferenceListSuite) aClearListOfSize(size int) *CrossReferenceList {
	list := suite.aListOfSize(size)

	list.Clear()

	return list
}

func (suite *CrossReferenceListSuite) someLocations(count int) []TileLocation {
	locations := make([]TileLocation, count)

	for index := 0; index < count; index++ {
		locations[index] = AtTile(uint16(index), uint16(index))
	}

	return locations
}

func (suite *CrossReferenceListSuite) TestNewCrossReferenceListReturnsListWithASizeOf1600(c *check.C) {
	list := NewCrossReferenceList()

	c.Check(list.size(), check.Equals, 1600)
}

func (suite *CrossReferenceListSuite) TestEncodeReturnsExpectedAmountOfBytes(c *check.C) {
	list := suite.aListOfSize(5)

	bytes := list.Encode()

	c.Check(len(bytes), check.Equals, data.LevelObjectCrossReferenceSize*5)
}

func (suite *CrossReferenceListSuite) TestEncodeSerializesAccordingToFormat(c *check.C) {
	list := suite.aListOfSize(1)

	entry0 := list.Entry(0)
	entry0.TileX = 0x0123
	entry0.TileY = 0x4567
	entry0.LevelObjectTableIndex = 0x89AB
	entry0.NextObjectIndex = 0xCDEF
	entry0.NextTileIndex = 0x0011

	bytes := list.Encode()

	c.Check(bytes, check.DeepEquals, []byte{0x23, 0x01, 0x67, 0x45, 0xAB, 0x89, 0xEF, 0xCD, 0x11, 0x00})
}

func (suite *CrossReferenceListSuite) TestDecodeCrossReferenceListSerializesAList(c *check.C) {
	list := suite.aClearListOfSize(3)

	bytes := list.Encode()

	newList := DecodeCrossReferenceList(bytes)

	c.Check(newList, check.DeepEquals, list)
}

func (suite *CrossReferenceListSuite) TestClearResetsTheList(c *check.C) {
	list := suite.aListOfSize(3)

	list.Clear()
	bytes := list.Encode()

	c.Check(bytes, check.DeepEquals, []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0,
		0xFF, 0xFF, 0xFF, 0xFF, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0,
		0xFF, 0xFF, 0xFF, 0xFF, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapReturnsAnIndex(c *check.C) {
	list := suite.aClearListOfSize(4)
	locations := []TileLocation{AtTile(1, 1)}

	index, _ := list.AddObjectToMap(1, suite.referencer, locations)

	c.Check(index, check.Not(check.Equals), CrossReferenceListIndex(0))
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapRegistersIndicesAtMap(c *check.C) {
	list := suite.aClearListOfSize(4)
	location1 := AtTile(1, 1)
	location2 := AtTile(2, 2)
	locations := []TileLocation{location1, location2}

	list.AddObjectToMap(1, suite.referencer, locations)

	c.Check(suite.referencer.ReferenceIndex(location1), check.Not(check.Equals), CrossReferenceListIndex(0))
	c.Check(suite.referencer.ReferenceIndex(location2), check.Not(check.Equals), CrossReferenceListIndex(0))
	c.Check(suite.referencer.ReferenceIndex(location1), check.Not(check.Equals), suite.referencer.ReferenceIndex(location2))
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapSetsPropertiesOfSingleEntry(c *check.C) {
	list := suite.aClearListOfSize(4)
	location1 := AtTile(7, 8)
	locations := []TileLocation{location1}

	index, _ := list.AddObjectToMap(20, suite.referencer, locations)

	firstEntry := list.Entry(index)
	c.Check(firstEntry, check.DeepEquals, &data.LevelObjectCrossReference{
		LevelObjectTableIndex: 20,
		NextObjectIndex:       0,
		NextTileIndex:         uint16(index),
		TileX:                 7,
		TileY:                 8})
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapSetsPropertiesOfMultipleEntries(c *check.C) {
	list := suite.aClearListOfSize(4)
	location1 := AtTile(3, 4)
	location2 := AtTile(5, 6)
	locations := []TileLocation{location1, location2}

	index, _ := list.AddObjectToMap(10, suite.referencer, locations)

	firstEntry := list.Entry(index)
	c.Check(firstEntry, check.DeepEquals, &data.LevelObjectCrossReference{
		LevelObjectTableIndex: 10,
		NextObjectIndex:       0,
		NextTileIndex:         uint16(index - 1),
		TileX:                 5,
		TileY:                 6})

	secondEntry := list.Entry(index - 1)
	c.Check(secondEntry, check.DeepEquals, &data.LevelObjectCrossReference{
		LevelObjectTableIndex: 10,
		NextObjectIndex:       0,
		NextTileIndex:         uint16(index),
		TileX:                 3,
		TileY:                 4})
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapKeepsReferencesOfObjectsInSameTile(c *check.C) {
	list := suite.aClearListOfSize(10)
	location1 := AtTile(16, 12)
	location2 := AtTile(1, 2)

	list.AddObjectToMap(40, suite.referencer, []TileLocation{AtTile(100, 100)})
	existingIndex, _ := list.AddObjectToMap(50, suite.referencer, []TileLocation{location1})
	index, _ := list.AddObjectToMap(60, suite.referencer, []TileLocation{location1, location2})

	firstEntry := list.Entry(index)
	c.Check(firstEntry.NextObjectIndex, check.Equals, uint16(0))

	secondEntry := list.Entry(index - 1)
	c.Check(secondEntry.NextObjectIndex, check.Equals, uint16(existingIndex))
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapReturnsErrorIfExhausted(c *check.C) {
	list := suite.aClearListOfSize(2)

	list.AddObjectToMap(40, suite.referencer, []TileLocation{AtTile(100, 100)})

	_, err := list.AddObjectToMap(50, suite.referencer, []TileLocation{AtTile(10, 10)})

	c.Check(err, check.NotNil)
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapRevertsIntermediateChangesOfEntriesIfExhausted(c *check.C) {
	list := suite.aClearListOfSize(3)

	list.AddObjectToMap(40, suite.referencer, []TileLocation{AtTile(100, 100)})

	_, err := list.AddObjectToMap(50, suite.referencer, []TileLocation{AtTile(10, 10), AtTile(20, 20)})
	c.Assert(err, check.NotNil)

	newIndex, newErr := list.AddObjectToMap(60, suite.referencer, []TileLocation{AtTile(12, 34)})
	c.Check(newErr, check.IsNil)
	c.Check(newIndex, check.Equals, CrossReferenceListIndex(2))
}

func (suite *CrossReferenceListSuite) TestAddObjectToMapRevertsMapReferencesIfExhausted(c *check.C) {
	list := suite.aClearListOfSize(3)

	list.AddObjectToMap(40, suite.referencer, []TileLocation{AtTile(100, 100)})

	_, err := list.AddObjectToMap(50, suite.referencer, []TileLocation{AtTile(10, 10), AtTile(20, 20)})
	c.Assert(err, check.NotNil)

	c.Check(suite.referencer.ReferenceIndex(AtTile(10, 10)), check.Equals, CrossReferenceListIndex(0))
	c.Check(suite.referencer.ReferenceIndex(AtTile(20, 20)), check.Equals, CrossReferenceListIndex(0))
}

func (suite *CrossReferenceListSuite) TestRemoveEntriesFromMapMakesEntriesAvailable(c *check.C) {
	list := suite.aClearListOfSize(3)

	firstIndex, _ := list.AddObjectToMap(40, suite.referencer, suite.someLocations(2))
	list.RemoveEntriesFromMap(firstIndex, suite.referencer)

	_, err := list.AddObjectToMap(50, suite.referencer, suite.someLocations(2))
	c.Check(err, check.IsNil)
}

func (suite *CrossReferenceListSuite) TestRemoveEntriesFromMapClearsMapReferences(c *check.C) {
	list := suite.aClearListOfSize(10)

	locations := suite.someLocations(2)
	firstIndex, _ := list.AddObjectToMap(40, suite.referencer, locations)
	list.RemoveEntriesFromMap(firstIndex, suite.referencer)

	c.Check(suite.referencer.ReferenceIndex(locations[0]), check.Equals, CrossReferenceListIndex(0))
	c.Check(suite.referencer.ReferenceIndex(locations[1]), check.Equals, CrossReferenceListIndex(0))
}

func (suite *CrossReferenceListSuite) TestRemoveEntriesFromMapKeepsEntriesOfOtherObjectsInTile(c *check.C) {
	list := suite.aClearListOfSize(10)

	locations := suite.someLocations(2)
	existingIndex, _ := list.AddObjectToMap(22, suite.referencer, locations[0:1])
	firstIndex, _ := list.AddObjectToMap(40, suite.referencer, locations)
	list.RemoveEntriesFromMap(firstIndex, suite.referencer)

	c.Check(suite.referencer.ReferenceIndex(locations[0]), check.Equals, existingIndex)
}

func (suite *CrossReferenceListSuite) TestRemoveEntriesFromMapModifiesEntriesOfOtherObjectsInTile(c *check.C) {
	list := suite.aClearListOfSize(10)

	locations := suite.someLocations(2)
	existingIndex, _ := list.AddObjectToMap(22, suite.referencer, locations[0:1])
	firstIndex, _ := list.AddObjectToMap(40, suite.referencer, locations)
	latestIndex, _ := list.AddObjectToMap(23, suite.referencer, locations[0:1])
	list.RemoveEntriesFromMap(firstIndex, suite.referencer)

	c.Check(list.Entry(latestIndex).NextObjectIndex, check.Equals, uint16(existingIndex))
}
