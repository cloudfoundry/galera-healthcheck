package monit_status
import (
	"github.com/cloudfoundry-incubator/galera-healthcheck/monit_wrapper"
	"regexp"
"errors"
	"fmt"
	"net/http"
)

type statusChecker struct {
	monitObject monit_wrapper.MonitWrapper
}

func New(monitObject monit_wrapper.MonitWrapper) *statusChecker {
	return &statusChecker{
		monitObject: monitObject,
	}
}


func (s *statusChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := s.status()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(status))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}


func (s *statusChecker) status() (string, error) {
	rawStatus := regexp.MustCompile(`mariadb_ctrl'\s*(.*)`).FindStringSubmatch(s.monitObject.Summary())[1]
	switch rawStatus {
	case "running":
		return "running", nil

	case "not monitored - stop pending":
		return "stopping", nil

	case "not monitored - start pending":
		return "starting", nil

	case "not monitored":
		return "stopped", nil

	case "execution failed":
		return "failing", nil
	}
	return "", errors.New(fmt.Sprintf("unknown state %s", rawStatus))
}

