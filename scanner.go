package gmitxt

import (
	"bufio"
	"bytes"
	"io"
)

// Scanner provides an interface for reading Gemini formatted text.  Text is
// assumed to be UTF-8 encoded.  This is a line-based scanner where new lines
// are delimited by either CRLF (\r\n DOS/Windows format) or LF (\n UNIX
// format).
//
// This uses the bufio.Scanner from Go's standard library to scan input text
// line by line.  Each successive call to the Scan method will step through the
// lines of the input text.  Each line can be identified by it's Gemini type by
// calling the Type method.
//
// Scanning stops unrecoverably at EOF, the first I/O error, or an input line
// too large to fit in the buffer.
//
// For reference, the text/gemini format is described here:
//
//     https://gemini.circumlunar.space/docs/specification.html
//
// Alternatively:
//
//     gemini://gemini.circumlunar.space/docs/specification.gmi
//
type Scanner struct {
	scan *bufio.Scanner // underlying bufio.Scanner used to scan lines
	line Line           // scanned Gemini line representation
	pre  bool           // are we in a preformatted text section?
	idx  int            // whitespace index to parse links
}

// NewScanner returns a new Scanner to read from r.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{scan: bufio.NewScanner(r)}
}

const (
	whitespace = " \t" // either space or tab
	tokHead3   = "###"
	tokHead2   = "##"
	tokHead1   = "#"
	tokLink    = "=>"
	tokPre     = "```"
	tokList    = "* "
	tokQuote   = ">"
)

// Scan advances the Scanner to the next line of text, which will then be
// available through the TextBytes or Text methods.  The Gemini line type
// scanned will be available through the Type method.  If the line type scanned
// was a link the URL for the link will be available via the URL and URLBytes
// methods. It returns false when the scan stops, either by reaching the end of
// the input or an error.  After Scan returns false, the Err method will return
// any error that occurred during scanning, except that if it was io.EOF, Err
// will return nil.
//
// Scan panics if the split function returns too many empty lines without
// advancing the input. This is a common error mode for scanners.
func (s *Scanner) Scan() bool {
	s.line.Text = nil
	s.line.URL = nil

	if !s.scan.Scan() {
		return false
	}

	s.line.Num++

	if s.pre {
		if bytes.HasPrefix(s.scan.Bytes(), []byte(tokPre)) {
			// End of preformatted text.
			s.line.Type = PreEnd
			s.pre = false

			return true
		}

		s.line.Type = PreBody
		s.line.Text = s.scan.Bytes()

		return true
	}

	switch {
	case bytes.HasPrefix(s.scan.Bytes(), []byte(tokHead3)):
		s.line.Type = Head3
		s.line.Text = trimLeftSpace(s.scan.Bytes()[3:])

		return true
	case bytes.HasPrefix(s.scan.Bytes(), []byte(tokHead2)):
		s.line.Type = Head2
		s.line.Text = trimLeftSpace(s.scan.Bytes()[2:])

		return true
	case bytes.HasPrefix(s.scan.Bytes(), []byte(tokHead1)):
		s.line.Type = Head1
		s.line.Text = trimLeftSpace(s.scan.Bytes()[1:])

		return true
	case bytes.HasPrefix(s.scan.Bytes(), []byte(tokLink)):
		s.line.Type = Link
		s.line.URL = trimLeftSpace(s.scan.Bytes()[2:])
		s.idx = bytes.IndexAny(s.line.URL, whitespace)

		if s.idx != -1 {
			s.line.Text = trimLeftSpace(s.line.URL[s.idx:])
			s.line.URL = s.line.URL[:s.idx]
		}

		return true
	case bytes.HasPrefix(s.scan.Bytes(), []byte(tokPre)):
		s.line.Type = PreStart
		s.line.Text = s.scan.Bytes()[3:]
		s.pre = true

		return true
	case bytes.HasPrefix(s.scan.Bytes(), []byte(tokList)):
		s.line.Type = List
		s.line.Text = s.scan.Bytes()[2:]

		return true
	case bytes.HasPrefix(s.scan.Bytes(), []byte(tokQuote)):
		s.line.Type = Quote
		s.line.Text = s.scan.Bytes()[1:]

		return true
	default:
		s.line.Type = Text
		s.line.Text = s.scan.Bytes()

		return true
	}
}

// trimLeftSpace is trims any whitespace to the left in the input byte slice.
// This returns nil if the input byte slice is all whitespace.
// This is a replacement for bytes.TrimLeft that doesn't allocate.
func trimLeftSpace(b []byte) []byte {
	if len(b) == 0 || !isWhitespace(b[0]) {
		return b
	}

	for idx, char := range b {
		if isWhitespace(char) {
			continue
		}

		return b[idx:]
	}

	return nil
}

// isWhitespace returns whether a byte character is a whitespace.  The Gemini
// specification defines whitespace as either a space or a tab character.
func isWhitespace(char byte) bool {
	if char == ' ' || char == '\t' {
		return true
	}

	return false
}

// Line returns the parsed Line of Gemini text that was just scanned by the Scan
// method.  The underlying data will be overwritten by subsequent calls to Scan.
// This does not allocate memory.
func (s *Scanner) Line() Line {
	return s.line
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *Scanner) Err() error {
	return s.scan.Err()
}

// Buffer sets the initial buffer to use when scanning and the maximum size of
// buffer that may be allocated during scanning. The maximum input line size is
// the larger of max and cap(buf). If max <= cap(buf), Scan will use this
// buffer only and do no allocation.
//
// By default, Scan uses an internal buffer and sets the maximum token size to
// bufio.MaxScanTokenSize (64 kilobytes).
//
// Buffer panics if it is called after scanning has started.
func (s *Scanner) Buffer(buf []byte, max int) {
	s.scan.Buffer(buf, max)
}
