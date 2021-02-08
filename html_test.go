package gmitxt_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"git.sr.ht/~kiba/gmitxt"
)

func TestToHTML(t *testing.T) {
	t.Log("Gemini headings to HTML")
	expectHTML(t, `#Head1
## This is a heading 2
### Head3
#### Head4?
`, `<h1>Head1</h1>
<h2>This is a heading 2</h2>
<h3>Head3</h3>
<h3># Head4?</h3>
`)

	t.Log("Gemini text to HTML")
	expectHTML(t, `I'm a line of text.

Look at me go!
`, `I'm a line of text.
<br>
<br>Look at me go!
`)

	t.Log("Gemini links to HTML")
	expectHTML(t, `=> gemini://gemini.circumlunar.space/ Project Gemini
=> testdata/example.gmi Example
=> https://example.tld/
`, `<a href="gemini://gemini.circumlunar.space/">Project Gemini</a>
<a href="testdata/example.gmi">Example</a>
<a href="https://example.tld/">https://example.tld/</a>
`)

	t.Log("Gemini preformatted text to HTML")
	expectHTML(t, "```txt\n\tPreformatted text.\n```ignore", `<pre class="txt">
	Preformatted text.
</pre>
`)
}

func BenchmarkToHTML(b *testing.B) {
	input, err := ioutil.ReadFile(example)
	if err != nil {
		b.Fatalf("could not read file %s: %v", example, err)
	}

	buf := &bytes.Buffer{}
	buf.Grow(1024 * 1024)

	for i := 0; i < b.N; i++ {
		scanner := gmitxt.NewScanner(bytes.NewReader(input))

		converter, err := gmitxt.NewConverter(scanner, buf)
		if err != nil {
			b.Fatalf("Unexpected error creating new HTML converter: %v", err)
		}

		if err := converter.Convert(); err != nil {
			b.Fatalf("Unexpected error converting gmi to HTML: %v", err)
		}
	}

	b.ReportAllocs()
}

func expectHTML(t *testing.T, input, expected string) {
	var out strings.Builder

	scanner := gmitxt.NewScanner(strings.NewReader(input))

	converter, err := gmitxt.NewConverter(scanner, &out)
	if err != nil {
		t.Fatalf("Unexpected error creating new HTML converter: %v", err)
	}

	if err := converter.Convert(); err != nil {
		t.Fatalf("Unexpected error converting gmi to HTML: %v", err)
	}

	if out.String() != expected {
		t.Errorf("Expected:\n%s\n...to convert to:\n%s\n...but got:\n%s",
			input, expected, out.String())
	}
}
