package base

type dictEntry struct {
	prev  *dictEntry
	depth int

	value byte
	key   word

	next [256]*dictEntry

	first byte
}

func rootDictEntry() *dictEntry {
	return &dictEntry{prev: nil, depth: 0, value: 0x00, key: reset}
}

func (entry *dictEntry) Add(value byte, key word) *dictEntry {
	newEntry := &dictEntry{
		prev:  entry,
		depth: entry.depth + 1,
		value: value,
		key:   key,
		first: entry.first}
	entry.next[value] = newEntry
	if entry.depth == 0 {
		newEntry.first = value
	}

	return newEntry
}

func (entry *dictEntry) Data() []byte {
	bytes := make([]byte, entry.depth, entry.depth)
	cur := entry
	for i := entry.depth - 1; i >= 0; i-- {
		bytes[i] = cur.value
		cur = cur.prev
	}

	return bytes
}

func (entry *dictEntry) FirstByte() byte {
	return entry.first
}
