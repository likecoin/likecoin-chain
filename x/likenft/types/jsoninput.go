package types

import (
	"encoding/json"
)

type JsonInput json.RawMessage // JSON encoded

// Normalize returns a sorted JSON without indentation
func (input JsonInput) Normalize() (json.RawMessage, error) {
	var v interface{}
	err := json.Unmarshal(input, &v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func (input JsonInput) MarshalJSON() ([]byte, error) {
	return json.RawMessage(input).MarshalJSON()
}

func (input *JsonInput) UnmarshalJSON(bz []byte) error {
	return (*json.RawMessage)(input).UnmarshalJSON(bz)
}

func (input JsonInput) Size() int {
	return len(input)
}

func (input JsonInput) Marshal() ([]byte, error) {
	return input.MarshalJSON()
}

func (input *JsonInput) Unmarshal(bz []byte) error {
	return input.UnmarshalJSON(bz)
}

func (input *JsonInput) MarshalTo(dAtA []byte) (int, error) {
	copy(dAtA, *input)
	return len(*input), nil
}

func (input JsonInput) String() string {
	return string(input)
}

func (input JsonInput) Validate() error {
	var v interface{}
	return json.Unmarshal(input, &v)
}

// for testing
func (input JsonInput) GetPath(path ...interface{}) (interface{}, bool) {
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
