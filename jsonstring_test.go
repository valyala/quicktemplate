package quicktemplate

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAppendJSONString(t *testing.T) {
	testAppendJSONString(t, ``)
	testAppendJSONString(t, `f`)
	testAppendJSONString(t, `"`)
	testAppendJSONString(t, `<`)
	testAppendJSONString(t, "\x00\n\r\t\b\f"+`"\`)
	testAppendJSONString(t, `"foobar`)
	testAppendJSONString(t, `foobar"`)
	testAppendJSONString(t, `foo "bar"
		baz`)
	testAppendJSONString(t, `this is a "тест"`)
	testAppendJSONString(t, `привет test ыва`)

	testAppendJSONString(t, `</script><script>alert('evil')</script>`)
	testAppendJSONString(t, "\u001b")
}

func testAppendJSONString(t *testing.T, s string) {
	expectedResult, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("unexpected error when encoding string %q: %s", s, err)
	}
	expectedResult = expectedResult[1 : len(expectedResult)-1]

	bb := AcquireByteBuffer()
	bb.B = AppendJSONString(bb.B[:0], s, false)
	result := string(bb.B)
	ReleaseByteBuffer(bb)

	if strings.Contains(result, "'") {
		t.Fatalf("json string shouldn't contain single quote: %q, src %q", result, s)
	}
	result = strings.Replace(result, `\u0027`, "'", -1)
	result = strings.Replace(result, ">", `\u003e`, -1)
	if result != string(expectedResult) {
		t.Fatalf("unexpected result %q. Expecting %q. original string %q", result, expectedResult, s)
	}
}
