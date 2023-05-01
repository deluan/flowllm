package pipelm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPipeLM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PipeLM Tests Suite")
}
