package srsmgmtrepo

import (
	"srsmgmt/config"
	"srsmgmt/internal/srsmgmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var db *gorm.DB

type Repo struct {
	Db     *gorm.DB
	Mock   sqlmock.Sqlmock
	Logger log.Logger
}

func (repo Repo) HealthCheck() (bool, error) {
	if sqlDB, err := db.DB(); err != nil {
		return false, err
	} else if err := sqlDB.Ping(); err != nil {
		return false, err
	}

	return true, nil
}

func New(logger log.Logger, gloggerMode glogger.LogLevel, cfg *config.Config) srsmgmt.Repository {
	var err error

	if gloggerMode < glogger.Silent {
		gloggerMode = glogger.Info
	}
	db, err = gorm.Open(postgres.Open(cfg.DbURI), &gorm.Config{
		Logger: glogger.Default.LogMode(gloggerMode),
	})
	if err != nil {
		level.Error(logger).Log("DB", "failed to connect database")
	}

	DbInit(logger, db)

	return Repo{
		Db:     db,
		Mock:   nil,
		Logger: logger,
	}
}

func DbInit(logger log.Logger, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		level.Error(logger).Log("DB", "failed to connect database: ", err)
	}

	db.AutoMigrate(&Stream{})

	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	logger.Log("DB", "Finished DB init")

}

func (repo Repo) GetMock() sqlmock.Sqlmock {
	return repo.Mock
}
