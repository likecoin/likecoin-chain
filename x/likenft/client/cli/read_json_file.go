package cli

import (
	"encoding/json"
	"io/ioutil"

	"github.com/likecoin/likechain/x/likenft/types"
)

func readClassInputJsonFile(path string) (*types.ClassInput, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	input := types.ClassInput{}
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func readNFTInputJsonFile(path string) (*types.NFTInput, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	input := types.NFTInput{}
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}
