package chunk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IdentifiedChunk struct {
	id    Identifier
	chunk *Chunk
}
type ChunkList []IdentifiedChunk

func (list ChunkList) IDs() []Identifier {
	ids := make([]Identifier, 0, len(list))
	for _, entry := range list {
		ids = append(ids, entry.id)
	}
	return ids
}

func (list ChunkList) Chunk(id Identifier) (chunk *Chunk, err error) {
	for _, entry := range list {
		if entry.id.Value() == id.Value() {
			chunk = entry.chunk
		}
	}
	if chunk == nil {
		err = fmt.Errorf("unknown id %v", id)
	}
	return chunk, err
}

type ProviderBackedStoreSuite struct {
	suite.Suite
	provider     ChunkList
	store        Store
	chunkCounter int
}

func TestProviderBackedStoreSuite(t *testing.T) {
	suite.Run(t, new(ProviderBackedStoreSuite))
}

func (suite *ProviderBackedStoreSuite) SetupTest() {
	suite.provider = nil
	suite.store = nil
}

func (suite *ProviderBackedStoreSuite) TestWithEmptyProvider() {
	suite.whenInstanceIsCreated()
	suite.thenIDsShouldBeEmpty()
}

func (suite *ProviderBackedStoreSuite) TestIDsDefaultsToProvider() {
	suite.givenProviderHas(ID(1), suite.aChunk())
	suite.givenProviderHas(ID(2), suite.aChunk())
	suite.whenInstanceIsCreated()
	suite.thenIDsShouldBe([]Identifier{ID(1), ID(2)})
}

func (suite *ProviderBackedStoreSuite) TestChunkDefaultsToChunksFromProvider() {
	chunkA := suite.aChunk()
	chunkB := suite.aChunk()
	suite.givenProviderHas(ID(1), chunkA)
	suite.givenProviderHas(ID(2), chunkB)
	suite.whenInstanceIsCreated()
	suite.thenReturnedChunkShouldBe(ID(1), chunkA)
	suite.thenReturnedChunkShouldBe(ID(2), chunkB)
}

func (suite *ProviderBackedStoreSuite) TestChunkReturnsErrorForUnknownChunk() {
	suite.whenInstanceIsCreated()
	suite.thenChunkShouldReturnErrorFor(ID(10))
}

func (suite *ProviderBackedStoreSuite) TestDelWillHaveStoreIgnoreChunkFromProvider() {
	suite.givenProviderHas(ID(1), suite.aChunk())
	suite.givenProviderHas(ID(2), suite.aChunk())
	suite.givenAnInstance()
	suite.whenChunkIsDeleted(ID(2))
	suite.thenIDsShouldBe([]Identifier{ID(1)})
	suite.thenChunkShouldReturnErrorFor(ID(2))
}

func (suite *ProviderBackedStoreSuite) TestDelWillHaveStoreIgnoreChunkFromProviderEvenIfReportedMultipleTimes() {
	suite.givenProviderHas(ID(1), suite.aChunk())
	suite.givenProviderHas(ID(2), suite.aChunk())
	suite.givenProviderHas(ID(2), suite.aChunk())
	suite.givenAnInstance()
	suite.whenChunkIsDeleted(ID(2))
	suite.thenIDsShouldBe([]Identifier{ID(1)})
	suite.thenChunkShouldReturnErrorFor(ID(2))
}

func (suite *ProviderBackedStoreSuite) TestPutOverridesProviderChunks() {
	suite.givenProviderHas(ID(2), suite.aChunk())
	suite.givenProviderHas(ID(1), suite.aChunk())
	suite.givenAnInstance()
	newChunk := suite.aChunk()
	suite.whenChunkIsPut(ID(2), newChunk)
	suite.thenIDsShouldBe([]Identifier{ID(2), ID(1)})
	suite.thenReturnedChunkShouldBe(ID(2), newChunk)
}

func (suite *ProviderBackedStoreSuite) TestPutAddsNewChunksAtEnd() {
	suite.givenProviderHas(ID(2), suite.aChunk())
	suite.givenProviderHas(ID(1), suite.aChunk())
	suite.givenAnInstance()
	newChunk := suite.aChunk()
	suite.whenChunkIsPut(ID(3), newChunk)
	suite.thenIDsShouldBe([]Identifier{ID(2), ID(1), ID(3)})
	suite.thenReturnedChunkShouldBe(ID(3), newChunk)
}

func (suite *ProviderBackedStoreSuite) givenProviderHas(id Identifier, chunk *Chunk) {
	suite.provider = append(suite.provider, IdentifiedChunk{id, chunk})
}

func (suite *ProviderBackedStoreSuite) givenAnInstance() {
	suite.whenInstanceIsCreated()
}

func (suite *ProviderBackedStoreSuite) whenInstanceIsCreated() {
	suite.store = NewProviderBackedStore(suite.provider)
}

func (suite *ProviderBackedStoreSuite) whenChunkIsDeleted(id Identifier) {
	suite.store.Del(id)
}

func (suite *ProviderBackedStoreSuite) whenChunkIsPut(id Identifier, chunk *Chunk) {
	suite.store.Put(id, chunk)
}

func (suite *ProviderBackedStoreSuite) thenIDsShouldBeEmpty() {
	assert.Equal(suite.T(), 0, len(suite.store.IDs()))
}

func (suite *ProviderBackedStoreSuite) thenIDsShouldBe(expected []Identifier) {
	assert.Equal(suite.T(), expected, suite.store.IDs())
}

func (suite *ProviderBackedStoreSuite) thenReturnedChunkShouldBe(id Identifier, expected *Chunk) {
	chunk, err := suite.store.Chunk(id)
	assert.Nil(suite.T(), err, "No error expected for ID %v", id)
	assert.Equal(suite.T(), expected, chunk, "Different chunk returned for ID %v", id)
}

func (suite *ProviderBackedStoreSuite) thenChunkShouldReturnErrorFor(id Identifier) {
	_, err := suite.store.Chunk(id)
	assert.Error(suite.T(), err, "Error expected for ID %v ", id) // nolint: vet
}

func (suite *ProviderBackedStoreSuite) aChunk() *Chunk {
	suite.chunkCounter++
	return &Chunk{BlockProvider: MemoryBlockProvider([][]byte{{byte(suite.chunkCounter)}})}
}
