package model

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

// LongPoll73InBound ... LP73 Incoming from EGM
type LongPoll73InBound struct {
	Address            byte
	Command            byte
	Length             byte
	RegistrationStatus byte
	AssetNumber        []byte
	RegistrationKey    []byte
	POSId              []byte
	CRC                []byte
}

// LongPoll73Outbound ... LP73 OutBound from To EGM from Interceptor
type LongPoll73Outbound struct {
	Address            byte
	Command            byte
	Length             byte
	RegistrationStatus byte
	AssetNumber        []byte
	RegistrationKey    []byte
	POSId              []byte
	CRC                []byte
}

// Parse ... Interface for all LP73
func (lp73 LongPoll73InBound) Parse(messageBytes []byte) error {

	if len(messageBytes) < 34 {
		return errors.New("Invalid message length")
	}

	if messageBytes[0] != 0x01 {
		return errors.New("Invalid message address")
	}
	lp73.Address = 0x01

	if messageBytes[1] != 0x73 {
		return errors.New("Invalid. Message command different from 0x72")
	}
	lp73.Command = 0x73

	lp73.RegistrationStatus = messageBytes[3]
	lp73.AssetNumber = messageBytes[4:8]
	lp73.RegistrationKey = messageBytes[8:28]
	lp73.POSId = messageBytes[28:32]

	return nil
}

// BuildAndByte ... build outbound LP73 messages
func (lp73 LongPoll73Outbound) BuildAndByte(params ...interface{}) ([]byte, error) {
	lp73.Address = 0x01
	lp73.Command = 0x73

	byt, ok := (params[0]).(byte)

	if !ok {
		return nil, errors.New("Invalid parameter. lp73 Outbound build method")
	}

	lp73.RegistrationStatus = byte(byt)
	lp73.POSId = []byte{0x00, 0x00, 0x00, 0x00}

	if byt == 0x00 || byt == 0x01 || byt == 0x40 {
		lp73.Length = 0x1D
		lp73.RegistrationKey = RegistrationKey

		binAssetNo := make([]byte, 4)
		binary.LittleEndian.PutUint32(binAssetNo, AssetNumber)
		lp73.AssetNumber = binAssetNo

		messageData := []byte{
			lp73.Address,
			lp73.Command,
		}

		messageData = append(messageData, lp73.RegistrationStatus)
		messageData = append(messageData, lp73.AssetNumber...)
		messageData = append(messageData, lp73.RegistrationKey...)
		messageData = append(messageData, lp73.POSId...)
		crc := CalculateCRC(0, messageData)
		messageData = append(messageData, crc...)

		return messageData, nil
	}

	if byt == 0x80 || byt == 0xFF {
		messageData := []byte{
			lp73.Address,
			lp73.Command,
		}

		messageData = append(messageData, lp73.RegistrationStatus)
		crc := CalculateCRC(0, messageData)
		messageData = append(messageData, crc...)

		return messageData, nil
	}

	return nil, errors.New("Invalid Resistration Code. Only valid registation codes are [0x00 0x01 0x40 0x80 0xFF]")
}
