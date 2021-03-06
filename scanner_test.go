package gmitxt_test

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"git.sr.ht/~kiba/gmitxt"
	"git.sr.ht/~kiba/gmitxt/internal/toast"
)

const example = "testdata/example.gmi"

func TestSanner(t *testing.T) {
	f, err := os.Open(example)
	if err != nil {
		t.Fatalf("could not open %s: %v", example, err)
	}
	defer f.Close()

	t.Logf("scanning: %s", example)

	s := gmitxt.NewScanner(f)
	expectStart(t, s)
	expectLine(t, s, 1, gmitxt.Head1, "This is my test Gemini ")
	expectLine(t, s, 2, gmitxt.Head1, "Heading #1")
	expectLine(t, s, 3, gmitxt.Head1, "")
	expectLine(t, s, 4, gmitxt.Head1, "")
	expectLine(t, s, 5, gmitxt.Head2, "This is a level two heading.")
	expectLine(t, s, 6, gmitxt.Head2, "Heading #2 ")
	expectLine(t, s, 7, gmitxt.Head2, "")
	expectLine(t, s, 8, gmitxt.Head2, "")
	expectLine(t, s, 9, gmitxt.Head3, "This is a level three heading.")
	expectLine(t, s, 10, gmitxt.Head3, "Heading #3 ")
	expectLine(t, s, 11, gmitxt.Head3, "")
	expectLine(t, s, 12, gmitxt.Head3, "")
	expectLine(t, s, 13, gmitxt.Text, "")
	expectLine(t, s, 14, gmitxt.Text, "This is a text line.")
	expectLine(t, s, 15, gmitxt.Text, "Another text line with trailing whitespace.   ") // nolint: lll
	expectLine(t, s, 16, gmitxt.Text, "")
	expectLine(t, s, 17, gmitxt.List, "List 1")
	expectLine(t, s, 18, gmitxt.Text, "*List 2")
	expectLine(t, s, 19, gmitxt.Text, "*")
	expectLine(t, s, 20, gmitxt.List, "")
	expectLine(t, s, 21, gmitxt.Text, "")
	expectLine(t, s, 22, gmitxt.Quote, " Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.") // nolint: lll
	expectLine(t, s, 23, gmitxt.Quote, "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.")                  // nolint: lll
	expectLine(t, s, 24, gmitxt.Quote, "")
	expectLine(t, s, 25, gmitxt.Text, "")
	expectLink(t, s, 26, "https://example.tld/", "")
	expectLink(t, s, 27, "gemini://example.tld/", "")
	expectLink(t, s, 28, "gemini://example.tld/", "Example link with a description") // nolint: lll
	expectLink(t, s, 29, "foo/bar/baz.txt", "A relative link ")
	expectLine(t, s, 30, gmitxt.PreStart, "go ")
	expectLine(t, s, 31, gmitxt.PreBody, "package main")
	expectLine(t, s, 32, gmitxt.PreBody, `import "fmt"`)
	expectLine(t, s, 33, gmitxt.PreBody, "func main() {")
	expectLine(t, s, 34, gmitxt.PreBody, `	fmt.Println("hello world")`)
	expectLine(t, s, 35, gmitxt.PreBody, "}")
	expectLine(t, s, 36, gmitxt.PreEnd, "")
	expectLine(t, s, 37, gmitxt.PreStart, "")
	expectLine(t, s, 38, gmitxt.PreBody, "Normal preformatted text")
	expectLine(t, s, 39, gmitxt.PreEnd, "")
	expectEnd(t, s, 39)
	expectEnd(t, s, 39)
	expectEnd(t, s, 39)

	t.Log("scanning with a tiny buffer")

	input := `# This is my test Gemini
## This is a level two heading.`
	buf := make([]byte, 0, 24)
	s = gmitxt.NewScanner(strings.NewReader(input))
	s.Buffer(buf, 25)
	expectStart(t, s)
	expectLine(t, s, 1, gmitxt.Head1, "This is my test Gemini")

	if s.Scan() {
		t.Errorf("Line %d: scanner should have stopped", s.Line())
	}

	if !errors.Is(s.Err(), bufio.ErrTooLong) {
		t.Errorf("Line %d: scanner should have error `%v`, got: %v",
			s.Line(), bufio.ErrTooLong, s.Err())
	}
}

func BenchmarkScanner(b *testing.B) {
	input, err := ioutil.ReadFile(example)
	if err != nil {
		b.Fatalf("could not read file %s: %v", example, err)
	}

	for i := 0; i < b.N; i++ {
		s := gmitxt.NewScanner(bytes.NewReader(input))
		for s.Scan() {
		}
	}

	b.ReportAllocs()
}

