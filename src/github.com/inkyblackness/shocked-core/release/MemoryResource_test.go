package release

import (
	"io/ioutil"

	check "gopkg.in/check.v1"
)

type MemoryResourceSuite struct {
}

var _ = check.Suite(&MemoryResourceSuite{})

func (suite *MemoryResourceSuite) SetUpSuite(c *check.C) {
}

func (suite *MemoryResourceSuite) TestNameReturnsNameOfResource(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", nil)

	c.Check(resource.Name(), check.Equals, "test1.res")
}

func (suite *MemoryResourceSuite) TestPathReturnsRelativePathOfResource(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", nil)

	c.Check(resource.Path(), check.Equals, "rel")
}

func (suite *MemoryResourceSuite) TestAsSourceReturnsReaderForData(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	reader, _ := resource.AsSource()
	data, _ := ioutil.ReadAll(reader)

	c.Assert(reader, check.NotNil)
	c.Check(data, check.DeepEquals, []byte{0x01, 0x02, 0x03})
}

func (suite *MemoryResourceSuite) TestAsSourceCreatesSeparateReaders(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	reader1, _ := resource.AsSource()
	reader2, _ := resource.AsSource()

	reader1.Read(make([]byte, 1))

	data, _ := ioutil.ReadAll(reader2)

	c.Check(data, check.DeepEquals, []byte{0x01, 0x02, 0x03})
}

func (suite *MemoryResourceSuite) TestAsSourceProhibitsAsSinkWhileOpen(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	resource.AsSource()
	_, err := resource.AsSink()

	c.Check(err, check.NotNil)
}

func (suite *MemoryResourceSuite) TestAsSinkProhibitsAsSourceWhileOpen(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	resource.AsSink()
	_, err := resource.AsSource()

	c.Check(err, check.NotNil)
}

func (suite *MemoryResourceSuite) TestAsSinkProhibitsSecondAsSinkWhileOpen(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	resource.AsSink()
	_, err := resource.AsSink()

	c.Check(err, check.NotNil)
}

func (suite *MemoryResourceSuite) TestAsSinkPossibleAfterAllSourcesClosed(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	source, _ := resource.AsSource()
	source.Close()

	sink, err := resource.AsSink()

	c.Assert(err, check.IsNil)
	c.Check(sink, check.NotNil)
}

func (suite *MemoryResourceSuite) TestAsSinkPossibleAfterOtherSinkClosed(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	sink1, _ := resource.AsSink()
	sink1.Close()

	sink2, err := resource.AsSink()

	c.Assert(err, check.IsNil)
	c.Check(sink2, check.NotNil)
}

func (suite *MemoryResourceSuite) TestAsSinkReturnsWriterForData(c *check.C) {
	resource := NewMemoryResource("test1.res", "rel", []byte{0x01, 0x02, 0x03})

	writer, _ := resource.AsSink()
	writer.Write([]byte{0x0A, 0x0B})
	writer.Close()

	reader, _ := resource.AsSource()
	data, _ := ioutil.ReadAll(reader)

	c.Check(data, check.DeepEquals, []byte{0x0A, 0x0B})
}
