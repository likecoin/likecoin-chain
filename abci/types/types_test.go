package types

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
			id := IDStr("MzMzMzMzMzMzMzMzMzMzMzMzMzM=")
			Convey("It should have correct DB key", func() {
				dbKey := id.DBKey("prefix", "suffix")
				expectedKey := []byte("prefix:id:_\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33_suffix")
				So(dbKey, ShouldResemble, expectedKey)
			})
		})
		Convey("For JSON marshaling and unmarshaling", func() {
			Convey("Given a valid LikeChainID", func() {
				id := IDStr("MzMzMzMzMzMzMzMzMzMzMzMzMzM=")
				Convey("It should be marshaled to JSON correctly", func() {
					bs, err := json.Marshal(&id)
					So(err, ShouldBeNil)
					So(bs, ShouldResemble, []byte(`"MzMzMzMzMzMzMzMzMzMzMzMzMzM="`))
					Convey("JSON unmarshaling should recover the same LikeChainID", func() {
						recoveredID := LikeChainID{}
						err = json.Unmarshal(bs, &recoveredID)
						So(err, ShouldBeNil)
						So(&recoveredID, ShouldResemble, id)
					})
				})
			})
			Convey("Given a valid string with invalid LikeChainID", func() {
				s := `"MzMzMzMzMzMzMzMzMzMzMzMzMzM"`
				Convey("JSON unmarshaling should return error", func() {
					recoveredID := LikeChainID{}
					err := json.Unmarshal([]byte(s), &recoveredID)
					So(err, ShouldNotBeNil)
				})
			})
			Convey("Given a non string", func() {
				s := `333333333333333333333333333`
				Convey("JSON unmarshaling should return error", func() {
					recoveredID := LikeChainID{}
					err := json.Unmarshal([]byte(s), &recoveredID)
					So(err, ShouldNotBeNil)
				})
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
			addr := Addr("0x3333333333333333333333333333333333333333")
			Convey("It should have correct DB key", func() {
				dbKey := addr.DBKey("prefix", "suffix")
				expectedKey := []byte("prefix:addr:_\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33_suffix")
				So(dbKey, ShouldResemble, expectedKey)
			})
			Convey("It should be marshaled to JSON correctly", func() {
				bs, err := json.Marshal(&addr)
				So(err, ShouldBeNil)
				So(bs, ShouldResemble, []byte(`"0x3333333333333333333333333333333333333333"`))
				Convey("JSON unmarshaling should recover the same Address", func() {
					recoveredAddr := Address{}
					err = json.Unmarshal(bs, &recoveredAddr)
					So(err, ShouldBeNil)
					So(&recoveredAddr, ShouldResemble, addr)
				})
			})
		})
		Convey("For JSON marshaling and unmarshaling", func() {
			Convey("Given a valid Address", func() {
				addr := Addr("0x3333333333333333333333333333333333333333")
				Convey("It should be marshaled to JSON correctly", func() {
					bs, err := json.Marshal(&addr)
					So(err, ShouldBeNil)
					So(bs, ShouldResemble, []byte(`"0x3333333333333333333333333333333333333333"`))
					Convey("JSON unmarshaling should recover the same Address", func() {
						recoveredAddr := Address{}
						err = json.Unmarshal(bs, &recoveredAddr)
						So(err, ShouldBeNil)
						So(&recoveredAddr, ShouldResemble, addr)
					})
				})
			})
			Convey("Given a valid string with invalid address", func() {
				s := `"0x333333333333333333333333333333333333333g"`
				Convey("JSON unmarshaling should return error", func() {
					recoveredAddr := Address{}
					err := json.Unmarshal([]byte(s), &recoveredAddr)
					So(err, ShouldNotBeNil)
				})
			})
			Convey("Given a non string", func() {
				s := `3333333333333333333333333333333333333333`
				Convey("JSON unmarshaling should return error", func() {
					recoveredAddr := Address{}
					err := json.Unmarshal([]byte(s), &recoveredAddr)
					So(err, ShouldNotBeNil)
				})
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
				bs, err := AminoCodec().MarshalBinaryLengthPrefixed(id)
				So(err, ShouldBeNil)
				Convey("Unmarshaling the bytes should gives out the same LikeChainID", func() {
					var iden Identifier
					err := AminoCodec().UnmarshalBinaryLengthPrefixed(bs, &iden)
					So(err, ShouldBeNil)
					So(iden, ShouldResemble, id)
				})
			})
		})
		Convey("Given an Address", func() {
			addr := Addr("0x0000000000000000000000000000000000000000")
			Convey("Marshaling the Address should succeed", func() {
				bs, err := AminoCodec().MarshalBinaryLengthPrefixed(addr)
				So(err, ShouldBeNil)
				Convey("Unmarshaling the bytes should gives out the same Address", func() {
					var iden Identifier
					err := AminoCodec().UnmarshalBinaryLengthPrefixed(bs, &iden)
					So(err, ShouldBeNil)
					So(iden, ShouldResemble, addr)
				})
			})
		})
		Convey("Given a BigInt", func() {
			v, _ := new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999999999", 10)
			n := BigInt{v}
			Convey("Marshaling the BigInt should succeed", func() {
				bs, err := AminoCodec().MarshalBinaryLengthPrefixed(n)
				So(err, ShouldBeNil)
				Convey("Unmarshaling the bytes should gives out the same BigInt", func() {
					var n2 BigInt
					err := AminoCodec().UnmarshalBinaryLengthPrefixed(bs, &n2)
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
		Convey("For JSON marshaling and unmarshaling", func() {
			Convey("Given a BigInt", func() {
				n, _ := NewBigIntFromString("999999999999999999999999999999999999999999999999999999999999")
				Convey("It should be marshaled to JSON correctly", func() {
					bs, err := json.Marshal(&n)
					So(err, ShouldBeNil)
					So(bs, ShouldResemble, []byte(`"999999999999999999999999999999999999999999999999999999999999"`))
					Convey("JSON unmarshaling should recover the same BigInt", func() {
						recoveredN := BigInt{}
						err = json.Unmarshal(bs, &recoveredN)
						So(err, ShouldBeNil)
						So(recoveredN, ShouldResemble, n)
					})
				})
			})
			Convey("Given a valid string of numbers without quotation", func() {
				s := `333333333333333333333333333333333333333`
				Convey("JSON unmarshaling should succeed", func() {
					n := BigInt{}
					err := json.Unmarshal([]byte(s), &n)
					So(err, ShouldBeNil)
					v, _ := new(big.Int).SetString("333333333333333333333333333333333333333", 10)
					So(n, ShouldResemble, BigInt{v})
				})
			})
			Convey("Given an invalid string of number", func() {
				s := `"333333333333333333333333333333333333333x"`
				Convey("JSON unmarshaling should return error", func() {
					n := BigInt{}
					err := json.Unmarshal([]byte(s), &n)
					So(err, ShouldNotBeNil)
				})
			})
		})
		Convey("For range", func() {
			Convey("Given a BigInt with value -1", func() {
				n := BigInt{big.NewInt(-1)}
				Convey("n.IsWithinRange() should return false", func() {
					So(n.IsWithinRange(), ShouldBeFalse)
				})
			})
			Convey("Given a BigInt with value 0", func() {
				n := BigInt{big.NewInt(0)}
				Convey("n.IsWithinRange() should return true", func() {
					So(n.IsWithinRange(), ShouldBeTrue)
				})
			})
			Convey("Given a BigInt with value 2^256-1", func() {
				v := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
				v.Sub(v, big.NewInt(1))
				n := BigInt{v}
				Convey("n.IsWithinRange() should return true", func() {
					So(n.IsWithinRange(), ShouldBeTrue)
				})
			})
			Convey("Given a BigInt with value 2^256", func() {
				v := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
				n := BigInt{v}
				Convey("n.IsWithinRange() should return false", func() {
					So(n.IsWithinRange(), ShouldBeFalse)
				})
			})
		})
	})
}
