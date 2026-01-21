package converter

import (
	"github.com/samber/lo"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func RepoWithdrawalListToService(list []*repoModel.Withdrawal) []*model.Withdrawal {
	out := make([]*model.Withdrawal, len(list))
	for i, w := range list {
		out[i] = lo.ToPtr(RepoWithdrawalToService(*w))
	}
	return out
}
