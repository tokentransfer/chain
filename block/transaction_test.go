package block

import (
	"encoding/hex"
	"testing"

	. "github.com/tokentransfer/check"
)

type TransactionSuite struct{}

func Test_Transaction(t *testing.T) {
	s := Suite(&TransactionSuite{})
	TestingRun(t, s)
}

func (suite *TransactionSuite) TestWithDataBlob(c *C) {
	blob := "670A9702080112146DA68A0C5DAAE0715AE6B62F00F548A2C6981C2F180120904E280A321442F32B004DA1093D51AE40A58F38E33BA4F46397424104D33A199D322FAFD28867E3B9FB2F5AC081D56CFF1AE803635730E1D01B77D837D9EE578346DD88B68D21B9A61B8F1EFE9B2574F08B4A471F864FA7EA7A29185C4A9C013130393733323032323333343035303739353231303632353935303137343531323933373934333732353130353532333435343735383835303536303839333834383133363831393933303637393A333436393933383739343134323339373435313139383730393837353930323037343133303634363938373830303430383636343430323536333539393136323135363432333939343831333912601A1B6F086F1A146AB60933B69A1CE3B278CBB58BCC4ADDC637885C280A1A236F086F1A146DA68A0C5DAAE0715AE6B62F00F548A2C6981C2F200128F0B1E883B1DE161A1C6F086F1A1442F32B004DA1093D51AE40A58F38E33BA4F4639728904E"

	data, err := hex.DecodeString(blob)
	c.Assert(err, IsNil)
	txWithData := &TransactionWithData{}
	err = txWithData.UnmarshalBinary(data)
	c.Assert(err, IsNil)

	c.Assert(txWithData, NotNil)
	// util.PrintJSON(">> txWithData", txWithData)
}

func (suite *TransactionSuite) TestTransactionBlob(c *C) {
	blob := "65080112146DA68A0C5DAAE0715AE6B62F00F548A2C6981C2F1801280A321442F32B004DA1093D51AE40A58F38E33BA4F46397424104D33A199D322FAFD28867E3B9FB2F5AC081D56CFF1AE803635730E1D01B77D837D9EE578346DD88B68D21B9A61B8F1EFE9B2574F08B4A471F864FA7EA7A29185C4A413C76F5EBE03A949D0776192F50C2284A65ECED1EB2DCA1C37719A30C650598D57DA6E63598223A881545D612F9B20B1CAC507AA8EAC458A39A347F0CE8A9720B00"

	data, err := hex.DecodeString(blob)
	c.Assert(err, IsNil)
	tx := &Transaction{}
	err = tx.UnmarshalBinary(data)
	c.Assert(err, IsNil)

	// util.PrintJSON(">> tx", tx)
	c.Assert(tx, NotNil)
	c.Assert(tx.Sequence, Equals, uint64(1))
}

func (suite *TransactionSuite) TestTransaction(c *C) {
	// fromKey, err := account.GenerateFamilySeed("masterpassphrase")
	// if err != nil {
	// 	panic(err)
	// }
	// fromPublic, err := fromKey.GetPublic()
	// if err != nil {
	// 	panic(err)
	// }
	// fromPublicKey, err := fromPublic.MarshalBinary()
	// if err != nil {
	// 	panic(err)
	// }
	// from, err := fromKey.GetAddress()
	// if err != nil {
	// 	panic(err)
	// }

	// to := account.NewAddress()
	// err = to.UnmarshalText([]byte("0x42f32B004Da1093d51AE40a58F38E33BA4f46397"))
	// if err != nil {
	// 	panic(err)
	// }
	// util.PrintJSON("from", from)
	// util.PrintJSON("to", to)

	// tx := &Transaction{
	// 	TransactionType: libblock.TransactionType(1),

	// 	Account:     from,
	// 	Sequence:    uint64(1),
	// 	Amount:      int64(10000),
	// 	Gas:         int64(10),
	// 	Destination: to,
	// 	Payload:     []byte{1, 2, 3, 4},
	// 	PublicKey:   libcore.PublicKey(fromPublicKey),
	// }

	// service := &crypto.CryptoService{}
	// service.Sign(fromKey, tx)

	// util.PrintJSON("signed tx", tx)

	// data, err := tx.MarshalBinary()
	// if err != nil {
	// 	panic(err)
	// }
	// hash := tx.GetHash()
	// util.PrintJSON(fmt.Sprintf("%s:%d", hash.String(), len(data)), libcore.Bytes(data))

	// newTx := &Transaction{}
	// err = newTx.UnmarshalBinary(data)
	// if err != nil {
	// 	panic(err)
	// }

	// util.PrintJSON("verify tx", newTx)
	// ok, err := service.Verify(newTx)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("verify", ok)

	// h, data, err := service.Raw(newTx)
	// if err != nil {
	// 	panic(err)
	// }
	// util.PrintJSON("new tx", newTx)
	// util.PrintJSON(fmt.Sprintf("%s:%d", h.String(), len(data)), libcore.Bytes(data))
}
