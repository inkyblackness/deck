package serial

import (
	"bytes"
	"fmt"
	"reflect"
)

func mapSingleData(valueType reflect.Type, value reflect.Value, coder Coder) {
	valueKind := valueType.Kind()

	if valueKind == reflect.Uint8 {
		temp := byte(value.Uint())
		coder.CodeByte(&temp)
		value.SetUint(uint64(temp))
	} else if valueKind == reflect.Int8 {
		temp := byte(value.Int())
		coder.CodeByte(&temp)
		value.SetInt(int64(temp))
	} else if valueKind == reflect.Uint16 {
		temp := uint16(value.Uint())
		coder.CodeUint16(&temp)
		value.SetUint(uint64(temp))
	} else if valueKind == reflect.Int16 {
		temp := uint16(value.Int())
		coder.CodeUint16(&temp)
		value.SetInt(int64(temp))
	} else if valueKind == reflect.Uint32 {
		temp := uint32(value.Uint())
		coder.CodeUint32(&temp)
		value.SetUint(uint64(temp))
	} else if valueKind == reflect.Int32 {
		temp := uint32(value.Int())
		coder.CodeUint32(&temp)
		value.SetInt(int64(temp))
	} else if valueKind == reflect.String {
		for _, temp := range bytes.NewBufferString(value.String()).Bytes() {
			coder.CodeByte(&temp)
		}

		var buf []byte
		temp := byte(0x00)
		coder.CodeByte(&temp)
		for temp != 0x00 {
			buf = append(buf, temp)
			coder.CodeByte(&temp)
		}
		value.SetString(bytes.NewBuffer(buf).String())
	} else if ((valueKind == reflect.Array) || (valueKind == reflect.Slice)) && (valueType.Elem().Kind() == reflect.Uint8) {
		temp := value.Slice(0, value.Len()).Bytes()
		coder.CodeBytes(temp)
	} else if valueKind == reflect.Array || valueKind == reflect.Slice {
		for j := 0; j < value.Len(); j++ {
			mapSingleData(valueType.Elem(), value.Index(j), coder)
		}
	} else if valueKind == reflect.Struct {
		mapStructData(valueType, value, coder)
	} else if valueKind == reflect.Ptr {
		mapSingleData(valueType.Elem(), reflect.Indirect(value), coder)
	} else if valueKind == reflect.Interface {
		MapData(value.Interface(), coder)
	} else {
		panic(fmt.Errorf("Unknown type <%v>", valueKind))
	}
}

func mapStructData(valueType reflect.Type, value reflect.Value, coder Coder) {
	fields := valueType.NumField()

	for i := 0; i < fields; i++ {
		structField := valueType.Field(i)
		fieldValue := value.Field(i)

		mapSingleData(structField.Type, fieldValue, coder)
	}
}

// MapData either encodes or decodes the given data structure with the provided Coder.
// Only those data types that can be serialized with the Coder are supported.
func MapData(dataStruct interface{}, coder Coder) {
	mapSingleData(reflect.TypeOf(dataStruct), reflect.ValueOf(dataStruct), coder)
}
