package types

import (
	"math"

	gocid "github.com/ipfs/go-cid"
)

type CID = gocid.Cid

type IscnDataType int

const (
	None           IscnDataType = iota // nil
	Number                             // uint32 / uint64 / float
	String                             // string
	NestedCID                          // CID (either parsed as CID, or unparsed as RawIscnMap)
	NestedIscnData                     // RawIscnMap
	// TODO: URL? Datetime? Bytes?
	Array
	Unknown
)

func (t IscnDataType) String() string {
	switch t {
	case None:
		return "None"
	case Number:
		return "Number"
	case String:
		return "String"
	case NestedCID:
		return "NestedCID"
	case NestedIscnData:
		return "NestedIscnData"
	case Array:
		return "Array"
	default:
		return "Unknown"
	}
}

// need a wrapper since interface{} can't be method receiver
type IscnDataField struct {
	value interface{}
}

type IscnDataArray []interface{}

type RawIscnMap = map[string]interface{}

type IscnData RawIscnMap

func (data IscnData) Get(field string) IscnDataField {
	return IscnDataField{RawIscnMap(data)[field]}
}

func (data IscnData) Set(field string, value interface{}) {
	RawIscnMap(data)[field] = value
}

func (arr IscnDataArray) Len() int {
	return len(arr)
}

func (arr IscnDataArray) Get(index int) IscnDataField {
	if index < 0 || index >= len(arr) {
		return IscnDataField{nil}
	}
	return IscnDataField{arr[index]}
}

func (arr IscnDataArray) Set(index int, v interface{}) bool {
	if index < 0 || index >= len(arr) {
		return false
	}
	arr[index] = v
	return true
}

func cidStrToCID(s string) (*CID, bool) {
	cid, err := gocid.Decode(s)
	if err != nil {
		return nil, false
	}
	return &cid, true
}

func ParseCID(v interface{}) (*CID, bool) {
	switch v.(type) {
	case *CID:
		return v.(*CID), true
	case CID:
		cid := v.(CID)
		return &cid, true
	case string:
		return cidStrToCID(v.(string))
	case RawIscnMap:
		m, ok := v.(RawIscnMap)
		if ok && CheckIscnType(m["/"]) == String {
			cidStr := m["/"].(string)
			return cidStrToCID(cidStr)
		}
	}
	return nil, false
}

func CheckIscnType(v interface{}) IscnDataType {
	switch v.(type) {
	case nil:
		return None
	case float32, float64, int, uint, int32, uint32, int64, uint64:
		return Number
	case string:
		return String
	case *CID:
		return NestedCID
	case RawIscnMap:
		m, ok := v.(RawIscnMap)
		if ok && CheckIscnType(m["/"]) == String {
			_, ok := ParseCID(m["/"])
			if ok {
				return NestedCID
			}
		}
		return NestedIscnData
	case []interface{}:
		// TODO: array of other types should be []interface{}?
		return Array
	}
	// TODO: any other cases? e.g. Time?
	return Unknown
}

func (field IscnDataField) Type() IscnDataType {
	return CheckIscnType(field.value)
}

func (field IscnDataField) AsAny() (interface{}, bool) {
	return field.value, true
}

func (field IscnDataField) AsString() (string, bool) {
	s, ok := field.value.(string)
	return s, ok
}

func (field IscnDataField) AsFloat64() (float64, bool) {
	switch field.value.(type) {
	case float64:
		return field.value.(float64), true
	case float32:
		return float64(field.value.(float32)), true
	case int:
		return float64(field.value.(int)), true
	case uint:
		return float64(field.value.(uint)), true
	case int32:
		return float64(field.value.(int32)), true
	case uint32:
		return float64(field.value.(uint32)), true
	case int64:
		n := field.value.(int64)
		if n >= 1<<53 {
			return 0, false
		}
		return float64(n), true
	case uint64:
		n := field.value.(uint64)
		if n >= 1<<53 {
			return 0, false
		}
		return float64(n), true
	default:
		return 0, false
	}
}

func (field IscnDataField) AsUint64() (uint64, bool) {
	var f float64
	switch field.value.(type) {
	case int:
		n := field.value.(int)
		if n < 0 {
			return 0, false
		}
		return uint64(n), true
	case uint:
		return uint64(field.value.(uint)), true
	case int32:
		n := field.value.(int32)
		if n < 0 {
			return 0, false
		}
		return uint64(n), true
	case uint32:
		return uint64(field.value.(uint32)), true
	case int64:
		n := field.value.(int64)
		if n < 0 {
			return 0, false
		}
		return uint64(n), true
	case uint64:
		return field.value.(uint64), true
	case float64:
		f = field.value.(float64)
	case float32:
		f = float64(field.value.(float32))
	default:
		return 0, false
	}
	if f > math.MaxUint64 {
		return 0, false
	}
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return 0, false
	}
	n := math.Floor(f)
	if n != f {
		return 0, false
	}
	return uint64(n), true
}

