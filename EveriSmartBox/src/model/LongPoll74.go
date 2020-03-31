package model

import (
	"encoding/binary"

	"github.com/albenik/bcd"
	"github.com/evri/CashlessPayments/EveriSmartBox/src/lib"
	"github.com/pkg/errors"
)

// LongPoll74Inbound ... LP74 Inbound from To EGM from Interceptor
type LongPoll74Inbound struct {
	Address               byte
	Command               byte
	Length                byte
	AssetNumber           []byte
	AssetNumberValue      uint32
	LockStatus            byte
	AvailableTransfer     byte
	HostCashoutStatus     byte
	AFTStatus             byte
	MaxBufferIndex        byte
	CashableBCD           []byte
	CashableAmount        uint64
	RestrictedBCD         []byte
	RestrictedAmount      uint64
	NonRestrictedBCD      []byte
	NonRestrictedAmount   uint64
	MaxTransferLimit      []byte
	MaxTransferLimitValue uint64
	RestrictedExp         [4]byte
	RestrictedPoolID      [2]byte
	CRC                   [2]byte
}

// LongPoll74Outbound ... LP74 Outbound To EGM from Interceptor
type LongPoll74Outbound struct {
	Address         byte
	Command         byte
	LockCode        byte
	TransactioNCode byte
	LockTimeout     []byte
	CRC             [2]byte
}

// LongPoll74IntrogateOutBound ... LP74 Introgate Outgoing message structure To EGM From Interceptor
type LongPoll74IntrogateOutBound struct {
	Address         byte
	Command         byte
	LockCode        byte
	TransactionCode byte
	LockTimeout     []byte
	CRC             [2]byte
}

// BuildAndByte ... USE LP74 to get machine asset number
func (lp74 LongPoll74Outbound) BuildAndByte(params ...interface{}) ([]byte, error) {

	lp74.Address = 0x01
	lp74.Command = 0x74
	lp74.LockCode = 0x00
	lp74.TransactioNCode = 0x00
	lp74.LockTimeout = []byte{0x10, 0x00}

	messageData := []byte{
		lp74.Address,
		lp74.Command,
		lp74.LockCode,
		lp74.TransactioNCode,
	}

	messageData = append(messageData, lp74.LockTimeout...)
	crc := CalculateCRC(0, messageData)
	messageData = append(messageData, crc...)
	return messageData, nil
}

// Parse ... Interface for all LP74
func (lp74 LongPoll74Inbound) Parse(messageBytes []byte, balanceChn chan<- uint64) error {
	if len(messageBytes) < 40 {
		return errors.New("Invalid message length")
	}

	if messageBytes[0] != 0x01 {
		return errors.New("Invalid message address")
	}
	lp74.Address = 0x01

	if messageBytes[1] != 0x74 {
		return errors.New("Invalid. Message command different from 0x74")
	}
	lp74.Command = 0x74

	lp74.AssetNumber = messageBytes[3:7]
	lp74.AssetNumberValue = binary.LittleEndian.Uint32(lp74.AssetNumber[:])
	AssetNumber = lp74.AssetNumberValue
	lp74.LockStatus = messageBytes[7]
	lp74.AvailableTransfer = messageBytes[8]
	lp74.HostCashoutStatus = messageBytes[9]
	lp74.AFTStatus = messageBytes[10]

	lp74.CashableBCD = messageBytes[12:17]
	lp74.CashableAmount = bcd.ToUint64(lp74.CashableBCD[:])

	lp74.RestrictedBCD = messageBytes[17:22]
	lp74.RestrictedAmount = bcd.ToUint64(lp74.RestrictedBCD[:])

	lp74.NonRestrictedBCD = messageBytes[22:27]
	lp74.NonRestrictedAmount = bcd.ToUint64(lp74.NonRestrictedBCD[:])

	lp74.MaxTransferLimit = messageBytes[27:32]
	lp74.MaxTransferLimitValue = bcd.ToUint64(lp74.MaxTransferLimit[:])

	lib.Log(`
	 Address				 %#X
	 Command				 %#X
	 Asset Number			 %d
	 Game Lock Status		 %#X
	 Available Transfer		 %#X
	 HostCashoutStatus		 %#X
	 AFT Status				 %#X
	 Cashable Amt			 %d
	 Restricted Amt			 %d
	 Non-Restricted Amt		 %d
	 Max Allowed			 %d
	 `,
		lp74.Address,
		lp74.Command,
		lp74.AssetNumberValue,
		lp74.LockStatus,
		lp74.AvailableTransfer,
		lp74.HostCashoutStatus,
		lp74.AFTStatus,
		lp74.CashableAmount,
		lp74.RestrictedAmount,
		lp74.NonRestrictedAmount,
		lp74.MaxTransferLimitValue)

	if balanceChn != nil {
		balanceChn <- lp74.CashableAmount
	}

	// 40- Byte array length for the lock response.
	if len(messageBytes) > 40 {
		IncomingMessageCh <- messageBytes[40:]
	}

	if lp74.LockStatus == 0x00 || lp74.LockStatus == 0xFF { // Unable to aquire lock
		MachineLockCh <- false
	}

	return nil
}

// InterogateBuildAndByte ... Intergate 74  to get the balance
func (lp74 LongPoll74IntrogateOutBound) InterogateBuildAndByte() ([]byte, error) {

	lp74.Address = 0x01
	lp74.Command = 0x74
	lp74.LockCode = 0xFF
	lp74.TransactionCode = 0x00
	lp74.LockTimeout = []byte{0x10, 0x00}

	messageData := []byte{
		lp74.Address,
		lp74.Command,
		lp74.LockCode,
		lp74.TransactionCode,
	}

	messageData = append(messageData, lp74.LockTimeout...)

	crc := CalculateCRC(0, messageData)
	messageData = append(messageData, crc...)

	return messageData, nil
}
