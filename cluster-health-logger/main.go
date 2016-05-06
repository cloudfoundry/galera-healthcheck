package main

import (
	"flag"
	"github.com/pivotal-cf-experimental/service-config"
	"os"
	"github.com/cloudfoundry-incubator/galera-healthcheck/cluster-health-logger/logwriter"
	"log"
	"database/sql"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/validator.v2"
)

func main() {
	var config logwriter.Config
	serviceConfig := service_config.New()

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	serviceConfig.AddFlags(flags)
	flags.Parse(os.Args[1:])
	err := serviceConfig.Read(&config)

	if err != nil {
		log.Fatal("Failed to read config", err)
	}

	err = validator.Validate(config)
	if err != nil {
		log.Fatal("Failed to validate config", err)
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			config.User,
			config.Password,
			"localhost",
			config.Port))

	if err != nil {
		log.Fatal("Failed to connect to mysql", err)
	}

	writer := logwriter.New(db, config.LogPath)

	for {
		writer.Write(time.Now().Format(time.RFC3339))
		time.Sleep(time.Duration(config.Interval)*time.Second)
	}
}