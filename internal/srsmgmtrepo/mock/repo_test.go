package mock

import (
	"context"
	"regexp"
	"srsmgmt/config"
	"srsmgmt/internal/srsmgmt"
	"srsmgmt/internal/srsmgmtrepo"

	"github.com/go-kit/log"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	glogger "gorm.io/gorm/logger"
)

func NewMock(logger log.Logger, cfg *config.Config) srsmgmt.Repository {
	sqlDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	if err != nil {
		logger.Log("[sqlmock new] %s", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DSN:        cfg.DbURI,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent),
	})
	if err != nil {
		logger.Log("[gorm open] %s", err)
	}

	srsmgmtrepo.DbInit(logger, db)

	return srsmgmtrepo.Repo{
		Db:     db,
		Mock:   mock,
		Logger: logger,
	}
}

func TestMockDB(t *testing.T) {
	logger := log.NewNopLogger()

	repo := NewMock(logger, config.GetConfig())

	s := srsmgmt.NewSrsMgmtService(repo, logger, nil, nil)
	if s == nil {
		t.Error("init service failed")
	}

	mock := repo.GetMock()

	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "streams" `)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	stream := srsmgmt.Stream{
		StreamID: uuid.FromStringOrNil("00000000-2222-0000-0000-000000000000"),
		Password: "12321",
	}
	_, err := s.CreateStream(context.Background(), stream)
	if err != nil {
		t.Errorf("Failed to CreateStream, got error: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "streams" `)).
		WithArgs("00000000-2222-0000-0000-000000000000").
		WillReturnRows(
			mock.NewRows([]string{"id", "app", "password", "status", "hls", "rtmpPush", "srtPush", "createdAt", "updatedAt", "clientId", "startedAt", "stopedAt"}).
				AddRow("00000000-0000-0000-0000-000000000000", "", "", 0, "", "", "", "0001-01-01 00:00:00 +0000 UTC", "0001-01-01 00:00:00 +0000 UTC", "", "0001-01-01 00:00:00 +0000 UTC", "0001-01-01 00:00:00 +0000 UTC"),
		)

	_, err = s.GetStream(context.Background(), stream)
	if err != nil {
		t.Errorf("Failed to GetStream, got error: %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
