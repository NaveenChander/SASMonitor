package outbound

import (
	"github.com/evri/CashlessPayments/SBController/config"
	model "github.com/evri/CashlessPayments/SBController/model"
)

var configuration = config.New()

// Adapter for generic functions
type Adapter interface {
	Load() (string, error)
	UnLoad() (string, error)
	UpdateJx(model.JXUpdateRequest) model.CommonResponse
}

type SBClient struct {
	request model.CommonRequest
}

// UcsClient to access UCS endpoints
type UcsClient struct {
	request model.CommonRequest
}

func Execute(request model.CommonRequest) Adapter {

	switch request.AssetNumber[0] {
	case configuration.SBASSETIDENTIFIER[0]:
		return SBClient{request: request}
	case configuration.UCSASSETIDENTIFIER[0]:
		return UcsClient{request: request}
	default:
		return nil
	}
}
