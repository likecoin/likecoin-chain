package types

import (
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLikeChainID(t *testing.T) {
	Convey("In the beginning", t, func() {
		Convey("Given a valid base64 string with 20 bytes length", func() {
			idStr := "AAAAAAAAAAAAAAAAAAAAAAAAAAA="
			idBytes, err := base64.StdEncoding.DecodeString(idStr)
			if err != nil {
				panic(err)
			}
			Convey("It should converts to a valid LikeChainID", func() {
				id, err := NewLikeChainIDFromString(idStr)
				So(err, ShouldBeNil)
				Convey("The LikeChainID should have same bytes with the base64 string", func() {
					bs := id.Bytes()
					So(bs, ShouldResemble, idBytes)
					Convey("The LikeChainID should be a base64 string representing the same bytes", func() {
						s := id.String()
						bs, err := base64.StdEncoding.DecodeString(s)
						So(err, ShouldBeNil)
						So(bs, ShouldResemble, idBytes)
					})
				})
			})
			Convey("It should converts to LikeChainID by bytes", func() {
				idByBytes := NewLikeChainID(idBytes)
				So(idByBytes[:], ShouldResemble, idBytes)
				idByBytes = ID(idBytes)
				So(idByBytes[:], ShouldResemble, idBytes)
			})
			Convey("It should converts to a valid LikeChainID without panic", func() {
				So(func() { IDStr(idStr) }, ShouldNotPanic)
			})
		})
		Convey("Given a valid base64 string with wrong length", func() {
			idStr := "AAAAAAAAAAAAAAAAAAAAAAAAAA=="
			_, err := base64.StdEncoding.DecodeString(idStr)
			if err != nil {
				panic(err)
			}
			Convey("Converting it to a LikeChainID should return error", func() {
				_, err := NewLikeChainIDFromString(idStr)
				So(err, ShouldNotBeNil)
				Convey("Converting it to a LikeChainID using IDStr should panic", func() {
					So(func() { IDStr(idStr) }, ShouldPanic)
				})
			})
		})
		Convey("Given a base64 string without padding", func() {
			idStr := "AAAAAAAAAAAAAAAAAAAAAAAAAAA"
			Convey("Converting it to a LikeChainID should return error", func() {
				_, err := NewLikeChainIDFromString(idStr)
				So(err, ShouldNotBeNil)
				Convey("Converting it to a LikeChainID using IDStr should panic", func() {
					So(func() { IDStr(idStr) }, ShouldPanic)
				})
			})
		})
		Convey("Given a base64 string with URL encoding", func() {
			idStr := "___________________________="
			Convey("Converting it to a LikeChainID should return error", func() {
				_, err := NewLikeChainIDFromString(idStr)
				So(err, ShouldNotBeNil)
				Convey("Converting it to a LikeChainID using IDStr should panic", func() {
					So(func() { IDStr(idStr) }, ShouldPanic)
				})
			})
		})
		Convey("Given an invalid base64 string", func() {
			idStr := "!"
			Convey("Converting it to a LikeChainID should return error", func() {
				_, err := NewLikeChainIDFromString(idStr)
				So(err, ShouldNotBeNil)
				Convey("Converting it to a LikeChainID using IDStr should panic", func() {
					So(func() { IDStr(idStr) }, ShouldPanic)
				})
			})
		})
		Convey("Given two identical LikeChainIDs", func() {
			id1 := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAA=")
			id2 := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAA=")
			Convey("Equals() should return true", func() {
				So(id1.Equals(id2), ShouldBeTrue)
				So(id2.Equals(id1), ShouldBeTrue)
			})
		})
		Convey("Given two different LikeChainIDs", func() {
			id1 := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAA=")
			id2 := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAE=")
			Convey("Equals() should return false", func() {
				So(id1.Equals(id2), ShouldBeFalse)
				So(id2.Equals(id1), ShouldBeFalse)
			})
		})
		Convey("Given a LikeChainID and an Address", func() {
			id := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAA=")
			addr := Addr("0000000000000000000000000000000000000000")
			Convey("Equals() should return false", func() {
				So(id.Equals(addr), ShouldBeFalse)
			})
		})
		Convey("Given a valid LikeChainID", func() {
			id := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAA=")
			Convey("It should have correct DB key", func() {
				dbKey := id.DBKey("prefix", "suffix")
				expectedKey := []byte(nil)
				expectedKey = append(expectedKey, []byte("prefix:id:_")...)
				expectedKey = append(expectedKey, id[:]...)
				expectedKey = append(expectedKey, []byte("_suffix")...)
				So(dbKey, ShouldResemble, expectedKey)
			})
		})
	})
}

func TestAddress(t *testing.T) {
	Convey("In the beginning", t, func() {
		Convey("Given a valid hex string with '0x' prefix and 20 bytes length", func() {
			addrHex := "0x0000000000000000000000000000000000000000"
			addrBytes, err := hex.DecodeString(addrHex[2:])
			if err != nil {
				panic(err)
			}
			Convey("It should converts to a valid Address", func() {
				addr, err := NewAddressFromHex(addrHex)
				So(err, ShouldBeNil)
				Convey("The Address should have same bytes with the hex string", func() {
					bs := addr.Bytes()
					So(bs, ShouldResemble, addrBytes)
					Convey("The Address should be a hex string with '0x' prefix representing the same bytes", func() {
						s := addr.String()
						So(s, ShouldStartWith, "0x")
						bs, err := hex.DecodeString(s[2:])
						So(err, ShouldBeNil)
						So(bs, ShouldResemble, addrBytes)
					})
				})
			})
			Convey("It should converts to Address by bytes", func() {
				addrByBytes := NewAddress(addrBytes)
				So(addrByBytes[:], ShouldResemble, addrBytes)
			})
			Convey("It should converts to a valid Address without panic", func() {
				So(func() { Addr(addrHex) }, ShouldNotPanic)
			})
		})
		Convey("Given a valid hex string with 20 bytes length without '0x' prefix", func() {
			addrHex := "0000000000000000000000000000000000000000"
			addrBytes, err := hex.DecodeString(addrHex)
			if err != nil {
				panic(err)
			}
			Convey("It should converts to a valid Address", func() {
				addr, err := NewAddressFromHex(addrHex)
				So(err, ShouldBeNil)
				Convey("The Address should have same bytes with the hex string", func() {
					bs := addr.Bytes()
					So(bs, ShouldResemble, addrBytes)
				})
			})
			Convey("It should converts to Address by bytes", func() {
				addrByBytes := NewAddress(addrBytes)
				So(addrByBytes[:], ShouldResemble, addrBytes)
			})
			Convey("It should converts to a valid Address without panic", func() {
				So(func() { Addr(addrHex) }, ShouldNotPanic)
			})
		})
		Convey("Given a valid hex string with wrong length", func() {
			addrHex := "0x00000000000000000000000000000000000000"
			_, err := hex.DecodeString(addrHex[2:])
			if err != nil {
				panic(err)
			}
			Convey("Converting it to a Address should return error", func() {
				_, err := NewAddressFromHex(addrHex)
				So(err, ShouldNotBeNil)
				Convey("Converting it to a addrHex using Addr should panic", func() {
					So(func() { Addr(addrHex) }, ShouldPanic)
				})
			})
		})
		Convey("Given an hex string with partial bytes", func() {
			addrHex := "000000000000000000000000000000000000000"
			Convey("Converting it to a Address should return error", func() {
				_, err := NewAddressFromHex(addrHex)
				So(err, ShouldNotBeNil)
				Convey("Converting it to a Address using Addr should panic", func() {
					So(func() { Addr(addrHex) }, ShouldPanic)
				})
			})
		})
		Convey("Given an hex string with invalid characters", func() {
			addrHex := "000000000000000000000000000000000000000g"
			Convey("Converting it to a Address should return error", func() {
				_, err := NewAddressFromHex(addrHex)
				So(err, ShouldNotBeNil)
				Convey("Converting it to a Address using Addr should panic", func() {
					So(func() { Addr(addrHex) }, ShouldPanic)
				})
			})
		})
		Convey("Given two identical Addresses", func() {
			addr1 := Addr("0000000000000000000000000000000000000000")
			addr2 := Addr("0x0000000000000000000000000000000000000000")
			Convey("Equals() should return true", func() {
				So(addr1.Equals(addr2), ShouldBeTrue)
				So(addr2.Equals(addr1), ShouldBeTrue)
			})
		})
		Convey("Given two different Addresses", func() {
			addr1 := Addr("0000000000000000000000000000000000000000")
			addr2 := Addr("0000000000000000000000000000000000000001")
			Convey("Equals() should return false", func() {
				So(addr1.Equals(addr2), ShouldBeFalse)
				So(addr2.Equals(addr1), ShouldBeFalse)
			})
		})
		Convey("Given an Address and a LikeChainID", func() {
			id := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAA=")
			addr := Addr("0000000000000000000000000000000000000000")
			Convey("Equals() should return false", func() {
				So(addr.Equals(id), ShouldBeFalse)
			})
		})
		Convey("Given a valid Address", func() {
			addr := Addr("0x0000000000000000000000000000000000000000")
			Convey("It should have correct DB key", func() {
				dbKey := addr.DBKey("prefix", "suffix")
				expectedKey := []byte(nil)
				expectedKey = append(expectedKey, []byte("prefix:addr:_")...)
				expectedKey = append(expectedKey, addr[:]...)
				expectedKey = append(expectedKey, []byte("_suffix")...)
				So(dbKey, ShouldResemble, expectedKey)
			})
		})
	})
}

func TestIdentifier(t *testing.T) {
	Convey("In the beginning", t, func() {
		Convey("Given a valid base64 string with 20 bytes length", func() {
			idStr := "AAAAAAAAAAAAAAAAAAAAAAAAAAA="
			id := IDStr(idStr)
			Convey("NewIdentifier should return a LikeChainID", func() {
				iden := NewIdentifier(idStr)
				So(iden, ShouldResemble, id)
			})
		})
		Convey("Given a valid hex string with 20 bytes length", func() {
			addrHex := "0000000000000000000000000000000000000000"
			addr := Addr(addrHex)
			Convey("NewIdentifier should return an address", func() {
				iden := NewIdentifier(addrHex)
				So(iden, ShouldResemble, addr)
			})
		})
		Convey("Given a string which is neither a LikeChainID nor an Address", func() {
			s := "abc"
			Convey("NewIdentifier should return nil", func() {
				iden := NewIdentifier(s)
				So(iden, ShouldBeNil)
			})
		})
	})
}

func TestAmino(t *testing.T) {
	Convey("In the beginning", t, func() {
		Convey("Given a LikeChainID", func() {
			id := IDStr("AAAAAAAAAAAAAAAAAAAAAAAAAAA=")
			Convey("Marshaling the LikeChainID should succeed", func() {
				bs, err := AminoCodec().MarshalBinary(id)
				So(err, ShouldBeNil)
				Convey("Unmarshaling the bytes should gives out the same LikeChainID", func() {
					var iden Identifier
					err := AminoCodec().UnmarshalBinary(bs, &iden)
					So(err, ShouldBeNil)
					So(iden, ShouldResemble, id)
				})
			})
		})
		Convey("Given an Address", func() {
			addr := Addr("0x0000000000000000000000000000000000000000")
			Convey("Marshaling the Address should succeed", func() {
				bs, err := AminoCodec().MarshalBinary(addr)
				So(err, ShouldBeNil)
				Convey("Unmarshaling the bytes should gives out the same Address", func() {
					var iden Identifier
					err := AminoCodec().UnmarshalBinary(bs, &iden)
					So(err, ShouldBeNil)
					So(iden, ShouldResemble, addr)
				})
			})
		})
		Convey("Given a BigInt", func() {
			v, _ := new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999999999", 10)
			n := BigInt{v}
			Convey("Marshaling the BigInt should succeed", func() {
				bs, err := AminoCodec().MarshalBinary(n)
				So(err, ShouldBeNil)
				Convey("Unmarshaling the bytes should gives out the same BigInt", func() {
					var n2 BigInt
					err := AminoCodec().UnmarshalBinary(bs, &n2)
					So(err, ShouldBeNil)
					So(n, ShouldResemble, n2)
				})
			})
		})
	})
}

func TestBigInt(t *testing.T) {
	Convey("In the beginning", t, func() {
		Convey("Given an int64", func() {
			i := int64(123)
			Convey("NewBigInt should return a BigInt representing that number", func() {
				n := NewBigInt(i)
				So(n.Int.Cmp(big.NewInt(123)), ShouldBeZeroValue)
			})
		})
		Convey("Given a valid string for positive number", func() {
			s := "123"
			Convey("NewBigIntFromString should return a BigInt representing that number", func() {
				n, ok := NewBigIntFromString(s)
				So(ok, ShouldBeTrue)
				So(n.Int.Cmp(big.NewInt(123)), ShouldBeZeroValue)
			})
		})
		Convey("Given a valid string for zero", func() {
			s := "0"
			Convey("NewBigIntFromString should return a BigInt representing zero", func() {
				n, ok := NewBigIntFromString(s)
				So(ok, ShouldBeTrue)
				So(n.Int.Cmp(big.NewInt(0)), ShouldBeZeroValue)
			})
		})
		Convey("Given a valid string for negative number", func() {
			s := "-123"
			Convey("NewBigIntFromString should fail", func() {
				_, ok := NewBigIntFromString(s)
				So(ok, ShouldBeFalse)
			})
		})
		Convey("Given an invalid string as number", func() {
			s := "fortytwo"
			Convey("NewBigIntFromString should fail", func() {
				_, ok := NewBigIntFromString(s)
				So(ok, ShouldBeFalse)
			})
		})
	})
}
