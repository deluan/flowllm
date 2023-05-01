package integration_tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegrationTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
