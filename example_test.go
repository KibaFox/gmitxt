package gmitxt_test

import (
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~kiba/gmitxt"
)

const geminiText = `# Example Gemini
This is a line of text.
=> gemini://gemini.circumlunar.space/ Gemini`

// Using Scanner to read Gemini text line by line.
func ExampleScanner() {
	scanner := gmitxt.NewScanner(strings.NewReader(geminiText))
	for scanner.Scan() {
		if scanner.Line().Type == gmitxt.Link {
			fmt.Printf("line %d: %s: url %s: %s\n",
				scanner.Line().Num,
				scanner.Line().Type,
				scanner.Line().URL,
				scanner.Line().Text,
			)
		} else {
			fmt.Printf("line %d: %s: %s\n",
				scanner.Line().Num,
				scanner.Line().Type,
				scanner.Line().Text,
			)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Output: line 1: Head1: Example Gemini
	// line 2: Text: This is a line of text.
	// line 3: Link: url gemini://gemini.circumlunar.space/: Gemini
}
