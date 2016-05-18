package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"testing"

	"github.com/valyala/quicktemplate"
)

func TestParseSwitchCaseSuccess(t *testing.T) {
	// single-case switch
	testParseSuccess(t, "{%func a()%}{%switch n%}{%case 1%}aaa{%endswitch%}{%endfunc%}")

	// multi-case switch
	testParseSuccess(t, "{%func a()%}{%switch%}\n\t  {%case foo()%}\nfoobar{%case bar()%}{%endswitch%}{%endfunc%}")

	// default statement
	testParseSuccess(t, "{%func a()%}{%switch%}{%default%}{%endswitch%}{%endfunc%}")

	// complex switch
	testParseSuccess(t, `{%func f()%}{% for %}
		{%switch foo() %}
		The text before the first case
		is converted into comment
		{%case "foobar" %}
			{% switch %}
			{% case bar() %}
				aaaa
			{% case baz() %}
				bbbb
			{% endswitch %}
		{% case "aaa" %}
			{% for i := 0; i < 10; i++ %}
				foobar
			{% endfor %}
		{% case "qwe" %}
			aaaa
			{% return %}
		{% case "www" %}
			break from the outer loop, not from the switch
			{% break %}
		{% default %}
			foobar
		{%endswitch%}
	{% endfor %}{%endfunc%}`)
}

func TestParseSwitchCaseFailure(t *testing.T) {
	// missing endswitch
	testParseFailure(t, "{%func a()%}{%switch%}{%endfunc%}")

	// empty switch
	testParseFailure(t, "{%func f()%}{%switch%}{%endswitch%}{%endfunc%}")

	// case outside switch
	testParseFailure(t, "{%func f()%}{%case%}{%endfunc%}")

	// the first tag inside switch is non-case
	testParseFailure(t, "{%func f()%}{%switch%}{%return%}{%endswitch%}{%endfunc%}")

	// break inside switch without outer loop
	testParseFailure(t, "{%func f()%}{%switch%}{%case 1%}{%break%}{%endswitch%}{%endfunc%}")

	// empty case
	testParseFailure(t, "{%func f()%}{%switch%}{%case%}aaa{%endswitch%}{%endfunc%}")

	// multiple default statements
	testParseFailure(t, "{%func f()%}{%switch%}{%case%}aaa{%default%}bbb{%default%}{%endswitch%}{%endfunc%}")
}

func TestParseBreakContinueReturn(t *testing.T) {
	testParseSuccess(t, `{% func a() %}{% for %}{% continue %}{% break %}{% return %}{% endfor %}{% endfunc %}`)
	testParseSuccess(t, `{% func a() %}{% for %}
		{% if f1() %}{% continue %}skip this{%s "and this" %}{% endif %}
		{% if f2() %}{% break %}{% for %}{% endfor %}skip this{% endif %}
		{% if f3() %}{% return %}foo{% if f4() %}{% for %}noop{% endfor %}{% endif %}bar skip this{% endif %}
		text
	{% endfor %}{% endfunc %}`)
}

func TestParseOutputTagSuccess(t *testing.T) {
	// identifier
	testParseSuccess(t, "{%func a()%}{%s foobar %}{%endfunc%}")

	// method call
	testParseSuccess(t, "{%func a()%}{%s foo.bar.baz(a, b, &A{d:e}) %}{%endfunc%}")

	// inline function call
	testParseSuccess(t, "{%func f()%}{%s func() string { return foo.bar(baz, aaa) }() %}{%endfunc%}")

	// map
	testParseSuccess(t, `{%func f()%}{%v map[int]string{1:"foo", 2:"bar"} %}{%endfunc%}`)

	// jsons-safe string
	testParseSuccess(t, `{% func f() %}{%j "foo\nbar" %}{%endfunc%}`)

	// url-encoded string
	testParseSuccess(t, `{% func A() %}{%u "fooab" %}{%endfunc%}`)

	// float precision
	testParseSuccess(t, `{% func A() %}{%f.2 3.1415926535 %}{%endfunc%}`)
	testParseSuccess(t, `{% func A() %}{%f.2= 3.1415926535 %}{%endfunc%}`)
	testParseSuccess(t, `{% func A() %}{%f.= 3.1415926535 %}{%endfunc%}`)
}

func TestParseOutputTagFailure(t *testing.T) {
	// empty tag
	testParseFailure(t, "{%func f()%}{%s %}{%endfunc%}")

	// multiple arguments
	testParseFailure(t, "{%func f()%}{%s a, b %}{%endfunc%}")

	// invalid code
	testParseFailure(t, "{%func f()%}{%s f(a, %}{%endfunc%}")
	testParseFailure(t, "{%func f()%}{%s Foo{Bar:1 %}{%endfunc%}")

	// unsupported code
	testParseFailure(t, "{%func f()%}{%s if (a) {} %}{%endfunc%}")
	testParseFailure(t, "{%func f()%}{%s for {} %}{%endfunc%}")
}

