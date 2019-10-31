package parser

import (
	"bytes"
	"reflect"
	"testing"
)

func TestScannerTagNameWithDotAndEqual(t *testing.T) {
	testScannerSuccess(t, "{%foo.bar.34 baz aaa%} awer{% aa= %}",
		[]tt{
			{ID: tagName, Value: "foo.bar.34"},
			{ID: tagContents, Value: "baz aaa"},
			{ID: text, Value: " awer"},
			{ID: tagName, Value: "aa="},
			{ID: tagContents, Value: ""},
		})
}

func TestScannerEndTagWithMinus(t *testing.T) {
	testScannerSuccess(t, "{%foo baz aaa-%} awer{% aa= -%}",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: "baz aaa"},
			{ID: text, Value: "awer"},
			{ID: tagName, Value: "aa="},
			{ID: tagContents, Value: ""},
		})
	testScannerSuccess(t, "{%foo baz aaa- %} awer{% aa= -% %}",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: "baz aaa-"},
			{ID: text, Value: " awer"},
			{ID: tagName, Value: "aa="},
			{ID: tagContents, Value: "-%"},
		})
	testScannerSuccess(t, "{%foo -%} \t\r\r\n\n\tawer{%aa= - %}  \nbb",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: ""},
			{ID: text, Value: "\n\tawer"},
			{ID: tagName, Value: "aa="},
			{ID: tagContents, Value: "-"},
			{ID: text, Value: "  \nbb"},
		})
	testScannerSuccess(t, "{%foo-%}\n  awer{%bar -%}\n",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: ""},
			{ID: text, Value: "  awer"},
			{ID: tagName, Value: "bar"},
			{ID: tagContents, Value: ""},
		})
	testScannerSuccess(t, "{%foo-%}{%bar%} \n {%baz%}",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: ""},
			{ID: tagName, Value: "bar"},
			{ID: tagContents, Value: ""},
			{ID: text, Value: " \n "},
			{ID: tagName, Value: "baz"},
			{ID: tagContents, Value: ""},
		})
}

func TestScannerBeginTagWithMinus(t *testing.T) {
	testScannerSuccess(t, "x\n   {%-foo%}  {%- bar%}",
		[]tt{
			{ID: text, Value: "x\n"},
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: ""},
			{ID: tagName, Value: "bar"},
			{ID: tagContents, Value: ""},
		})
	testScannerSuccess(t, "{%- foo%}",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: ""},
		})
	testScannerSuccess(t, "{%-foo baz aaa%} awer\n\t {%-aa= %}",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: "baz aaa"},
			{ID: text, Value: " awer\n"},
			{ID: tagName, Value: "aa="},
			{ID: tagContents, Value: ""},
		})
	testScannerSuccess(t, "{%-foo baz aaa%} awer {- % aa= %}",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: "baz aaa"},
			{ID: text, Value: " awer {- % aa= %}"},
		})
	testScannerSuccess(t, "\n\n   {%- foo baz aaa -%} awer {-{%- aa= xxx%}",
		[]tt{
			{ID: text, Value: "\n\n"},
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: "baz aaa"},
			{ID: text, Value: "awer {-"},
			{ID: tagName, Value: "aa="},
			{ID: tagContents, Value: "xxx"},
		})
	testScannerSuccess(t, "{%-foo%}{%bar%} \n {%-baz%}",
		[]tt{
			{ID: tagName, Value: "foo"},
			{ID: tagContents, Value: ""},
			{ID: tagName, Value: "bar"},
			{ID: tagContents, Value: ""},
			{ID: text, Value: " \n"},
			{ID: tagName, Value: "baz"},
			{ID: tagContents, Value: ""},
		})
}

func TestScannerStripspaceSuccess(t *testing.T) {
	testScannerSuccess(t, "  aa\n\t {%stripspace%} \t\n  f\too \n   b  ar \n\r\t {%  bar baz  asd %}\n\nbaz \n\t \taaa  \n{%endstripspace%} bb  ", []tt{
		{ID: text, Value: "  aa\n\t "},
		{ID: text, Value: "f\toob  ar"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: "baz  asd"},
		{ID: text, Value: "bazaaa"},
		{ID: text, Value: " bb  "},
	})
	testScannerSuccess(t, "{%stripspace  %}{% stripspace fobar %} {%space%}  a\taa\n\r\t bb  b  {%endstripspace  %}  {%endstripspace  baz%}", []tt{
		{ID: text, Value: " "},
		{ID: text, Value: "a\taabb  b"},
	})

	// sripspace wins over collapsespace
	testScannerSuccess(t, "{%stripspace%} {%collapsespace%}foo\n\t bar{%endcollapsespace%} \r\n\t {%endstripspace%}", []tt{
		{ID: text, Value: "foobar"},
	})
}