// BenchmarkBufioScanner benchmarks the bufio.Scanner for comparison.
func BenchmarkBufioScanner(b *testing.B) {
	input, err := ioutil.ReadFile(example)
	if err != nil {
		b.Fatalf("could not read file %s: %v", example, err)
	}

	for i := 0; i < b.N; i++ {
		s := bufio.NewScanner(bytes.NewReader(input))
		for s.Scan() {
		}
	}

	b.ReportAllocs()
}

// BenchmarkToastParser benchmarks the toast.cafe/x/gmi parser for comparison.
func BenchmarkToastParser(b *testing.B) {
	input, err := ioutil.ReadFile(example)
	if err != nil {
		b.Fatalf("could not read file %s: %v", example, err)
	}

	for i := 0; i < b.N; i++ {
		p := toast.NewParser(bytes.NewReader(input))
		p.Parse() // nolint: errcheck // ignore error for benchmark
	}

	b.ReportAllocs()
}

func expectEnd(t *testing.T, s *gmitxt.Scanner, num uint32) {
	if s.Scan() {
		t.Errorf("Line %d: scanner should be finished", s.Line())
	}

	if s.Line().Num != num {
		t.Errorf("Line number was expected to be %d, but got: %d",
			num, s.Line().Num)
	}

	if s.Err() != nil {
		t.Errorf("Line %d: encountered error unexpected error: %v",
			s.Line().Num, s.Err())
	}

	if len(s.Line().Text) != 0 {
		t.Errorf("Line %d: end text should be an empty string, got: `%s`",
			s.Line().Num, s.Line().Text)
	}

	if len(s.Line().URL) != 0 {
		t.Errorf("Line %d: end url should be an empty string, got: `%s`",
			s.Line().Num, s.Line().URL)
	}
}

func expectStart(t *testing.T, s *gmitxt.Scanner) {
	if s.Line().Num != 0 {
		t.Errorf("Initial scanner should start at line number 0, got: %d",
			s.Line().Num)
	}

	if s.Err() != nil {
		t.Errorf("Initial scanner Err should be nil, got: %v", s.Err())
	}

	if len(s.Line().Text) != 0 {
		t.Errorf("Initial scanner text should be an empty string, got: `%x`",
			s.Line().Text)
	}

	if len(s.Line().URL) != 0 {
		t.Errorf("Initial scanner url be an empty string, got: `%x`",
			s.Line().URL)
	}
}

func expectLine(
	t *testing.T,
	s *gmitxt.Scanner,
	num uint32,
	typ gmitxt.LineType,
	expected string,
) {
	s.Scan()

	if s.Line().Num != num {
		t.Errorf("Line number was expected to be %d, but got: %d",
			num, s.Line().Num)
	}

	if s.Err() != nil {
		t.Errorf("Line %d: encountered error unexpected error: %v",
			s.Line().Num, s.Err())
	}

	if s.Line().Type != typ {
		t.Errorf("Line %d: type was not detected as %s, got: %s",
			s.Line().Num, typ, s.Line().Type)
	}

	if !bytes.Equal(s.Line().Text, []byte(expected)) {
		t.Errorf("Line %d: bytes do not match %x got: %x",
			s.Line().Num, []byte(expected), s.Line().Text)
	}

	if len(s.Line().URL) != 0 {
		t.Errorf("Line %d: url should be empty; got: %x",
			s.Line().Num, s.Line().URL)
	}
}

func expectLink(t *testing.T, s *gmitxt.Scanner, num uint32, url, text string) {
	s.Scan()

	if s.Line().Num != num {
		t.Errorf("Line number was expected to be %d, but got: %d",
			num, s.Line().Num)
	}

	if s.Err() != nil {
		t.Errorf("Line %d: encountered error unexpected error: %v",
			s.Line().Num, s.Err())
	}

	if s.Line().Type != gmitxt.Link {
		t.Errorf("Line %d: type was not detected as Link, got: %s",
			s.Line().Num, s.Line().Type)
	}

	if !bytes.Equal(s.Line().URL, []byte(url)) {
		t.Errorf("Line %d: url do not match %x got: %x",
			s.Line().Num, []byte(url), s.Line().URL)
	}

	if !bytes.Equal(s.Line().Text, []byte(text)) {
		t.Errorf("Line %d: text do not match %v got: %v",
			s.Line().Num, []byte(text), s.Line().Text)
	}
}
