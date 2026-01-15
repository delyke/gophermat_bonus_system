package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func ServiceOrderToRepo(ro model.Order) repoModel.Order {
	return repoModel.Order{
		UUID:       ro.UUID,
		OrderID:    ro.OrderID,
		Status:     ServiceStatusToRepo(ro.Status),
		Accrual:    ro.Accrual,
		UserUUID:   ro.UserUUID,
		UploadedAt: ro.UploadedAt,
	}
}
