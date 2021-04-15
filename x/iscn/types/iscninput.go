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
	if input == nil {
		return 4 // `null`
	}
	return len(input)
}

func (input IscnInput) Marshal() ([]byte, error) {
	return input.MarshalJSON()
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
