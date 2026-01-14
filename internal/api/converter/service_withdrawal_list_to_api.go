package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func ServiceWithdrawalListToApi(in []*model.Withdrawal) []bonusV1.WithdrawalDto {
	out := make([]bonusV1.WithdrawalDto, len(in))
	for i := range in {
		out[i] = ServiceWithdrawalToApi(*in[i])
	}
	return out
}
