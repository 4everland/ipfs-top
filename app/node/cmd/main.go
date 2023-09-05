package main

import (
	"flag"
	"github.com/4everland/ipfs-servers/app/node/internal/conf"
	"github.com/4everland/ipfs-servers/app/node/internal/server"
	"github.com/4everland/ipfs-servers/third_party/logx"
	"github.com/4everland/ipfs-servers/third_party/pprofx"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	_ "go.uber.org/automaxprocs"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	loglevel string

	id, _ = os.Hostname()
)

func init() {
	Name = "ipfs.dht.node"
	Version = "0.0.1"
	flag.StringVar(&flagconf, "conf", "configs", "config path, eg: -conf config.yaml")
	// debug, info, warn, error, fatal
	flag.StringVar(&loglevel, "level", "info", "log level,  eg: -level info")
}

func newApp(logger log.Logger, ns *server.NodeServer, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		pprofx.Server(ns, gs),
	)
}

func main() {
	flag.Parse()
	logger := log.With(logx.NewLogger(log.NewStdLogger(os.Stdout), logx.WithFilterLevel(log.ParseLevel(loglevel))),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}

}
