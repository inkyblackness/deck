package dos

import (
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/serial"
)

const (
	// MagicHeader is the header value in an object properties file.
	MagicHeader = uint32(0x2D)
)

type typeEntry struct {
	genericOffset  uint32
	genericLength  uint32
	specificOffset uint32
	specificLength uint32
	commonOffset   uint32
}

var errSizeMismatch = fmt.Errorf("Size mismatch")

func codeObjectData(coder serial.PositioningCoder, entry *typeEntry, data *objprop.ObjectData) {
	if uint32(len(data.Generic)) != entry.genericLength {
		panic(errSizeMismatch)
	}
	if uint32(len(data.Specific)) != entry.specificLength {
		panic(errSizeMismatch)
	}
	if uint32(len(data.Common)) != objprop.CommonPropertiesLength {
		panic(errSizeMismatch)
	}

	coder.SetCurPos(entry.genericOffset)
	coder.CodeBytes(data.Generic)
	coder.SetCurPos(entry.specificOffset)
	coder.CodeBytes(data.Specific)
	coder.SetCurPos(entry.commonOffset)
	coder.CodeBytes(data.Common)
}

func expectedDataLength(descriptors []objprop.ClassDescriptor) uint32 {
	length := uint32(0)

	length += uint32(4)
	for _, classDesc := range descriptors {
		length += classDesc.TotalDataLength()
	}

	return length
}

func calculateEntryValues(descriptors []objprop.ClassDescriptor) map[res.ObjectID]*typeEntry {
	startOffset := uint32(4)
	entries := make(map[res.ObjectID]*typeEntry)
	var entryList []*typeEntry

	for classIndex, classDesc := range descriptors {
		genericOffset := startOffset
		specificOffset := startOffset + classDesc.GenericDataLength*classDesc.TotalTypeCount()

		for subclassIndex, subclassDesc := range classDesc.Subclasses {
			for typeIndex := uint32(0); typeIndex < subclassDesc.TypeCount; typeIndex++ {
				entry := &typeEntry{
					genericOffset:  genericOffset,
					genericLength:  classDesc.GenericDataLength,
					specificOffset: specificOffset,
					specificLength: subclassDesc.SpecificDataLength}
				entryKey := res.MakeObjectID(res.ObjectClass(classIndex), res.ObjectSubclass(subclassIndex), res.ObjectType(typeIndex))

				entries[entryKey] = entry
				entryList = append(entryList, entry)
				specificOffset += subclassDesc.SpecificDataLength
				genericOffset += classDesc.GenericDataLength
			}
			startOffset = specificOffset
		}
	}
	for _, entry := range entryList {
		entry.commonOffset = startOffset
		startOffset += objprop.CommonPropertiesLength
	}

	return entries
}
