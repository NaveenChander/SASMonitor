package model

import (
	"github.com/albenik/bcd"
	"github.com/pkg/errors"
)

type LongPoll1BInbound struct {
	Address     byte
	Command     byte
	Progressive byte
	Level       byte
	Amount      uint64
	PartialPay  []byte
	ResetId     byte
	CRC         []byte
}

func (lp1B LongPoll1BInbound) Parse(messageBytes []byte, jxAmountChn chan<- uint64) error {

	if len(messageBytes) < 24 {
		return errors.New("Invalid Message")
	}
	if messageBytes[0] != 0x1B {
		return errors.New("Invalid message address")
	}

	lp1B.Address = 0x01

	lp1B.Amount = bcd.ToUint64(messageBytes[3:8])

	jxAmountChn <- lp1B.Amount

	return nil
}
