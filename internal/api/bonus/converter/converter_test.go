package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func TestFloat64PtrToOptFloat64(t *testing.T) {
	result := Float64PtrToOptFloat64(nil)
	require.False(t, result.IsSet())

	value := 12.5
	result = Float64PtrToOptFloat64(&value)
	require.True(t, result.IsSet())
	require.Equal(t, 12.5, result.Value)
}

func TestServiceOrderStatusToApi(t *testing.T) {
	require.Equal(t, bonusV1.OrderStatusNEW, ServiceOrderStatusToAPI(model.OrderNew))
	require.Equal(t, bonusV1.OrderStatusPROCESSING, ServiceOrderStatusToAPI(model.OrderProcessing))
	require.Equal(t, bonusV1.OrderStatusINVALID, ServiceOrderStatusToAPI(model.OrderInvalid))
	require.Equal(t, bonusV1.OrderStatusPROCESSED, ServiceOrderStatusToAPI(model.OrderProcessed))
	require.Equal(t, bonusV1.OrderStatusINVALID, ServiceOrderStatusToAPI(model.OrderStatus("unknown")))
}

func TestServiceOrderToApi(t *testing.T) {
	now := time.Now()
	accrual := 10.5
	order := model.Order{
		OrderID:    "20000006",
		Status:     model.OrderProcessed,
		Accrual:    &accrual,
		UploadedAt: now,
	}

	result := ServiceOrderToAPI(order)

	require.Equal(t, "20000006", result.Number)
	require.Equal(t, bonusV1.OrderStatusPROCESSED, result.Status)
	require.Equal(t, now, result.UploadedAt)
	require.True(t, result.Accrual.IsSet())
	require.Equal(t, 10.5, result.Accrual.Value)
}

func TestServiceOrdersListToApi(t *testing.T) {
	now := time.Now()
	orders := []*model.Order{
		{
			OrderID:    "20000006",
			Status:     model.OrderNew,
			UploadedAt: now,
		},
	}

	result := ServiceOrdersListToAPI(orders)

	require.Len(t, result, 1)
	require.Equal(t, "20000006", result[0].Number)
	require.Equal(t, bonusV1.OrderStatusNEW, result[0].Status)
	require.Equal(t, now, result[0].UploadedAt)
}

func TestServiceWithdrawalToApi(t *testing.T) {
	now := time.Now()
	withdrawal := model.Withdrawal{
		OrderID:     "20000006",
		Amount:      5.5,
		ProcessedAt: now,
	}

	result := ServiceWithdrawalToAPI(withdrawal)

	require.Equal(t, "20000006", result.Order)
	require.Equal(t, 5.5, result.Sum)
	require.Equal(t, now, result.ProcessedAt)
}

func TestServiceWithdrawalListToApi(t *testing.T) {
	now := time.Now()
	withdrawals := []*model.Withdrawal{
		{
			OrderID:     "20000006",
			Amount:      5.5,
			ProcessedAt: now,
		},
	}

	result := ServiceWithdrawalListToAPI(withdrawals)

	require.Len(t, result, 1)
	require.Equal(t, "20000006", result[0].Order)
	require.Equal(t, 5.5, result[0].Sum)
	require.Equal(t, now, result[0].ProcessedAt)
}
