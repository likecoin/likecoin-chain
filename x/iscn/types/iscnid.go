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
const IscnIdRegexpPattern = "iscn://([-_.:a-zA-Z0-9]+)/([-_a-zA-Z0-9]+)(?:/([0-9]+))?$"

var IscnIdRegexp *regexp.Regexp

func Base64ToTracingId(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func (iscnId IscnId) TracingIdString() string {
	return base64.RawURLEncoding.EncodeToString(iscnId.TracingId)
}

func (iscnId IscnId) String() string {
	return fmt.Sprintf("iscn://%s/%s/%d", iscnId.RegistryId, iscnId.TracingIdString(), iscnId.Version)
}

func ParseIscnID(s string) (IscnId, error) {
	id := IscnId{}
	matches := IscnIdRegexp.FindStringSubmatch(s)
	if matches == nil {
		return id, fmt.Errorf("invalid ISCN ID format")
	}
	id.RegistryId = matches[1]
	tracingId, err := base64.RawURLEncoding.DecodeString(matches[2])
	if err != nil {
		return id, err
	}
	if len(tracingId) == 0 {
		return id, fmt.Errorf("empty tracing ID")
	}
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
	tracingId := hasher.Sum(nil)
	return IscnId{
		RegistryId: registryId,
		TracingId:  tracingId,
		Version:    1,
	}
}

func init() {
	IscnIdRegexp = regexp.MustCompile(IscnIdRegexpPattern)
}
