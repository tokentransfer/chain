package account

import (
	"github.com/tokentransfer/chain/account/eth"

	libaccount "github.com/tokentransfer/interfaces/account"
	libcore "github.com/tokentransfer/interfaces/core"
)

const (
	ETH     libaccount.KeyType = 0
	JINGTUM libaccount.KeyType = 1
)

func GenerateFamilySeed(password string) (libaccount.KeyType, libaccount.Key, error) {
	key, err := eth.GenerateFamilySeed(password)
	if err != nil {
		return 0, nil, err
	}
	return ETH, key, nil
}

func NewKeyFromSecret(secret string) (libaccount.KeyType, libaccount.Key, error) {
	key := &eth.Key{}
	err := key.UnmarshalText([]byte(secret))
	if err != nil {
		return 0, nil, err
	}
	return ETH, key, nil
}

func NewKeyFromBytes(data []byte) (libaccount.KeyType, libaccount.Key, error) {
	key := &eth.Key{}
	err := key.UnmarshalBinary(data)
	if err != nil {
		return 0, nil, err
	}
	return ETH, key, nil
}

func NewPublicFromHex(s string) (libaccount.KeyType, libaccount.PublicKey, error) {
	key := &eth.Public{}
	err := key.UnmarshalText([]byte(s))
	if err != nil {
		return 0, nil, err
	}
	return ETH, key, nil
}

func NewPublicFromBytes(data []byte) (libaccount.KeyType, libaccount.PublicKey, error) {
	key := &eth.Public{}
	err := key.UnmarshalBinary(data)
	if err != nil {
		return 0, nil, err
	}
	return ETH, key, nil
}

func NewAccountFromAddress(address string) (libaccount.KeyType, libcore.Address, error) {
	a := &eth.Address{}
	err := a.UnmarshalText([]byte(address))
	if err != nil {
		return 0, nil, err
	}
	return ETH, a, nil
}

func NewAccountFromBytes(data []byte) (libaccount.KeyType, libcore.Address, error) {
	a := &eth.Address{}
	err := a.UnmarshalBinary(data)
	if err != nil {
		return 0, nil, err
	}
	return ETH, a, nil
}
