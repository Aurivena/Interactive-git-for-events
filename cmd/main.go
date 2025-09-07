package main

import (
	"arch/internal/delivery/http"
	"arch/internal/domain/entity"
	"arch/internal/initialization"
	"arch/internal/server"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.Info("start init server")
	if err := initialization.LoadConfiguration(); err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info("end init server")
}

// @title           Arch
// @version         1.0.0
// @description     архитектура

// @host      		localhost:1941
// @BasePath  		/api/
func main() {
	var serverInstance server.Server
	routes, businessDatabase := initialization.InitLayers()
	go run(serverInstance, routes, &initialization.ConfigService.Server)
	stop()
	serverInstance.Stop(context.Background(), businessDatabase)
}

func run(server server.Server, routes *http.Http, config *entity.ServerConfig) {
	ginEngine := routes.InitHTTPHttps(config)
	certificates := initialization.ConfigService.Certificates

	if err := server.Run(config.Port, ginEngine, certificates); err != nil {
		if err.Error() != "http: Server closed" {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}
}

func stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	<-quit
}
