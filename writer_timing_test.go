package quicktemplate

import (
	"testing"
)

func BenchmarkQWriterVString(b *testing.B) {
	v := createTestS(100)
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.V(v)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterVInt(b *testing.B) {
	v := 1233455
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.V(v)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterQ1(b *testing.B) {
	benchmarkQWriterQ(b, 1)
}

func BenchmarkQWriterQ10(b *testing.B) {
	benchmarkQWriterQ(b, 10)
}

func BenchmarkQWriterQ100(b *testing.B) {
	benchmarkQWriterQ(b, 100)
}

func BenchmarkQWriterQ1K(b *testing.B) {
	benchmarkQWriterQ(b, 1000)
}

func BenchmarkQWriterQ10K(b *testing.B) {
	benchmarkQWriterQ(b, 10000)
}

func benchmarkQWriterQ(b *testing.B, size int) {
	s := createTestS(size)
	b.SetBytes(int64(size))
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.Q(s)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterJ1(b *testing.B) {
	benchmarkQWriterJ(b, 1)
}

func BenchmarkQWriterJ10(b *testing.B) {
	benchmarkQWriterJ(b, 10)
}

func BenchmarkQWriterJ100(b *testing.B) {
	benchmarkQWriterJ(b, 100)
}

func BenchmarkQWriterJ1K(b *testing.B) {
	benchmarkQWriterJ(b, 1000)
}

func BenchmarkQWriterJ10K(b *testing.B) {
	benchmarkQWriterJ(b, 10000)
}

func benchmarkQWriterJ(b *testing.B, size int) {
	s := createTestS(size)
	b.SetBytes(int64(size))
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.J(s)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterU1(b *testing.B) {
	benchmarkQWriterU(b, 1)
}

func BenchmarkQWriterU10(b *testing.B) {
	benchmarkQWriterU(b, 10)
}

func BenchmarkQWriterU100(b *testing.B) {
	benchmarkQWriterU(b, 100)
}

func BenchmarkQWriterU1K(b *testing.B) {
	benchmarkQWriterU(b, 1000)
}

func BenchmarkQWriterU10K(b *testing.B) {
	benchmarkQWriterU(b, 10000)
}

func benchmarkQWriterU(b *testing.B, size int) {
	s := createTestS(size)
	b.SetBytes(int64(size))
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.U(s)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterFfloat(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		f := 123.456
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.F(f)
			bb.Reset()
			f++
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterFint(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		f := 123.0
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.F(f)
			bb.Reset()
			f++
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterD(b *testing.B) {
	n := 123456
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.D(n)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func BenchmarkQWriterZ1(b *testing.B) {
	benchmarkQWriterZ(b, 1)
}

func BenchmarkQWriterZ10(b *testing.B) {
	benchmarkQWriterZ(b, 10)
}

func BenchmarkQWriterZ100(b *testing.B) {
	benchmarkQWriterZ(b, 100)
}

func BenchmarkQWriterZ1K(b *testing.B) {
	benchmarkQWriterZ(b, 1000)
}

func BenchmarkQWriterZ10K(b *testing.B) {
	benchmarkQWriterZ(b, 10000)
}

func BenchmarkQWriterS1(b *testing.B) {
	benchmarkQWriterS(b, 1)
}

func BenchmarkQWriterS10(b *testing.B) {
	benchmarkQWriterS(b, 10)
}

func BenchmarkQWriterS100(b *testing.B) {
	benchmarkQWriterS(b, 100)
}

func BenchmarkQWriterS1K(b *testing.B) {
	benchmarkQWriterS(b, 1000)
}

func BenchmarkQWriterS10K(b *testing.B) {
	benchmarkQWriterS(b, 10000)
}

func benchmarkQWriterZ(b *testing.B, size int) {
	z := createTestZ(size)
	b.SetBytes(int64(size))
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.Z(z)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func benchmarkQWriterS(b *testing.B, size int) {
	s := createTestS(size)
	b.SetBytes(int64(size))
	b.RunParallel(func(pb *testing.PB) {
		var w QWriter
		bb := AcquireByteBuffer()
		w.w = bb
		for pb.Next() {
			w.S(s)
			bb.Reset()
		}
		ReleaseByteBuffer(bb)
	})
}

func createTestS(size int) string {
	return string(createTestZ(size))
}

var sample = []byte(`0123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM`)

func createTestZ(size int) []byte {
	var b []byte
	for i := 0; i < size; i++ {
		b = append(b, sample[i%len(sample)])
	}
	return b
}
