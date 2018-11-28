package txs

import (
	"testing"

	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEIP712Signature(t *testing.T) {
	Convey("In the beginning", t, func() {
		Convey("Given an EIP-712 sign data", func() {
			n, _ := types.NewBigIntFromString("108925094107761721718559880609268176661768870649949697295083485164148668344506")
			addr := types.Addr("0x1337133713371337133713371337133713371337")
			signData := EIP712SignData{
				Name: "Some Struct",
				Fields: []EIP712Field{
					{"uint64_field", EIP712Uint64(0x1337c0dedeadbeef)},
					{"uint256_field", EIP712Uint256(n)},
					{"string_field", EIP712String("This is SPARTAAAAAAA")},
					{"empty_string_field", EIP712String("")},
					{"bytes32_field", EIP712Bytes32([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32})},
					{"address_field", EIP712Address(*addr)},
				},
			}
			Convey("The hash should be computed correctly", func() {
				hash, err := signData.Hash()
				So(err, ShouldBeNil)
				So(common.Bytes2Hex(hash), ShouldResemble, "67b9798a291aa3cc2f9445a424798888d608c30239f7c28a24cdfda63115561c")
			})
		})
	})
}
