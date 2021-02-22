package crypto

import (
	"crypto/sha256"
	"errors"

	"github.com/tokentransfer/chain/account"

	libaccount "github.com/tokentransfer/interfaces/account"
	libcore "github.com/tokentransfer/interfaces/core"
	libcrypto "github.com/tokentransfer/interfaces/crypto"
)

type CryptoService struct {
}

func (service *CryptoService) GetSize() int {
	return 32
}

func (service *CryptoService) Hash(msg []byte) (libcore.Hash, error) {
	h := sha256.New()
	h.Write(msg)
	b := h.Sum(nil)
	return libcore.Hash(b), nil
}

func (service *CryptoService) Raw(h libcrypto.Hashable) (libcore.Hash, []byte, error) {
	data, err := h.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	hash, err := service.Hash(data)
	if err != nil {
		return nil, nil, err
	}
	h.SetHash(hash)
	return hash, data, nil
}

func (service *CryptoService) Sign(p libaccount.Key, s libcrypto.Signable) error {
	publicKey, err := p.GetPublic()
	if err != nil {
		return err
	}
	publicBytes, err := publicKey.MarshalBinary()
	if err != nil {
		return err
	}
	s.SetPublicKey(libcore.PublicKey(publicBytes))

	data, err := s.Raw()
	if err != nil {
		return err
	}
	hash, err := service.Hash(data)
	if err != nil {
		return err
	}
	signature, err := p.Sign(hash, data)
	if err != nil {
		return err
	}
	s.SetSignature(signature)

	h, _, err := service.Raw(s)
	if err != nil {
		return err
	}
	s.SetHash(h)

	return nil
}

func (service *CryptoService) Verify(s libcrypto.Signable) (bool, error) {
	publicBytes := []byte(s.GetPublicKey())

	p := account.NewPublicKey()
	err := p.UnmarshalBinary(publicBytes)
	if err != nil {
		return false, err
	}
	a, err := p.GenerateAddress()
	if err != nil {
		return false, err
	}
	if !libcore.Equals(a, s.GetAccount()) {
		return false, errors.New("error signature")
	}
	data, err := s.Raw()
	if err != nil {
		return false, err
	}
	hash, err := service.Hash(data)
	if err != nil {
		return false, err
	}
	signature := s.GetSignature()
	return p.Verify(hash, data, signature)
}
