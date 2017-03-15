package interpreters

import (
	"sort"
)

type sortedEntries struct {
	entries map[string]*entry
	keys    []string
}

func sortKeys(entries map[string]*entry) []string {
	list := sortedEntries{entries: entries, keys: make([]string, 0, len(entries))}
	for key := range entries {
		list.keys = append(list.keys, key)
	}
	sort.Sort(&list)
	return list.keys
}

func (entries *sortedEntries) Len() int {
	return len(entries.keys)
}

func (entries *sortedEntries) Less(i, j int) bool {
	entryA := entries.entries[entries.keys[i]]
	entryB := entries.entries[entries.keys[j]]

	return entryA.start < entryB.start
}

func (entries *sortedEntries) Swap(i, j int) {
	entries.keys[i], entries.keys[j] = entries.keys[j], entries.keys[i]
}
