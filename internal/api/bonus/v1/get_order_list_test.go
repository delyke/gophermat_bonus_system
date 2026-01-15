package v1

import (
	"errors"
	"time"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (s *ApiSuite) TestGetOrdersNumberListInternalServerError() {
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("ListByUploadedDesc", s.ctx).Return(nil, errors.New("unexpected error"))

	res, err := s.api.GetOrdersNumberList(s.ctx)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.GetOrdersNumberListInternalServerError{}, res)
}

func (s *ApiSuite) TestGetOrdersNumberListNoContent() {
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("ListByUploadedDesc", s.ctx).Return([]*model.Order{}, nil)

	res, err := s.api.GetOrdersNumberList(s.ctx)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.GetOrdersNumberListNoContent{}, res)
}

func (s *ApiSuite) TestGetOrdersNumberListOK() {
	uploadedAt := time.Now()
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("ListByUploadedDesc", s.ctx).Return([]*model.Order{
		{
			OrderID:    "20000006",
			Status:     model.OrderProcessing,
			UploadedAt: uploadedAt,
		},
	}, nil)

	res, err := s.api.GetOrdersNumberList(s.ctx)

	s.Require().NoError(err)
	okRes, ok := res.(*bonusV1.GetOrdersNumberListResponse)
	s.Require().True(ok)
	s.Require().Len(*okRes, 1)
	s.Require().Equal("20000006", (*okRes)[0].Number)
	s.Require().Equal(bonusV1.OrderStatusPROCESSING, (*okRes)[0].Status)
	s.Require().Equal(uploadedAt, (*okRes)[0].UploadedAt)
}
