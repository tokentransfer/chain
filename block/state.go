package block

import (
	"errors"

	"github.com/tokentransfer/chain/account"
	"github.com/tokentransfer/chain/core"
	"github.com/tokentransfer/chain/core/pb"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
)

type State struct {
	Hash libcore.Hash

	Account    libcore.Address
	Sequence   uint64
	BlockIndex uint64

	StateType libblock.StateType
}

func (s *State) GetHash() libcore.Hash {
	return s.Hash
}

func (s *State) SetHash(h libcore.Hash) {
	s.Hash = h
}

func (s *State) GetStateType() libblock.StateType {
	return s.StateType
}

func (s *State) GetBlockIndex() uint64 {
	return s.BlockIndex
}

func (s *State) SetBlockIndex(index uint64) {
	s.BlockIndex = index
}

func (s *State) GetAccount() libcore.Address {
	return s.Account
}

func (s *State) GetIndex() uint64 {
	return s.Sequence
}

type AccountState struct {
	State

	Amount core.Amount
}

func (s *AccountState) GetStateKey() string {
	return core.GetAccountKey(s.Account, s.Amount.Currency, s.Amount.Issuer, "-")
}

func (s *AccountState) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_ACCOUNT_STATE {
		return errors.New("error State data")
	}

	state := msg.(*pb.AccountState)

	_, account, err := account.NewAccountFromBytes(state.Account[:])
	if err != nil {
		return err
	}

	a, err := core.NewAmount(state.Amount)
	if err != nil {
		return err
	}

	s.StateType = libblock.StateType(core.CORE_ACCOUNT_STATE)
	s.BlockIndex = state.BlockIndex
	s.Account = account
	s.Sequence = state.Sequence
	s.Amount = *a
	return nil
}

func (s *AccountState) MarshalBinary() ([]byte, error) {
	a, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.AccountState{
		StateType:  uint32(core.CORE_ACCOUNT_STATE),
		BlockIndex: s.BlockIndex,
		Account:    a,
		Sequence:   s.Sequence,
		Amount:     s.Amount.String(),
	})
}

func (s *AccountState) Raw(ignoreSigningFields bool) ([]byte, error) {
	a, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.AccountState{
		StateType: uint32(core.CORE_ACCOUNT_STATE),
		Account:   a,
		Sequence:  s.Sequence,
		Amount:    s.Amount.String(),
	})
}

type CurrencyState struct {
	State

	Name        string
	Symbol      string
	Decimals    uint32
	TotalSupply core.Amount
}

func (s *CurrencyState) GetStateKey() string {
	return core.GetCurrencyKey(s.TotalSupply.Currency, s.TotalSupply.Issuer, "-")
}

func (s *CurrencyState) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_CURRENCY_STATE {
		return errors.New("error state data")
	}
	state := msg.(*pb.CurrencyState)

	_, issuer, err := account.NewAccountFromBytes(state.Account[:])
	if err != nil {
		return err
	}

	a, err := core.NewAmount(state.TotalSupply)
	if err != nil {
		return err
	}

	s.StateType = libblock.StateType(core.CORE_CURRENCY_STATE)
	s.BlockIndex = state.BlockIndex
	s.Account = issuer
	s.Sequence = state.Sequence
	s.Name = state.Name
	s.Symbol = state.Symbol
	s.Decimals = state.Decimals
	s.TotalSupply = *a

	return nil
}

func (s *CurrencyState) MarshalBinary() ([]byte, error) {
	issuer, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.CurrencyState{
		StateType:   uint32(core.CORE_CURRENCY_STATE),
		BlockIndex:  s.BlockIndex,
		Account:     issuer,
		Sequence:    s.Sequence,
		Name:        s.Name,
		Symbol:      s.Symbol,
		Decimals:    s.Decimals,
		TotalSupply: s.TotalSupply.String(),
	})
}

func (s *CurrencyState) Raw(ignoreSigningFields bool) ([]byte, error) {
	issuer, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.CurrencyState{
		StateType:   uint32(core.CORE_CURRENCY_STATE),
		Account:     issuer,
		Sequence:    s.Sequence,
		Name:        s.Name,
		Symbol:      s.Symbol,
		Decimals:    s.Decimals,
		TotalSupply: s.TotalSupply.String(),
	})
}

func ReadState(data []byte) (libblock.State, error) {
	if len(data) == 0 {
		return nil, errors.New("error entry")
	}
	meta := core.GetMeta(data)
	switch meta {
	case core.CORE_ACCOUNT_STATE:
		s := &AccountState{}
		err := s.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
		return s, nil
	case core.CORE_CURRENCY_STATE:
		s := &CurrencyState{}
		err := s.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
		return s, nil
	default:
		return nil, errors.New("error data")
	}
}

func CloneState(state libblock.State) (libblock.State, error) {
	data, err := state.MarshalBinary()
	if err != nil {
		return nil, err
	}
	s, err := ReadState(data)
	if err != nil {
		return nil, err
	}
	return s, nil
}
