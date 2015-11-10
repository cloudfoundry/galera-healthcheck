package monit_status_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry-incubator/galera-healthcheck/monit_status"
	"github.com/cloudfoundry-incubator/galera-healthcheck/monit_wrapper/fakes"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Monit Status", func() {

	var (
		fakeMonit *fakes.FakeMonitWrapper
	)

	BeforeEach(func() {
		fakeMonit = &fakes.FakeMonitWrapper{}
	})

	It("database is Running", func(){
		fakeMonit.SummaryReturns(`
The Monit daemon 5.2.4 uptime: 4d 16h 29m

Process 'mariadb_ctrl'              running
Process 'galera-healthcheck'        running
Process 'gra-log-purger-executable' running
System 'system_435c8a75-abf5-4068-a840-9492c27890bb' running
		`)
		statusChecker := monit_status.New(fakeMonit)

		r, err := http.NewRequest("GET", "/monit_status", nil)
		Expect(err).ToNot(HaveOccurred())

		w := httptest.NewRecorder()
		statusChecker.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(Equal("running"))
	})

//	It("database is Starting", func(){
//		fakeMonit.SummaryReturns(`
//The Monit daemon 5.2.4 uptime: 4d 16h 29m
//
//Process 'mariadb_ctrl'              not monitored - start pending
//Process 'galera-healthcheck'        running
//Process 'gra-log-purger-executable' running
//System 'system_435c8a75-abf5-4068-a840-9492c27890bb' running
//		`)
//		statusChecker := monit_status.New(fakeMonit)
//		Expect(statusChecker.Status()).To(Equal("starting"))
//	})
//
//	It("database is Stopping", func(){
//		fakeMonit.SummaryReturns(`
//The Monit daemon 5.2.4 uptime: 4d 16h 29m
//
//Process 'mariadb_ctrl'              not monitored - stop pending
//Process 'galera-healthcheck'        running
//Process 'gra-log-purger-executable' running
//System 'system_435c8a75-abf5-4068-a840-9492c27890bb' running
//		`)
//		statusChecker := monit_status.New(fakeMonit)
//
//		Expect(statusChecker.Status()).To(Equal("stopping"))
//	})
//
//	It("database is failing", func(){
//		fakeMonit.SummaryReturns(`
//The Monit daemon 5.2.4 uptime: 4d 16h 29m
//
//Process 'mariadb_ctrl'              execution failed
//Process 'galera-healthcheck'        running
//Process 'gra-log-purger-executable' running
//System 'system_435c8a75-abf5-4068-a840-9492c27890bb' running
//		`)
//		statusChecker := monit_status.New(fakeMonit)
//
//		Expect(statusChecker.Status()).To(Equal("failing"))
//	})
//
	It("returns an error for unknown states", func() {
		fakeMonit.SummaryReturns(`
The Monit daemon 5.2.4 uptime: 4d 16h 29m

Process 'mariadb_ctrl'              NO IDEA
Process 'galera-healthcheck'        running
Process 'gra-log-purger-executable' running
System 'system_435c8a75-abf5-4068-a840-9492c27890bb' running
		`)
		statusChecker := monit_status.New(fakeMonit)

		r, err := http.NewRequest("GET", "/monit_status", nil)
		Expect(err).ToNot(HaveOccurred())

		w := httptest.NewRecorder()
		statusChecker.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(http.StatusInternalServerError))
		Expect(w.Body.String()).To(Equal("unknown state NO IDEA"))
	})
})
