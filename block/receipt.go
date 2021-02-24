package block

import (
	"errors"

	"github.com/tokentransfer/chain/core"
	"github.com/tokentransfer/chain/core/pb"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
)

type Receipt struct {
	Hash              libcore.Hash
	TransactionIndex  uint32
	TransactionResult libblock.TransactionResult

	States []libblock.State
}

func (r *Receipt) GetHash() libcore.Hash {
	return r.Hash
}

func (r *Receipt) SetHash(h libcore.Hash) {
	r.Hash = h
}

func (r *Receipt) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}

	if meta != core.CORE_RECEIPT {
		return errors.New("error receipt data")
	}
	receipt := msg.(*pb.Receipt)
	r.TransactionResult = libblock.TransactionResult(receipt.TransactionResult)
	r.TransactionIndex = receipt.TransactionIndex

	list := receipt.GetStates()
	l := len(list)
	states := make([]libblock.State, l)
	for i := 0; i < l; i++ {
		data := list[i]
		state, err := ReadState(data)
		if err != nil {
			return err
		}
		states[i] = state
	}
	r.States = states

	return nil
}

func (r *Receipt) MarshalBinary() ([]byte, error) {
	l := len(r.States)
	states := make([][]byte, l)
	for i := 0; i < l; i++ {
		state := r.States[i]
		data, err := state.MarshalBinary()
		if err != nil {
			return nil, err
		}
		states[i] = data
	}
	receipt := &pb.Receipt{
		TransactionResult: uint32(r.TransactionResult),
		TransactionIndex:  r.TransactionIndex,
		States:            states,
	}
	return core.Marshal(receipt)
}

func (r *Receipt) Raw(ignoreSigningFields bool) ([]byte, error) {
	l := len(r.States)
	states := make([][]byte, l)
	for i := 0; i < l; i++ {
		state := r.States[i]
		data, err := state.Raw(ignoreSigningFields)
		if err != nil {
			return nil, err
		}
		states[i] = data
	}
	receipt := &pb.Receipt{
		TransactionResult: uint32(r.TransactionResult),
		States:            states,
	}
	return core.Marshal(receipt)
}

func (r *Receipt) GetTransactionIndex() uint32 {
	return r.TransactionIndex
}

func (r *Receipt) SetTransactionIndex(index uint32) {
	r.TransactionIndex = index
}

func (r *Receipt) GetTransactionResult() libblock.TransactionResult {
	return r.TransactionResult
}

func (r *Receipt) GetStates() []libblock.State {
	return r.States
}
