package generator

import "testing"

func BenchmarkGenerateSmallLengthKeys(b *testing.B) {
	const N = 10

	gen := NewGenerator(N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.Generate()
	}

	b.ReportAllocs()
}

func BenchmarkGenerateBigLengthKeys(b *testing.B) {
	const N = 1000

	gen := NewGenerator(N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.Generate()
	}

	b.ReportAllocs()
}
