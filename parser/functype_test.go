package parser

import (
	"testing"
)

func TestParseFuncCallSuccess(t *testing.T) {
	// func without args
	testParseFuncCallSuccess(t, "f()", "streamf(qw422016)")

	// func with args
	testParseFuncCallSuccess(t, "Foo(a, b)", "StreamFoo(qw422016, a, b)")

	// method without args
	testParseFuncCallSuccess(t, "a.f()", "a.streamf(qw422016)")

	// method with args
	testParseFuncCallSuccess(t, "a.f(xx)", "a.streamf(qw422016, xx)")

	// chained method
	testParseFuncCallSuccess(t, "foo.bar.Baz(x, y)", "foo.bar.StreamBaz(qw422016, x, y)")

	// complex args
	testParseFuncCallSuccess(t, `as.ffs.SS(
		func(x int, y string) {
			panic("foobar")
		},
		map[string]int{
			"foo":1,
			"bar":2,
		},
		qawe)`,
		`as.ffs.StreamSS(qw422016, 
		func(x int, y string) {
			panic("foobar")
		},
		map[string]int{
			"foo":1,
			"bar":2,
		},
		qawe)`)
}

func TestParseFuncCallFailure(t *testing.T) {
	testParseFuncCallFailure(t, "")

	// non-func
	testParseFuncCallFailure(t, "foobar")
	testParseFuncCallFailure(t, "a, b, c")
	testParseFuncCallFailure(t, "{}")
	testParseFuncCallFailure(t, "(a)")
	testParseFuncCallFailure(t, "(f())")

	// inline func
	testParseFuncCallFailure(t, "func() {}()")
	testParseFuncCallFailure(t, "func a() {}()")

	// nonempty tail after func call
	testParseFuncCallFailure(t, "f(); f1()")
	testParseFuncCallFailure(t, "f()\nf1()")
	testParseFuncCallFailure(t, "f()\n for {}")
}

func testParseFuncCallFailure(t *testing.T, s string) {
	_, err := parseFuncCall([]byte(s))
	if err == nil {
		t.Fatalf("expecting non-nil error when parsing %q", s)
	}
}

func testParseFuncCallSuccess(t *testing.T, s, callStream string) {
	f, err := parseFuncCall([]byte(s))
	if err != nil {
		t.Fatalf("unexpected error when parsing %q: %s", s, err)
	}
	cs := f.CallStream("qw422016")
	if cs != callStream {
		t.Fatalf("unexpected CallStream: %q. Expecting %q. s=%q", cs, callStream, s)
	}
}

func TestParseFuncDefSuccess(t *testing.T) {
	// private func without args
	testParseFuncDefSuccess(t, "xx()", "xx() string",
		"streamxx(qw422016 *qt422016.Writer)", "streamxx(qw422016)",
		"writexx(qq422016 qtio422016.Writer)", "writexx(qq422016)")

	// public func with a single arg
	testParseFuncDefSuccess(t, "F(a int)", "F(a int) string",
		"StreamF(qw422016 *qt422016.Writer, a int)", "StreamF(qw422016, a)",
		"WriteF(qq422016 qtio422016.Writer, a int)", "WriteF(qq422016, a)")

	// public method without args
	testParseFuncDefSuccess(t, "(f *foo) M()", "(f *foo) M() string",
		"(f *foo) StreamM(qw422016 *qt422016.Writer)", "f.StreamM(qw422016)",
		"(f *foo) WriteM(qq422016 qtio422016.Writer)", "f.WriteM(qq422016)")

	// private method with three args
	testParseFuncDefSuccess(t, "(f *Foo) bar(x, y string, z int)", "(f *Foo) bar(x, y string, z int) string",
		"(f *Foo) streambar(qw422016 *qt422016.Writer, x, y string, z int)", "f.streambar(qw422016, x, y, z)",
		"(f *Foo) writebar(qq422016 qtio422016.Writer, x, y string, z int)", "f.writebar(qq422016, x, y, z)")

	// method with complex args
	testParseFuncDefSuccess(t, "(t TPL) Head(h1, h2 func(x, y int), h3 map[int]struct{})", "(t TPL) Head(h1, h2 func(x, y int), h3 map[int]struct{}) string",
		"(t TPL) StreamHead(qw422016 *qt422016.Writer, h1, h2 func(x, y int), h3 map[int]struct{})", "t.StreamHead(qw422016, h1, h2, h3)",
		"(t TPL) WriteHead(qq422016 qtio422016.Writer, h1, h2 func(x, y int), h3 map[int]struct{})", "t.WriteHead(qq422016, h1, h2, h3)")

	// method with variadic arguments
	testParseFuncDefSuccess(t, "(t TPL) Head(name string, num int, otherNames ...string)", "(t TPL) Head(name string, num int, otherNames ...string) string",
		"(t TPL) StreamHead(qw422016 *qt422016.Writer, name string, num int, otherNames ...string)", "t.StreamHead(qw422016, name, num, otherNames...)",
		"(t TPL) WriteHead(qq422016 qtio422016.Writer, name string, num int, otherNames ...string)", "t.WriteHead(qq422016, name, num, otherNames...)")
}

func TestParseFuncDefFailure(t *testing.T) {
	testParseFuncDefFailure(t, "")

	// invalid syntax
	testParseFuncDefFailure(t, "foobar")
	testParseFuncDefFailure(t, "f() {")
	testParseFuncDefFailure(t, "for {}")

	// missing func name
	testParseFuncDefFailure(t, "()")
	testParseFuncDefFailure(t, "(a int, b string)")

	// missing method name
	testParseFuncDefFailure(t, "(x XX) ()")
	testParseFuncDefFailure(t, "(x XX) (y, z string)")

	// func with return values
	testParseFuncDefFailure(t, "f() string")
	testParseFuncDefFailure(t, "f() (int, string)")
	testParseFuncDefFailure(t, "(x XX) f() string")
	testParseFuncDefFailure(t, "(x XX) f(a int) (int, string)")
}

func testParseFuncDefFailure(t *testing.T, s string) {
	f, err := parseFuncDef([]byte(s))
	if err == nil {
		t.Fatalf("expecting error when parsing %q. got %#v", s, f)
	}
}

func testParseFuncDefSuccess(t *testing.T, s, defString, defStream, callStream, defWrite, callWrite string) {
	f, err := parseFuncDef([]byte(s))
	if err != nil {
		t.Fatalf("cannot parse %q: %s", s, err)
	}
	ds := f.DefString()
	if ds != defString {
		t.Fatalf("unexpected DefString: %q. Expecting %q. s=%q", ds, defString, s)
	}
	ds = f.DefStream("qw422016")
	if ds != defStream {
		t.Fatalf("unexpected DefStream: %q. Expecting %q. s=%q", ds, defStream, s)
	}
	cs := f.CallStream("qw422016")
	if cs != callStream {
		t.Fatalf("unexpected CallStream: %q. Expecting %q. s=%q", cs, callStream, s)
	}
	dw := f.DefWrite("qq422016")
	if dw != defWrite {
		t.Fatalf("unexpected DefWrite: %q. Expecting %q. s=%q", dw, defWrite, s)
	}
	cw := f.CallWrite("qq422016")
	if cw != callWrite {
		t.Fatalf("unexpected CallWrite: %q. Expecting %q. s=%q", cw, callWrite, s)
	}
}
