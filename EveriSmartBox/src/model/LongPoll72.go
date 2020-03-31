package model

import (
	"encoding/binary"
	"time"

	"github.com/albenik/bcd"
	"github.com/pkg/errors"
	"github.com/vjeantet/jodaTime"
)

// LongPoll72InBound ... LP72 Incoming message structure from EGM
type LongPoll72InBound struct {
	Address             byte
	Command             byte
	Length              byte
	TransferCode        byte
	TransactionIndex    byte
	TransferType        byte
	CashableBCD         [5]byte
	RestrictedBCD       [5]byte
	NonRestrictedBCD    [5]byte
	TransferFlags       byte
	AssetNumber         [4]byte
	RegistrationKey     [20]byte
	TransactionIDLength byte
	TransactionID       []byte
	Expiration          [4]byte
	PoolID              [2]byte
	ReceiptDataLength   byte
	LockTimeout         []byte
	CRC                 [2]byte
}

// LongPoll72OutBound ... LP72 Outgoing message structure To EGM From Interceptor
type LongPoll72OutBound struct {
	Address             byte
	Command             byte
	Length              byte
	TransferCode        byte
	TransactionIndex    byte
	TransferType        byte
	CashableBCD         []byte
	RestrictedBCD       []byte
	NonRestrictedBCD    []byte
	TransferFlags       byte
	AssetNumber         []byte
	RegistrationKey     []byte
	TransactionIDLength byte
	TransactionID       []byte
	Expiration          []byte
	PoolID              []byte
	ReceiptDataLength   byte
	LockTimeout         []byte
	CRC                 []byte
}

// LongPoll72IntrogateOutBound ... LP72 Introgate Outgoing message structure To EGM From Interceptor
type LongPoll72IntrogateOutBound struct {
	Address          byte
	Command          byte
	Length           byte
	TransferCode     byte
	TransactionIndex byte
}

// Parse ... Interface for all LP72
func (lp72 LongPoll72InBound) Parse(messageBytes []byte) error {

	if messageBytes[0] != 0x01 {
		return errors.New("Invalid message address")
	}
	lp72.Address = 0x01

	if messageBytes[1] != 0x72 {
		return errors.New("Invalid. Message command different from 0x72")
	}
	lp72.Command = 0x72

	// 63- Byte array length for the load response.
	if len(messageBytes) > 63 {
		IncomingMessageCh <- messageBytes[63:]
	}

	return nil
}

// BuildAndByte .. build LP72 outbound
func (LP72 LongPoll72OutBound) BuildAndByte(params ...interface{}) ([]byte, error) {

	cashableMoneyInCents := params[0].(uint64)
	restrictedMoneyInCents := params[1].(uint64)
	nonRestrictedMoneyInCents := params[2].(uint64)

	binAssetNo := make([]byte, 4)
	binary.LittleEndian.PutUint32(binAssetNo, AssetNumber)

	LP72.Address = 0x01
	LP72.Command = 0x72
	LP72.TransferCode = 0x00
	LP72.TransactionIndex = 0x00
	LP72.CashableBCD = bcd.FromUint(cashableMoneyInCents, 5)
	LP72.RestrictedBCD = bcd.FromUint(restrictedMoneyInCents, 5)
	LP72.NonRestrictedBCD = bcd.FromUint(nonRestrictedMoneyInCents, 5)
	LP72.TransferFlags = 0x00
	LP72.AssetNumber = binAssetNo
	LP72.RegistrationKey = RegistrationKey
	LP72.Expiration = []byte{0x01, 0x01, 0x20, 0x25}
	LP72.PoolID = []byte{0x00, 0x00}
	LP72.ReceiptDataLength = 0x00
	LP72.LockTimeout = []byte{0x44, 0xCD}

	date := jodaTime.Format("YYYYMMddHHmmssSSS", time.Now())

	LP72.TransactionID = []byte(date)
	LP72.TransactionIDLength = byte(len(LP72.TransactionID))

	if params[3].(uint64) != 0 {
		LP72.TransferType = 0x80
	} else {
		LP72.TransferType = 0x00
	}

	messageLen := (1 + 1 + 1 + 5 + 5 + 5 + 1 + 4 + 20 + 1 + len(LP72.TransactionID) + 4 + 2 + 1 + 2)
	LP72.Length = byte(messageLen)

	messageData := []byte{
		LP72.Address,
		LP72.Command,
		LP72.Length,
		LP72.TransferCode,
		LP72.TransactionIndex,
		LP72.TransferType,
	}
	messageData = append(messageData, LP72.CashableBCD...)
	messageData = append(messageData, LP72.RestrictedBCD...)
	messageData = append(messageData, LP72.NonRestrictedBCD...)
	messageData = append(messageData, LP72.TransferFlags)
	messageData = append(messageData, LP72.AssetNumber...)
	messageData = append(messageData, LP72.RegistrationKey...)
	messageData = append(messageData, LP72.TransactionIDLength)
	messageData = append(messageData, LP72.TransactionID...)
	messageData = append(messageData, LP72.Expiration...)
	messageData = append(messageData, LP72.PoolID...)
	messageData = append(messageData, LP72.ReceiptDataLength)
	messageData = append(messageData, LP72.LockTimeout...)

	crc := CalculateCRC(0, messageData)

	messageData = append(messageData, crc...)
	return messageData, nil

}
