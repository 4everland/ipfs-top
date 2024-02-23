package data

import (
	"context"
	"github.com/4everland/ipfs-top/api/routing"
	"github.com/4everland/ipfs-top/app/provide/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var ProviderSet = wire.NewSet(
	NewGormClient,
	NewDataSetRepo,
	NewNodes,
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

type NamedRoutingClient struct {
	Name   string
	Client routing.RoutingClient
}

func NewNodes(config *conf.Data) []*NamedRoutingClient {
	nodes := make([]*NamedRoutingClient, 0)
	for _, endpoint := range config.Target {
		if endpoint.ServerType != conf.ServerType_GRPC {
			continue
		}
		tlsOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		conn, err := grpc2.Dial(
			context.Background(),
			grpc2.WithEndpoint(endpoint.Endpoint),
			grpc2.WithMiddleware(
				tracing.Client(),
				recovery.Recovery(),
			),
			grpc2.WithTimeout(time.Minute),
			grpc2.WithOptions(tlsOption),
		)
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, &NamedRoutingClient{
			Name:   endpoint.Endpoint,
			Client: routing.NewRoutingClient(conn),
		})
	}

	return nodes
}