func TestParseTemplateCodeSuccess(t *testing.T) {
	// empty code
	testParseSuccess(t, "{% code %}")
	testParseSuccess(t, "{% func f() %}{% code %}{% endfunc %}")

	// comment
	testParseSuccess(t, `{% code // foobar %}`)
	testParseSuccess(t, `{% func f() %}{% code // foobar %}{% endfunc %}`)
	testParseSuccess(t, `{% code
		// foo
		// bar
	%}`)
	testParseSuccess(t, `{% func f() %}{% code
		// foo
		// bar
	%}{% endfunc %}`)
	testParseSuccess(t, `{%
		code
		/*
			foo
			bar
		*/
	%}`)
	testParseSuccess(t, `{% func f() %}{%
		code
		/*
			foo
			bar
		*/
	%}{% endfunc %}`)

	testParseSuccess(t, `{% code var a int %}`)
	testParseSuccess(t, `{% func f() %}{% code var a int %}{% endfunc %}`)
	testParseSuccess(t, `{% func f() %}{% code a := 0 %}{% endfunc %}`)
	testParseSuccess(t, `{% func f() %}{% code type A struct{} %}{% endfunc %}`)

	// declarations
	testParseSuccess(t, `{%code
		// comment
		type Foo struct {}
		var b = &Foo{}

		func (f *Foo) Bar() {}

		// yet another comment
		func Bar(baz int) string {
			return fmt.Sprintf("%d", baz)
		}
	%}`)
}

func TestParseTemplateCodeFailure(t *testing.T) {
	// import inside the code
	testParseFailure(t, `{% code import "foo" %}`)

	// incomplete code
	testParseFailure(t, `{% code type A struct { %}`)
	testParseFailure(t, `{% code func F() { %}`)

	// invalid code
	testParseFailure(t, `{%code { %}`)
	testParseFailure(t, `{%code {} %}`)
	testParseFailure(t, `{%code ( %}`)
	testParseFailure(t, `{%code () %}`)
}

func TestParseImportSuccess(t *testing.T) {
	// single line import
	testParseSuccess(t, `{% import "github.com/foo/bar" %}`)

	// multiline import
	testParseSuccess(t, `{% import (
		"foo"
		xxx "bar/baz/aaa"

		"yyy.com/zzz"
	) %}`)

	// multiple imports
	testParseSuccess(t, `{% import "foo" %}
		baaas
		{% import (
			"bar"
			"baasd"
		)
		%}
		sddf
	`)
}

func TestParseImportFailure(t *testing.T) {
	// empty import
	testParseFailure(t, `{%import %}`)

	// invalid import
	testParseFailure(t, `{%import foo %}`)

	// non-import code
	testParseFailure(t, `{%import {"foo"} %}`)
	testParseFailure(t, `{%import "foo"
		type A struct {}
	%}`)
	testParseFailure(t, `{%import type a struct{} %}`)
}

func TestParseFailure(t *testing.T) {
	// unknown tag
	testParseFailure(t, "{% foobar %}")

	// unexpected tag outside func
	testParseFailure(t, "aaa{% for %}bbb{%endfor%}")
	testParseFailure(t, "{% return %}")
	testParseFailure(t, "{% break %}")
	testParseFailure(t, "{% if 1==1 %}aaa{%endif%}")
	testParseFailure(t, "{%s abc %}")
	testParseFailure(t, "{%v= aaaa(xx) %}")
	testParseFailure(t, "{%= a() %}")

	// import after func and/or code
	testParseFailure(t, `{%code var i = 0 %}{%import "fmt"%}`)
	testParseFailure(t, `{%func f()%}{%endfunc%}{%import "fmt"%}`)

	// missing endfunc
	testParseFailure(t, "{%func a() %}aaaa")

	// empty func name
	testParseFailure(t, "{% func () %}aaa{% endfunc %}")
	testParseFailure(t, "{% func (a int, b string) %}aaa{% endfunc %}")

	// empty func arguments
	testParseFailure(t, "{% func aaa %}aaa{% endfunc %}")

	// func with anonymous argument
	testParseFailure(t, "{% func a(x int, string) %}{%endfunc%}")

	// func with incorrect arguments' list
	testParseFailure(t, "{% func x(foo, bar) %}{%endfunc%}")
	testParseFailure(t, "{% func x(foo bar baz) %}{%endfunc%}")

	// empty if condition
	testParseFailure(t, "{% func a() %}{% if    %}aaaa{% endif %}{% endfunc %}")

	// missing endif
	testParseFailure(t, "{%func a() %}{%if foo %}aaa{% endfunc %}")

	// missing endfor
	testParseFailure(t, "{%func a()%}{%for %}aaa{%endfunc%}")

	// break outside for
	testParseFailure(t, "{%func a()%}{%break%}{%endfunc%}")

	// invalid if condition
	testParseFailure(t, "{%func a()%}{%if a = b %}{%endif%}{%endfunc%}")
	testParseFailure(t, "{%func f()%}{%if a { %}{%endif%}{%endfunc%}")

	// invalid for
	testParseFailure(t, "{%func a()%}{%for a = b %}{%endfor%}{%endfunc%}")
	testParseFailure(t, "{%func f()%}{%for { %}{%endfor%}{%endfunc%}")

	// invalid code inside func
	testParseFailure(t, "{%func f()%}{%code } %}{%endfunc%}")
	testParseFailure(t, "{%func f()%}{%code { %}{%endfunc%}")

	// interface inside func
	testParseFailure(t, "{%func f()%}{%interface A { Foo() } %}{%endfunc%}")

	// interface without name
	testParseFailure(t, "{%interface  { Foo() } %}")

	// empty interface
	testParseFailure(t, "{% interface Foo {} %}")

	// invalid interface
	testParseFailure(t, "{%interface aaaa %}")
	testParseFailure(t, "{%interface aa { Foo() %}")

	// unnamed method
	testParseFailure(t, "{%func (s *S) () %}{%endfunc%}")

	// empty method arguments
	testParseFailure(t, "{%func (s *S) Foo %}{%endfunc %}")

	// method with return values
	testParseFailure(t, "{%func (s *S) Foo() string %}{%endfunc%}")
	testParseFailure(t, "{%func (s *S) Bar() (int, string) %}{%endfunc%}")
}

