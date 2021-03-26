package jingtum

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/btcsuite/btcd/btcec"

	libaccount "github.com/tokentransfer/interfaces/account"
	libcore "github.com/tokentransfer/interfaces/core"
)

const (
	ADDRESS_PREFIX uint8 = 0
	SECRET_PREFIX  uint8 = 33
)

var (
	ORDER = btcec.S256().N
	ZERO  = big.NewInt(0)
	ONE   = big.NewInt(1)
)

func GenerateFamilySeed(password string) (*Key, error) {
	b := Sha512Quarter([]byte(password))
	k := &Key{}
	err := k.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return k, nil
}

type Key struct {
	privateKey Private
	publicKey  *Public
	address    *Address
}

func (p *Key) UnmarshalBinary(data []byte) error {
	return p.privateKey.UnmarshalBinary(data)
}

func (p *Key) MarshalBinary() ([]byte, error) {
	return p.privateKey.MarshalBinary()
}

func (p *Key) UnmarshalText(b []byte) error {
	return p.privateKey.UnmarshalText(b)
}

func (p *Key) MarshalText() ([]byte, error) {
	return p.privateKey.MarshalText()
}

func (p *Key) Sign(hash libcore.Hash, msg []byte) (libcore.Signature, error) {
	return p.privateKey.Sign(hash, msg)
}

func (p *Key) Verify(hash libcore.Hash, msg []byte, signature libcore.Signature) (bool, error) {
	publicKey, err := p.GetPublic()
	if err != nil {
		return false, err
	}
	return publicKey.Verify(hash, msg, signature)
}

func (p *Key) GetPrivate() (libaccount.PrivateKey, error) {
	return &p.privateKey, nil
}

func (p *Key) GetPublic() (libaccount.PublicKey, error) {
	if p.publicKey == nil {
		pk, err := p.privateKey.GeneratePublic()
		if err != nil {
			return nil, err
		}
		p.publicKey = pk.(*Public)
	}
	return p.publicKey, nil
}

func (p *Key) GetAddress() (libcore.Address, error) {
	_, err := p.GetPublic()
	if err != nil {
		return nil, err
	}
	if p.address == nil {
		a, err := p.publicKey.GenerateAddress()
		if err != nil {
			return nil, err
		}
		p.address = a.(*Address)
	}
	return p.address, nil
}

type Private struct {
	seed []byte
	*btcec.PrivateKey
}

func newKey(data []byte) *btcec.PrivateKey {
	inc := big.NewInt(0).SetBytes(data)
	inc.Lsh(inc, 32)
	for k := big.NewInt(0); ; inc.Add(inc, ONE) {
		k.SetBytes(Sha512Half(inc.Bytes()))
		if k.Cmp(ZERO) > 0 && k.Cmp(ORDER) < 0 {
			privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), k.Bytes())
			return privKey
		}
	}
}

func (p *Private) UnmarshalBinary(data []byte) error {
	l := len(data)
	if l != 16 {
		return errors.New("16 bytes required")
	}
	p.seed = data
	p.PrivateKey = newKey(data)
	return nil
}

func (p *Private) MarshalBinary() ([]byte, error) {
	priv := p.generateKey(uint32(0))
	return priv.D.Bytes(), nil
}

func (p *Private) UnmarshalText(b []byte) error {
	ret, err := Base58Decode(string(b))
	if err != nil {
		return err
	}
	seed := ret[1:17]
	return p.UnmarshalBinary(seed)
}

func (p *Private) MarshalText() ([]byte, error) {
	data := p.seed
	buffer := bufCat0(SECRET_PREFIX, data)

	checksum := Sha256(Sha256(buffer))[0:4]
	ret := bufCat1(buffer, checksum)
	s := Base58Encode(ret)
	return []byte(s), nil
}

func (p *Private) Sign(hash libcore.Hash, msg []byte) (libcore.Signature, error) {
	hashBytes := []byte(hash)
	priv := p.generateKey(uint32(0))
	sig, err := priv.Sign(hashBytes)
	if err != nil {
		return nil, err
	}
	return libcore.Signature(sig.Serialize()), nil
}

func (p *Private) generateKey(sequence uint32) *btcec.PrivateKey {
	seed := make([]byte, btcec.PubKeyBytesLenCompressed+4)
	copy(seed, p.PubKey().SerializeCompressed())
	binary.BigEndian.PutUint32(seed[btcec.PubKeyBytesLenCompressed:], sequence)
	key := newKey(seed)
	key.D.Add(key.D, p.D).Mod(key.D, ORDER)
	key.X, key.Y = key.ScalarBaseMult(key.D.Bytes())
	return key
}

func (p *Private) GeneratePublic() (libaccount.PublicKey, error) {
	priv := p.generateKey(uint32(0))
	return &Public{priv.PubKey()}, nil
}

func (p *Private) GetSecret() (string, error) {
	b, err := p.MarshalText()
	return string(b), err
}

type Public struct {
	*btcec.PublicKey
}

func (p *Public) UnmarshalBinary(data []byte) error {
	pk, err := btcec.ParsePubKey(data, btcec.S256())
	if err != nil {
		return err
	}
	p.PublicKey = pk
	return nil
}

func (p *Public) MarshalBinary() ([]byte, error) {
	return p.SerializeCompressed(), nil
}

func (p *Public) UnmarshalText(b []byte) error {
	data, err := hex.DecodeString(string(b))
	if err != nil {
		return err
	}
	return p.UnmarshalBinary(data)
}

func (p *Public) MarshalText() ([]byte, error) {
	data, err := p.MarshalBinary()
	if err != nil {
		return nil, err
	}
	s := hex.EncodeToString(data)
	return []byte(s), nil
}

func (p *Public) Verify(hash libcore.Hash, msg []byte, signature libcore.Signature) (bool, error) {
	signatureBytes := []byte(signature)
	sig, err := btcec.ParseDERSignature(signatureBytes, btcec.S256())
	if err != nil {
		return false, err
	}
	hashBytes := []byte(hash)
	return sig.Verify(hashBytes, p.PublicKey), nil
}

func (p *Public) GenerateAddress() (libcore.Address, error) {
	data, err := p.MarshalBinary()
	if err != nil {
		return nil, err
	}
	seed := Sha256RipeMD160(data)

	ra := Address{}
	err = ra.UnmarshalBinary(seed)
	if err != nil {
		return nil, err
	}
	return &ra, nil
}

type Address [20]byte

func (a *Address) UnmarshalText(b []byte) error {
	ret, err := Base58Decode(string(b))
	if err != nil {
		return err
	}
	copy(a[:], ret[1:21])
	return nil
}

func (a Address) MarshalText() ([]byte, error) {
	seed := a[:]
	buffer := bufCat0(ADDRESS_PREFIX, seed)
	checksum := Sha256(Sha256(buffer))[0:4]
	ret := bufCat1(buffer, checksum)

	s := Base58Encode(ret)
	return []byte(s), nil
}

func (a *Address) UnmarshalBinary(data []byte) error {
	copy(a[:], data)
	return nil
}

func (a Address) MarshalBinary() ([]byte, error) {
	return a[:], nil
}

func (a Address) String() string {
	b, _ := a.MarshalText()
	return string(b)
}
