package tendermint

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	ethSecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/types"
)

func signatureToEthereumSigDer(sig, msgHash []byte, ethAddr common.Address) []byte {
	if sig[0] != 0x30 {
		log.
			WithField("sig", common.ToHex(sig)).
			Panic("sig[0] is not 0x30")
	}
	sigLen := int(sig[1])
	if len(sig) != sigLen+2 {
		log.
			WithField("sig", common.ToHex(sig)).
			Panic("sig len not correct")
	}
	if sig[2] != 0x02 {
		log.
			WithField("sig", common.ToHex(sig)).
			Panic("r part not start with 0x02")
	}
	rLen := int(sig[3])
	sStart := 4 + rLen
	r := sig[4:sStart]
	if sig[sStart] != 0x02 {
		log.
			WithField("sig", common.ToHex(sig)).
			Panic("s part not start with 0x02")
	}
	sLen := int(sig[sStart+1])
	sEnd := sStart + 2 + sLen
	s := sig[sStart+2 : sEnd]
	result := make([]byte, 64)
	if rLen < 32 {
		copy(result[32-rLen:], r)
	} else {
		copy(result, r[rLen-32:])
	}
	if sLen < 32 {
		copy(result[64-sLen:], s)
	} else {
		copy(result[32:], s[sLen-32:])
	}
	return signatureToEthereumSig64(result, msgHash, ethAddr)
}

func signatureToEthereumSig64(sig, msgHash []byte, ethAddr common.Address) []byte {
	result := make([]byte, 65)
	copy(result, sig)
	for v := byte(0); v <= 3; v++ {
		result[64] = v
		recoveredPubKey, err := ethCrypto.SigToPub(msgHash, result)
		if err != nil {
			continue
		}
		recoveredAddr := ethCrypto.PubkeyToAddress(*recoveredPubKey)
		if recoveredAddr == ethAddr {
			result[64] = 27 + v
			return result
		}
	}
	log.
		WithField("sig", common.ToHex(sig)).
		WithField("msg_hash", common.ToHex(msgHash)).
		WithField("eth_addr", ethAddr.Hex()).
		Panic("Cannot find v to recover the address from signature")
	return nil
}

// SignatureToEthereumSig transforms an encoded Tendermint secp256k1 signature to an Ethereum one
func SignatureToEthereumSig(sig, msgHash []byte, ethAddr common.Address) []byte {
	switch len(sig) {
	case 64:
		return signatureToEthereumSig64(sig, msgHash, ethAddr)
	default:
		return signatureToEthereumSigDer(sig, msgHash, ethAddr)
	}
}

// PubKeyToEthAddr transforms a Tendermint secp256k1 public key to an Ethereum address
func PubKeyToEthAddr(tmPubKey *secp256k1.PubKeySecp256k1) common.Address {
	x, y := ethSecp256k1.DecompressPubkey(tmPubKey[:])
	ethPubKey := ecdsa.PublicKey{
		Curve: ethCrypto.S256(),
		X:     x,
		Y:     y,
	}
	return ethCrypto.PubkeyToAddress(ethPubKey)
}

// MapValidatorIndexToEthAddr takes the validator list, returns a mapping from validator index to Ethereum address of
// the validator
func MapValidatorIndexToEthAddr(validators []types.Validator) map[int]common.Address {
	tmToEthAddr := make(map[int]common.Address)
	for i, v := range validators {
		pubKey := v.PubKey.(secp256k1.PubKeySecp256k1)
		tmToEthAddr[i] = PubKeyToEthAddr(&pubKey)
		log.
			WithField("validator_index", i).
			WithField("validator_addr", tmToEthAddr[i].Hex()).
			Debug("Mapped validator with index")
	}
	return tmToEthAddr
}
