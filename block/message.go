package block

import (
	"github.com/tokentransfer/chain/core"
	"github.com/tokentransfer/chain/core/pb"

	libcore "github.com/tokentransfer/interfaces/core"
)

type Message struct {
	Hash libcore.Hash

	SigningPubKey *libcore.PublicKey
	TxnSignature  *libcore.Signature
	MessageKey    libcore.Bytes
}

func (tx *Message) GetHash() libcore.Hash {
	return tx.Hash
}

func (tx *Message) SetHash(h libcore.Hash) {
	tx.Hash = h
}

func (tx *Message) SetMessage(m libcore.Bytes) {
	tx.MessageKey = m
}

func (tx *Message) UnmarshalBinary(d []byte) error {
	_, msg, err := core.Unmarshal(d)
	if err != nil {
		return err
	}
	t := msg.(*pb.MessageKey)

	tx.MessageKey = t.GetMessageKey()

	pk := libcore.PublicKey(t.GetPublicKey())
	sig := libcore.Signature(t.GetSignature())

	tx.SigningPubKey = &pk
	tx.TxnSignature = &sig

	return nil
}

func (tx *Message) MarshalBinary() ([]byte, error) {
	t := &pb.MessageKey{
		MessageKey: tx.MessageKey,
	}
	if tx.SigningPubKey != nil {
		t.PublicKey = *tx.SigningPubKey
	}
	if tx.TxnSignature != nil {
		t.Signature = *tx.TxnSignature
	}
	data, err := core.Marshal(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (tx *Message) Raw(ignoreSigningFields bool) ([]byte, error) {
	if ignoreSigningFields {
		t := &pb.MessageKey{
			MessageKey: tx.MessageKey,
		}
		if tx.SigningPubKey != nil {
			t.PublicKey = *tx.SigningPubKey
		}
		data, err := core.Marshal(t)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	return tx.MarshalBinary()
}

func (tx *Message) GetAccount() libcore.Address {
	publicBytes := []byte(tx.GetPublicKey())

	_, p, err := as.NewPublicFromBytes(publicBytes)
	if err != nil {
		return nil
	}
	a, err := p.GenerateAddress()
	if err != nil {
		return nil
	}
	return a
}

func (tx *Message) GetPublicKey() libcore.PublicKey {
	if tx.SigningPubKey != nil {
		return *tx.SigningPubKey
	}
	return nil
}

func (tx *Message) SetPublicKey(p libcore.PublicKey) {
	if p != nil {
		tx.SigningPubKey = &p
	} else {
		tx.SigningPubKey = nil
	}
}

func (tx *Message) GetSignature() libcore.Signature {
	if tx.TxnSignature != nil {
		return *tx.TxnSignature
	}
	return nil
}

func (tx *Message) SetSignature(s libcore.Signature) {
	if s != nil {
		tx.TxnSignature = &s
	} else {
		tx.TxnSignature = nil
	}
}
