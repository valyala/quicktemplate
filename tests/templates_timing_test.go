package tests

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"testing"
	texttemplate "text/template"

	"github.com/valyala/quicktemplate"
	"github.com/valyala/quicktemplate/testdata/templates"
)

var tpl = template.Must(template.ParseFiles("../testdata/templates/bench.tpl"))
var textTpl = texttemplate.Must(texttemplate.ParseFiles("../testdata/templates/bench.tpl"))

func init() {
	// make sure that both template engines generate the same result
	rows := getBenchRows(3)

	bb1 := &quicktemplate.ByteBuffer{}
	if err := tpl.Execute(bb1, rows); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	bb2 := &quicktemplate.ByteBuffer{}
	templates.WriteBenchPage(bb2, rows)

	if !bytes.Equal(bb1.B, bb2.B) {
		log.Fatalf("results mismatch:\n%q\n%q", bb1, bb2)
	}
}

func BenchmarkQuickTemplate1(b *testing.B) {
	benchmarkQuickTemplate(b, 1)
}

func BenchmarkQuickTemplate10(b *testing.B) {
	benchmarkQuickTemplate(b, 10)
}

func BenchmarkQuickTemplate100(b *testing.B) {
	benchmarkQuickTemplate(b, 100)
}

func benchmarkQuickTemplate(b *testing.B, rowsCount int) {
	rows := getBenchRows(rowsCount)
	b.RunParallel(func(pb *testing.PB) {
		bb := quicktemplate.AcquireByteBuffer()
		for pb.Next() {
			templates.WriteBenchPage(bb, rows)
			bb.Reset()
		}
		quicktemplate.ReleaseByteBuffer(bb)
	})
}

func BenchmarkHTMLTemplate1(b *testing.B) {
	benchmarkHTMLTemplate(b, 1)
}

func BenchmarkHTMLTemplate10(b *testing.B) {
	benchmarkHTMLTemplate(b, 10)
}

func BenchmarkHTMLTemplate100(b *testing.B) {
	benchmarkHTMLTemplate(b, 100)
}

func benchmarkHTMLTemplate(b *testing.B, rowsCount int) {
	rows := getBenchRows(rowsCount)
	b.RunParallel(func(pb *testing.PB) {
		bb := quicktemplate.AcquireByteBuffer()
		for pb.Next() {
			if err := tpl.Execute(bb, rows); err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			bb.Reset()
		}
		quicktemplate.ReleaseByteBuffer(bb)
	})
}

func BenchmarkTextTemplate1(b *testing.B) {
	benchmarkTextTemplate(b, 1)
}

func BenchmarkTextTemplate10(b *testing.B) {
	benchmarkTextTemplate(b, 10)
}

func BenchmarkTextTemplate100(b *testing.B) {
	benchmarkTextTemplate(b, 100)
}

func benchmarkTextTemplate(b *testing.B, rowsCount int) {
	rows := getBenchRows(rowsCount)
	b.RunParallel(func(pb *testing.PB) {
		bb := quicktemplate.AcquireByteBuffer()
		for pb.Next() {
			if err := textTpl.Execute(bb, rows); err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			bb.Reset()
		}
		quicktemplate.ReleaseByteBuffer(bb)
	})
}

func getBenchRows(n int) []templates.BenchRow {
	rows := make([]templates.BenchRow, n)
	for i := 0; i < n; i++ {
		rows[i] = templates.BenchRow{
			ID:      i,
			Message: fmt.Sprintf("message %d", i),
			Print:   ((i & 1) == 0),
		}
	}
	return rows
}
