package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

const IscnIdRegexpPattern = "iscn://([-_.:=+,a-zA-Z0-9]+)/([-_.:=+,a-zA-Z0-9]+)(?:/([0-9]+))?$"

var IscnIdRegexp *regexp.Regexp = regexp.MustCompile(IscnIdRegexpPattern)

func (prefix IscnIdPrefix) String() string {
	return fmt.Sprintf("iscn://%s/%s", prefix.RegistryName, prefix.ContentId)
}

func NewIscnId(registryName, contentId string, version uint64) IscnId {
	return IscnId{
		Prefix: IscnIdPrefix{
			RegistryName: registryName,
			ContentId:    contentId,
		},
		Version: 1,
	}
}

func (iscnId IscnId) PrefixId() IscnId {
	prefixId := iscnId
	prefixId.Version = 0
	return prefixId
}

func (iscnId IscnId) String() string {
	return fmt.Sprintf("%s/%d", iscnId.Prefix.String(), iscnId.Version)
}

func (iscnId *IscnId) PrefixEqual(iscnId2 *IscnId) bool {
	return iscnId.Prefix.Equal(iscnId2.Prefix)
}

func ParseIscnId(s string) (IscnId, error) {
	id := IscnId{}
	matches := IscnIdRegexp.FindStringSubmatch(s)
	if matches == nil {
		return id, fmt.Errorf("invalid ISCN ID format")
	}
	id.Prefix.RegistryName = matches[1]
	contentId := matches[2]
	id.Prefix.ContentId = contentId
	if len(matches) > 3 && len(matches[3]) > 0 {
		version, err := strconv.ParseUint(matches[3], 10, 64)
		if err != nil {
			return id, err
		}
		id.Version = version
	} else {
		id.Version = 0
	}
	return id, nil
}

func (iscnId IscnId) MarshalJSON() ([]byte, error) {
	return json.Marshal(iscnId.String())
}

func (iscnId *IscnId) UnmarshalJSON(bz []byte) error {
	var s string
	err := json.Unmarshal(bz, &s)
	if err != nil {
		return err
	}
	parsed, err := ParseIscnId(s)
	if err != nil {
		return err
	}
	*iscnId = parsed
	return nil
}

func GenerateNewIscnIdWithSeed(registryName string, seed []byte) IscnId {
	hasher := tmhash.New()
	hasher.Write([]byte(registryName))
	hasher.Write(seed)
	contentIdBytes := hasher.Sum(nil)
	contentId := base64.RawURLEncoding.EncodeToString(contentIdBytes)
	return NewIscnId(registryName, contentId, 1)
}
