package withdrawal

import (
	"errors"

	"github.com/google/uuid"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *ServiceSuite) TestGetListUnauthorized() {
	_, err := s.service.GetList(s.ctx)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUnauthorized)
}

func (s *ServiceSuite) TestGetListError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	unexpectedErr := errors.New("list error")
	s.bonusRepository.On("Withdrawals").Return(s.withdrawalRepository)
	s.withdrawalRepository.On("ListByUserByDateDesc", s.ctx, userUUID).Return(nil, unexpectedErr)

	_, err := s.service.GetList(s.ctx)

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
}

func (s *ServiceSuite) TestGetListSuccess() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	withdrawals := []*model.Withdrawal{
		{
			UserUUID: userUUID,
			OrderID:  "20000006",
			Amount:   10,
		},
	}
	s.bonusRepository.On("Withdrawals").Return(s.withdrawalRepository)
	s.withdrawalRepository.On("ListByUserByDateDesc", s.ctx, userUUID).Return(withdrawals, nil)

	result, err := s.service.GetList(s.ctx)

	s.Require().NoError(err)
	s.Require().Equal(withdrawals, result)
}
