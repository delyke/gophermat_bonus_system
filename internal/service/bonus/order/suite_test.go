package order

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/mocks"
	workerMocks "github.com/delyke/gophermat_bonus_system/internal/service/workers/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx             context.Context //nolint:containedctx
	bonusRepository *mocks.BonusRepository
	accrualWorker   *workerMocks.AccrualWorker
	service         *service
	faker           *gofakeit.Faker
	orderRepository *mocks.OrderRepository
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.bonusRepository = mocks.NewBonusRepository(s.T())
	s.accrualWorker = workerMocks.NewAccrualWorker(s.T())
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.service = NewService(s.bonusRepository, s.accrualWorker)
	s.faker = gofakeit.New(0)
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

func (s *ServiceSuite) CreateFakeOrder() *model.Order {
	oUuid, err := uuid.Parse(s.faker.UUID())
	if err != nil {
		panic(err)
	}

	oUserUuid, err := uuid.Parse(s.faker.UUID())
	if err != nil {
		panic(err)
	}

	return &model.Order{
		UUID:       oUuid,
		OrderID:    s.faker.Phrase(),
		Status:     model.OrderProcessing,
		Accrual:    nil,
		UserUUID:   oUserUuid,
		UploadedAt: time.Time{},
	}
}
