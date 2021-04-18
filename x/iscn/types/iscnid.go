package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

const TracingIdLength = tmhash.Size
const IscnIdRegexpPattern = "iscn://([-_.:=+,a-zA-Z0-9]+)/([-_.:=+,a-zA-Z0-9]+)(?:/([0-9]+))?$"

var IscnIdRegexp *regexp.Regexp = regexp.MustCompile(IscnIdRegexpPattern)

func (iscnId IscnId) Prefix() string {
	return fmt.Sprintf("%s/%s", iscnId.RegistryId, iscnId.TracingId)
}

func (iscnId IscnId) String() string {
	return fmt.Sprintf("iscn://%s/%s/%d", iscnId.RegistryId, iscnId.TracingId, iscnId.Version)
}

func (iscnId *IscnId) PrefixEqual(iscnId2 *IscnId) bool {
	return iscnId.RegistryId == iscnId2.RegistryId && iscnId.TracingId == iscnId2.TracingId
}

func (iscnId *IscnId) Equal(iscnId2 *IscnId) bool {
	return iscnId.PrefixEqual(iscnId2) && iscnId.Version == iscnId2.Version
}

func ParseIscnID(s string) (IscnId, error) {
	id := IscnId{}
	matches := IscnIdRegexp.FindStringSubmatch(s)
	if matches == nil {
		return id, fmt.Errorf("invalid ISCN ID format")
	}
	id.RegistryId = matches[1]
	tracingId := matches[2]
	id.TracingId = tracingId
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
	parsed, err := ParseIscnID(s)
	if err != nil {
		return err
	}
	*iscnId = parsed
	return nil
}

func GenerateNewIscnIdWithSeed(registryId string, seed []byte) IscnId {
	hasher := tmhash.New()
	hasher.Write([]byte(registryId))
	hasher.Write(seed)
	tracingIdBytes := hasher.Sum(nil)
	tracingId := base64.RawURLEncoding.EncodeToString(tracingIdBytes)
	return IscnId{
		RegistryId: registryId,
		TracingId:  tracingId,
		Version:    1,
	}
}