func TestParserSuccess(t *testing.T) {
	// empty template
	testParseSuccess(t, "")

	// template without code and funcs
	testParseSuccess(t, "foobar\nbaz")

	// template with code
	testParseSuccess(t, "{%code var a struct {}\nconst n = 123%}")

	// import
	testParseSuccess(t, `{%import "foobar"%}`)
	testParseSuccess(t, `{% import (
	"foo"
	"bar"
)%}`)
	testParseSuccess(t, `{%import "foo"%}{%import "bar"%}`)

	// func
	testParseSuccess(t, "{%func a()%}{%endfunc%}")

	// func with with condition
	testParseSuccess(t, "{%func a(x bool)%}{%if x%}foobar{%endif%}{%endfunc%}")

	// func with complex arguments
	testParseSuccess(t, "{%func f(h1, h2 func(x, y int) string, d int)%}{%endfunc%}")

	// for
	testParseSuccess(t, "{%func a()%}{%for%}aaa{%endfor%}{%endfunc%}")

	// return
	testParseSuccess(t, "{%func a()%}{%return%}{%endfunc%}")

	// nested for
	testParseSuccess(t, "{%func a()%}{%for i := 0; i < 10; i++ %}{%for j := 0; j < i; j++%}aaa{%endfor%}{%endfor%}{%endfunc%}")

	// plain containing arbitrary tags
	testParseSuccess(t, "{%func f()%}{%plain%}This {%endfunc%} is ignored{%endplain%}{%endfunc%}")

	// comment with arbitrary tags
	testParseSuccess(t, "{%func f()%}{%comment%}This {%endfunc%} is ignored{%endcomment%}{%endfunc%}")

	// complex if
	testParseSuccess(t, "{%func a()%}{%if n, err := w.Write(p); err != nil %}{%endif%}{%endfunc%}")

	// complex for
	testParseSuccess(t, "{%func a()%}{%for i, n := 0, len(s); i < n && f(i); i++ %}{%endfor%}{%endfunc%}")

	// complex code inside func
	testParseSuccess(t, `{%func f()%}{%code
		type A struct{}
		var aa []A
		for i := 0; i < 10; i++ {
			aa = append(aa, &A{})
			if i == 42 {
				break
			}
		}
		return
	%}{%endfunc%}`)

	// break inside for loop
	testParseSuccess(t, `{%func f()%}{%for%}
		{% if a() %}
			{% break %}
		{% else %}
			{% return %}
		{% endif %}
	{%endfor%}{%endfunc%}`)

	// interface
	testParseSuccess(t, "{%interface Foo { Bar()\nBaz() } %}")
	testParseSuccess(t, "{%iface Foo { Bar()\nBaz() } %}")

	// method
	testParseSuccess(t, "{%func (s *S) Foo(bar, baz string) %}{%endfunc%}")
}

func testParseFailure(t *testing.T, str string) {
	r := bytes.NewBufferString(str)
	w := &bytes.Buffer{}
	if err := parse(w, r, "memory/foobar.tpl", "memory"); err == nil {
		t.Fatalf("expecting error when parsing %q", str)
	}
}

func testParseSuccess(t *testing.T, str string) {
	r := bytes.NewBufferString(str)
	w := &bytes.Buffer{}
	if err := parse(w, r, "memory/foobar.tpl", "memory"); err != nil {
		t.Fatalf("unexpected error when parsing %q: %s", str, err)
	}
}

func TestParseFile(t *testing.T) {
	filename := "testdata/test.qtpl"
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("cannot open file %q: %s", filename, err)
	}
	defer f.Close()

	packageName, err := getPackageName(filename)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	w := quicktemplate.AcquireByteBuffer()
	if err := parse(w, f, filename, packageName); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	code, err := format.Source(w.B)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	quicktemplate.ReleaseByteBuffer(w)

	expectedFilename := filename + ".compiled"
	expectedCode, err := ioutil.ReadFile(expectedFilename)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !bytes.Equal(code, expectedCode) {
		t.Fatalf("unexpected code: %q\nExpecting %q", code, expectedCode)
	}

}
