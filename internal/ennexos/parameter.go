package ennexos

import (
	"encoding/json"
	"os"
)

type ParameterFile struct {
	Header     ParameterHeader
	Parameters []Parameter
}

type ParameterHeader struct {
	Description string `json:"description"`
}

type ParameterFormatFixed struct {
	DecimalShift int     `json:"decimal_shift"`
	LowerBound   *int64  `json:"lower_bound,omitempty"`
	UpperBound   *int64  `json:"upper_bound,omitempty"`
	Default      *int64  `json:"default,omitempty"`
	Unit         *string `json:"unit,omitempty"`
}

type ParameterFormatTaglist struct {
	Default *int64           `json:"default,omitempty"`
	Values  map[int64]string `json:"values"`
}

type ParameterFormat struct {
	Fixed   *ParameterFormatFixed   `json:"fixed,omitempty"`
	Taglist *ParameterFormatTaglist `json:"taglist,omitempty"`
	Unknown *string                 `json:"unknown,omitempty"`
}

type Parameter struct {
	Address    uint32          `json:"address"`
	Name       string          `json:"name"`
	AccessRead bool            `json:"access_read"`
	DataType   DataType        `json:"data_type"`
	DataFormat ParameterFormat `json:"data_format"`
	Group      string          `json:"group"`
	Channel    string          `json:"channel"`
}

func ReadParameterFile(path string) (*ParameterFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result ParameterFile

	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
