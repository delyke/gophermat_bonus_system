package converter

import (
	"github.com/delyke/gophermat_bonus_system/internal/model"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func RepoUserToService(user *repoModel.User) *model.User {
	return &model.User{
		UUID:     user.UUID,
		Login:    user.Login,
		Password: user.Password,
		Balance:  user.Balance,
	}
}
