package main

import (
	"context"
	"errors"
	"imager/internal/service"
	"imager/internal/storage"
	"imager/internal/storage/kafka"
	"imager/internal/storage/postgres"
	"imager/internal/web"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	kafkaLib "github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

var (
	Port         = "80"
	Debug        = "true"
	ConfigPath   = "../config/config.yml"
	PostgresConn = "postgres://dev:qqq@localhost:5432/test?sslmode=disable"
	KafkaBrokers = "localhost:9092"
	KafkaTopic   = "pictures"
	KafkaGroupID = "my-test-group"
	Templates    = "templates/*.html"
)

func main() {
	zlog.Init()
	env()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()
	err := cfg.Load(ConfigPath)
	if err != nil {
		panic(err)
	}

	kfk := kafka.New(strings.Split(KafkaBrokers, ","), KafkaTopic, KafkaGroupID)
	pg := postgres.New(PostgresConn)
	str := storage.New(pg, kfk)
	srv := service.New(str)

	out := make(chan kafkaLib.Message)
	kfk.StartConsuming(ctx, out)
	srv.StartConsuming(ctx, out)

	router := ginext.New()
	web.SetRoutes(router, Templates, srv)

	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		panic(err)
	}

	server := web.NewServer(router, cfg)

	go func() {
		if err := server.Serve(listener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				zlog.Logger.Error().Err(err).Send()
			}
		}
	}()

	zlog.Logger.Info().Msg("server started")

	<-sig
	zlog.Logger.Info().Msg("gracefull shutdown")
	_ = server.Shutdown(context.Background())
	cancel()
	time.Sleep(time.Second * 1)
	srv.Shutdown()
	str.Shutdown()
}
