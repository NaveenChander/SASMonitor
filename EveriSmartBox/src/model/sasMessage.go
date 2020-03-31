package model

import (
	"encoding/binary"
)

type treeNodeType int

// AssetNumber ... EGM's asset number
var AssetNumber uint32 = 2222

// RegistrationKey ... EGM's Registration key
var RegistrationKey = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20}

const (
	address       treeNodeType = 0
	command       treeNodeType = 1
	commandLength treeNodeType = 2
	data          treeNodeType = 3
	crc           treeNodeType = 4
)

const (
	initializeRegistration  byte = 0x00
	registrationEGM         byte = 0x01
	requestOperatorAck      byte = 0x40
	unregisteredEGM         byte = 0x80
	readCurrentRegistration byte = 0xFF
)

//TreeNode ... Tree structure
type TreeNode struct {
	nodeData       byte
	children       []TreeNode
	nodeType       treeNodeType
	SASMessageType SASIncomingMessage
}

// IncomingMessageCh ... Buffered channel for incoming messages from Serial port
var IncomingMessageCh = make(chan []byte, 100)

// MachineLoadFuncCh ... bool to ack machine locked
var MachineLoadFuncCh = make(chan bool)

// MachineUnloadFuncCh ... bool to ack machine locked
var MachineUnloadFuncCh = make(chan bool)

// MeachineIntrogationCh ... bool to ack Introgation message
var MeachineIntrogationCh = make(chan bool)

// MachineLockCh ... bool to ack machine locked
var MachineLockCh = make(chan bool)

// JackpotReceivedCh ... ack when Jackpot Received to ack machine locked
var JackpotReceivedCh = make(chan bool)

// SASMessage ... Blank interface to encapsulate all SASMessages
type SASMessage interface{}

// SASIncomingMessage ... Interface for all incoming messages
type SASIncomingMessage interface {
	Parse(messageBytes []byte) error
}

// SASOutboundMessage ... Interface for all outgoing messages
type SASOutboundMessage interface {
	BuildAndByte(params ...interface{}) ([]byte, error)
}

// Event51 ... Raise 51
type Event51 struct {
	Command byte
}

// Event69 ... Raise 69
type Event69 struct {
	Command byte
}

// Event6f ... Raise 6f
type Event6f struct {
	Command byte
}

// Parse ... Interface for all Event51
func (Event51) Parse(messageBytes []byte) error {
	JackpotReceivedCh <- true
	return nil
}

// Parse ... Interface for all Event69
func (Event69) Parse(messageBytes []byte) error {
	MachineLoadFuncCh <- true
	return nil
}

// Parse ... Interface for all Event6f
func (Event6f) Parse(messageBytes []byte) error {
	MachineLockCh <- true
	return nil
}

// CalculateCRC ... CalculateCRC
func CalculateCRC(crcValue uint16, data []byte) []byte {

	for i := 0; i < len(data); i++ {

		dataByte := uint16(data[i])
		byte1 := uint16(0x0F)
		byte2 := uint16(0x1081)

		crcValue = (uint16((crcValue >> 4) ^ (((crcValue ^ dataByte) & byte1) * byte2)))
		crcValue = (uint16((crcValue >> 4) ^ (((crcValue ^ (dataByte >> 4)) & byte1) * byte2)))
	}

	returnValue := make([]byte, 2)
	binary.LittleEndian.PutUint16(returnValue, crcValue)

	return returnValue
}

// // BuildInBoundSASTreeRequest ... BuildSASMessageTree
// func BuildInBoundSASTreeRequest() []TreeNode {
// 	return []TreeNode{
// 		{
// 			nodeData: 0x01,
// 			nodeType: data,
// 			children: []TreeNode{
// 				{
// 					nodeData: 0x00,
// 					children: nil,
// 					nodeType: data,
// 				},
// 				{
// 					nodeData: 0x72,
// 					nodeType: command,
// 					children: []TreeNode{},
// 				},
// 				{
// 					nodeData: 0x73,
// 					nodeType: command,
// 					children: []TreeNode{},
// 				},
// 				{
// 					nodeData: 0x74,
// 					nodeType: command,
// 					children: []TreeNode{},
// 				},
// 				{
// 					nodeData: 0xA8,
// 					nodeType: command,
// 					children: []TreeNode{},
// 				},
// 				{
// 					nodeData: 0x94,
// 					nodeType: command,
// 					children: []TreeNode{},
// 				},
// 				{
// 					nodeData: 0x1B,
// 					nodeType: command,
// 					children: []TreeNode{},
// 				},
// 			},
// 		},
// 		{
// 			nodeData: 0x6f,
// 			nodeType: data,
// 			children: nil,
// 		},
// 		{
// 			nodeData: 0x69,
// 			nodeType: data,
// 			children: nil,
// 		},
// 		{
// 			nodeData: 0x51,
// 			nodeType: data,
// 			children: nil,
// 		},
// 	}
// }
