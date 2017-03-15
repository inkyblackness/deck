package interpreters

import (
	check "gopkg.in/check.v1"
)

type InstanceSuite struct {
	data []byte
	inst *Instance
}

var _ = check.Suite(&InstanceSuite{})

func (suite *InstanceSuite) SetUpTest(c *check.C) {
	sub1 := New().
		With("subField0", 0, 1).
		With("subField1", 1, 2)

	sub2 := New().
		With("subFieldA", 0, 2).
		With("subFieldB", 2, 1)

	desc := New().
		With("field0", 0, 1).
		With("field1", 1, 1).
		With("field2", 2, 2).
		With("field3", 4, 4).
		With("misaligned", 9, 2).
		With("beyond", 256, 16).
		Refining("sub1", 3, 3, sub1, Always).
		Refining("sub2", 6, 3, sub2, func(inst *Instance) bool { return inst.Get("field0") == 0 }).
		Refining("sub3", 7, 1, New(), Always)

	suite.data = []byte{0x01, 0x5A, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}
	suite.inst = desc.For(suite.data)
}

func (suite *InstanceSuite) TestGetReturnsZeroForUnknownKey(c *check.C) {
	result := suite.inst.Get("unknown")

	c.Check(result, check.Equals, uint32(0))
}

func (suite *InstanceSuite) TestGetReturnsValueLittleEndian(c *check.C) {
	result := suite.inst.Get("field2")

	c.Check(result, check.Equals, uint32(0x0403))
}

func (suite *InstanceSuite) TestGetReturnsZeroForKeyBeyondSize(c *check.C) {
	result := suite.inst.Get("beyond")

	c.Check(result, check.Equals, uint32(0))
}

func (suite *InstanceSuite) TestSetIgnoresMisalignedFields(c *check.C) {
	suite.inst.Set("misaligned", 0xEEFF)

	c.Check(suite.data[8:10], check.DeepEquals, []byte{0x09, 0x0A})
}

func (suite *InstanceSuite) TestSetStoresValue(c *check.C) {
	suite.inst.Set("field3", 0xAABBCCDD)

	c.Check(suite.data[4:8], check.DeepEquals, []byte{0xDD, 0xCC, 0xBB, 0xAA})
}

func (suite *InstanceSuite) TestRefinedForUnknownKeyReturnsDummyInstance(c *check.C) {
	refined := suite.inst.Refined("unknown")

	c.Assert(refined, check.NotNil)
	c.Check(refined.Get("something"), check.Equals, uint32(0))
}

func (suite *InstanceSuite) TestRefinedReturnsInstanceForSubsection(c *check.C) {
	refined := suite.inst.Refined("sub1")

	c.Check(refined.Get("subField1"), check.Equals, uint32(0x0605))
}

func (suite *InstanceSuite) TestRefinedAllowsModificationOfOriginalData(c *check.C) {
	refined := suite.inst.Refined("sub1")
	refined.Set("subField0", 0xAB)

	c.Check(suite.data[3], check.Equals, byte(0xAB))
}

func (suite *InstanceSuite) TestRefinedReturnsInstanceEvenIfNotActive(c *check.C) {
	refined := suite.inst.Refined("sub2")

	c.Check(refined.Get("subFieldA"), check.Equals, uint32(0x0807))
}

func (suite *InstanceSuite) TestKeysReturnsListOfKeysSortedByStartIndex(c *check.C) {
	keys := suite.inst.Keys()

	c.Check(keys, check.DeepEquals, []string{"field0", "field1", "field2", "field3", "misaligned", "beyond"})
}

func (suite *InstanceSuite) TestActiveRefinementsReturnsListOfActiveKeysSortedByStartIndex(c *check.C) {
	keys := suite.inst.ActiveRefinements()

	c.Check(keys, check.DeepEquals, []string{"sub1", "sub3"})
}

func (suite *InstanceSuite) TestActiveRefinementsCanChange(c *check.C) {
	suite.inst.Set("field0", 0)
	keys := suite.inst.ActiveRefinements()

	c.Check(keys, check.DeepEquals, []string{"sub1", "sub2", "sub3"})
}

func (suite *InstanceSuite) TestRawReturnsData(c *check.C) {
	data := suite.inst.Refined("sub2").Raw()

	c.Check(data, check.DeepEquals, []byte{0x07, 0x08, 0x09})
}

func (suite *InstanceSuite) TestRawAllowsModificationOfOriginalData(c *check.C) {
	suite.inst.Refined("sub2").Raw()[1] = 0xEF

	c.Check(suite.data[7], check.Equals, byte(0xEF))
}

func (suite *InstanceSuite) TestUndefinedReturnsUndefinedBits(c *check.C) {
	data := suite.inst.Undefined()

	c.Check(data, check.DeepEquals, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF})
}

func (suite *InstanceSuite) TestUndefinedConsidersActiveRefinements(c *check.C) {
	suite.inst.Set("field0", 0)
	data := suite.inst.Undefined()

	c.Check(data, check.DeepEquals, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF})
}
