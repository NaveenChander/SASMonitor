package model

type LongPollA8Inbound struct {
	Address byte
	Command byte
	ACKCode byte
	CRC     []byte
}

// LongPoll74IntrogateOutBound ... LP74 Introgate Outgoing message structure To EGM From Interceptor
type LongPollA8OutBound struct {
	Address   byte
	Command   byte
	ResetCode byte
	CRC       [2]byte
}

func (lp80 LongPollA8Inbound) Parse(messageBytes []byte) error {
	return nil
}

// BuildAndByte ... USE LP74 to get machine asset number
func (lp80 LongPollA8OutBound) BuildAndByte() ([]byte, error) {

	lp80.Address = 0x01
	lp80.Command = 0xA8
	lp80.ResetCode = 0x01

	messageData := []byte{
		lp80.Address,
		lp80.Command,
		lp80.ResetCode,
	}

	crc := CalculateCRC(0, messageData)
	messageData = append(messageData, crc...)
	return messageData, nil
}
