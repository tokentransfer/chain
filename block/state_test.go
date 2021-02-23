package block

import (
	"fmt"
	"testing"

	"github.com/tokentransfer/chain/account"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
	. "gopkg.in/check.v1"
)

type StateSuite struct{}

func Test_State(t *testing.T) {
	s := Suite(&StateSuite{})
	TestingRun(t, s)
	// TestingT(t)
}

func (suite *StateSuite) TestState(c *C) {
	a := account.NewAddress()
	err := a.UnmarshalText([]byte("jngGY9W1F2Ky8wCzeHTahbtxjadU9wNFRz"))
	if err != nil {
		panic(err)
	}

	s1 := &AccountState{
		State: State{
			StateType:  1,
			BlockIndex: uint64(2),
		},
		Account:  a,
		Sequence: uint64(1),
		Amount:   int64(100),
	}
	// util.PrintJSON("state", s1)

	s2 := &AccountState{
		State: State{
			StateType:  1,
			BlockIndex: uint64(2),
		},
		Account:  a,
		Sequence: uint64(1),
		Amount:   int64(100),
	}
	// util.PrintJSON("state", s2)

	state := &Receipt{
		States: []libblock.State{s1, s2},
	}
	// util.PrintJSON("state", state)

	sbs, err := state.MarshalBinary()
	if err != nil {
		panic(err)
	}
	s := &Receipt{}
	err = s.UnmarshalBinary(sbs)
	if err != nil {
		panic(err)
	}
	// util.PrintJSON("state", s)
	fmt.Println("hex", libcore.EncodeToString(sbs))
}
