package vectorstores_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVectorStores(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VectorStores Suite")
}
