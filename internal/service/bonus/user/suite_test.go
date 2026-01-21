package user

import (
	"context"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/repository/mocks"
	serviceMocks "github.com/delyke/gophermat_bonus_system/internal/service/mocks"
	workerMocks "github.com/delyke/gophermat_bonus_system/internal/service/workers/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx                  context.Context //nolint:containedctx
	bonusRepository      *mocks.BonusRepository
	accrualWorker        *workerMocks.AccrualWorker
	service              *service
	faker                *gofakeit.Faker
	userRepository       *mocks.UserRepository
	tokenIssuer          *serviceMocks.TokenIssuer
	withdrawalRepository *mocks.WithdrawalRepository
	logger               *logger.Logger
	originalArgs         []string
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.bonusRepository = mocks.NewBonusRepository(s.T())
	s.accrualWorker = workerMocks.NewAccrualWorker(s.T())
	s.userRepository = mocks.NewUserRepository(s.T())
	s.tokenIssuer = serviceMocks.NewTokenIssuer(s.T())
	s.logger = logger.NewNop()
	s.service = NewService(s.bonusRepository, s.tokenIssuer, s.logger)
	s.withdrawalRepository = mocks.NewWithdrawalRepository(s.T())
	s.faker = gofakeit.New(0)
	s.originalArgs = os.Args
	os.Args = []string{os.Args[0]}
	err := config.Load()
	s.Require().NoError(err)
}

func (s *ServiceSuite) TearDownTest() {
	if s.originalArgs != nil {
		os.Args = s.originalArgs
	}
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
