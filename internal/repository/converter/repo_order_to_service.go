package converter

import (
	model "github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func RepoOrderToService(ro repoModel.Order) model.Order {
	return model.Order{
		UUID:       ro.UUID,
		OrderID:    ro.OrderID,
		Status:     RepoStatusToService(ro.Status),
		Accrual:    ro.Accrual,
		UserUUID:   ro.UserUUID,
		UploadedAt: ro.UploadedAt,
	}
}
