package v1

import (
	"context"

	"github.com/samber/lo"

	"github.com/delyke/gophermat_bonus_system/internal/api/converter"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (a *api) GetOrdersNumberList(ctx context.Context) (bonusV1.GetOrdersNumberListRes, error) {
	orders, err := a.bonusService.Orders().ListByUploadedDesc(ctx)
	if err != nil {
		return &bonusV1.GetOrdersNumberListInternalServerError{}, nil
	}
	if len(orders) == 0 {
		return &bonusV1.GetOrdersNumberListNoContent{}, nil
	}
	apiOrders := converter.ServiceOrdersListToApi(orders)
	return lo.ToPtr(bonusV1.GetOrdersNumberListResponse(apiOrders)), nil
}
