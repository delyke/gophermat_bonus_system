package converter

import accrualV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/accrual/v1"

func OptFloat64ToFloat64(in accrualV1.OptFloat64) *float64 {
	if in.Set {
		return &in.Value
	}
	return nil
}
