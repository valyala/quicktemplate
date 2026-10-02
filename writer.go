package quicktemplate

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

// Writer implements auxiliary writer used by quicktemplate functions.
//
// Use AcquireWriter for creating new writers.
type Writer struct {
	e QWriter
	n QWriter
}

// W returns the underlying writer passed to AcquireWriter.
func (qw *Writer) W() io.Writer {
	return qw.n.w
}

// E returns QWriter with enabled html escaping.
func (qw *Writer) E() *QWriter {
	return &qw.e
}

// N returns QWriter without html escaping.
func (qw *Writer) N() *QWriter {
	return &qw.n
}

// WritePadding writes s to w with padding without html escaping.
func (qw *Writer) WritePadding(s string, width int) {
	qw.n.WritePadding(s, width)
}

// WritePadLeft writes s to w with left padding without html escaping.
func (qw *Writer) WritePadLeft(s string, width int) {
	qw.n.WritePadLeft(s, width)
}

// WritePadRight writes s to w with right padding without html escaping.
func (qw *Writer) WritePadRight(s string, width int) {
	qw.n.WritePadRight(s, width)
}

// WritePaddingZ writes z to w with padding without html escaping.
func (qw *Writer) WritePaddingZ(z []byte, width int) {
	qw.n.WritePaddingZ(z, width)
}

// WritePadLeftZ writes z to w with left padding without html escaping.
func (qw *Writer) WritePadLeftZ(z []byte, width int) {
	qw.n.WritePadLeftZ(z, width)
}

// WritePadRightZ writes z to w with right padding without html escaping.
func (qw *Writer) WritePadRightZ(z []byte, width int) {
	qw.n.WritePadRightZ(z, width)
}

// AcquireWriter returns new writer from the pool.
//
// Return unneeded writer to the pool by calling ReleaseWriter
// in order to reduce memory allocations.
func AcquireWriter(w io.Writer) *Writer {
	v := writerPool.Get()
	if v == nil {
		qw := &Writer{}
		qw.e.w = &htmlEscapeWriter{}
		v = qw
	}
	qw := v.(*Writer)
	qw.e.w.(*htmlEscapeWriter).w = w
	qw.n.w = w
	return qw
}

// ReleaseWriter returns the writer to the pool.
//
// Do not access released writer, otherwise data races may occur.
func ReleaseWriter(qw *Writer) {
	hw := qw.e.w.(*htmlEscapeWriter)
	hw.w = nil
	qw.e.Reset()
	qw.e.w = hw

	qw.n.Reset()

	writerPool.Put(qw)
}

var writerPool sync.Pool

// QWriter is auxiliary writer used by Writer.
type QWriter struct {
	w   io.Writer
	err error
	b   []byte
}

// Write implements io.Writer.
func (w *QWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.w.Write(p)
	if err != nil {
		w.err = err
	}
	return n, err
}

// Reset resets QWriter to the original state.
func (w *QWriter) Reset() {
	w.w = nil
	w.err = nil
}

// S writes s to w.
func (w *QWriter) S(s string) {
	w.Write(unsafeStrToBytes(s))
}

// Z writes z to w.
func (w *QWriter) Z(z []byte) {
	w.Write(z)
}

// SZ is a synonym to Z.
func (w *QWriter) SZ(z []byte) {
	w.Write(z)
}

// D writes n to w.
func (w *QWriter) D(n int) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = strconv.AppendInt(bb.B, int64(n), 10)
	} else {
		w.b = strconv.AppendInt(w.b[:0], int64(n), 10)
		w.Write(w.b)
	}
}

// DL writes n to w
func (w *QWriter) DL(n int64) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = strconv.AppendInt(bb.B, n, 10)
	} else {
		w.b = strconv.AppendInt(w.b[:0], n, 10)
		w.Write(w.b)
	}
}

// DUL writes n to w
func (w *QWriter) DUL(n uint64) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = strconv.AppendUint(bb.B, n, 10)
	} else {
		w.b = strconv.AppendUint(w.b[:0], n, 10)
		w.Write(w.b)
	}
}

// F writes f to w.
func (w *QWriter) F(f float64) {
	n := int(f)
	if float64(n) == f {
		// Fast path - just int.
		w.D(n)
		return
	}

	// Slow path.
	w.FPrec(f, -1)
}

// FPrec writes f to w using the given floating point precision.
func (w *QWriter) FPrec(f float64, prec int) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = strconv.AppendFloat(bb.B, f, 'f', prec, 64)
	} else {
		w.b = strconv.AppendFloat(w.b[:0], f, 'f', prec, 64)
		w.Write(w.b)
	}
}

