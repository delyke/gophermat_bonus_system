package v1

import (
	"errors"
	"io"
	"strings"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

type errorReader struct {
	err error
}

func (r errorReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

func (s *APISuite) TestOrderNumberLoadReadError() {
	req := bonusV1.OrderNumberLoadReq{Data: errorReader{err: io.ErrUnexpectedEOF}}

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadBadRequest{}, res)
}

func (s *APISuite) TestOrderNumberLoadEmptyBody() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("   ")}

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadBadRequest{}, res)
}

func (s *APISuite) TestOrderNumberLoadBadCredentials() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrBadCredentials)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadBadRequest{}, res)
}

func (s *APISuite) TestOrderNumberLoadLuhnInvalid() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrOrderIDLuhnInvalid)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadUnprocessableEntity{}, res)
}

func (s *APISuite) TestOrderNumberLoadUnauthorized() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrUnauthorized)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadUnauthorized{}, res)
}

func (s *APISuite) TestOrderNumberLoadConflict() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrOrderBelongsToAnotherUser)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadConflict{}, res)
}

func (s *APISuite) TestOrderNumberLoadAlreadyUploaded() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrOrderAlreadyUploaded)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadOK{}, res)
}

func (s *APISuite) TestOrderNumberLoadInternalServerError() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", errors.New("unexpected error"))

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadInternalServerError{}, res)
}

func (s *APISuite) TestOrderNumberLoadAccepted() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("uuid", nil)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadAccepted{}, res)
}
