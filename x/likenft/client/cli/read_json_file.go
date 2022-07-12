package cli

import (
	"encoding/json"
	"io/ioutil"
)

func readJsonFile[T any](path string) (*T, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var input T
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}
