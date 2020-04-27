package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/multiformats/go-multibase"
)

var RegistryID = "1" // TODO: move into config

type IscnID []byte

func (id IscnID) String() string {
	return fmt.Sprintf("%s/%s", RegistryID, CidMbaseEncoder.Encode(id))
}

func (id IscnID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *IscnID) UnmarshalJSON(bz []byte) error {
	parts := strings.Split("/", string(bz))
	if len(parts) != 2 {
		return fmt.Errorf("invalid Iscn ID format")
	}
	enc, bz, err := multibase.Decode(parts[1])
	if err != nil {
		return err
	}
	if enc != CidMbaseEncoder.Encoding() {
		return fmt.Errorf("invalid Iscn ID multibase encoding")
	}
	*id = bz
	return nil
}
