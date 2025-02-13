package db

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/plasmatrip/avito_merch/internal/config"
	"github.com/plasmatrip/avito_merch/internal/logger"
	mock_db "github.com/plasmatrip/avito_merch/internal/storage/db/mock"
)

type DBTestSuite struct {
	suite.Suite
	config *config.Config
	logger *logger.Logger
}

func (suite *DBTestSuite) SetupSuite() {
	var err error
	suite.config = &config.Config{
		Host:         "localhost:8080",
		Database:     "postgres://avito_merch:password@localhost:5432/avito_merch?sslmode=disable",
		LogLevel:     "debug",
		TokenSecret:  "T0kenS3cRE7",
		ReadTimeout:  5,
		WriteTimeout: 10,
		IdleTimeout:  60,
	}
	suite.logger, err = logger.NewLogger("debug")
	if err != nil {
		panic(err)
	}
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (suite *DBTestSuite) TestPing() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock_db.NewMockDB(ctrl)

	suite.NotNil(m)

	m.EXPECT().
		Ping(gomock.Any()).
		Return(nil)

	suite.NotPanics(func() {
		m.Ping(context.Background())
	})
}

func (suite *DBTestSuite) TestBuyItem() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock_db.NewMockDB(ctrl)

	suite.NotNil(m)

	m.EXPECT().
		BuyItem(gomock.Any(), gomock.Any(), gomock.Any()).
		Return()

	suite.NotPanics(func() {
		m.BuyItem(context.Background(), uuid.Nil, "")
	})
}
