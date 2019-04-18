package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

/*TestScannerConfigurationTagNameWithDotAndEqual ..."
 */
func TestScannerConfigurationTagNameWithDotAndEqual(t *testing.T) {
	testScannerConfigurationSuccess(t, "[{foo.bar.34 baz aaa}] awer[{ aa= }]",
		[]confToken{
			{ID: tagName, Value: "foo.bar.34"},
			{ID: tagContents, Value: "baz aaa"},
			{ID: text, Value: " awer"},
			{ID: tagName, Value: "aa="},
			{ID: tagContents, Value: ""},
		})
}

/*TestScannerConfigurationStripspaceSuccess ...
 */
func TestScannerConfigurationStripspaceSuccess(t *testing.T) {
	testScannerConfigurationSuccess(t, "  aa\n\t [{stripspace}] \t\n  f\too \n   b  ar \n\r\t [{  bar baz  asd }]\n\nbaz \n\t \taaa  \n[{endstripspace}] bb  ", []confToken{
		{ID: text, Value: "  aa\n\t "},
		{ID: text, Value: "f\toob  ar"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: "baz  asd"},
		{ID: text, Value: "bazaaa"},
		{ID: text, Value: " bb  "},
	})
	testScannerConfigurationSuccess(t, "[{stripspace  }][{ stripspace fobar }] [{space}]  a\taa\n\r\t bb  b  [{endstripspace  }]  [{endstripspace  baz}]", []confToken{
		{ID: text, Value: " "},
		{ID: text, Value: "a\taabb  b"},
	})

	// sripspace wins over collapsespace
	testScannerConfigurationSuccess(t, "[{stripspace}] [{collapsespace}]foo\n\t bar[{endcollapsespace}] \r\n\t [{endstripspace}]", []confToken{
		{ID: text, Value: "foobar"},
	})
}

/*TestScannerConfigurationStripspaceFailure ...
 */
func TestScannerConfigurationStripspaceFailure(t *testing.T) {
	// incomplete stripspace tag
	testScannerConfigurationFailure(t, "[{stripspace   ")

	// incomplete endstripspace tag
	testScannerConfigurationFailure(t, "[{stripspace}]aaa[{endstripspace")

	// missing endstripspace
	testScannerConfigurationFailure(t, "[{stripspace}] foobar")

	// missing stripspace
	testScannerConfigurationFailure(t, "aaa[{endstripspace}]")

	// missing the second endstripspace
	testScannerConfigurationFailure(t, "[{stripspace}][{stripspace}]aaaa[{endstripspace}]")
}

/*TestScannerConfigurationCollapsespaceSuccess ...
 */
func TestScannerConfigurationCollapsespaceSuccess(t *testing.T) {
	testScannerConfigurationSuccess(t, "  aa\n\t [{collapsespace}] \t\n  foo \n   bar[{  bar baz  asd }]\n\nbaz \n   \n[{endcollapsespace}] bb  ", []confToken{
		{ID: text, Value: "  aa\n\t "},
		{ID: text, Value: " foo bar"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: "baz  asd"},
		{ID: text, Value: " baz "},
		{ID: text, Value: " bb  "},
	})
	testScannerConfigurationSuccess(t, "[{collapsespace  }][{ collapsespace fobar }] [{space}]  aaa\n\r\t bbb  [{endcollapsespace  }]  [{endcollapsespace  baz}]", []confToken{
		{ID: text, Value: " "},
		{ID: text, Value: " "},
		{ID: text, Value: " aaa bbb "},
		{ID: text, Value: " "},
	})
}

/*TestScannerCollapsespaceFailure ...
 */
func TestScannerConfigurationCollapsespaceFailure(t *testing.T) {
	// incomplete collapsespace tag
	testScannerConfigurationFailure(t, "[{collapsespace   ")

	// incomplete endcollapsespace tag
	testScannerConfigurationFailure(t, "[{collapsespace}]aaa[{endcollapsespace")

	// missing endcollapsespace
	testScannerConfigurationFailure(t, "[{collapsespace}] foobar")

	// missing collapsespace
	testScannerConfigurationFailure(t, "aaa[{endcollapsespace}]")

	// missing the second endcollapsespace
	testScannerConfigurationFailure(t, "[{collapsespace}][{collapsespace}]aaaa[{endcollapsespace}]")
}

/*TestScannerPlainSuccess ...
 */
