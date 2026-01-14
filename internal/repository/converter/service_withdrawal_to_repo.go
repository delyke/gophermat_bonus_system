package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func ServiceWithdrawalToRepo(in model.Withdrawal) repoModel.Withdrawal {
	return repoModel.Withdrawal{
		UUID:        in.UUID,
		UserUUID:    in.UserUUID,
		OrderID:     in.OrderID,
		Amount:      in.Amount,
		ProcessedAt: in.ProcessedAt,
	}
}
