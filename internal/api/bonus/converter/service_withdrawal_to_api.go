package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func ServiceWithdrawalToAPI(in model.Withdrawal) bonusV1.WithdrawalDto {
	return bonusV1.WithdrawalDto{
		Order:       in.OrderID,
		Sum:         in.Amount,
		ProcessedAt: in.ProcessedAt,
	}
}
