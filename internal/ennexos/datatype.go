package ennexos

import (
	"encoding/json"
	"fmt"
	"os"
)

type DataType int

var _ json.Marshaler = (*DataType)(nil)
var _ json.Unmarshaler = (*DataType)(nil)

//go:generate go run golang.org/x/tools/cmd/stringer -linecomment -type=DataType -output=datatype_string.go
const (
	DataTypeUnspecified DataType = iota

	Int16  // s16
	Int32  // s32
	Int64  // s64
	Uint16 // u16
	Uint32 // u32
	Uint64 // u64
	Text   // text
)

var allDataTypes = []DataType{
	DataTypeUnspecified,
	Int16,
	Int32,
	Int64,
	Uint16,
	Uint32,
	Uint64,
	Text,
}

var dataTypeByName = map[string]DataType{}

func init() {
	for _, i := range allDataTypes {
		dataTypeByName[i.String()] = i
	}
}

func (t DataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *DataType) UnmarshalJSON(data []byte) error {
	var name string

	if err := json.Unmarshal(data, &name); err != nil {
		return nil
	}

	value, ok := dataTypeByName[name]
	if !ok {
		return fmt.Errorf("%w: unsupported data type %q", os.ErrInvalid, name)
	}

	*t = value

	return nil
}
