package block

import (
	"github.com/tokentransfer/chain/core"

	libblock "github.com/tokentransfer/interfaces/block"
)

// TransactionType
const (
	TRANSACTION libblock.TransactionType = libblock.TransactionType(core.CORE_TRANSACTION)
)

// StateType
const (
	ACCOUNT_STATE  libblock.StateType = libblock.StateType(core.CORE_ACCOUNT_STATE)
	CURRENCY_STATE libblock.StateType = libblock.StateType(core.CORE_CURRENCY_STATE)
)

func init() {
	// TransactionType
	TRANSACTION.Register("Transaction", func(t libblock.TransactionType) libblock.Transaction {
		tx := &Transaction{}
		tx.TransactionType = t
		return tx
	})

	// StateType
	ACCOUNT_STATE.Register("AccountState", func(t libblock.StateType) libblock.State {
		info := &AccountState{}
		info.StateType = t
		return info
	})
	CURRENCY_STATE.Register("CurrencyState", func(t libblock.StateType) libblock.State {
		info := &CurrencyState{}
		info.StateType = t
		return info
	})
}
