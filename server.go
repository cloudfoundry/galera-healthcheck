package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudfoundry-incubator/galera-healthcheck/healthcheck"
	. "github.com/cloudfoundry-incubator/galera-healthcheck/logger"
	_ "github.com/go-sql-driver/mysql"
)

var serverBind = flag.String(
        "bind",
        "",
        "Specifies the IP address to bind to",
)

var serverPort = flag.Int(
	"port",
	8080,
	"Specifies the port of the healthcheck server",
)

var mysqlUser = flag.String(
	"user",
	"root",
	"Specifies the MySQL user to connect as",
)

var mysqlPassword = flag.String(
	"password",
	"",
	"Specifies the MySQL password to connect with",
)

var mysqlHost = flag.String(
        "host",
        "localhost",
        "Specifies the MySQL host to connect to",
)

var mysqlPort = flag.Int(
	"mysqlPort",
	3306,
	"Specifies the MySQL port to connect to",
)

var availableWhenDonor = flag.Bool(
	"availWhenDonor",
	true,
	"Specifies if the healthcheck allows availability when in donor state",
)

var availableWhenReadOnly = flag.Bool(
	"availWhenReadOnly",
	false,
	"Specifies if the healthcheck allows availability when in read only mode",
)

var pidfile = flag.String(
	"pidfile",
	"",
	"Location for the pidfile",
)

var connectionCutterPath = flag.String(
	"connectionCutterPath",
	"",
	"Location for the script which cuts mysql connections",
)

var healthchecker *healthcheck.Healthchecker

func handler(w http.ResponseWriter, r *http.Request) {
	result, msg := healthchecker.Check()
	if result {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	fmt.Fprintf(w, "Galera Cluster Node status: %s", msg)
	LogWithTimestamp(msg)
}

func main() {
	flag.Parse()

	err := ioutil.WriteFile(*pidfile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		panic(err)
	}

	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", *mysqlUser, *mysqlPassword, *mysqlHost, *mysqlPort))
	config := healthcheck.HealthcheckerConfig{
		*availableWhenDonor,
		*availableWhenReadOnly,
	}

	healthchecker = healthcheck.New(db, config)

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *serverBind, *serverPort), nil)
}
