package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func ServiceOrderToAPI(serviceOrder model.Order) bonusV1.OrderDto {
	return bonusV1.OrderDto{
		Number:     serviceOrder.OrderID,
		Status:     ServiceOrderStatusToAPI(serviceOrder.Status),
		Accrual:    Float64PtrToOptFloat64(serviceOrder.Accrual),
		UploadedAt: serviceOrder.UploadedAt,
	}
}
