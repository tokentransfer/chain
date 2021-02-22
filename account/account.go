package account

import (
	"github.com/tokentransfer/chain/account/eth"

	libaccount "github.com/tokentransfer/interfaces/account"
	libcore "github.com/tokentransfer/interfaces/core"
)

func GenerateFamilySeed(password string) (libaccount.Key, error) {
	return eth.GenerateFamilySeed(password)
}

func NewKey() libaccount.Key {
	return &eth.Key{}
}

func NewPublicKey() libaccount.PublicKey {
	return &eth.Public{}
}

func NewAddress() libcore.Address {
	return &eth.Address{}
}
