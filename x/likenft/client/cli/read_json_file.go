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

func readCmdUpdateClassInput(path string) (*CmdUpdateClassInput, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	input := CmdUpdateClassInput{}
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func readCmdMintNFTInput(path string) (*CmdMintNFTInput, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	input := CmdMintNFTInput{}
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}
