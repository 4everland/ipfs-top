package main

import (
	"flag"
	"github.com/4everland/ipfs-top/app/provide/internal/conf"
	"github.com/4everland/ipfs-top/app/provide/internal/server"
	"github.com/4everland/ipfs-top/third_party/pprofx"
	"github.com/go-kratos/kratos/v2/transport"
	"os"

	"github.com/4everland/golog"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	_ "go.uber.org/automaxprocs"
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
	Name = "ipfs.provide"
	Version = "0.0.1"
	flag.StringVar(&flagconf, "conf", "configs", "config path, eg: -conf config.yaml")
	// debug, info, warn, error, fatal
	flag.StringVar(&loglevel, "level", "info", "log level,  eg: -level info")
}

func newApp(logger log.Logger, es *server.EventServer, rs *server.ReproviderServer) *kratos.App {
	var servers = []transport.Server{rs}
	if es != nil {
		servers = append(servers, es)
	}
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		pprofx.Server(servers...),
	)
}

func main() {
	flag.Parse()
	logger := golog.NewFormatStdLogger(os.Stdout, golog.WithFilterLevel(log.LevelInfo), golog.WithServerName(Name, Version))
	if err := golog.InitOTLPTracer(Name+":"+Version, golog.RatioFromEnv()); err != nil {
		log.NewHelper(logger).Warnf("InitOTLPTracer err: %s", err)
	}

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
	app, cleanup, err := wireApp(bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}

}