func TestScannerStripspaceFailure(t *testing.T) {
	// incomplete stripspace tag
	testScannerFailure(t, "{%stripspace   ")

	// incomplete endstripspace tag
	testScannerFailure(t, "{%stripspace%}aaa{%endstripspace")

	// missing endstripspace
	testScannerFailure(t, "{%stripspace%} foobar")

	// missing stripspace
	testScannerFailure(t, "aaa{%endstripspace%}")

	// missing the second endstripspace
	testScannerFailure(t, "{%stripspace%}{%stripspace%}aaaa{%endstripspace%}")
}

func TestScannerCollapsespaceSuccess(t *testing.T) {
	testScannerSuccess(t, "  aa\n\t {%collapsespace%} \t\n  foo \n   bar{%  bar baz  asd %}\n\nbaz \n   \n{%endcollapsespace%} bb  ", []tt{
		{ID: text, Value: "  aa\n\t "},
		{ID: text, Value: " foo bar"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: "baz  asd"},
		{ID: text, Value: " baz "},
		{ID: text, Value: " bb  "},
	})
	testScannerSuccess(t, "{%collapsespace  %}{% collapsespace fobar %} {%space%}  aaa\n\r\t bbb  {%endcollapsespace  %}  {%endcollapsespace  baz%}", []tt{
		{ID: text, Value: " "},
		{ID: text, Value: " "},
		{ID: text, Value: " aaa bbb "},
		{ID: text, Value: " "},
	})
}

func TestScannerCollapsespaceFailure(t *testing.T) {
	// incomplete collapsespace tag
	testScannerFailure(t, "{%collapsespace   ")

	// incomplete endcollapsespace tag
	testScannerFailure(t, "{%collapsespace%}aaa{%endcollapsespace")

	// missing endcollapsespace
	testScannerFailure(t, "{%collapsespace%} foobar")

	// missing collapsespace
	testScannerFailure(t, "aaa{%endcollapsespace%}")

	// missing the second endcollapsespace
	testScannerFailure(t, "{%collapsespace%}{%collapsespace%}aaaa{%endcollapsespace%}")
}

func TestScannerPlainSuccess(t *testing.T) {
	testScannerSuccess(t, "{%plain%}{%endplain%}", nil)
	testScannerSuccess(t, "{%plain%}{%foo bar%}asdf{%endplain%}", []tt{
		{ID: text, Value: "{%foo bar%}asdf"},
	})
	testScannerSuccess(t, "{%plain%}{%foo{%endplain%}", []tt{
		{ID: text, Value: "{%foo"},
	})
	testScannerSuccess(t, "aa{%plain%}bbb{%cc%}{%endplain%}{%plain%}dsff{%endplain%}", []tt{
		{ID: text, Value: "aa"},
		{ID: text, Value: "bbb{%cc%}"},
		{ID: text, Value: "dsff"},
	})
	testScannerSuccess(t, "mmm{%plain%}aa{% bar {%%% }baz{%endplain%}nnn", []tt{
		{ID: text, Value: "mmm"},
		{ID: text, Value: "aa{% bar {%%% }baz"},
		{ID: text, Value: "nnn"},
	})
	testScannerSuccess(t, "{% plain dsd %}0{%comment%}123{%endcomment%}45{% endplain aaa %}", []tt{
		{ID: text, Value: "0{%comment%}123{%endcomment%}45"},
	})
}

func TestScannerPlainFailure(t *testing.T) {
	testScannerFailure(t, "{%plain%}sdfds")
	testScannerFailure(t, "{%plain%}aaaa%{%endplain")
	testScannerFailure(t, "{%plain%}{%endplain%")
}

