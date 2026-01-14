package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func RepoStatusToService(status repoModel.OrderStatus) model.OrderStatus {
	switch status {
	case repoModel.OrderNew:
		return model.OrderNew
	case repoModel.OrderProcessing:
		return model.OrderProcessing
	case repoModel.OrderInvalid:
		return model.OrderInvalid
	case repoModel.OrderProcessed:
		return model.OrderProcessed
	}
	return model.OrderInvalid
}
