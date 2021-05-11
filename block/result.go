package block

import (
	libblock "github.com/tokentransfer/interfaces/block"
)

const (
	TrSUCCESS  libblock.TransactionResult = 0 // The transaction was applied.
	TrCOMPOUND libblock.TransactionResult = 1 // The compound error.

	TrINVALID   libblock.TransactionResult = -1 // Invalid.
	TrTIMEOUT   libblock.TransactionResult = -2 // No replies received and timeout.
	TrEXCEPTION libblock.TransactionResult = -3 // Unexpected program state.
	TrALREADY   libblock.TransactionResult = -4 // Already in the ledger.

	TrNO_ENTRY   libblock.TransactionResult = -101 // No matching entry found.
	TrNO_ACCOUNT libblock.TransactionResult = -102 // The account does not exist.

	TrINSUFF_GAS libblock.TransactionResult = -110 // Insufficient balance to pay gas.

	TrBAD_SEQUENCE    libblock.TransactionResult = -120 // Bad sequence.
	TrBAD_PARAMETER   libblock.TransactionResult = -121 // Bad parameter.
	TrBAD_TRANSACTION libblock.TransactionResult = -122 // Bad transaction.
	TrBAD_SIGNATURE   libblock.TransactionResult = -123 // Bad signature.
	TrBAD_AMOUNT      libblock.TransactionResult = -124 // Bad amount.
	TrBAD_ACCOUNT     libblock.TransactionResult = -125 // Bad account.
	TrBAD_SECRET      libblock.TransactionResult = -126 // Bad secret.
)

var resultNames = map[libblock.TransactionResult]struct {
	Name string
	Desc string
}{
	TrSUCCESS:  {"trSUCCESS", `The transaction was applied.`},
	TrCOMPOUND: {"trCOMPOUND", `The compound error.`},

	TrINVALID:   {"trINVALID", `Invalid.`},
	TrTIMEOUT:   {"trTIMEOUT", `No replies received and timeout.`},
	TrEXCEPTION: {"trEXCEPTION", `Unexpected program state.`},
	TrALREADY:   {"trALREADY", `Already in the ledger.`},

	TrNO_ENTRY:   {"trNO_ENTRY", `No matching entry found.`},
	TrNO_ACCOUNT: {"trNO_ACCOUNT", `The account does not exist.`},

	TrINSUFF_GAS: {"trINSUFF_GAS", `Insufficient balance to pay gas.`},

	TrBAD_SEQUENCE:    {"trBAD_SEQUENCE", `Bad sequence.`},
	TrBAD_PARAMETER:   {"trBAD_PARAMETER", `Bad parameter.`},
	TrBAD_TRANSACTION: {"trBAD_TRANSACTION", `Bad transaction.`},
	TrBAD_SIGNATURE:   {"trBAD_SIGNATURE", `Bad signature.`},
	TrBAD_AMOUNT:      {"trBAD_AMOUNT", `Bad amount.`},
	TrBAD_ACCOUNT:     {"trBAD_ACCOUNT", `Bad account.`},
	TrBAD_SECRET:      {"trBAD_SECRET", `Bad secret.`},
}

func init() {
	for result, name := range resultNames {
		result.Register(name.Name, name.Desc)
	}
}
