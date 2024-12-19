package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mhdiiilham/dating-app/delivery/restful"
	"github.com/mhdiiilham/dating-app/pkg/common"
	"github.com/mhdiiilham/dating-app/pkg/credential"
	"github.com/mhdiiilham/dating-app/repository"
	"github.com/mhdiiilham/dating-app/usecase/authentication"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})

	config := common.ReadConfig()
	log.Infof("read configuration success. version: %s", config.Version)

	dbConn := common.ConnectDb(config)

	accessTokenDuration := 2 * time.Hour
	jwtClient := credential.NewJwtGenerator(config.AppName, accessTokenDuration, config.JWTSecret)
	passwordHasher := credential.Hasher{}

	userRepository := repository.NewUser(dbConn)

	authenticatorService := authentication.NewService(userRepository, jwtClient, passwordHasher)

	echo.NotFoundHandler = restful.NotFoundHandler()
	e := echo.New()
	apiV1 := e.Group("api/v1")
	authenticationRoutes := apiV1.Group("/authentications")
	authenticationRoutes.POST("/signup", restful.HandleUserSignUp(authenticatorService))
	authenticationRoutes.POST("/signin", restful.HandleUserSignIn(authenticatorService))

	go func() {
		if err := e.Start(config.GetServerPort()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	logrus.Info("waiting shutdown signal")
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, done := context.WithTimeout(context.Background(), 30*time.Second)
	defer done()

	log.Infof("shutting down echo server; error=%v", e.Shutdown(ctx))
	log.Info("server is shuting down...")
	log.Infof("closing db connection; %v", dbConn.Close())
}
