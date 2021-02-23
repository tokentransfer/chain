package jingtum

import (
	"fmt"
	"math/big"
	"testing"

	. "github.com/tokentransfer/check"
)

type JingtumSuite struct{}

func Test_Jingtum(t *testing.T) {
	s := Suite(&JingtumSuite{})
	TestingRun(t, s)
}

func (suite *JingtumSuite) TestAccount(c *C) {
	// s := "snoPBjXtMeMyMHUVTgbuqAfg1SUTb" //jHb9CJAWyB4jr91VRWn96DkukG4bwdtyTh
	s := "snxT8kKatzLzLZHAdYcwbVQJc79vN" //jngGY9W1F2Ky8wCzeHTahbtxjadU9wNFRz
	k := &Key{}
	err := k.UnmarshalText([]byte(s))

	// k, err := GenerateFamilySeed("masterpassphrase")
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

	data, err := a.MarshalBinary()
	if err != nil {
		panic(err)
	}
	i := new(big.Int)
	i.SetBytes(data)
	fmt.Println(address, i.Uint64())
}
