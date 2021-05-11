package account

import (
	"github.com/tokentransfer/chain/account/jingtum"

	libaccount "github.com/tokentransfer/interfaces/account"
	libcore "github.com/tokentransfer/interfaces/core"
)

const (
	// ETH     libaccount.KeyType = 0
	JINGTUM libaccount.KeyType = 1
)

func init() {
	// KeyType
	// ETH.Register("ETH", func(t libaccount.KeyType) libaccount.Key {
	// 	key, _ := eth.GenerateFamilySeed("masterpassphrase")
	// 	return key
	// })
	JINGTUM.Register("Jingtum", func(t libaccount.KeyType) libaccount.Key {
		key, _ := jingtum.GenerateFamilySeed("masterpassphrase")
		return key
	})
}

type AccountService struct {
}

func (service *AccountService) GenerateFamilySeed(password string) (libaccount.KeyType, libaccount.Key, error) {
	key, err := jingtum.GenerateFamilySeed(password)
	if err != nil {
		return 0, nil, err
	}
	return JINGTUM, key, nil
}

func (service *AccountService) NewKeyFromSecret(secret string) (libaccount.KeyType, libaccount.Key, error) {
	key := &jingtum.Key{}
	err := key.UnmarshalText([]byte(secret))
	if err != nil {
		return 0, nil, err
	}
	return JINGTUM, key, nil
}

func (service *AccountService) NewKeyFromBytes(data []byte) (libaccount.KeyType, libaccount.Key, error) {
	key := &jingtum.Key{}
	err := key.UnmarshalBinary(data)
	if err != nil {
		return 0, nil, err
	}
	return JINGTUM, key, nil
}

func (service *AccountService) NewPublicFromHex(s string) (libaccount.KeyType, libaccount.PublicKey, error) {
	key := &jingtum.Public{}
	err := key.UnmarshalText([]byte(s))
	if err != nil {
		return 0, nil, err
	}
	return JINGTUM, key, nil
}

func (service *AccountService) NewPublicFromBytes(data []byte) (libaccount.KeyType, libaccount.PublicKey, error) {
	key := &jingtum.Public{}
	err := key.UnmarshalBinary(data)
	if err != nil {
		return 0, nil, err
	}
	return JINGTUM, key, nil
}

func (service *AccountService) NewAccountFromAddress(address string) (libaccount.KeyType, libcore.Address, error) {
	a := &jingtum.Address{}
	err := a.UnmarshalText([]byte(address))
	if err != nil {
		return 0, nil, err
	}
	return JINGTUM, a, nil
}

func (service *AccountService) NewAccountFromBytes(data []byte) (libaccount.KeyType, libcore.Address, error) {
	a := &jingtum.Address{}
	err := a.UnmarshalBinary(data)
	if err != nil {
		return 0, nil, err
	}
	return JINGTUM, a, nil
}
