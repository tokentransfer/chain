package block

import (
	"testing"

	. "github.com/tokentransfer/check"
)

type ReceiptSuite struct{}

func Test_Receipt(t *testing.T) {
	s := Suite(&ReceiptSuite{})
	TestingRun(t, s)
}
