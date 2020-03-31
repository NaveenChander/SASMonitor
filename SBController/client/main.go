package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evri/CashlessPayments/SBController/config"
	model "github.com/evri/CashlessPayments/SBController/model"
	outbound "github.com/evri/CashlessPayments/SBController/outbound"
	mux "github.com/gorilla/mux"
)

var configuration = config.New()

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/cashin", cashin).Methods("POST")
	router.HandleFunc("/cashout", cashout).Methods("POST")
	router.HandleFunc("/api/SBController/jackpot/{assetNumber}", jxupdate).Methods("PUT")

	fmt.Println("Server started in 9000 port")

	http.ListenAndServe(":9000", router)

}

func cashin(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")

	var cashinRequest model.CommonRequest
	_ = json.NewDecoder(req.Body).Decode(&cashinRequest)

	adapter := outbound.Execute(cashinRequest)

	response, _ := adapter.Load()

	// TODO - Handle error

	var cahsinResponse model.CommonResponse

	cahsinResponse.Status = response

	json.NewEncoder(resp).Encode(&cahsinResponse)

}

func cashout(resp http.ResponseWriter, req *http.Request) {

	var cashoutRequest model.CommonRequest
	_ = json.NewDecoder(req.Body).Decode(&cashoutRequest)

	adapter := outbound.Execute(cashoutRequest)

	response, _ := adapter.UnLoad()

	// TODO - Handle error
	var commonResponse model.CommonResponse

	commonResponse.Status = response

	json.NewEncoder(resp).Encode(&commonResponse)

}

func jxupdate(resp http.ResponseWriter, req *http.Request) {

	fmt.Println("JX release machine request")

	resp.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)

	assetNumber := params["assetNumber"]

	var jxRequest model.JXUpdateRequest
	_ = json.NewDecoder(req.Body).Decode(&jxRequest)

	var commonRequest model.CommonRequest

	commonRequest.AssetNumber = assetNumber

	adapter := outbound.Execute(commonRequest)

	response := adapter.UpdateJx(jxRequest)

	json.NewEncoder(resp).Encode(&response)

}