func TestScannerConfigurationPlainSuccess(t *testing.T) {
	testScannerConfigurationSuccess(t, "[{plain}][{endplain}]", nil)
	testScannerConfigurationSuccess(t, "[{plain}][{foo bar}]asdf[{endplain}]", []confToken{
		{ID: text, Value: "[{foo bar}]asdf"},
	})
	testScannerConfigurationSuccess(t, "[{plain}][{foo[{endplain}]", []confToken{
		{ID: text, Value: "[{foo"},
	})
	testScannerConfigurationSuccess(t, "aa[{plain}]bbb[{cc}][{endplain}][{plain}]dsff[{endplain}]", []confToken{
		{ID: text, Value: "aa"},
		{ID: text, Value: "bbb[{cc}]"},
		{ID: text, Value: "dsff"},
	})
	testScannerConfigurationSuccess(t, "mmm[{plain}]aa[{ bar [{%% }baz[{endplain}]nnn", []confToken{
		{ID: text, Value: "mmm"},
		{ID: text, Value: "aa[{ bar [{%% }baz"},
		{ID: text, Value: "nnn"},
	})
	testScannerConfigurationSuccess(t, "[{ plain dsd }]0[{comment}]123[{endcomment}]45[{ endplain aaa }]", []confToken{
		{ID: text, Value: "0[{comment}]123[{endcomment}]45"},
	})
}

/*TestScannerPlainFailure ...
 */
func TestScannerConfigurationPlainFailure(t *testing.T) {
	testScannerConfigurationFailure(t, "[{plain}]sdfds")
	testScannerConfigurationFailure(t, "[{plain}]aaaa%[{endplain")
	testScannerConfigurationFailure(t, "[{plain}][{endplain%")
}

/*TestScannerCommentSuccess ...
 */
func TestScannerConfigurationCommentSuccess(t *testing.T) {
	testScannerConfigurationSuccess(t, "[{comment}][{endcomment}]", nil)
	testScannerConfigurationSuccess(t, "[{comment}]foo[{endcomment}]", nil)
	testScannerConfigurationSuccess(t, "[{comment}]foo[{endcomment}][{comment}]sss[{endcomment}]", nil)
	testScannerConfigurationSuccess(t, "[{comment}]foo[{bar}][{endcomment}]", nil)
	testScannerConfigurationSuccess(t, "[{comment}]foo[{bar [{endcomment}]", nil)
	testScannerConfigurationSuccess(t, "[{comment}]foo[{bar&^[{endcomment}]", nil)
	testScannerConfigurationSuccess(t, "[{comment}]foo[{ bar\n\rs%[{endcomment}]", nil)
	testScannerConfigurationSuccess(t, "xx[{x}]www[{ comment aux data }]aaa[{ comment }][{ endcomment }]yy", []confToken{
		{ID: text, Value: "xx"},
		{ID: tagName, Value: "x"},
		{ID: tagContents, Value: ""},
		{ID: text, Value: "www"},
		{ID: text, Value: "yy"},
	})
}

/*TestScannerCommentFailure ...
 */
func TestScannerConfigurationCommentFailure(t *testing.T) {
	testScannerConfigurationFailure(t, "[{comment}]...no endcomment")
	testScannerConfigurationFailure(t, "[{ comment }]foobar[{ endcomment")
}

func TestScannerConfigurationSuccess(t *testing.T) {
	testScannerConfigurationSuccess(t, "", nil)
	testScannerConfigurationSuccess(t, "a}]{foo}bar", []confToken{
		{ID: text, Value: "a}]{foo}bar"},
	})
	testScannerConfigurationSuccess(t, "[{ foo bar baz(a, b, 123) }]", []confToken{
		{ID: tagName, Value: "foo"},
		{ID: tagContents, Value: "bar baz(a, b, 123)"},
	})
	testScannerConfigurationSuccess(t, "foo[{bar}]baz", []confToken{
		{ID: text, Value: "foo"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: ""},
		{ID: text, Value: "baz"},
	})
	testScannerConfigurationSuccess(t, "{{[{\n\r\tfoo bar\n\rbaz%%\n   \r }]}", []confToken{
		{ID: text, Value: "{{"},
		{ID: tagName, Value: "foo"},
		{ID: tagContents, Value: "bar\n\rbaz%%"},
		{ID: text, Value: "}"},
	})
	testScannerConfigurationSuccess(t, "[{}]", []confToken{
		{ID: tagName, Value: ""},
		{ID: tagContents, Value: ""},
	})
	testScannerConfigurationSuccess(t, "[{%aaa bb}]", []confToken{
		{ID: tagName, Value: ""},
		{ID: tagContents, Value: "%aaa bb"},
	})
	testScannerConfigurationSuccess(t, "foo[{ bar }][{ baz aa (123)}]321", []confToken{
		{ID: text, Value: "foo"},
		{ID: tagName, Value: "bar"},
		{ID: tagContents, Value: ""},
		{ID: tagName, Value: "baz"},
		{ID: tagContents, Value: "aa (123)"},
		{ID: text, Value: "321"},
	})
}

/*TestScannerConfigurationFailure ...
 */
func TestScannerConfigurationFailure(t *testing.T) {
	testScannerConfigurationFailure(t, "a[{")
	testScannerConfigurationFailure(t, "a[{foo")
	testScannerConfigurationFailure(t, "a[{% }foo")
	testScannerConfigurationFailure(t, "a[{ foo %")
	testScannerConfigurationFailure(t, "b[{ fo() }]bar")
	testScannerConfigurationFailure(t, "aa[{ foo bar")
}

