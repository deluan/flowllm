package vectorstores_test

import (
	. "github.com/deluan/pipelm/vectorstores"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CosineSimilarity", func() {
	It("should return 0 when both input vectors are empty", func() {
		var a []float32
		var b []float32
		Expect(CosineSimilarity(a, b)).To(Equal(float32(0)))
	})

	It("should return 0 when one of the input vectors is empty", func() {
		a := []float32{1, 2, 3}
		var b []float32
		Expect(CosineSimilarity(a, b)).To(Equal(float32(0)))
	})

	It("should return 1 when both input vectors are the same", func() {
		a := []float32{1, 2, 3}
		b := []float32{1, 2, 3}
		Expect(CosineSimilarity(a, b)).To(BeNumerically("~", float32(1), 1e-6))
	})

	It("should return 0 when input vectors are orthogonal", func() {
		a := []float32{1, 0, 0}
		b := []float32{0, 1, 0}
		Expect(CosineSimilarity(a, b)).To(BeNumerically("~", float32(0), 1e-6))
	})

	It("should return -1 when input vectors are opposite", func() {
		a := []float32{1, 2, 3}
		b := []float32{-1, -2, -3}
		Expect(CosineSimilarity(a, b)).To(BeNumerically("~", float32(-1), 1e-6))
	})

	It("should return the correct similarity for arbitrary vectors", func() {
		a := []float32{1, 0, -1}
		b := []float32{2, 2, 0}
		expected := 0.500
		Expect(CosineSimilarity(a, b)).To(BeNumerically("~", expected, 1e-6))
	})

	It("should handle vectors with different lengths", func() {
		a := []float32{1, 2, 3, 4, 5}
		b := []float32{1, 2, 3}
		Expect(CosineSimilarity(a, b)).To(BeNumerically("~", float32(1), 1e-6))
	})

	It("should handle vectors with large values", func() {
		a := []float32{1e5, 2e5, 3e5}
		b := []float32{1e5, 2e5, 3e5}
		Expect(CosineSimilarity(a, b)).To(BeNumerically("~", float32(1), 1e-6))
	})
})
