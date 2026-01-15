package v1

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (s *ApiSuite) TestGetWithdrawalsUnauthorized() {
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("GetList", s.ctx).Return(nil, model.ErrUnauthorized)

	res, err := s.api.GetWithdrawals(s.ctx)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.GetWithdrawalsUnauthorized{}, res)
}

func (s *ApiSuite) TestGetWithdrawalsInternalServerError() {
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("GetList", s.ctx).Return(nil, errors.New("unexpected error"))

	res, err := s.api.GetWithdrawals(s.ctx)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.GetWithdrawalsInternalServerError{}, res)
}

func (s *ApiSuite) TestGetWithdrawalsNoContent() {
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("GetList", s.ctx).Return([]*model.Withdrawal{}, nil)

	res, err := s.api.GetWithdrawals(s.ctx)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.GetWithdrawalsNoContent{}, res)
}

func (s *ApiSuite) TestGetWithdrawalsOK() {
	processedAt := time.Now()
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("GetList", s.ctx).Return([]*model.Withdrawal{
		{
			UserUUID:    uuid.New(),
			OrderID:     "20000006",
			Amount:      10,
			ProcessedAt: processedAt,
		},
	}, nil)

	res, err := s.api.GetWithdrawals(s.ctx)

	s.Require().NoError(err)
	okRes, ok := res.(*bonusV1.GetWithdrawalsResponse)
	s.Require().True(ok)
	s.Require().Len(*okRes, 1)
	s.Require().Equal("20000006", (*okRes)[0].Order)
	s.Require().Equal(10.0, (*okRes)[0].Sum)
	s.Require().Equal(processedAt, (*okRes)[0].ProcessedAt)
}
