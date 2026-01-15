package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func RepoWithdrawalToService(in repoModel.Withdrawal) model.Withdrawal {
	return model.Withdrawal{
		UUID:        in.UUID,
		UserUUID:    in.UserUUID,
		OrderID:     in.OrderID,
		Amount:      in.Amount,
		ProcessedAt: in.ProcessedAt,
	}
}
