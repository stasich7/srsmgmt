package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"srsmgmt/config"
	"srsmgmt/internal/srsmgmt"
	"srsmgmt/internal/srsmgmtrepo"
	pb "srsmgmt/pb"
	"srsmgmt/pkg/playlist"
	"srsmgmt/pkg/srsclient"
	"srsmgmt/pkg/srsconfig"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.GetConfig()

	logLevel := level.AllowError()
	if cfg.IsDebug || true {
		logLevel = level.AllowDebug()
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.TimestampFormat(time.Now, time.RFC3339))
		logger = log.With(logger, "caller", log.Caller(5))
		logger = level.NewFilter(logger, logLevel)
	}

	repo := srsmgmtrepo.New(logger, 0, cfg)
	plist := playlist.New()

	var srsClient srsclient.SrsClient
	{
		srsClient = srsclient.New(cfg.SRSAddr, logger)
		srsClient = srsclient.LoggingMiddleware(logger)(srsClient)
	}

	srsConfigClient := srsClient.(srsconfig.ConfigReloader)
	srsConfig := srsconfig.New(cfg.SRSConfPath, cfg.TplStorage, srsConfigClient)

	var s srsmgmt.Service
	{
		s = srsmgmt.NewSrsMgmtService(repo, logger, plist, srsClient, srsConfig)
		s = srsmgmt.LoggingMiddleware(logger)(s)
	}

	var (
		httpHandler = srsmgmt.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
		grpcServer  = srsmgmt.MakeGRPCHandler(s, log.With(logger, "component", "GRPC"))
	)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", cfg.HTTPAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", cfg.HTTPAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", cfg.GRPCAddr)
			baseServer := grpc.NewServer()
			pb.RegisterSrsMgmtServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	if cfg.IsPprof {
		go func() {
			logger.Log("pprof", "enabled on :6060")
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	logger.Log("exit", g.Run())

}
