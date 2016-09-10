package model

import (
	check "gopkg.in/check.v1"
)

type BufferedTextureStoreSuite struct {
	store   *BufferedTextureStore
	queries map[string]int
}

var _ = check.Suite(&BufferedTextureStoreSuite{})

func (suite *BufferedTextureStoreSuite) SetUpTest(c *check.C) {
	suite.queries = make(map[string]int)
	suite.store = NewBufferedTextureStore(func(id TextureKey) {
		suite.queries[id.String()]++
	})
}

func (suite *BufferedTextureStoreSuite) TestTextureReturnsNilForUnknownTexture(c *check.C) {
	texture := suite.store.Texture(GameTextureKeyFor(10))

	c.Check(texture, check.IsNil)
}

func (suite *BufferedTextureStoreSuite) TestTextureRequestsDataForUnknownTexture(c *check.C) {
	suite.store.Texture(GameTextureKeyFor(20))

	c.Check(suite.queries, check.DeepEquals, map[string]int{"20": 1})
}

func (suite *BufferedTextureStoreSuite) TestTextureRequestsDataForAnUnknownTextureOnlyOnce(c *check.C) {
	suite.store.Texture(GameTextureKeyFor(30))
	suite.store.Texture(GameTextureKeyFor(30))
	suite.store.Texture(GameTextureKeyFor(30))

	c.Check(suite.queries, check.DeepEquals, map[string]int{"30": 1})
}

func (suite *BufferedTextureStoreSuite) TestSetTextureRegistersAnInstance(c *check.C) {
	instance := aTexture()
	suite.store.SetTexture(GameTextureKeyFor(40), instance)
	texture := suite.store.Texture(GameTextureKeyFor(40))

	c.Check(texture, check.Equals, instance)
}

func (suite *BufferedTextureStoreSuite) TestSetTextureDisposesPreviousInstance(c *check.C) {
	oldInstance := aTestingTexture()
	newInstance := aTexture()
	suite.store.SetTexture(GameTextureKeyFor(50), oldInstance)
	suite.store.SetTexture(GameTextureKeyFor(50), newInstance)
	texture := suite.store.Texture(GameTextureKeyFor(50))

	c.Assert(texture, check.Equals, newInstance)
	c.Check(oldInstance.disposed, check.Equals, true)
}

func (suite *BufferedTextureStoreSuite) TestTextureDoesNotRequestsDataForAlreadyKnownTexture(c *check.C) {
	suite.store.SetTexture(GameTextureKeyFor(60), aTexture())
	suite.store.Texture(GameTextureKeyFor(60))

	c.Check(suite.queries, check.DeepEquals, map[string]int{})
}

func (suite *BufferedTextureStoreSuite) TestResetDisposesAllTextures(c *check.C) {
	instance1 := aTestingTexture()
	instance2 := aTestingTexture()
	suite.store.SetTexture(GameTextureKeyFor(60), instance1)
	suite.store.SetTexture(GameTextureKeyFor(61), instance2)

	suite.store.Reset()

	c.Check(instance1.disposed, check.Equals, true)
	c.Check(instance2.disposed, check.Equals, true)
}

func (suite *BufferedTextureStoreSuite) TestResetCausesNewQueriesToBeMade(c *check.C) {
	suite.store.SetTexture(GameTextureKeyFor(60), aTexture())

	suite.store.Reset()

	texture := suite.store.Texture(GameTextureKeyFor(60))

	c.Assert(texture, check.IsNil)
	c.Check(suite.queries, check.DeepEquals, map[string]int{"60": 1})
}
