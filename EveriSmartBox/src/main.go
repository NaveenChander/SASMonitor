package main

import (
	"bufio"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/evri/CashlessPayments/EveriSmartBox/src/config"
	"github.com/evri/CashlessPayments/EveriSmartBox/src/core"
	"github.com/evri/CashlessPayments/EveriSmartBox/src/lib"
)

var quitCh = make(chan struct{})
var configuration = config.New()

func main() {
	lib.Log("Hello %s/%s\n", runtime.GOOS, runtime.GOARCH)
	// 01 72 45 00 00 00 00 00 00 10 00 00 00 00 00 00 00 00 00 00 00 00 AE 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 12 54 65 73 74 20 54 72 61 6E 73 61 63 74 69 6F 6E 31 39 05 30 20 03 0C 00 00
	go core.InitializePort(configuration.ActivePort, configuration.PassivePort)
	// go core.DequeueIncomingMessage()
	go getConsoleCommands()

loop:
	for {
		select {
		case <-quitCh:
			lib.Log("\n  Application Exiting \n")
			break loop
		}
	}
	time.Sleep(1000 * time.Millisecond) // adding 1000 milli seconds delay to allow all goRoutines to complete
}

func getConsoleCommands() {

loop:
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			go core.DisableEGM()
		case "2":
			go core.EnableEGM()
		case "3":
			go core.LoadEGMwithFunds(1000, 0, 0)
		case "4":
			// go loadWalletWithCredits(assetNumber, registrationKey)
		case "5":
			// registerAFTwithEGM(assetNumber, registrationKey)
		case "q":
			lib.Log("q Entered")
			quitCh <- struct{}{}
			break loop
		default:
			lib.Log("Invalid Entry.. Received Command is %s ", input)
		}

	}
}

func random(min, max int) int {
	rand.Seed(1000)
	return rand.Intn(max-min) + min
}
