package quicktemplate

import (
	"testing"
)

func TestWriter(t *testing.T) {
	bb := AcquireByteBuffer()
	qw := AcquireWriter(bb)
	w := qw.W()
	bbNew, ok := w.(*ByteBuffer)
	if !ok {
		t.Fatalf("W() must return ByteBuffer, not %T", w)
	}
	if bbNew != bb {
		t.Fatalf("unexpected ByteBuffer returned: %p. Expecting %p", bbNew, bb)
	}

	wn := qw.N()
	we := qw.E()

	wn.S("<a></a>")
	wn.D(123)
	we.DUL(18446744073709551615)
	wn.Z([]byte("'"))
	wn.Q("foo")
	wn.J("ds")
	wn.F(1.23)
	wn.U("абв")
	wn.V(struct{}{})
	wn.SZ([]byte("aaa"))
	wn.QZ([]byte("asadf"))
	wn.JZ([]byte("asd"))
	wn.UZ([]byte("abc"))

	we.S("<a></a>")
	we.D(321)
	we.DUL(18446744073709551615)
	we.Z([]byte("'"))
	we.Q("foo")
	we.J("ds")
	we.F(1.23)
	we.U("абв")
	we.V(struct{}{})
	we.SZ([]byte("aaa"))
	we.QZ([]byte("asadf"))
	we.JZ([]byte("asd"))
	we.UZ([]byte("abc"))

	ReleaseWriter(qw)

	expectedS := "<a></a>12318446744073709551615'\"foo\"ds1.23%D0%B0%D0%B1%D0%B2{}aaa\"asadf\"asdabc" +
		"&lt;a&gt;&lt;/a&gt;32118446744073709551615&#39;&quot;foo&quot;ds1.23%D0%B0%D0%B1%D0%B2{}aaa&quot;asadf&quot;asdabc"
	if string(bb.B) != expectedS {
		t.Fatalf("unexpected output:\n%q\nExpecting\n%q", bb.B, expectedS)
	}

	ReleaseByteBuffer(bb)
}

func TestQWriterS(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "\x00foo<>&'\" bar\n\t</script>=;\\/+%йцу\x00foo&lt;&gt;&amp;&#39;&quot; bar\n\t&lt;/script&gt;=;\\/+%йцу"
		wn.S(s)
		we.S(s)
		return expectedS
	})
}

func TestQWriterZ(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "\x00foo<>&'\" bar\n\t</script>=;\\/+%йцу\x00foo&lt;&gt;&amp;&#39;&quot; bar\n\t&lt;/script&gt;=;\\/+%йцу"
		wn.Z([]byte(s))
		we.Z([]byte(s))
		return expectedS
	})
}

func TestQWriterSZ(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "\x00foo<>&'\" bar\n\t</script>=;\\/+%йцу\x00foo&lt;&gt;&amp;&#39;&quot; bar\n\t&lt;/script&gt;=;\\/+%йцу"
		wn.SZ([]byte(s))
		we.SZ([]byte(s))
		return expectedS
	})
}

func TestQWriterQ(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "\"\\u0000foo\\u003c>&\\u0027\\\" bar\\n\\t\\u003c/script>=;\\\\/+%йцу\"&quot;\\u0000foo\\u003c&gt;&amp;\\u0027\\&quot; bar\\n\\t\\u003c/script&gt;=;\\\\/+%йцу&quot;"
		wn.Q(s)
		we.Q(s)
		return expectedS
	})
}

func TestQWriterQZ(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "\"\\u0000foo\\u003c>&\\u0027\\\" bar\\n\\t\\u003c/script>=;\\\\/+%йцу\"&quot;\\u0000foo\\u003c&gt;&amp;\\u0027\\&quot; bar\\n\\t\\u003c/script&gt;=;\\\\/+%йцу&quot;"
		wn.QZ([]byte(s))
		we.QZ([]byte(s))
		return expectedS
	})
}

func TestQWriterJ(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "\\u0000foo\\u003c>&\\u0027\\\" bar\\n\\t\\u003c/script>=;\\\\/+%йцу\\u0000foo\\u003c&gt;&amp;\\u0027\\&quot; bar\\n\\t\\u003c/script&gt;=;\\\\/+%йцу"
		wn.J(s)
		we.J(s)
		return expectedS
	})
}

func TestQWriterJZ(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "\\u0000foo\\u003c>&\\u0027\\\" bar\\n\\t\\u003c/script>=;\\\\/+%йцу\\u0000foo\\u003c&gt;&amp;\\u0027\\&quot; bar\\n\\t\\u003c/script&gt;=;\\\\/+%йцу"
		wn.JZ([]byte(s))
		we.JZ([]byte(s))
		return expectedS
	})
}

func TestQWriterU(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "%00foo%3C%3E%26%27%22+bar%0A%09%3C%2Fscript%3E%3D%3B%5C%2F%2B%25%D0%B9%D1%86%D1%83%00foo%3C%3E%26%27%22+bar%0A%09%3C%2Fscript%3E%3D%3B%5C%2F%2B%25%D0%B9%D1%86%D1%83"
		wn.U(s)
		we.U(s)
		return expectedS
	})
}

func TestQWriterUZ(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "%00foo%3C%3E%26%27%22+bar%0A%09%3C%2Fscript%3E%3D%3B%5C%2F%2B%25%D0%B9%D1%86%D1%83%00foo%3C%3E%26%27%22+bar%0A%09%3C%2Fscript%3E%3D%3B%5C%2F%2B%25%D0%B9%D1%86%D1%83"
		wn.UZ([]byte(s))
		we.UZ([]byte(s))
		return expectedS
	})
}

func TestQWriterV(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		s := "\u0000" + `foo<>&'" bar
	</script>=;\/+%йцу`
		expectedS := "{\x00foo<>&'\" bar\n\t</script>=;\\/+%йцу}{\x00foo&lt;&gt;&amp;&#39;&quot; bar\n\t&lt;/script&gt;=;\\/+%йцу}"
		ss := struct{ S string }{s}
		wn.V(ss)
		we.V(ss)
		return expectedS
	})
}

func TestQWriterF(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		f := 1.9234
		wn.F(f)
		we.F(f)
		return "1.92341.9234"
	})
}

func TestQWriterFPrec(t *testing.T) {
	testQWriter(t, func(wn, we *QWriter) string {
		f := 1.9254
		wn.FPrec(f, 2)
		we.FPrec(f, 3)
		wn.FPrec(f, 0)
		we.FPrec(f, 1)
		return "1.931.92521.9"
	})
}

func testQWriter(t *testing.T, f func(wn, we *QWriter) (expectedS string)) {
	bb := AcquireByteBuffer()
	qw := AcquireWriter(bb)
	wn := qw.N()
	we := qw.E()

	expectedS := f(wn, we)

	ReleaseWriter(qw)

	if string(bb.B) != expectedS {
		t.Fatalf("unexpected output:\n%q\nExpecting\n%q", bb.B, expectedS)
	}

	ReleaseByteBuffer(bb)
}
