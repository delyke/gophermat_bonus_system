package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *ServiceSuite) TestCreateLuhnValidateError() {
	_, err := s.service.Create(s.ctx, []byte("g495599"))
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrBadCredentials)
}

func (s *ServiceSuite) TestCreateLuhnValidateNumberError() {
	_, err := s.service.Create(s.ctx, []byte("495599"))
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderIDLuhnInvalid)
}

func (s *ServiceSuite) TestCreateUUIDFromContextFailure() {
	_, err := s.service.Create(s.ctx, []byte("20000006"))
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUnauthorized)
}

func (s *ServiceSuite) TestCreateOrderExistsOnCurrentUser() {
	userID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userID,
		Login:    "max",
	})
	s.bonusRepository.On("Orders").Return(s.orderRepository)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return(nil, model.ErrOrderIDAlreadyExists)
	s.orderRepository.On("GetByNumber", s.ctx, "20000006").Return(&model.Order{
		UserUUID: userID,
	}, nil)
	_, err := s.service.Create(s.ctx, []byte("20000006"))

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyUploaded)
}

func (s *ServiceSuite) TestCreateOrderUnexpectedError() {
	userID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userID,
		Login:    "max",
	})
	unexpErr := errors.New("unexpected error")
	s.bonusRepository.On("Orders").Return(s.orderRepository)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return(nil, unexpErr)
	_, err := s.service.Create(s.ctx, []byte("20000006"))

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpErr)
}

func (s *ServiceSuite) TestCreateOrderExistsOnAnotherUser() {
	userID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userID,
		Login:    "max",
	})
	s.bonusRepository.On("Orders").Return(s.orderRepository)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return(nil, model.ErrOrderIDAlreadyExists)
	s.orderRepository.On("GetByNumber", s.ctx, "20000006").Return(&model.Order{
		UserUUID: uuid.New(),
	}, nil)
	_, err := s.service.Create(s.ctx, []byte("20000006"))

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderBelongsToAnotherUser)
}

func (s *ServiceSuite) TestCreateOrderGetByIdUnexpected() {
	userID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userID,
		Login:    "max",
	})
	s.bonusRepository.On("Orders").Return(s.orderRepository)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return(nil, model.ErrOrderIDAlreadyExists)
	gByNumErr := errors.New("unexpected error")
	s.orderRepository.On("GetByNumber", s.ctx, "20000006").Return(nil, gByNumErr)
	_, err := s.service.Create(s.ctx, []byte("20000006"))

	s.Require().Error(err)
	s.Require().ErrorIs(err, gByNumErr)
}

func (s *ServiceSuite) TestCreateOrderSuccess() {
	userID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userID,
		Login:    "max",
	})
	orderUUID := uuid.New()
	s.bonusRepository.On("Orders").Return(s.orderRepository)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return(&model.Order{
		UUID: orderUUID,
	}, nil)
	s.accrualWorker.On("Enqueue", mock.Anything).Return(nil)
	cUUID, err := s.service.Create(s.ctx, []byte("20000006"))

	s.Require().NoError(err)
	s.Require().Equal(cUUID, orderUUID.String())
}
