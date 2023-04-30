package splitters_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSplitters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Splitters Tests Suite")
}
