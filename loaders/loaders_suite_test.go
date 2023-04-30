package loaders_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLoaders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Loaders Suite")
}
