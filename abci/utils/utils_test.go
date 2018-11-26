package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func TestIsValidBigIntegerString(t *testing.T) {
	Convey("For an integer string", t, func() {
		Convey("If input is greater than or equal to 0, it should pass", func() {
			So(IsValidBigIntegerString("0"), ShouldBeTrue)
			So(IsValidBigIntegerString("10000000000000000000000000000000000000000000000000000000000"), ShouldBeTrue)
		})

		Convey("If input is less than 0, it should fail", func() {
			So(IsValidBigIntegerString("-1"), ShouldBeFalse)
		})

		Convey("If input is not a valid integer, it should fail", func() {
			So(IsValidBigIntegerString("0x10"), ShouldBeFalse)
			So(IsValidBigIntegerString("1.1"), ShouldBeFalse)
			So(IsValidBigIntegerString("NaN"), ShouldBeFalse)
		})
	})
}

func TestDbRawKey(t *testing.T) {
	Convey("For DbRawKey", t, func() {
		bs := []byte{0, 1, 2, 3, 255}
		Convey("If the prefix is empty, the key should have no prefix", func() {
			So(DbRawKey(bs, "", "suffix"), ShouldResemble, []byte("\x00\x01\x02\x03\xff_suffix"))
		})
		Convey("If the suffix is empty, the key should have no suffix", func() {
			So(DbRawKey(bs, "prefix", ""), ShouldResemble, []byte("prefix_\x00\x01\x02\x03\xff"))
		})
		Convey("If both the prefix and suffix are empty, the key should be the same as the original bytes", func() {
			So(DbRawKey(bs, "", ""), ShouldResemble, bs)
		})
		Convey("If both the prefix and suffix are not empty, the key should have prefix and suffix", func() {
			So(DbRawKey(bs, "prefix", "suffix"), ShouldResemble, []byte("prefix_\x00\x01\x02\x03\xff_suffix"))
		})
	})
}

func TestDbTxHashKey(t *testing.T) {
	Convey("For DbTxHashKey", t, func() {
		txHash := []byte{0, 1, 2, 3, 255}
		Convey("If the suffix is empty, the key should have no suffix", func() {
			So(DbTxHashKey(txHash, ""), ShouldResemble, []byte("tx:hash:_\x00\x01\x02\x03\xff"))
		})
		Convey("If the suffix is not empty, the key should have suffix", func() {
			So(DbTxHashKey(txHash, "suffix"), ShouldResemble, []byte("tx:hash:_\x00\x01\x02\x03\xff_suffix"))
		})
	})
}

func TestHashRawTx(t *testing.T) {
	Convey("HashRawTx test cases", t, func() {
		for i := 0; i < 100; i++ {
			len := cmn.RandIntn(100)
			bs := cmn.RandBytes(len)
			So(HashRawTx(bs), ShouldResemble, tmhash.Sum(bs))
		}
		So(HashRawTx([]byte{}), ShouldResemble, tmhash.Sum([]byte{}))
	})
}

func TestHex2Bytes(t *testing.T) {
	Convey("For Hex2Bytes", t, func() {
		Convey("For a valid hex string with 0x prefix", func() {
			s := "0x12345678"
			Convey("Hex2Bytes should parse the string into byte array correctly", func() {
				bs, err := Hex2Bytes(s)
				So(err, ShouldBeNil)
				So(bs, ShouldResemble, []byte{0x12, 0x34, 0x56, 0x78})
			})
		})
		Convey("For a valid hex string without 0x prefix", func() {
			s := "12345678"
			Convey("Hex2Bytes should parse the string into byte array correctly", func() {
				bs, err := Hex2Bytes(s)
				So(err, ShouldBeNil)
				So(bs, ShouldResemble, []byte{0x12, 0x34, 0x56, 0x78})
			})
		})
		Convey("For a hex string with half bytes", func() {
			s := "1234567"
			Convey("Hex2Bytes should return error", func() {
				_, err := Hex2Bytes(s)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("For an invalid hex string", func() {
			s := "gg"
			Convey("Hex2Bytes should return error", func() {
				_, err := Hex2Bytes(s)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestEncodeAndDecodeUint64(t *testing.T) {
	Convey("For EncodeUint64 and DecodeUint64", t, func() {
		Convey("For any uint64, DecodeUint64 should decode to the original input from EncodeUint64", func() {
			for i := 0; i < 100; i++ {
				n := cmn.RandUint64()
				So(DecodeUint64(EncodeUint64(n)), ShouldEqual, n)
			}
			So(DecodeUint64(EncodeUint64(0)), ShouldEqual, 0)
			So(DecodeUint64(EncodeUint64(uint64(0xFFFFFFFFFFFFFFFF))), ShouldEqual, uint64(0xFFFFFFFFFFFFFFFF))
		})
		Convey("For any byte array with length 8, EncodeUint64 should decode to the original input from DecodeUint64", func() {
			for i := 0; i < 100; i++ {
				bs := cmn.RandBytes(8)
				So(EncodeUint64(DecodeUint64(bs)), ShouldResemble, bs)
			}
			bs := []byte{0, 0, 0, 0, 0, 0, 0, 0}
			So(EncodeUint64(DecodeUint64(bs)), ShouldResemble, bs)
			bs = []byte{255, 255, 255, 255, 255, 255, 255, 255}
			So(EncodeUint64(DecodeUint64(bs)), ShouldResemble, bs)
		})
	})
}

func TestJoinKeys(t *testing.T) {
	Convey("For JoinKeys", t, func() {
		Convey("If we join 'hello' and 'world', we should get 'hello_world'", func() {
			keys := [][]byte{
				[]byte("hello"),
				[]byte("world"),
			}
			So(JoinKeys(keys), ShouldResemble, []byte("hello_world"))
		})
		Convey("If we join 'hello', we should get 'hello'", func() {
			keys := [][]byte{
				[]byte("hello"),
			}
			So(JoinKeys(keys), ShouldResemble, []byte("hello"))
		})
		Convey("If we join 'hello' and '', we should get 'hello_'", func() {
			keys := [][]byte{
				[]byte("hello"),
				[]byte{},
			}
			So(JoinKeys(keys), ShouldResemble, []byte("hello_"))
		})
		Convey("If we join nothing, we should get nil", func() {
			keys := [][]byte{}
			So(JoinKeys(keys), ShouldBeNil)
		})
	})
}

func TestPrefixKey(t *testing.T) {
	Convey("For PrefixKey", t, func() {
		Convey("If we prefix 'hello' and 'world', we should get 'hello_world_'", func() {
			keys := [][]byte{
				[]byte("hello"),
				[]byte("world"),
			}
			So(PrefixKey(keys), ShouldResemble, []byte("hello_world_"))
		})
		Convey("If we join 'hello', we should get 'hello_'", func() {
			keys := [][]byte{
				[]byte("hello"),
			}
			So(PrefixKey(keys), ShouldResemble, []byte("hello_"))
		})
		Convey("If we join 'hello' and '', we should get 'hello__'", func() {
			keys := [][]byte{
				[]byte("hello"),
				[]byte{},
			}
			So(PrefixKey(keys), ShouldResemble, []byte("hello__"))
		})
		Convey("If we join nothing, we should get nil", func() {
			keys := [][]byte{}
			So(PrefixKey(keys), ShouldBeNil)
		})
	})
}
