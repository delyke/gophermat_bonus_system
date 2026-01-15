package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func ServiceOrderToApi(serviceOrder model.Order) bonusV1.OrderDto {
	return bonusV1.OrderDto{
		Number:     serviceOrder.OrderID,
		Status:     ServiceOrderStatusToApi(serviceOrder.Status),
		Accrual:    Float64PtrToOptFloat32(serviceOrder.Accrual),
		UploadedAt: serviceOrder.UploadedAt,
	}
}
