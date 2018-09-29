package main

import (
	"encoding/base64"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/likecoin/likechain/abci/utils"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	ethAddressRegexString      = `^0x[0-9a-fA-F]{40}$`
	ethAddressUpperRegexString = `^0x[0-9A-F]{40}$`
	ethAddressLowerRegexString = `^0x[0-9a-f]{40}$`
	ethSignatureRegexString    = `^0x[0-9a-f]{130}$`
)

var (
	ethAddressRegex      = regexp.MustCompile(ethAddressRegexString)
	ethaddressRegexUpper = regexp.MustCompile(ethAddressUpperRegexString)
	ethAddressRegexLower = regexp.MustCompile(ethAddressLowerRegexString)
	ethSignatureRegex    = regexp.MustCompile(ethSignatureRegexString)
)

func validateBigInteger(
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

// Ref: https://github.com/go-playground/validator/blob/v9/baked_in.go#L420
func isEthereumAddress(
	v *validator.Validate,
	topStruct reflect.Value,
	currentStructOrField reflect.Value,
	field reflect.Value,
	fieldType reflect.Type,
	fieldKind reflect.Kind,
	param string,
) bool {
	if addr, ok := field.Interface().(string); ok {
		if !ethAddressRegex.MatchString(addr) {
			return false
		}

		if ethaddressRegexUpper.MatchString(addr) ||
			ethAddressRegexLower.MatchString(addr) {
			return true
		}
	}

	return true
}

func isEthereumSignature(
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

func isIdentity(
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
	return isEthereumAddress(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
}

func bindValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("biginteger", validateBigInteger)
		v.RegisterValidation("eth_addr", isEthereumAddress)
		v.RegisterValidation("eth_sig", isEthereumSignature)
		v.RegisterValidation("identity", isIdentity)
	}
}
