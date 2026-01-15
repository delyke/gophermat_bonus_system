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

func (s *ApiSuite) TestOrderNumberLoadReadError() {
	req := bonusV1.OrderNumberLoadReq{Data: errorReader{err: io.ErrUnexpectedEOF}}

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadBadRequest{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadEmptyBody() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("   ")}

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadBadRequest{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadBadCredentials() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrBadCredentials)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadBadRequest{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadLuhnInvalid() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrOrderIdLuhnInvalid)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadUnprocessableEntity{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadUnauthorized() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrUnauthorized)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadUnauthorized{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadConflict() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrOrderBelongsToAnotherUser)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadConflict{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadAlreadyUploaded() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", model.ErrOrderAlreadyUploaded)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadOK{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadInternalServerError() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("", errors.New("unexpected error"))

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadInternalServerError{}, res)
}

func (s *ApiSuite) TestOrderNumberLoadAccepted() {
	req := bonusV1.OrderNumberLoadReq{Data: strings.NewReader("20000006")}
	s.bonusService.On("Orders").Return(s.orderService)
	s.orderService.On("Create", s.ctx, []byte("20000006")).Return("uuid", nil)

	res, err := s.api.OrderNumberLoad(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.OrderNumberLoadAccepted{}, res)
}
