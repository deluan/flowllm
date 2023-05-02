package flowllm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFlowLLM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FlowLLM Tests Suite")
}
