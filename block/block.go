package block

import (
	"errors"
	"log"

	"github.com/tokentransfer/chain/core"
	"github.com/tokentransfer/chain/core/pb"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
)

type Block struct {
	Hash libcore.Hash

	BlockIndex      uint64
	ParentHash      libcore.Hash
	TransactionHash libcore.Hash
	ReceiptHash     libcore.Hash
	Timestamp       int64

	Transactions []libblock.TransactionWithData
	Receipts     []libblock.Receipt
}

func (b *Block) GetIndex() uint64 {
	return b.BlockIndex
}

func (b *Block) GetHash() libcore.Hash {
	return b.Hash
}

func (b *Block) SetHash(h libcore.Hash) {
	b.Hash = h
}

func (b *Block) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}

	if meta != core.CORE_BLOCK {
		return errors.New("error block data")
	}
	block := msg.(*pb.Block)
	b.BlockIndex = block.BlockIndex
	b.ParentHash = libcore.Hash(block.ParentHash)
	b.TransactionHash = libcore.Hash(block.TransactionHash)
	b.ReceiptHash = libcore.Hash(block.ReceiptHash)
	b.Timestamp = block.Timestamp

	l := len(block.Transactions)
	transactions := make([]libblock.TransactionWithData, l)
	for i := 0; i < l; i++ {
		data, err := core.Marshal(block.Transactions[i])
		if err != nil {
			log.Println(err)
			return err
		}
		tx := &TransactionWithData{}
		err = tx.UnmarshalBinary(data)
		if err != nil {
			log.Println(err)
			return err
		}
		transactions[i] = tx
	}
	b.Transactions = transactions

	l = len(block.Receipts)
	receipts := make([]libblock.Receipt, l)
	for i := 0; i < l; i++ {
		data, err := core.Marshal(block.Receipts[i])
		if err != nil {
			return err
		}
		r := &Receipt{}
		err = r.UnmarshalBinary(data)
		if err != nil {
			return err
		}
		receipts[i] = r
	}
	b.Receipts = receipts

	return nil
}

func (b *Block) MarshalBinary() ([]byte, error) {
	block := &pb.Block{
		BlockIndex:      b.BlockIndex,
		ParentHash:      []byte(b.ParentHash),
		TransactionHash: []byte(b.TransactionHash),
		ReceiptHash:     []byte(b.ReceiptHash),
		Timestamp:       b.Timestamp,
	}

	l := len(b.Transactions)
	transactions := make([]*pb.TransactionWithData, l)
	for i := 0; i < l; i++ {
		tx := b.Transactions[i]
		data, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		_, msg, err := core.Unmarshal(data)
		if err != nil {
			return nil, err
		}
		transactions[i] = msg.(*pb.TransactionWithData)
	}
	block.Transactions = transactions

	l = len(b.Receipts)
	receipts := make([]*pb.Receipt, l)
	for i := 0; i < l; i++ {
		r := b.Receipts[i]
		data, err := r.MarshalBinary()
		if err != nil {
			return nil, err
		}
		_, msg, err := core.Unmarshal(data)
		if err != nil {
			return nil, err
		}
		receipts[i] = msg.(*pb.Receipt)
	}
	block.Receipts = receipts

	return core.Marshal(block)
}

func (b *Block) Raw(ignoreSigningFields bool) ([]byte, error) {
	block := &pb.Block{
		BlockIndex:      b.BlockIndex,
		ParentHash:      []byte(b.ParentHash),
		TransactionHash: []byte(b.TransactionHash),
		ReceiptHash:     []byte(b.ReceiptHash),
		Timestamp:       b.Timestamp,
	}

	l := len(b.Transactions)
	transactions := make([]*pb.TransactionWithData, l)
	for i := 0; i < l; i++ {
		tx := b.Transactions[i]
		data, err := tx.Raw(ignoreSigningFields)
		if err != nil {
			return nil, err
		}
		_, msg, err := core.Unmarshal(data)
		if err != nil {
			return nil, err
		}
		transactions[i] = msg.(*pb.TransactionWithData)
	}
	block.Transactions = transactions

	l = len(b.Receipts)
	receipts := make([]*pb.Receipt, l)
	for i := 0; i < l; i++ {
		r := b.Receipts[i]
		data, err := r.Raw(ignoreSigningFields)
		if err != nil {
			return nil, err
		}
		_, msg, err := core.Unmarshal(data)
		if err != nil {
			return nil, err
		}
		receipts[i] = msg.(*pb.Receipt)
	}
	block.Receipts = receipts

	return core.Marshal(block)
}

func (b *Block) GetParentHash() libcore.Hash {
	return libcore.Hash(b.ParentHash)
}

func (b *Block) GetTransactionHash() libcore.Hash {
	return libcore.Hash(b.TransactionHash)
}

func (b *Block) GetReceiptHash() libcore.Hash {
	return libcore.Hash(b.ReceiptHash)
}

func (b *Block) GetTransactions() []libblock.TransactionWithData {
	l := len(b.Transactions)
	ret := make([]libblock.TransactionWithData, l)
	for i := 0; i < l; i++ {
		ret[i] = b.Transactions[i]
	}
	return ret
}

func (b *Block) GetReceipts() []libblock.Receipt {
	l := len(b.Receipts)
	ret := make([]libblock.Receipt, l)
	for i := 0; i < l; i++ {
		ret[i] = b.Receipts[i]
	}
	return ret
}
