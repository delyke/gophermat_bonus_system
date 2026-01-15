package v1

import (
	"errors"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (s *ApiSuite) TestBalanceWithdrawalUnauthorized() {
	req := &bonusV1.BalanceWithdrawalRequest{Order: "20000006", Sum: 10}
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("Create", s.ctx, []byte("20000006"), 10.0).Return(model.ErrUnauthorized)

	res, err := s.api.BalanceWithdrawal(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.BalanceWithdrawalUnauthorized{}, res)
}

func (s *ApiSuite) TestBalanceWithdrawalNotEnoughBalance() {
	req := &bonusV1.BalanceWithdrawalRequest{Order: "20000006", Sum: 10}
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("Create", s.ctx, []byte("20000006"), 10.0).Return(model.ErrNotEnoughBalance)

	res, err := s.api.BalanceWithdrawal(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.BalanceWithdrawalPaymentRequired{}, res)
}

func (s *ApiSuite) TestBalanceWithdrawalInvalidOrder() {
	req := &bonusV1.BalanceWithdrawalRequest{Order: "20000006", Sum: 10}
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("Create", s.ctx, []byte("20000006"), 10.0).Return(model.ErrOrderIdLuhnInvalid)

	res, err := s.api.BalanceWithdrawal(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.BalanceWithdrawalUnprocessableEntity{}, res)
}

func (s *ApiSuite) TestBalanceWithdrawalBadCredentials() {
	req := &bonusV1.BalanceWithdrawalRequest{Order: "20000006", Sum: 10}
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("Create", s.ctx, []byte("20000006"), 10.0).Return(model.ErrBadCredentials)

	res, err := s.api.BalanceWithdrawal(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.BalanceWithdrawalUnprocessableEntity{}, res)
}

func (s *ApiSuite) TestBalanceWithdrawalInternalServerError() {
	req := &bonusV1.BalanceWithdrawalRequest{Order: "20000006", Sum: 10}
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("Create", s.ctx, []byte("20000006"), 10.0).Return(errors.New("unexpected error"))

	res, err := s.api.BalanceWithdrawal(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.BalanceWithdrawalInternalServerError{}, res)
}

func (s *ApiSuite) TestBalanceWithdrawalOK() {
	req := &bonusV1.BalanceWithdrawalRequest{Order: "20000006", Sum: 10}
	s.bonusService.On("Withdrawals").Return(s.withdrawalService)
	s.withdrawalService.On("Create", s.ctx, []byte("20000006"), 10.0).Return(nil)

	res, err := s.api.BalanceWithdrawal(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.BalanceWithdrawalOK{}, res)
}
