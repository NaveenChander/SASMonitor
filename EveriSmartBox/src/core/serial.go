package core

import (
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/evri/CashlessPayments/EveriSmartBox/src/lib"
	"github.com/evri/CashlessPayments/EveriSmartBox/src/model"
	serial "go.bug.st/serial.v1"
)

//Port ... Serail Port to Read / Write
var port serial.Port
var portName string

//InitializePort ... Create Serial Object
func InitializePort(activePortName string, passivePortName string) {

	portName = activePortName

	mode := &serial.Mode{
		BaudRate: 19200,
		Parity:   serial.SpaceParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	var err error
	port, err = serial.Open(activePortName, mode)
	time.Sleep(time.Second) // Sleep for 1 second and wait for the port to open

	if err != nil {
		lib.Log("Exiting Application unable to open Port %s", activePortName)
		os.Exit(500)
	}
	defer port.ResetInputBuffer()
	defer port.ResetInputBuffer()
	defer port.Close() // Close Active Port before existing

	port, err = serial.Open(passivePortName, mode)
	time.Sleep(time.Second) // Sleep for 1 second and wait for the port to open

	if err != nil {
		lib.Log("Exiting Application unable to open Port %s", passivePortName)
		os.Exit(500)
	}

	defer port.ResetInputBuffer()
	defer port.ResetInputBuffer()
	defer port.Close() // Close Passive Port before existing

	readPort()

}

// WriteData ... Write to serial Port
func WriteData(data []byte) {
	lib.Log("\n----------- BEGIN SEND -------------\n%s ", hex.Dump(data[:]))
	port.ResetOutputBuffer()
	port.Write(data[:])
	lib.Log("----------- END SEND -------------")
}

func readPort() {

	buff := make([]byte, 256)

	for {
		n, err := port.Read(buff)
		dataRead := buff[:n]

		if err != nil {
			lib.Log("Exiting Application Error reading from %s", portName)
			os.Exit(500)
		}

		if buff[0] == 0x00 {
			continue
		}

		// if dataRead[0] != 0x51 || dataRead[0] == 0x69 || dataRead[0] == 0x6f {
		// 	model.IncomingMessageCh <- []byte{dataRead[0]}
		// 	dataRead = make([]byte, 256)
		// }

		ch := make(chan []byte)
		exitLoopCh := make(chan bool)

		timer := time.NewTicker(time.Millisecond * 100)
		go func() {

		InnerForLoop:
			for {
				buff2 := make([]byte, 256)
				x, _ := port.Read(buff2)
				ch <- buff2[:x]
				select {
				case <-timer.C:
					timer.Stop()
					exitLoopCh <- true
					break InnerForLoop
				default:
				}
			}
		}()
	outerForLoop:
		for {
			select {
			case newBuff := <-ch:
				dataRead = append(dataRead, newBuff...)
			case <-exitLoopCh:
				break outerForLoop
			default:
			}
		}

		close(ch)

		if dataRead[0] != 0x00 {
			model.IncomingMessageCh <- dataRead
		}

		if err != nil {
			log.Fatal(err)
			break
		}
	}
}
