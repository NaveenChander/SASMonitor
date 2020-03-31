package outbound

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	model "github.com/evri/CashlessPayments/SBController/model"
)

func (ucs UcsClient) Load() (string, error) {

	jsonData := map[string]string{
		"amount":     strconv.Itoa(int(ucs.request.Cents)),
		"terminalId": ucs.request.AssetNumber,
	}
	jsonValue, _ := json.Marshal(jsonData)

	url := configuration.UCSHOST + configuration.UCSLOAD

	url = strings.Replace(url, "{assetNumber}", ucs.request.AssetNumber, -1)
	url = strings.Replace(url, "{playerAccountID}", ucs.request.PlayerCardNumber, -1)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		return "", err
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		log.Fatalf("Response from UCS: %s", data)
		// TODO - determine the response
		return "success", nil
	}
}

func (ucs UcsClient) UnLoad() (string, error) {

	// read amount from response
	jsonData := map[string]string{
		"amount":     "0",
		"terminalId": ucs.request.AssetNumber,
	}
	jsonValue, _ := json.Marshal(jsonData)

	url := configuration.UCSHOST + configuration.UCSUNLOAD

	url = strings.Replace(url, "{assetNumber}", ucs.request.AssetNumber, -1)
	url = strings.Replace(url, "{playerAccountID}", ucs.request.PlayerCardNumber, -1)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		return "", err
	}

	data, _ := ioutil.ReadAll(response.Body)

	log.Fatalf("Response from UCS: %s", data)
	// TODO - determine the response

	var ucsBalance model.UCSBalance

	_ = json.Unmarshal(data, &ucsBalance)

	return ucsBalance.Amount, nil

}

func (sb UcsClient) UpdateJx(jxRequest model.JXUpdateRequest) model.CommonResponse {
	return model.CommonResponse{Status: "NotImplemented"}
}
