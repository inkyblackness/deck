package chunk

// ProviderBackedStore is a chunk store that defaults to a provider for
// returning chunks.
type ProviderBackedStore struct {
	ids       []Identifier
	retriever map[uint16]func() (*Chunk, error)
}

func cloneIDs(source []Identifier) []Identifier {
	cloned := make([]Identifier, len(source))
	copy(cloned, source)
	return cloned
}

// NewProviderBackedStore returns a new store based on the specified provider.
// The created store will have a snapshot of the current provider; Changes in
// the list of available IDs are not reflected.
func NewProviderBackedStore(provider Provider) *ProviderBackedStore {
	store := &ProviderBackedStore{
		ids:       cloneIDs(provider.IDs()),
		retriever: make(map[uint16]func() (*Chunk, error))}

	defaultToProvider := func(id Identifier) func() (*Chunk, error) {
		return func() (*Chunk, error) { return provider.Chunk(id) }
	}
	for _, id := range store.ids {
		store.retriever[id.Value()] = defaultToProvider(id)
	}

	return store
}

// IDs returns a list of available IDs this store currently contains.
func (store *ProviderBackedStore) IDs() []Identifier {
	return cloneIDs(store.ids)
}

// Chunk returns a chunk for the given identifier.
func (store *ProviderBackedStore) Chunk(id Identifier) (*Chunk, error) {
	retriever, existing := store.retriever[id.Value()]
	if !existing {
		return nil, ErrChunkDoesNotExist(id)
	}
	return retriever()
}

// Del removes the chunk with given identifier from the store.
func (store *ProviderBackedStore) Del(id Identifier) {
	key := id.Value()
	if _, existing := store.retriever[key]; existing {
		delete(store.retriever, key)
		newIDs := make([]Identifier, 0, len(store.ids)-1)
		for _, oldID := range store.ids {
			if oldID.Value() != key {
				newIDs = append(newIDs, oldID)
			}
		}
		store.ids = newIDs
	}
}

// Put (re-)assigns an identifier with data. If no chunk with given ID exists,
// then it is created. Existing chunks are overwritten with the provided data.
func (store *ProviderBackedStore) Put(id Identifier, chunk *Chunk) {
	key := id.Value()
	if _, existing := store.retriever[key]; !existing {
		store.ids = append(store.ids, id)
	}
	store.retriever[key] = func() (*Chunk, error) { return chunk, nil }
}
