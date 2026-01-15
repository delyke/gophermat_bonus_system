package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func ServiceStatusToRepo(status model.OrderStatus) repoModel.OrderStatus {
	switch status {
	case model.OrderNew:
		return repoModel.OrderNew
	case model.OrderProcessing:
		return repoModel.OrderProcessing
	case model.OrderInvalid:
		return repoModel.OrderInvalid
	case model.OrderProcessed:
		return repoModel.OrderProcessed
	}
	return repoModel.OrderInvalid
}
