package monit_status_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMonitStatus(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Monit Status Suite")
}