func TestScannerCommentSuccess(t *testing.T) {
	testScannerSuccess(t, "{%comment%}{%endcomment%}", nil)
	testScannerSuccess(t, "{%comment%}foo{%endcomment%}", nil)
	testScannerSuccess(t, "{%comment%}foo{%endcomment%}{%comment%}sss{%endcomment%}", nil)
	testScannerSuccess(t, "{%comment%}foo{%bar%}{%endcomment%}", nil)
	testScannerSuccess(t, "{%comment%}foo{%bar {%endcomment%}", nil)
	testScannerSuccess(t, "{%comment%}foo{%bar&^{%endcomment%}", nil)
	testScannerSuccess(t, "{%comment%}foo{% bar\n\rs%{%endcomment%}", nil)
	testScannerSuccess(t, "xx{%x%}www{% comment aux data %}aaa{% comment %}{% endcomment %}yy", []tt{
		{ID: text, Value: "xx"},
		{ID: tagName, Value: "x"},
		{ID: tagContents, Value: ""},
		{ID: text, Value: "www"},
		{ID: text, Value: "yy"},
	})
}

func TestScannerCommentFailure(t *testing.T) {
	testScannerFailure(t, "{%comment%}...no endcomment")
	testScannerFailure(t, "{% comment %}foobar{% endcomment")
}

func TestScannerSuccess(t *testing.T) {
	testScannerSuccess(t, "", nil)
	testScannerSuccess(t, "a%}{foo}bar", []tt{
		{ID: text, Value: "a%}{foo}bar"},
	})
	testScannerSuccess(t, "{% foo bar baz(a, b, 123) %}", []tt{
		{ID: tagName, Value: "foo"},
		{ID: tagContents, Value: "bar baz(a, b, 123)"},
	})
	testScannerSuccess(t, "foo{%bar%}baz", []tt{
		{ID: text, Value: "foo"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: ""},
		{ID: text, Value: "baz"},
	})
	testScannerSuccess(t, "{{{%\n\r\tfoo bar\n\rbaz%%\n   \r %}}", []tt{
		{ID: text, Value: "{{"},
		{ID: tagName, Value: "foo"},
		{ID: tagContents, Value: "bar\n\rbaz%%"},
		{ID: text, Value: "}"},
	})
	testScannerSuccess(t, "{%%}", []tt{
		{ID: tagName, Value: ""},
		{ID: tagContents, Value: ""},
	})
	testScannerSuccess(t, "{%%aaa bb%}", []tt{
		{ID: tagName, Value: ""},
		{ID: tagContents, Value: "%aaa bb"},
	})
	testScannerSuccess(t, "foo{% bar %}{% baz aa (123)%}321", []tt{
		{ID: text, Value: "foo"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: ""},
		{ID: tagName, Value: "baz"},
		{ID: tagContents, Value: "aa (123)"},
		{ID: text, Value: "321"},
	})
}

func TestScannerFailure(t *testing.T) {
	testScannerFailure(t, "a{%")
	testScannerFailure(t, "a{%foo")
	testScannerFailure(t, "a{%% }foo")
	testScannerFailure(t, "a{% foo %")
	testScannerFailure(t, "b{% fo() %}bar")
	testScannerFailure(t, "aa{% foo bar")

	testScannerFailure(t, "a{%-")
	testScannerFailure(t, "a{%-foo")
	testScannerFailure(t, "a{%-% }foo")
	testScannerFailure(t, "a{%- foo %")
	testScannerFailure(t, "b{%- fo() %}bar")
	testScannerFailure(t, "aa{%- foo bar")
}

func testScannerFailure(t *testing.T, str string) {
	r := bytes.NewBufferString(str)
	s := newScanner(r, "memory")
	var tokens []tt
	for s.Next() {
		tokens = append(tokens, tt{
			ID:    s.Token().ID,
			Value: string(s.Token().Value),
		})
	}
	if err := s.LastError(); err == nil {
		t.Fatalf("expecting error when scanning %q. got tokens %v", str, tokens)
	}
}

func testScannerSuccess(t *testing.T, str string, expectedTokens []tt) {
	r := bytes.NewBufferString(str)
	s := newScanner(r, "memory")
	var tokens []tt
	for s.Next() {
		tokens = append(tokens, tt{
			ID:    s.Token().ID,
			Value: string(s.Token().Value),
		})
	}
	if err := s.LastError(); err != nil {
		t.Fatalf("unexpected error: %s. str=%q", err, str)
	}
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Fatalf("unexpected tokens %v. Expecting %v. str=%q", tokens, expectedTokens, str)
	}
}

type tt struct {
	ID    int
	Value string
}
