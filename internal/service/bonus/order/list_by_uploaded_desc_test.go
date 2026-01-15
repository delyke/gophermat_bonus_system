package order

import (
	"errors"

	"github.com/google/uuid"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *ServiceSuite) TestListByUploadedDescService_NotAuthorized() {
	_, err := s.service.ListByUploadedDesc(s.ctx)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUnauthorized)
}

func (s *ServiceSuite) TestListByUploadedDescService_RepoError() {
	userID, err := uuid.Parse(s.faker.UUID())
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userID,
		Login:    "max",
	})
	unExpErr := errors.New("unexpected error")
	s.bonusRepository.On("Orders").Return(s.orderRepository)
	s.orderRepository.On("GetListByUploadedDesc", s.ctx, userID).Return(nil, unExpErr)

	_, err = s.service.ListByUploadedDesc(s.ctx)
	s.Require().Error(err)
	s.Require().ErrorIs(err, unExpErr)
}

func (s *ServiceSuite) TestListByUploadedDescService_Success() {
	userID, err := uuid.Parse(s.faker.UUID())
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userID,
		Login:    "max",
	})
	orders := make([]*model.Order, 3)
	orders = append(orders, s.CreateFakeOrder())
	orders = append(orders, s.CreateFakeOrder())
	orders = append(orders, s.CreateFakeOrder())
	s.bonusRepository.On("Orders").Return(s.orderRepository)
	s.orderRepository.On("GetListByUploadedDesc", s.ctx, userID).Return(orders, nil)

	ords, err := s.service.ListByUploadedDesc(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(orders, ords)
}
