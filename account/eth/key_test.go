package eth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/tokentransfer/check"
	libcore "github.com/tokentransfer/interfaces/core"
)

type KeySuite struct{}

func Test_Key(t *testing.T) {
	s := Suite(&KeySuite{})
	TestingRun(t, s)
}

func (suite *KeySuite) TestRandom(c *C) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}
	k := &Key{}
	err = k.UnmarshalBinary(seed)
	if err != nil {
		panic(err)
	}
	a, err := k.GetAddress()
	if err != nil {
		panic(err)
	}
	address, err := a.GetAddress()
	if err != nil {
		panic(err)
	}
	p, err := k.GetPrivate()
	if err != nil {
		panic(err)
	}
	secret, err := p.GetSecret()
	if err != nil {
		panic(err)
	}
	fmt.Println("address", address)
	fmt.Println("secret", secret)
}

func (suite *KeySuite) TestGenerate(c *C) {
	//address 0x6da68a0c5dAAE0715AE6b62F00f548A2C6981c2f
	//secret beec9ec61c17b04cb9e4a9b7017e749f92835e2743e95f94cde218d667b14109
	k, err := GenerateFamilySeed("masterpassphrase")
	c.Assert(err, IsNil)
	a, err := k.GetAddress()
	c.Assert(err, IsNil)
	address, err := a.GetAddress()
	c.Assert(err, IsNil)
	p, err := k.GetPrivate()
	c.Assert(err, IsNil)
	secret, err := p.GetSecret()
	c.Assert(err, IsNil)

	// fmt.Println("address", address)
	// fmt.Println("secret", secret)
	c.Assert(address, Equals, "0x6da68a0c5dAAE0715AE6b62F00f548A2C6981c2f")
	c.Assert(secret, Equals, "beec9ec61c17b04cb9e4a9b7017e749f92835e2743e95f94cde218d667b14109")
}

func (suite *KeySuite) TestAccount(c *C) {
	//{"address":"0x42f32B004Da1093d51AE40a58F38E33BA4f46397","private":"4774628228852ee570d188f92cd10df3282bb5d895fc701733f43fca6bfb9852","public":"04d811caac49ba458fda498e5bc385bc9cc6e67aa6b19ba754c6cd75953ef06310e8607798ce5810a0b32fbd41fe8915de52fd511e7660038ff7067a0e94fc9481"}
	//address 0x42f32B004Da1093d51AE40a58F38E33BA4f46397
	//secret 4774628228852ee570d188f92cd10df3282bb5d895fc701733f43fca6bfb9852
	s := "4774628228852ee570d188f92cd10df3282bb5d895fc701733f43fca6bfb9852"
	k := &Key{}
	err := k.UnmarshalText([]byte(s))
	c.Assert(err, IsNil)
	a, err := k.GetAddress()
	c.Assert(err, IsNil)
	address, err := a.GetAddress()
	c.Assert(err, IsNil)
	p, err := k.GetPrivate()
	c.Assert(err, IsNil)
	secret, err := p.GetSecret()
	c.Assert(err, IsNil)

	// fmt.Println("address", address)
	// fmt.Println("secret", secret)
	c.Assert(address, Equals, "0x42f32B004Da1093d51AE40a58F38E33BA4f46397")
	c.Assert(secret, Equals, "4774628228852ee570d188f92cd10df3282bb5d895fc701733f43fca6bfb9852")
}

func (suite *KeySuite) TestSignAndVerify(c *C) {
	msgBlob := "65080112146DA68A0C5DAAE0715AE6B62F00F548A2C6981C2F1801280A321442F32B004DA1093D51AE40A58F38E33BA4F46397424104D33A199D322FAFD28867E3B9FB2F5AC081D56CFF1AE803635730E1D01B77D837D9EE578346DD88B68D21B9A61B8F1EFE9B2574F08B4A471F864FA7EA7A29185C4A413C76F5EBE03A949D0776192F50C2284A65ECED1EB2DCA1C37719A30C650598D57DA6E63598223A881545D612F9B20B1CAC507AA8EAC458A39A347F0CE8A9720B00"
	msg, err := hex.DecodeString(msgBlob)
	c.Assert(err, IsNil)

	h := make([]byte, 32)
	n, err := rand.Read(h)
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 32)

	k, err := GenerateFamilySeed("masterpassphrase")
	c.Assert(err, IsNil)

	fmt.Println("hash", libcore.Hash(h).String())
	sig, err := k.Sign(h, msg)
	c.Assert(err, IsNil)

	pk, err := k.GetPublic()
	c.Assert(err, IsNil)

	ok, err := pk.Verify(h, msg, sig)
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
}
