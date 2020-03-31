package jackpot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/evri/CashlessPayments/EveriSmartBox/src/config"
)

var configuration = config.New()

// ProcessJackpot ... Process Jackpot
func ProcessJackpot(amount uint64) {

	fmt.Println("In process jackpot.....")

	jsonData := map[string]string{
		"amount":      strconv.Itoa(int(amount)),
		"assetNumber": "12222",
	}
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(configuration.JXURL, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Println("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

}
