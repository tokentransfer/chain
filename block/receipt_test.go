package block

import (
	"testing"

	. "gopkg.in/check.v1"
)

type ReceiptSuite struct{}

func Test_Receipt(t *testing.T) {
	s := Suite(&ReceiptSuite{})
	TestingRun(t, s)
	// TestingT(t)
}
