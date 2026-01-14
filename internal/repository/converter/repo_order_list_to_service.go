package converter

import (
	"github.com/samber/lo"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func RepoOrderListToService(orderList []*repoModel.Order) []*model.Order {
	out := make([]*model.Order, len(orderList))
	for i, order := range orderList {
		out[i] = lo.ToPtr(RepoOrderToService(*order))
	}
	return out
}
