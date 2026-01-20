package withdrawal

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *ServiceSuite) TestCreateUnauthorized() {
	err := s.service.Create(s.ctx, []byte("20000006"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUnauthorized)
}

func (s *ServiceSuite) TestCreateLuhnValidateError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})

	err := s.service.Create(s.ctx, []byte("g495599"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderIDLuhnInvalid)
}

func (s *ServiceSuite) TestCreateLuhnInvalid() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})

	err := s.service.Create(s.ctx, []byte("495599"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderIDLuhnInvalid)
}

func (s *ServiceSuite) TestCreateGetBalanceError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	unexpectedErr := errors.New("withdrawal error")
	s.bonusRepository.On("CreateWithdrawalWithBalance", s.ctx, mock.Anything).Return(unexpectedErr)

	err := s.service.Create(s.ctx, []byte("20000006"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
}

func (s *ServiceSuite) TestCreateNotEnoughBalance() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	s.bonusRepository.On("CreateWithdrawalWithBalance", s.ctx, mock.Anything).Return(model.ErrNotEnoughBalance)

	err := s.service.Create(s.ctx, []byte("20000006"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrNotEnoughBalance)
}

func (s *ServiceSuite) TestCreateWithdrawalError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	unexpectedErr := errors.New("create withdrawal error")
	s.bonusRepository.On("CreateWithdrawalWithBalance", s.ctx, mock.Anything).Return(unexpectedErr)
	err := s.service.Create(s.ctx, []byte("20000006"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
}

func (s *ServiceSuite) TestCreateSuccess() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})

	s.bonusRepository.On("CreateWithdrawalWithBalance", s.ctx, mock.Anything).Return(nil)

	err := s.service.Create(s.ctx, []byte("20000006"), 100)

	s.Require().NoError(err)
}
