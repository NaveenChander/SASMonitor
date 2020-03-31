package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/evri/CashlessPayments/EveriSmartBox/src/lib"
	"github.com/evri/CashlessPayments/EveriSmartBox/src/model"
	"github.com/evri/CashlessPayments/EveriSmartBox/src/outbound/jackpot"
)

var balanceChn chan uint64

var jxAmountChn = make(chan uint64)

// DisableEGM ... Send Command to Disable EGM
func DisableEGM() {
	fmt.Println("There....")
	message := []byte{0x01, 0x01, 0x51, 0x08}
	WriteData(message)
}

// EnableEGM ... Send Command to Enable EGM
func EnableEGM() {
	message := []byte{0x01, 0x02, 0xCA, 0x3A}
	WriteData(message)
}

// LoadEGMwithFunds ... Load EGM with funds
func LoadEGMwithFunds(cashableMoneyInCents uint64, restrictedMoneyInCents uint64, nonRestrictedMoneyInCents uint64) error {

	// Step 1: Send Lock command
	lockEGMBeforeFundTransfer()

	// Step 2: Wait for a proper response ->6f  Else return error within 15 seconds make this timeout configurable
	shouldTimeout := true
	select {
	case res := <-model.MachineLockCh:
		if res == false { // if unable to aquire a lock retry
			lib.Log("Retrying to Lock EGM")
			lockEGMBeforeFundTransfer()
		} else {
			shouldTimeout = false
		}
	case <-time.After(15 * time.Second):
		if shouldTimeout {
			lib.Log("Timed out no response from machine on Lock EGM for fund transfer")
			return errors.New("TIMEOUT RECEIVED -> Lock EGM for fund transfer Failed ")
		}
	}

	// Step 3: Send the fund transfer amount LP72
	var lp72 model.LongPoll72OutBound
	var transCode uint64 = 0

	messageData, _ := lp72.BuildAndByte(cashableMoneyInCents, restrictedMoneyInCents, nonRestrictedMoneyInCents, transCode)
	WriteData(messageData)

	// Step 4: Wait for an ackknowledgement
	shouldTimeout = true
	select {
	case res := <-model.MachineLoadFuncCh:
		if res == false { // Failed to transfer funds retry
			lib.Log("Retrying to Load funds")
			WriteData(messageData)
		} else {
			shouldTimeout = false
		}
	case <-time.After(20 * time.Second):
		if shouldTimeout {
			lib.Log("Timed out no response from machine on Load Funds to EGM")
			return errors.New("TIMEOUT RECEIVED -> Load funds Failed ")
		}
	}

	// Step 5: Respond with an introgation message and then confirm wallet
	messageData = []byte{0x01, 0x72, 0x02, 0xFF, 0x00, 0x0F, 0x22}
	WriteData(messageData)

	return nil
}

// UnLoadFunds ...Unload funds from EGM
func UnLoadFunds() error {

	balanceChn = make(chan uint64)

	// step 1: Lock EGM

	lockEGMBeforeFundTransfer()

	// Step 2: Wait for a proper response ->6f  Else return error within 15 seconds make this timeout configurable
	shouldTimeout := true
	select {
	case res := <-model.MachineLockCh:
		if res == false { // if unable to aquire a lock retry
			lib.Log("Retrying to Lock EGM")
			lockEGMBeforeFundTransfer()
		} else {
			shouldTimeout = false
		}
	case <-time.After(15 * time.Second):
		if shouldTimeout {
			lib.Log("Timed out no response from machine on Lock EGM for fund transfer")
			return errors.New("TIMEOUT RECEIVED -> Lock EGM for fund transfer Failed ")
		}
	}

	// Balance from 74 response
	balance := <-balanceChn

	balanceChn = nil

	fmt.Println("Balance............", balance)

	// Step 3: Send the fund transfer amount LP72
	var lp72 model.LongPoll72OutBound
	var transCode uint64 = 80
	var restrictedMoneyInCents uint64 = 0
	var nonRestrictedMoneyInCents uint64 = 0

	messageData, _ := lp72.BuildAndByte(balance, restrictedMoneyInCents, nonRestrictedMoneyInCents, transCode)
	WriteData(messageData)

	// Step 4: Wait for an ackknowledgement
	shouldTimeout = true
	select {
	case res := <-model.MachineLoadFuncCh:
		if res == false { // Failed to transfer funds retry
			lib.Log("Retrying to Load funds")
			WriteData(messageData)
		} else {
			shouldTimeout = false
		}
	case <-time.After(20 * time.Second):
		if shouldTimeout {
			lib.Log("Timed out no response from machine on Load Funds to EGM")
			return errors.New("TIMEOUT RECEIVED -> Load funds Failed ")
		}
	}

	// Step 5: Respond with an introgation message and then confirm wallet
	messageData = []byte{0x01, 0x72, 0x02, 0xFF, 0x00, 0x0F, 0x22}
	WriteData(messageData)
	return nil
}

func lockEGMBeforeFundTransfer() {
	var lp74 model.LongPoll74Outbound
	data, _ := lp74.BuildAndByte(nil)
	WriteData(data)
}

func interogateAFT() {
	var interogate model.LongPoll74IntrogateOutBound
	data, _ := interogate.InterogateBuildAndByte()
	WriteData(data)
}

// Handset repay ... Send Command to repay
func HandsetRepay(dispenseType string) {
	message := []byte{0x01, 0x94, 0x75, 0xCB}

	if dispenseType == "KeyToCredit" {
		var lpA8 model.LongPollA8OutBound
		data, _ := lpA8.BuildAndByte()
		WriteData(data)

		time.Sleep(1 * time.Second)

		WriteData(message)

	}
	WriteData(message)
}

// DequeueIncomingMessage ... Receive incoming channel
func DequeueIncomingMessage() {

	for {
		select {
		case n := <-model.IncomingMessageCh:
			lib.Log("\n*********************** - Begin READ - *************************\n%s\n*********************** - End READ   - *************************", hex.Dump(n[:]))
			SASParser(n)
			break
		}
	}
}

// SASParser ... Parse incoming SAS message build object
func SASParser(incoming []byte) model.SASMessage {
	switch incoming[0] {
	case 0x01:
		switch incoming[1] {
		case 0x72:
			var lp72 model.LongPoll72InBound
			lp72.Parse(incoming)
			return lp72
		case 0x73:
			var lp73 model.LongPoll73InBound
			lp73.Parse(incoming)
			return lp73
		case 0x74:
			var lp74 model.LongPoll74Inbound

			go lp74.Parse(incoming, balanceChn)

			return lp74
		}
		break
	case 0x6f:
		var e6f model.Event6f
		model.Event6f.Parse(e6f, incoming)
		break
	case 0x69:
		var e69 model.Event69
		model.Event69.Parse(e69, incoming)
		break
	case 0x51:

		var e51 model.Event51
		go model.Event51.Parse(e51, incoming)

		jackpotResp := <-model.JackpotReceivedCh

		if jackpotResp == true {
			messageData := []byte{0x01, 0x1B}
			WriteData(messageData)
		}
		break
	case 0x1B:

		var jx1B model.LongPoll1BInbound
		go jx1B.Parse(incoming, jxAmountChn)

		jxAmount := <-jxAmountChn

		jackpot.ProcessJackpot(jxAmount)
		break
	default:
		break
	}

	return nil
}
