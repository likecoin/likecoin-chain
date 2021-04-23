package types

import (
	"encoding/json"
)

type IscnInput json.RawMessage // JSON encoded

// Normalize returns a sorted JSON without indentation
func (input IscnInput) Normalize() (json.RawMessage, error) {
	var v interface{}
	err := json.Unmarshal(input, &v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func (input IscnInput) MarshalJSON() ([]byte, error) {
	return json.RawMessage(input).MarshalJSON()
}

func (input *IscnInput) UnmarshalJSON(bz []byte) error {
	return (*json.RawMessage)(input).UnmarshalJSON(bz)
}

func (input IscnInput) Size() int {
	return len(input)
}

func (input IscnInput) Marshal() ([]byte, error) {
	return input.Normalize()
}

func (input *IscnInput) Unmarshal(bz []byte) error {
	return input.UnmarshalJSON(bz)
}

func (input *IscnInput) MarshalTo(dAtA []byte) (int, error) {
	copy(dAtA, *input)
	return len(*input), nil
}

func (input IscnInput) String() string {
	return string(input)
}

func (input IscnInput) Validate() error {
	var v interface{}
	return json.Unmarshal(input, &v)
}

// for testing
func (input IscnInput) GetPath(path ...interface{}) (interface{}, bool) {
	var v interface{}
	err := json.Unmarshal(input, &v)
	if err != nil {
		return nil, false
	}

	for _, subpath := range path {
		switch subpath.(type) {
		case string:
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil, false
			}
			v, ok = m[subpath.(string)]
			if !ok {
				return nil, false
			}
		case int:
			arr, ok := v.([]interface{})
			if !ok {
				return nil, false
			}
			index := subpath.(int)
			if index < 0 || index >= len(arr) {
				return nil, false
			}
			v = arr[index]
		default:
			panic("invalid subpath type")
		}
	}

	return v, true
}