func testScannerConfigurationFailure(t *testing.T, str string) {
	r := bytes.NewBufferString(str)
	s := newScannerWithTagConf(r, "memory", "[{", "}]")
	var tokens []confToken
	for s.Next() {
		id := s.Token().ID
		val := string(s.Token().Value)
		fmt.Printf("%d %s", id, val)
		tokens = append(tokens, confToken{
			ID:    s.Token().ID,
			Value: string(s.Token().Value),
		})
	}
	if err := s.LastError(); err == nil {
		t.Fatalf("expecting error when scanning %q. got tokens %v", str, tokens)
	}
}

func testScannerConfigurationSuccess(t *testing.T, str string, expectedTokens []confToken) {
	r := bytes.NewBufferString(str)
	s := newScannerWithTagConf(r, "memory", "[{", "}]")
	var tokens []confToken
	for s.Next() {
		tokens = append(tokens, confToken{
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

// /*TestScannerConfigurationSuccessTest ...
//  */
// func TestScannerConfigurationSuccess(t *testing.T) {
// 	testScannerConfigurationSuccess(t, "", nil)
// 	testScannerConfigurationSuccess(t, "a%]{foo}bar", []confToken{
// 		{ID: text, Value: "a%]{foo}bar"},
// 	})
// 	testScannerConfigurationSuccess(t, "[% foo bar baz(a, b, 123) %]", []confToken{
// 		{ID: tagName, Value: "foo"},
// 		{ID: tagContents, Value: "bar baz(a, b, 123)"},
// 	})
// 	testScannerConfigurationSuccess(t, "foo[%bar%]baz", []confToken{
// 		{ID: text, Value: "foo"},
// 		{ID: tagName, Value: "bar"},
// 		{ID: tagContents, Value: ""},
// 		{ID: text, Value: "baz"},
// 	})
// 	testScannerConfigurationSuccess(t, "{{[%\n\r\tfoo bar\n\rbaz%%\n   \r }]}", []confToken{
// 		{ID: text, Value: "{{"},
// 		{ID: tagName, Value: "foo"},
// 		{ID: tagContents, Value: "bar\n\rbaz%%"},
// 		{ID: text, Value: "}"},
// 	})
// 	testScannerConfigurationSuccess(t, "[%%]", []confToken{
// 		{ID: tagName, Value: ""},
// 		{ID: tagContents, Value: ""},
// 	})
// 	testScannerConfigurationSuccess(t, "[%%aaa bb%]", []confToken{
// 		{ID: tagName, Value: ""},
// 		{ID: tagContents, Value: "%aaa bb"},
// 	})
// 	testScannerConfigurationSuccess(t, "foo[% bar %][% baz aa (123)%]321", []confToken{
// 		{ID: text, Value: "foo"},
// 		{ID: tagName, Value: "bar"},
// 		{ID: tagContents, Value: ""},
// 		{ID: tagName, Value: "baz"},
// 		{ID: tagContents, Value: "aa (123)"},
// 		{ID: text, Value: "321"},
// 	})
// 	testScannerConfigurationSuccess(t, "  aa\n\t [%stripspace%] \t\n  f\too \n   b  ar \n\r\t [%  bar baz  asd %]\n\nbaz \n\t \taaa  \n[%endstripspace%] bb  ", []confToken{
// 		{ID: text, Value: "  aa\n\t "},
// 		{ID: text, Value: "f\toob  ar"},
// 		{ID: tagName, Value: "bar"},
// 		{ID: tagContents, Value: "baz  asd"},
// 		{ID: text, Value: "bazaaa"},
// 		{ID: text, Value: " bb  "},
// 	})
// }

/*TestCrash ...
 */
func TestCrash(t *testing.T) {

	var source = `This is a base page template. All the other template pages implement this interface.

{% interface
Page {
	Title()
	Body()
}
}]


Page prints a page implementing Page interface.
[{ func PageTemplate(p Page) }]
<html>
	<head>
		<title>{%= p.Title() }]</title>
	</head>
	<body>
		<div>
			<a href="/">return to main page</a>
		</div>
		[{= p.Body() }]
	</body>
</html>
[{ endfunc }]


Base page implementation. Other pages may inherit from it if they need
overriding only certain Page methods
[{ code type BasePage struct {} }]
[{ func (p *BasePage) Title() }]This is a base title[{ endfunc }]
[{ func (p *BasePage) Body() }]This is a base body[{ endfunc }]
}`
	r := bytes.NewBufferString(source)
	s := newScannerWithTagConf(r, "memory", "[{", "}]")
	var tokens []confToken
	for s.Next() {
		tokens = append(tokens, confToken{
			ID:    s.Token().ID,
			Value: string(s.Token().Value),
		})
	}
	if err := s.LastError(); err != nil {
		t.Fatalf("unexpected error: %s. str=%q", err, source)
	}
}

type confToken struct {
	ID    int
	Value string
}
