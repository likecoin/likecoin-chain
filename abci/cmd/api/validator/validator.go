package validator

import (
	"encoding/base64"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/likecoin/likechain/abci/utils"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	ethAddressRegexString   = `^0x[0-9a-fA-F]{40}$`
	ethSignatureRegexString = `^0x[0-9a-f]{130}$`
)

var (
	ethAddressRegex   = regexp.MustCompile(ethAddressRegexString)
	ethSignatureRegex = regexp.MustCompile(ethSignatureRegexString)
)

// ValidateBigInteger validates big integer
func ValidateBigInteger(
	v *validator.Validate,
	topStruct reflect.Value,
	currentStructOrField reflect.Value,
	field reflect.Value,
	fieldType reflect.Type,
	fieldKind reflect.Kind,
	param string,
) bool {
	if s, ok := field.Interface().(string); ok {
		return utils.IsValidBigIntegerString(s)
	}
	return false
}

// IsEthereumAddress is a modified implementation of validator.v9 IsEtheremAddress
// Ref: https://github.com/go-playground/validator/blob/v9/baked_in.go#L420
func IsEthereumAddress(
	v *validator.Validate,
	topStruct reflect.Value,
	currentStructOrField reflect.Value,
	field reflect.Value,
	fieldType reflect.Type,
	fieldKind reflect.Kind,
	param string,
) bool {
	if addr, ok := field.Interface().(string); ok {
		return ethAddressRegex.MatchString(addr)
	}

	return false
}

// IsEthereumSignature validates Ethereum signature
func IsEthereumSignature(
	v *validator.Validate,
	topStruct reflect.Value,
	currentStructOrField reflect.Value,
	field reflect.Value,
	fieldType reflect.Type,
	fieldKind reflect.Kind,
	param string,
) bool {
	if sig, ok := field.Interface().(string); ok {
		if !ethSignatureRegex.MatchString(sig) {
			return false
		}
	}

	return true
}

// IsIdentity validates LikeChain identity
func IsIdentity(
	v *validator.Validate,
	topStruct reflect.Value,
	currentStructOrField reflect.Value,
	field reflect.Value,
	fieldType reflect.Type,
	fieldKind reflect.Kind,
	param string,
) bool {
	if validator.IsBase64(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) {
		data, err := base64.StdEncoding.DecodeString(field.String())
		if err != nil {
			return false
		}
		return len(data) == 20
	}
	return IsEthereumAddress(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
}

// Bind binds all custom validators
func Bind() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("biginteger", ValidateBigInteger)
		v.RegisterValidation("eth_addr", IsEthereumAddress)
		v.RegisterValidation("eth_sig", IsEthereumSignature)
		v.RegisterValidation("identity", IsIdentity)
	}
}
