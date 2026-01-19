package v1

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	serviceMocks "github.com/delyke/gophermat_bonus_system/internal/service/mocks"
)

type APISuite struct {
	suite.Suite
	ctx               context.Context //nolint:containedctx
	api               *api
	bonusService      *serviceMocks.BonusService
	userService       *serviceMocks.UserService
	orderService      *serviceMocks.OrderService
	withdrawalService *serviceMocks.WithdrawalService
	originalArgs      []string
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.bonusService = serviceMocks.NewBonusService(s.T())
	s.userService = serviceMocks.NewUserService(s.T())
	s.orderService = serviceMocks.NewOrderService(s.T())
	s.withdrawalService = serviceMocks.NewWithdrawalService(s.T())
	s.api = NewAPI(s.bonusService, logger.NewNop())
	s.originalArgs = os.Args
	os.Args = []string{os.Args[0]}
	err := config.Load()
	s.Require().NoError(err)
}

func (s *APISuite) TearDownTest() {
	if s.originalArgs != nil {
		os.Args = s.originalArgs
	}
}

func TestApiSuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
