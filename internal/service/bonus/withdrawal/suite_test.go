package withdrawal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx                  context.Context //nolint:containedctx
	bonusRepository      *mocks.BonusRepository
	service              *service
	userRepository       *mocks.UserRepository
	withdrawalRepository *mocks.WithdrawalRepository
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.bonusRepository = mocks.NewBonusRepository(s.T())
	s.userRepository = mocks.NewUserRepository(s.T())
	s.withdrawalRepository = mocks.NewWithdrawalRepository(s.T())
	s.service = NewService(s.bonusRepository)
	err := logger.Init(
		"debug",
		true,
	)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
