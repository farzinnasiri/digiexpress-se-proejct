package dlocator

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/digiexpress/dlocator/internal/pkg/api"
	"github.com/digiexpress/dlocator/internal/pkg/config"
	"github.com/digiexpress/dlocator/internal/pkg/log"
	"github.com/sirupsen/logrus"
)

const (
	AppServerCount  = 1
	DefaultLogLevel = logrus.WarnLevel
)

func makeGracefulChannel() chan os.Signal {
	gracefulStop := make(chan os.Signal, 1)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	return gracefulStop
}

type App struct {
	Config     *config.AppConfig
	grpcServer *api.DLocatorGrpcServer
}

func NewApp(config *config.AppConfig, grpcServer *api.DLocatorGrpcServer) *App {
	app := &App{
		Config:     config,
		grpcServer: grpcServer,
	}
	app.setupAppLogger()

	return app
}

func (app *App) makeCtx(appCtx context.Context) (context.Context, context.CancelFunc) {
	log.Debug("creating server context...")

	var (
		gracefulStop = makeGracefulChannel()
		ctx, cancel  = context.WithCancel(appCtx)
	)

	// start the cancellation goroutine
	go func() {
		log.Debug("listening to graceful stop")
		<-gracefulStop
		log.Debug("graceful stop has been triggered, sending cancellation signal...")
		cancel()
		log.Info("cancellation signal has been sent")
	}()

	log.Info("server context created")

	return ctx, cancel
}

func (app *App) startGrpcServer(wg *sync.WaitGroup) {
	log.Debugf("starting grpc server on port %d...", app.Config.Grpc.Port)

	go func() {
		defer func() {
			wg.Done()
			log.Info("grpc server is shutdown")
		}()

		if err := app.grpcServer.Serve(); err != nil {
			log.Panic("failed to serve grpc", err)
		}
	}()

	log.Infof("grpc server started on port %d", app.Config.Grpc.Port)
}

func (app *App) Run(appCtx context.Context) {
	serverCtx, serverCancel := app.makeCtx(appCtx)
	defer serverCancel()

	var serverWaitGroup sync.WaitGroup

	serverWaitGroup.Add(AppServerCount)
	app.startGrpcServer(&serverWaitGroup)

	<-serverCtx.Done()
	app.grpcServer.Stop()

	serverWaitGroup.Wait()
}

func (app *App) setupAppLogger() {
	app.setGlobalLoggerFormatter()
	app.setGlobalLoggerLevel()
}

func (app *App) setGlobalLoggerLevel() {
	if app.Config.Logging.Level != "" {
		level, err := logrus.ParseLevel(app.Config.Logging.Level)
		if err != nil {
			log.Error("failed to parse log level", fmt.Errorf("got level: %w", err))
			log.Infof("setting the default log level: %s", DefaultLogLevel)
			logrus.SetLevel(DefaultLogLevel)
		} else {
			logrus.SetLevel(level)
		}
	} else {
		log.Infof("setting the default log level: %s", DefaultLogLevel)
		logrus.SetLevel(DefaultLogLevel)
	}
}

func (app *App) setGlobalLoggerFormatter() {
	logrus.SetFormatter(
		&logrus.JSONFormatter{
			DisableTimestamp:  false,
			PrettyPrint:       app.Config.Logging.PrettyPrint,
			DisableHTMLEscape: true,
			DataKey:           "context",
			TimestampFormat:   time.RFC3339,
			CallerPrettyfier:  nil,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "msg",
				logrus.FieldKeyFunc:  "caller",
			},
		},
	)
}
