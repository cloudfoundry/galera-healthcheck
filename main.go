package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	_ "github.com/go-sql-driver/mysql"

	"github.com/cloudfoundry-incubator/galera-healthcheck/api"
	"github.com/cloudfoundry-incubator/galera-healthcheck/config"
	"github.com/cloudfoundry-incubator/galera-healthcheck/healthcheck"
	"github.com/cloudfoundry-incubator/galera-healthcheck/monit_client"
	"github.com/cloudfoundry-incubator/galera-healthcheck/mysqld_cmd"
	"github.com/cloudfoundry-incubator/galera-healthcheck/sequence_number"
)

func main() {
	rootConfig, err := config.NewConfig(os.Args)

	logger := rootConfig.Logger

	err = rootConfig.Validate()
	if err != nil {
		logger.Fatal("Failed to validate config", err)
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			rootConfig.DB.User,
			rootConfig.DB.Password,
			rootConfig.DB.Host,
			rootConfig.DB.Port))

	if err != nil {
		logger.Fatal("db-initialize", err, lager.Data{
			"dbHost": rootConfig.DB.Host,
			"dbPort": rootConfig.DB.Port,
			"dbUser": rootConfig.DB.User,
		})
	} else {
		logger.Info("db-initialize", lager.Data{
			"dbHost": rootConfig.DB.Host,
			"dbPort": rootConfig.DB.Port,
			"dbUser": rootConfig.DB.User,
		})
	}

	mysqldCmd := mysqld_cmd.NewMysqldCmd(logger, *rootConfig)
	monitClient := monit_client.New(rootConfig.Monit, logger)
	healthchecker := healthcheck.New(db, *rootConfig, logger)
	sequenceNumberchecker := sequence_number.New(db, mysqldCmd, *rootConfig, logger)
	stateSnapshotter := &healthcheck.DBStateSnapshotter{
		DB:     db,
		Logger: logger,
	}

	router, err := api.NewRouter(
		logger,
		rootConfig,
		monitClient,
		sequenceNumberchecker,
		healthchecker,
		healthchecker,
		stateSnapshotter,
	)
	if err != nil {
		logger.Fatal("Failed to create router", err)
	}

	address := fmt.Sprintf("%s:%d", rootConfig.Host, rootConfig.Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("tcp-listen", err, lager.Data{
			"address": address,
		})
	}

	url := fmt.Sprintf("http://%s/", address)
	logger.Info("Serving healthcheck endpoint", lager.Data{
		"url": url,
	})

	if err := http.Serve(l, router); err != nil {
		logger.Fatal("http-server", err)
	}
	logger.Info("graceful-exit")
}
