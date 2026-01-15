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
	s.Require().ErrorIs(err, model.ErrBadCredentials)
}

func (s *ServiceSuite) TestCreateLuhnInvalid() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})

	err := s.service.Create(s.ctx, []byte("495599"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderIdLuhnInvalid)
}

func (s *ServiceSuite) TestCreateGetBalanceError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	unexpectedErr := errors.New("balance error")
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(0.0, unexpectedErr)

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
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(50.0, nil)

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
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.bonusRepository.On("Withdrawals").Return(s.withdrawalRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(200.0, nil)
	s.withdrawalRepository.On("Create", s.ctx, mock.Anything).Return(unexpectedErr)

	err := s.service.Create(s.ctx, []byte("20000006"), 100)

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
}

func (s *ServiceSuite) TestCreateSetBalanceError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	unexpectedErr := errors.New("set balance error")
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.bonusRepository.On("Withdrawals").Return(s.withdrawalRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(200.0, nil)
	s.withdrawalRepository.On("Create", s.ctx, mock.Anything).Return(nil)
	s.userRepository.On("SetBalanceByUUID", s.ctx, userUUID, 100.0).Return(unexpectedErr)

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
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.bonusRepository.On("Withdrawals").Return(s.withdrawalRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(200.0, nil)
	s.withdrawalRepository.On("Create", s.ctx, mock.Anything).Return(nil)
	s.userRepository.On("SetBalanceByUUID", s.ctx, userUUID, 100.0).Return(nil)

	err := s.service.Create(s.ctx, []byte("20000006"), 100)

	s.Require().NoError(err)
}
