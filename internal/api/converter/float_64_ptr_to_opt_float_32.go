package converter

import bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"

func Float64PtrToOptFloat32(v *float64) bonusV1.OptFloat64 {
	if v == nil {
		return bonusV1.OptFloat64{Set: false}
	}

	return bonusV1.OptFloat64{
		Value: *v,
		Set:   true,
	}
}
