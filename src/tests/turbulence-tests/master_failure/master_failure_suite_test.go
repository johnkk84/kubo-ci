package master_failure_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMasterFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Master failure suite")
}

var _ = BeforeSuite(func() {
	CheckRequiredEnvs([]string{
		"BOSH_DEPLOYMENT",
		"BOSH_ENVIRONMENT",
		"BOSH_CLIENT",
		"BOSH_CLIENT_SECRET",
		"BOSH_CA_CERT",
		"TURBULENCE_HOST",
		"TURBULENCE_PORT",
		"TURBULENCE_USERNAME",
		"TURBULENCE_PASSWORD",
		"TURBULENCE_CA_CERT",
	})
	director := NewDirector()
	director.EnableResurrection(true)
})
