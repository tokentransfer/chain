package jingtum

import (
	"crypto/sha256"
	"crypto/sha512"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

var (
	ALPHABET = base58.NewAlphabet("jpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65rkm8oFqi1tuvAxyz")
)

func Sha256(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}

func Sha512Half(b []byte) []byte {
	h := sha512.New()
	h.Write(b)
	return h.Sum(nil)[:32]
}

func Sha512Quarter(b []byte) []byte {
	h := sha512.New()
	h.Write(b)
	return h.Sum(nil)[:16]
}

func bufCat0(item1 uint8, buf2 []byte) []byte {
	return append([]byte{item1}, buf2...)
}

func bufCat1(buf1 []byte, buf2 []byte) []byte {
	var buf []byte
	buf = append(buf, buf1...)
	buf = append(buf, buf2...)
	return buf
}

func Sha256RipeMD160(data []byte) []byte {
	sha := sha256.New()
	sha.Write(data)

	ripe := ripemd160.New()
	ripe.Write(sha.Sum(nil))
	return ripe.Sum(nil)
}

func Base58Encode(data []byte) string {
	return base58.EncodeAlphabet(data, ALPHABET)
}

func Base58Decode(s string) ([]byte, error) {
	return base58.DecodeAlphabet(s, ALPHABET)
}