func (field IscnDataField) AsCID() (*CID, bool) {
	return ParseCID(field.value)
}

func (field IscnDataField) AsArray() (IscnDataArray, bool) {
	arr, ok := field.value.([]interface{})
	return arr, ok
}

func (field IscnDataField) AsRawMap() (RawIscnMap, bool) {
	rawMap, ok := field.value.(RawIscnMap)
	return rawMap, ok
}

func (field IscnDataField) AsIscnData() (IscnData, bool) {
	rawMap, ok := field.value.(RawIscnMap)
	if !ok {
		return nil, false
	}
	return IscnData(rawMap), true
}

// Basically for ISCN ID only
func (field IscnDataField) AsBytes() ([]byte, bool) {
	bz, ok := field.value.([]byte)
	return bz, ok
}

type IscnSchema func(RawIscnMap) bool

func (f IscnSchema) ConstructIscnData(rawMap RawIscnMap) (IscnData, bool) {
	if !f(rawMap) {
		return nil, false
	}
	return IscnData(rawMap), true
}

func NewSchemaValidator(fs ...IscnSchema) IscnSchema {
	return func(rawMap RawIscnMap) bool {
		for _, f := range fs {
			if !f(rawMap) {
				return false
			}
		}
		return true
	}
}

func Field(field string, fs ...ValueValidator) IscnSchema {
	return func(rawMap RawIscnMap) bool {
		v := rawMap[field]
		for _, f := range fs {
			if !f(v) {
				return false
			}
		}
		return true
	}
}

type ValueValidator func(interface{}) bool

func InType(typ ...IscnDataType) ValueValidator {
	return func(v interface{}) bool {
		t := CheckIscnType(v)
		for _, target := range typ {
			if t == target {
				return true
			}
		}
		return false
	}
}

func InSchema(schema IscnSchema) ValueValidator {
	return func(v interface{}) bool {
		if CheckIscnType(v) != NestedIscnData {
			return false
		}
		return schema(v.(RawIscnMap))
	}
}

func IsArrayOf(fs ...ValueValidator) ValueValidator {
	return func(v interface{}) bool {
		if CheckIscnType(v) != Array {
			return false
		}
		arr := v.([]interface{})
		for _, v = range arr {
			for _, f := range fs {
				if !f(v) {
					return false
				}
			}
		}
		return true
	}
}

func Any(fs ...ValueValidator) ValueValidator {
	return func(v interface{}) bool {
		for _, f := range fs {
			if f(v) {
				return true
			}
		}
		return false
	}
}

func IsUint32(v interface{}) bool {
	isInt := false
	isFloat := false
	var n int64
	var f float64
	switch v.(type) {
	case uint32:
		return true
	case uint:
		return v.(uint) <= math.MaxUint32
	case uint64:
		return v.(uint64) <= math.MaxUint32
	case int:
		n = int64(v.(int))
		isInt = true
	case int32:
		n = int64(v.(int32))
		isInt = true
	case int64:
		n = int64(v.(int64))
		isInt = true
	case float32:
		f = float64(v.(float32))
		isFloat = true
	case float64:
		f = v.(float64)
		isFloat = true
	default:
		return false
	}
	if isInt {
		return n >= 0 && n <= math.MaxUint32
	}
	if isFloat {
		if f > math.MaxUint32 {
			return false
		}
		if math.IsInf(f, 0) || math.IsNaN(f) {
			return false
		}
		return math.Floor(f) == f
	}
	return false
}

func IsUint64(v interface{}) bool {
	isInt := false
	isFloat := false
	var n int64
	var f float64
	switch v.(type) {
	case uint, uint32, uint64:
		return true
	case int:
		n = int64(v.(int))
		isInt = true
	case int32:
		n = int64(v.(int32))
		isInt = true
	case int64:
		n = int64(v.(int64))
		isInt = true
	case float32:
		f = float64(v.(float32))
		isFloat = true
	case float64:
		f = v.(float64)
		isFloat = true
	default:
		return false
	}
	if isInt {
		return n >= 0
	}
	if isFloat {
		if f > math.MaxUint64 {
			return false
		}
		if math.IsInf(f, 0) || math.IsNaN(f) {
			return false
		}
		return math.Floor(f) == f
	}
	return false
}

func IsCIDWithCodec(codec uint64) ValueValidator {
	return func(v interface{}) bool {
		if CheckIscnType(v) != NestedCID {
			return false
		}
		cid, ok := ParseCID(v)
		return ok && cid.Prefix().GetCodec() == codec
	}
}

func IsURI(v interface{}) bool {
	// TODO
	return false
}

func IsTimestamp(v interface{}) bool {
	// TODO
	return false
}
