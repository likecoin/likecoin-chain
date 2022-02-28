package cli

import (
	"encoding/json"
	"io/ioutil"
)

func readCmdNewClassInput(path string) (*CmdNewClassInput, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	input := CmdNewClassInput{}
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}
