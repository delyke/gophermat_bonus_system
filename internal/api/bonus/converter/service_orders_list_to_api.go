package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func ServiceOrdersListToAPI(in []*model.Order) []bonusV1.OrderDto {
	out := make([]bonusV1.OrderDto, len(in))
	for i := range in {
		out[i] = ServiceOrderToAPI(*in[i])
	}
	return out
}
