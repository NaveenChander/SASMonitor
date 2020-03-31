package model

// Cashin Request
type CommonRequest struct {
	Cents            uint64 `json:"cents"`
	Reference        uint64 `json:"reference"`
	AssetNumber      string `json:"assetNumber"`
	PlayerCardNumber string `json:"playerCardNumber"`
}

type JXUpdateRequest struct {
	DispenseType string `json:"dispensetype"`
	Amount       uint64 `json:"amount"`
}
