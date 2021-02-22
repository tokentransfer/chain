package block

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/skywell/skywell-go/util"
	"github.com/tokentransfer/chain/account"
	"github.com/tokentransfer/chain/core"
	"github.com/tokentransfer/chain/crypto"
	"github.com/tokentransfer/chain/node"

	libblock "github.com/tokentransfer/interfaces/block"
	. "gopkg.in/check.v1"
)

type BlockSuite struct{}

func Test_Block(t *testing.T) {
	s := Suite(&BlockSuite{})
	TestingRun(t, s)
	// TestingT(t)
}

func (s *BlockSuite) TestBlob(c *C) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	blob := "64080A12209BB0595780A120423974A6EE9658ACCEA0833D3F04F4D7A9BE1FDFD2EE4C01CB1A202E91BE87194D808F9A7E2EEFE1968A6E55C6F385A593AD38D7EFCDB5057AA2D02220091212A515AA5DB777930FD7EE3C5B3AF189EF666B345301E920E759BD8D0F7328C6CAC2F7D0ABF1B11632FF010A9D0108011214B5F762798A53D543A014CAF8B297CFF8F2F937E81801280A3214334DB1FE97E118B2B225C33A55D6185FAE431F3D42210330E7FC9D56BB25D6893BA3F317AE5BCF33B3291BD63DB32654A313222F7FD0204A4630440220774DE24E4E1EC75C10AAF859D4DF457CA6C30A4B4A598F2692C12106D98EE245022019C24199E3F6DEF99E80F5938A60522028A65A700DB8C8ED7600DD5ED5917E5D125D1A1B6F086F1A14334DB1FE97E118B2B225C33A55D6185FAE431F3D280A1A236F086F1A14B5F762798A53D543A014CAF8B297CFF8F2F937E82001288080E983B1DE161A196F086F1A14334DB1FE97E118B2B225C33A55D6185FAE431F3D3A5D1A1B6F086F1A14334DB1FE97E118B2B225C33A55D6185FAE431F3D280A1A236F086F1A14B5F762798A53D543A014CAF8B297CFF8F2F937E82001288080E983B1DE161A196F086F1A14334DB1FE97E118B2B225C33A55D6185FAE431F3D"
	data, err := hex.DecodeString(blob)
	c.Assert(err, IsNil)

	b := &Block{}
	err = b.UnmarshalBinary(data)
	c.Assert(err, IsNil)
	// util.PrintJSON(">> block", b)
}

func generateTransaction(seq uint64, value int64, gas int64) *Transaction {
	fromKey, err := account.GenerateFamilySeed("masterpassphrase")
	if err != nil {
		panic(err)
	}
	from, err := fromKey.GetAddress()
	if err != nil {
		panic(err)
	}

	to := account.NewAddress()
	err = to.UnmarshalText([]byte("0x42f32B004Da1093d51AE40a58F38E33BA4f46397"))
	if err != nil {
		panic(err)
	}

	tx := &Transaction{
		TransactionType: libblock.TransactionType(1),

		Account:     from,
		Sequence:    seq,
		Amount:      value,
		Gas:         gas,
		Destination: to,
		Payload:     []byte{1, 2, 3, 4},
	}

	service := &crypto.CryptoService{}
	service.Sign(fromKey, tx)

	return tx
}

func (suite *BlockSuite) TestBlock(c *C) {
	cryptoService := &crypto.CryptoService{}

	config, err := core.NewConfig("localhost", 7001)
	if err != nil {
		panic(err)
	}

	merkleService := &node.MerkleService{
		CryptoService: cryptoService,
	}
	err = merkleService.Init(config)
	if err != nil {
		panic(err)
	}
	consensusService := &node.ConsensusService{
		CryptoService: cryptoService,
		MerkleService: merkleService,
	}

	genesis, err := consensusService.GenerateBlock(nil)
	if err != nil {
		panic(err)
	}
	util.PrintJSON("genesis", genesis)

	_, err = consensusService.VerifyBlock(genesis)
	if err != nil {
		panic(err)
	}
	err = consensusService.AddBlock(genesis)
	if err != nil {
		panic(err)
	}

	tx1 := generateTransaction(1, 100, 10)
	ok, err := merkleService.VerifyTransaction(tx1)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic(errors.New("error transaction"))
	}
	tx1WithData, err := merkleService.ProcessTransaction(tx1)
	if err != nil {
		panic(err)
	}
	err = merkleService.PutTransaction(tx1WithData)
	if err != nil {
		panic(err)
	}

	tx2 := generateTransaction(2, 200, 10)
	ok, err = merkleService.VerifyTransaction(tx2)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic(errors.New("error transaction"))
	}
	tx2WithData, err := merkleService.ProcessTransaction(tx2)
	if err != nil {
		panic(err)
	}
	err = merkleService.PutTransaction(tx2WithData)
	if err != nil {
		panic(err)
	}

	fmt.Println("transaction root", merkleService.GetTransactionRoot())
	fmt.Println("receipt root", merkleService.GetReceiptRoot())

	err = merkleService.Cancel()
	if err != nil {
		panic(err)
	}
	fmt.Println("transaction root", merkleService.GetTransactionRoot())
	fmt.Println("receipt root", merkleService.GetReceiptRoot())

	block1, err := consensusService.GenerateBlock([]libblock.TransactionWithData{
		tx1WithData,
		tx2WithData,
	})
	if err != nil {
		panic(err)
	}
	util.PrintJSON("block 1", block1)

	_, err = consensusService.VerifyBlock(block1)
	if err != nil {
		panic(err)
	}
	err = consensusService.AddBlock(block1)
	if err != nil {
		panic(err)
	}
}