// Q writes quoted json-safe s to w.
func (w *QWriter) Q(s string) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = AppendJSONString(bb.B, s, true)
	} else {
		w.b = AppendJSONString(w.b[:0], s, true)
		w.Write(w.b)
	}
}

var strQuote = []byte(`"`)

// QZ writes quoted json-safe z to w.
func (w *QWriter) QZ(z []byte) {
	w.Q(unsafeBytesToStr(z))
}

// J writes json-safe s to w.
//
// Unlike Q it doesn't qoute resulting s.
func (w *QWriter) J(s string) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = AppendJSONString(bb.B, s, false)
	} else {
		w.b = AppendJSONString(w.b[:0], s, false)
		w.Write(w.b)
	}
}

// JZ writes json-safe z to w.
//
// Unlike Q it doesn't qoute resulting z.
func (w *QWriter) JZ(z []byte) {
	w.J(unsafeBytesToStr(z))
}

// V writes v to w.
func (w *QWriter) V(v interface{}) {
	fmt.Fprintf(w, "%v", v)
}

// U writes url-encoded s to w.
func (w *QWriter) U(s string) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = appendURLEncode(bb.B, s)
	} else {
		w.b = appendURLEncode(w.b[:0], s)
		w.Write(w.b)
	}
}

// UZ writes url-encoded z to w.
func (w *QWriter) UZ(z []byte) {
	w.U(unsafeBytesToStr(z))
}

var spaces = func() []byte {
	b := make([]byte, 512)
	for i := range b {
		b[i] = ' '
	}
	return b
}()

func appendSpaces(b []byte, n int) []byte {
	for n > len(spaces) {
		b = append(b, spaces...)
		n -= len(spaces)
	}
	if n > 0 {
		b = append(b, spaces[:n]...)
	}
	return b
}

func (w *QWriter) writeSpaces(n int) {
	for n > len(spaces) {
		w.Write(spaces)
		n -= len(spaces)
	}
	if n > 0 {
		w.Write(spaces[:n])
	}
}

// WritePadLeft writes s to w with spaces padded on the left to reach width.
func (w *QWriter) WritePadLeft(s string, width int) {
	if width < 0 {
		width = -width
	}
	padLen := width - len(s)
	if padLen <= 0 {
		w.S(s)
		return
	}
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = appendSpaces(bb.B, padLen)
		bb.B = append(bb.B, s...)
	} else {
		w.writeSpaces(padLen)
		w.S(s)
	}
}

// WritePadRight writes s to w with spaces padded on the right to reach width.
func (w *QWriter) WritePadRight(s string, width int) {
	if width < 0 {
		width = -width
	}
	padLen := width - len(s)
	if padLen <= 0 {
		w.S(s)
		return
	}
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = append(bb.B, s...)
		bb.B = appendSpaces(bb.B, padLen)
	} else {
		w.S(s)
		w.writeSpaces(padLen)
	}
}

// WritePadding writes s to w with padding.
// If width > 0, spaces are added to the left (right-aligned).
// If width < 0, spaces are added to the right (left-aligned).
// If width == 0, s is written without padding.
func (w *QWriter) WritePadding(s string, width int) {
	if width > 0 {
		w.WritePadLeft(s, width)
	} else if width < 0 {
		w.WritePadRight(s, -width)
	} else {
		w.S(s)
	}
}

// WritePadLeftZ writes z to w with spaces padded on the left to reach width.
func (w *QWriter) WritePadLeftZ(z []byte, width int) {
	if width < 0 {
		width = -width
	}
	padLen := width - len(z)
	if padLen <= 0 {
		w.Z(z)
		return
	}
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = appendSpaces(bb.B, padLen)
		bb.B = append(bb.B, z...)
	} else {
		w.writeSpaces(padLen)
		w.Z(z)
	}
}

// WritePadRightZ writes z to w with spaces padded on the right to reach width.
func (w *QWriter) WritePadRightZ(z []byte, width int) {
	if width < 0 {
		width = -width
	}
	padLen := width - len(z)
	if padLen <= 0 {
		w.Z(z)
		return
	}
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = append(bb.B, z...)
		bb.B = appendSpaces(bb.B, padLen)
	} else {
		w.Z(z)
		w.writeSpaces(padLen)
	}
}

// WritePaddingZ writes z to w with padding.
// If width > 0, spaces are added to the left (right-aligned).
// If width < 0, spaces are added to the right (left-aligned).
// If width == 0, z is written without padding.
func (w *QWriter) WritePaddingZ(z []byte, width int) {
	if width > 0 {
		w.WritePadLeftZ(z, width)
	} else if width < 0 {
		w.WritePadRightZ(z, -width)
	} else {
		w.Z(z)
	}
}
