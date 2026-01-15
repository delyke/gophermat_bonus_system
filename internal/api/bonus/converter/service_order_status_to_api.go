package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func ServiceOrderStatusToAPI(status model.OrderStatus) bonusV1.OrderStatus {
	switch status {
	case model.OrderNew:
		return bonusV1.OrderStatusNEW
	case model.OrderProcessing:
		return bonusV1.OrderStatusPROCESSING
	case model.OrderInvalid:
		return bonusV1.OrderStatusINVALID
	case model.OrderProcessed:
		return bonusV1.OrderStatusPROCESSED
	}
	return bonusV1.OrderStatusINVALID
}
