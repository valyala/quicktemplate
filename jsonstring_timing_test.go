package quicktemplate

import (
	"strconv"
	"testing"
)

func BenchmarkAppendJSONString(b *testing.B) {
	b.Run("no-special-chars", func(b *testing.B) {
		benchmarkAppendJSONString(b, "foo bar baz abc defkjlkj lkjdfs klsdjflfdjoqjo lkj ss")
	})
	b.Run("with-special-chars", func(b *testing.B) {
		benchmarkAppendJSONString(b, `foo bar baz abc defkjlkj lkjdf" klsdjflfdjoqjo\lkj ss`)
	})
}

func benchmarkAppendJSONString(b *testing.B, s string) {
	b.ReportAllocs()
	b.SetBytes(int64(len(s)))
	b.RunParallel(func(pb *testing.PB) {
		var buf []byte
		for pb.Next() {
			buf = AppendJSONString(buf[:0], s, true)
		}
	})
}

func BenchmarkAppendJSONStringViaStrconv(b *testing.B) {
	b.Run("no-special-chars", func(b *testing.B) {
		benchmarkAppendJSONStringViaStrconv(b, "foo bar baz abc defkjlkj lkjdfs klsdjflfdjoqjo lkj ss")
	})
	b.Run("with-special-chars", func(b *testing.B) {
		benchmarkAppendJSONStringViaStrconv(b, `foo bar baz abc defkjlkj lkjdf" klsdjflfdjoqjo\lkj ss`)
	})
}

func benchmarkAppendJSONStringViaStrconv(b *testing.B, s string) {
	b.ReportAllocs()
	b.SetBytes(int64(len(s)))
	b.RunParallel(func(pb *testing.PB) {
		var buf []byte
		for pb.Next() {
			buf = strconv.AppendQuote(buf[:0], s)
		}
	})
}
