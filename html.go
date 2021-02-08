package gmitxt

import (
	"errors"
	"fmt"
	"io"
)

// Converter scans through Gemini text and converts it to HTML.
type Converter struct {
	s *Scanner
	w io.Writer
}

// ErrNoScanner occurs when NewConverter is called without a Gemini Scanner.
var ErrNoScanner = errors.New(
	"programmer error: scanner not initialized for HTML encoder")

// NewConverter makes a new Gemini to HTML converter that will scan through the
// Gemini text with the provided scanner and write HTML to the provided writer.
func NewConverter(s *Scanner, w io.Writer) (*Converter, error) {
	if s == nil {
		return nil, ErrNoScanner
	}

	return &Converter{
		s: s,
		w: w,
	}, nil
}

// Convert writes the HTML version of the Gemini text to the stream.
func (c *Converter) Convert() error {
	var previous LineType

	for c.s.Scan() {
		if err := c.writeHTML(previous); err != nil {
			return err
		}

		previous = c.s.Type()
	}

	if err := c.s.Err(); err != nil {
		return fmt.Errorf("scan error converting to HTML: %w", err)
	}

	return nil
}

func (c *Converter) err(err error) error {
	return NewWriteErr(err, c.s.Type(), c.s.Line())
}

func (c *Converter) write(bytes []byte) error {
	if _, err := c.w.Write(bytes); err != nil {
		return c.err(err)
	}

	return nil
}

func (c *Converter) writeStr(str string) error {
	if _, err := io.WriteString(c.w, str); err != nil {
		return c.err(err)
	}

	return nil
}

// ErrNotSupportedType occurs when writeHTML is called on on a line of Gemini
// text that is not supported.
var ErrNotSupportedType = errors.New("programmer error writing HTML for " +
	"unsupported Gemini text line type")

func (c *Converter) writeHTML(previous LineType) error {
	switch c.s.Type() {
	case Head1, Head2, Head3:
		return c.writeHead()
	case Text:
		return c.writeText(previous)
	case Link:
		return c.writeLink()
	case PreStart, PreBody, PreEnd:
		return c.writePre()
	case List:
		return c.writeList(previous)
	case Quote:
		return c.writeQuote(previous)
	default:
		return NewWriteErr(ErrNotSupportedType, c.s.Type(), c.s.Line())
	}
}

// ErrNotHeader occurs when writeHead is called on on a line of Gemini text
// that is not a header.
var ErrNotHeader = errors.New("programmer error writing a header " +
	"for a Gemini line that's not a header")

func (c *Converter) writeHead() error {
	var start, end string

	switch c.s.Type() { // nolint: exhaustive // only handle headings here
	case Head1:
		start, end = "<h1>", "</h1>\n"
	case Head2:
		start, end = "<h2>", "</h2>\n"
	case Head3:
		start, end = "<h3>", "</h3>\n"
	default:
		return c.err(ErrNotHeader)
	}

	if err := c.writeStr(start); err != nil {
		return err
	}

	if err := c.write(c.s.TextBytes()); err != nil {
		return err
	}

	return c.writeStr(end)
}

// ErrNotText occurs when writeText is called on on a line of Gemini text
// that is not a text line.
var ErrNotText = errors.New("programmer error writing a text line " +
	"for a Gemini line that's not a text line")

func (c *Converter) writeText(previous LineType) error {
	if c.s.Type() != Text {
		return c.err(ErrNotText)
	}

	if previous == Text {
		if err := c.writeStr("<br>"); err != nil {
			return err
		}
	}

	if err := c.write(c.s.TextBytes()); err != nil {
		return err
	}

	return c.writeStr("\n")
}

// ErrNotLink occurs when writeLink is called on on a line of Gemini text
// that is not a link.
var ErrNotLink = errors.New("programmer error writing a link " +
	"for a Gemini line that's not a link")

func (c *Converter) writeLink() error {
	if c.s.Type() != Link {
		return c.err(ErrNotLink)
	}

	if err := c.writeStr("<a href=\""); err != nil {
		return err
	}

	if err := c.write(c.s.URLBytes()); err != nil {
		return err
	}

	if err := c.writeStr("\">"); err != nil {
		return err
	}

	if len(c.s.TextBytes()) > 0 {
		if err := c.write(c.s.TextBytes()); err != nil {
			return err
		}
	} else {
		if err := c.write(c.s.URLBytes()); err != nil {
			return err
		}
	}

	return c.writeStr("</a>\n")
}

// ErrNotPre occurs when writePre is called on on a line of Gemini text
// that is not a preformatted text line.
var ErrNotPre = errors.New("programmer error writing preformatted text " +
	"for a Gemini line that's not a preformatted text line")

func (c *Converter) writePre() error {
	switch c.s.Type() { // nolint: exhaustive // only handle pre here
	case PreStart:
		if len(c.s.TextBytes()) > 0 {
			if err := c.writeStr("<pre class=\""); err != nil {
				return err
			}

			if err := c.write(c.s.TextBytes()); err != nil {
				return err
			}

			return c.writeStr("\">\n")
		}

		return c.writeStr("<pre>\n")
	case PreBody:
		if err := c.write(c.s.TextBytes()); err != nil {
			return err
		}

		return c.writeStr("\n")
	case PreEnd:
		return c.writeStr("</pre>\n")
	default:
		return c.err(ErrNotPre)
	}
}

// ErrNotList occurs when writeList is called on on a line of Gemini text
// that is not a list line.
var ErrNotList = errors.New("programmer error writing a list " +
	"for a Gemini line that's not a list")

func (c *Converter) writeList(previous LineType) error {
	if c.s.Type() != List {
		return c.err(ErrNotList)
	}

	return nil // TODO
}

// ErrNotQuote occurs when writeQuote is called on on a line of Gemini text
// that is not a quote line.
var ErrNotQuote = errors.New("programmer error writing a quote " +
	"for a Gemini line that's not a quote")

func (c *Converter) writeQuote(previous LineType) error {
	if c.s.Type() != Quote {
		return c.err(ErrNotQuote)
	}

	return nil // TODO
}

// WriteErr wraps an error that occurs while writing.  It contains an underlying
// write error, the last scanned Gemini text line type, and the last
// scanned line number.
type WriteErr struct {
	// Err is the underlying error that occurred while writing.
	Err error
	// Type is the Gemini text line type that was scanned as input.
	Type LineType
	// Line is the line number of the Gemini text scanned.
	Line int
}

// NewWriteErr creates a new WriteErr from an error and uses the provided
// Scanner to set the last scanned LineType and Position.
func NewWriteErr(err error, typ LineType, line int) *WriteErr {
	return &WriteErr{
		Err:  err,
		Type: typ,
		Line: line,
	}
}

// Unwrap returns the underlying write error.
func (e WriteErr) Unwrap() error {
	return e.Err
}

// Error returns an error message with the line number and type scanned with the
// underlying error message.
func (e WriteErr) Error() string {
	return fmt.Sprintf("line %d: error writing %s: %s", e.Line, e.Type, e.Err)
}
