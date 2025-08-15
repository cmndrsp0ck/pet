package dialog

import (
	"testing"

	"github.com/go-test/deep"
)

func TestSearchForParams(t *testing.T) {
	command := "<a=1> <b> hello"

	want := [][2]string{
		{"a", "1"},
		{"b", ""},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_WithNoParams(t *testing.T) {
	command := "no params"

	got := SearchForParams(command)

	if got != nil {
		t.Fatalf("wanted nil, got '%v'", got)
	}
}

func TestSearchForParams_WithMultipleParams(t *testing.T) {
	command := "<a=1> <b> <c=3>"

	want := [][2]string{
		{"a", "1"},
		{"b", ""},
		{"c", "3"},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_WithEmptyCommand(t *testing.T) {
	command := ""

	got := SearchForParams(command)

	if got != nil {
		t.Fatalf("wanted nil, got '%v'", got)
	}
}

func TestSearchForParams_WithNewline(t *testing.T) {
	command := "<a=1> <b> hello\n<c=3>"

	want := [][2]string{
		{"a", "1"},
		{"b", ""},
		{"c", "3"},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_ValueWithSpaces(t *testing.T) {
	command := "example_function --flag=<param=Lots of Bananas>"

	want := [][2]string{
		{"param", "Lots of Bananas"},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_InvalidParamFormat(t *testing.T) {
	command := "<a=1 <b> hello"
	want := [][2]string{
		{"b", ""},
	}
	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_InvalidParamFormatWithoutSpaces(t *testing.T) {
	command := "<a=1<b>hello"
	want := [][2]string{
		{"b", ""},
	}
	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_ConfusingBrackets(t *testing.T) {
	command := "cat <<EOF > <file=path/to/file>\nEOF"
	want := [][2]string{
		{"file", "path/to/file"},
	}
	got := SearchForParams(command)
	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKey(t *testing.T) {
	command := "<a=1> <a=2> <a=3>"
	want := [][2]string{
		{"a", "3"},
	}
	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues(t *testing.T) {
	command := "<a=1> <a=2> <a=3>"
	want := [][2]string{
		{"a", "3"},
	}
	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_MultipleLines(t *testing.T) {
	command := "<a=1> <a=2> <a=3>\n<b=4>"
	want := [][2]string{
		{"a", "3"},
		{"b", "4"},
	}
	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat(t *testing.T) {
	command := "<a=1> <a=2 <a=3>"
	want := [][2]string{
		{"a", "3"},
	}
	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat_MultipleLines(t *testing.T) {
	command := "<a=1> <a=2> <a=3 \n<b=4>"
	want := [][2]string{
		{"a", "2"},
		{"b", "4"},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat_MultipleLines2(t *testing.T) {
	command := "<a=1> <a=2> <a=3>\n<b=4"
	want := [][2]string{
		{"a", "3"},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_EqualsInDefaultValueIgnored(t *testing.T) {
	command := "echo \"<param=Hello == World!===>\""
	want := [][2]string{
		{"param", "Hello == World!==="},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleDefaultValuesDoNotBreakFunction(t *testing.T) {
	command := "echo \"<param=|_Hello_||_Hello world_||_How are you?_|> <second=Hello>, <third>\""
	want := [][2]string{
		{"param", "|_Hello_||_Hello world_||_How are you?_|"},
		{"second", "Hello"},
		{"third", ""},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestInsertParams(t *testing.T) {
	command := "<a=1> <a> <b> hello"

	params := map[string]string{
		"a": "test",
		"b": "case",
	}

	got := insertParams(command, params)
	want := "test test case hello"
	if want != got {
		t.Fatalf("wanted '%s', got '%s'", want, got)
	}
}

func TestInsertParams_unique_parameters(t *testing.T) {
	command := "curl -X POST \"<host=http://localhost:9200>/<index>\" -H 'Content-Type: application/json'"

	params := map[string]string{
		"host":  "localhost:9200",
		"index": "test",
	}

	got := insertParams(command, params)
	want := "curl -X POST \"localhost:9200/test\" -H 'Content-Type: application/json'"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestInsertParams_complex(t *testing.T) {
	command := "something <host=http://localhost:9200>/<test>/_delete_by_query/<host>"

	params := map[string]string{
		"host": "localhost:9200",
		"test": "case",
	}

	got := insertParams(command, params)
	want := "something localhost:9200/case/_delete_by_query/localhost:9200"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestInsertParams_EqualsInDefaultValueIgnored(t *testing.T) {
	command := "echo \"<param=Hello == World!===>\""

	params := map[string]string{
		"param": "something == something",
	}

	got := insertParams(command, params)
	want := "echo \"something == something\""
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

// Test escape functionality
func TestSearchForParams_WithEscapedParams(t *testing.T) {
	command := "<a=1> \\<b=2\\> <c=3>"

	want := [][2]string{
		{"a", "1"},
		{"c", "3"},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_WithEscapedParamsOnly(t *testing.T) {
	command := "\\<a=1\\> \\<b=2\\>"

	got := SearchForParams(command)

	if got != nil {
		t.Fatalf("wanted nil, got '%v'", got)
	}
}

func TestSearchForParams_WithPartialEscape(t *testing.T) {
	command := "\\<a=1> <b=2\\>"

	// Should find both a=1 and b=2\ since neither is fully escaped
	// (a=1 is missing trailing \, b=2\ is missing leading \)
	want := [][2]string{
		{"a", "1"},
		{"b", "2\\"},
	}

	got := SearchForParams(command)

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestInsertParams_WithEscapedParams(t *testing.T) {
	command := "<a=1> \\<b=literal\\> <c=3>"

	params := map[string]string{
		"a": "replaced_a",
		"c": "replaced_c",
	}

	got := insertParams(command, params)
	want := "replaced_a <b=literal> replaced_c"

	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestInsertParams_WithOnlyEscapedParams(t *testing.T) {
	command := "\\<a=literal\\> \\<b=also_literal\\>"

	params := map[string]string{}

	got := insertParams(command, params)
	want := "<a=literal> <b=also_literal>"

	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestRemoveEscapeChars(t *testing.T) {
	command := "echo \\<hello\\> world \\<foo\\>"

	got := removeEscapeChars(command)
	want := "echo <hello> world <foo>"

	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestInsertParams_MixedEscapedAndNormal(t *testing.T) {
	command := "curl -X POST <host>/api \\<version=v1\\> -d '<data>'"

	params := map[string]string{
		"host": "localhost:8080",
		"data": "test_payload",
	}

	got := insertParams(command, params)
	want := "curl -X POST localhost:8080/api <version=v1> -d 'test_payload'"

	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
