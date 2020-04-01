package types

import (
	"encoding/json"
)

type Stakeholder struct {
	Type  string `json:"type" yaml:"type"`
	Id    string `json:"id" yaml:"id"`
	Stake uint32 `json:"stake" yaml:"stake"`
}

type Period struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}

type Right struct {
	Holder    string `json:"holder" yaml:"holder"`
	Type      string `json:"type" yaml:"type"`
	Terms     string `json:"terms" yaml:"terms"`
	Period    Period `json:"period" yaml:"period"`
	Territory string `json:"territory" yaml:"territory"`
}

type IscnContent struct {
	Type        string   `json:"type" yaml:"type"`
	Source      string   `json:"source" yaml:"source"`
	Fingerprint string   `json:"fingerprint" yaml:"fingerprint"`
	Feature     string   `json:"feature" yaml:"feature"`
	Edition     string   `json:"edition" yaml:"edition"`
	Title       string   `json:"title" yaml:"title"`
	Description string   `json:"description" yaml:"description"`
	Tags        []string `json:"tags" yaml:"tags"`
	// TODO: allow custom fields
}

type IscnRecord struct {
	Stakeholders []Stakeholder `json:"stakeholders" yaml:"stakeholders"`
	Timestamp    int64         `json:"timestamp" yaml:"timestamp"`
	Parent       string        `json:"parent" yaml:"parent"`
	Version      uint32        `json:"version" yaml:"version"`
	Right        []Right       `json:"right" yaml:"right"`
	Content      IscnContent   `json:"content" yaml:"content"`
}

func (iscnRecord IscnRecord) String() string {
	// TODO: timestamp should be ISO-8601
	bz, err := json.Marshal(iscnRecord)
	if err != nil {
		panic(err)
	}
	return string(bz)
}
