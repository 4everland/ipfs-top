package data

import (
	"github.com/4everland/ipfs-top/app/provide/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var ProviderSet = wire.NewSet(
	NewGormClient,
	NewDataSetRepo,
)

type dbLog struct {
	logger *log.Helper
}

func (dblog *dbLog) Printf(format string, args ...interface{}) {
	dblog.logger.Infof(format, args...)
}

func newDbLog(logger *log.Helper) *dbLog {
	return &dbLog{logger: logger}
}

func NewGormClient(conf *conf.Data, l log.Logger) *gorm.DB {
	dataLog := log.NewHelper(log.With(l, "module", "db"))
	db, err := gorm.Open(postgres.Open(conf.PinSet.SourcesDsn), &gorm.Config{
		Logger: logger.New(newDbLog(dataLog), logger.Config{
			Colorful: false,
			LogLevel: logger.Warn,

			SlowThreshold:             time.Second * 2,
			IgnoreRecordNotFoundError: true,
		}),
	})

	if err != nil {
		dataLog.Fatalf("failed opening connection to db: %s, err: %v", conf.PinSet.SourcesDsn, err)
	}

	return db
}
