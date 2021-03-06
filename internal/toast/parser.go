package toast

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Parser handles parsing of a text/gemini document
type Parser struct {
	s     *bufio.Scanner
	Lines []Line
	toc   []*Heading // see toc.go
	pft   bool
}

// NewParser instantiates a new parser from a reader
func NewParser(r io.Reader) *Parser {
	return &Parser{bufio.NewScanner(r), nil, nil, false}
}

// Parse will parse the provided input until termination
// The returned error will be non-nil only if a significant error occured
func (p *Parser) Parse() error {
	for p.s.Scan() {
		p.Lines = append(p.Lines, p.parseLine(p.s.Text()))
	}
	return p.s.Err()
}

func get(s string, i int) (r rune) {
	if i >= len(s) {
		return
	}
	for ii, c := range s {
		if i == ii {
			return c
		}
	}
	return // should never happen
}

func (p *Parser) parseLine(l string) Line {
	a, b, c := get(l, 0), get(l, 1), get(l, 2)

	// preformatted mode
	if p.pft {
		if a == '`' && b == '`' && c == '`' {
			p.pft = false
			var s string
			return (*PreformatToggleLine)(&s)
		}
		return (*PreformatLine)(&l) // the only time trailing whitespace is not stripped, though line endings still are
	}

	// normal mode
	switch a {
	case '=':
		if b == '>' {
			return parseLink(l)
		}
	case '`':
		if b == '`' && c == '`' {
			p.pft = true
			s := strings.TrimPrefix(l, "```")
			s = strings.TrimSpace(s)
			return (*PreformatToggleLine)(&s)
		}
	case '#':
		return parseHeading(l)
	case '*':
		if b == ' ' {
			return parseUList(l)
		}
	case '>':
		return parseQuote(l)
	}
	return (*TextLine)(&l)
}

// whitespace is removed from the end of the line too
// because the spec doesn't really say what the line endings are meant to be for text/gemini itself
// and whitespace at the end of a line is basically meaningless anyway in the presence of empty text lines

func parseLink(l string) Line {
	l = strings.TrimPrefix(l, "=>")
	l = strings.TrimSpace(l)

	// is there whitespace in what's left?
	w := strings.IndexFunc(l, unicode.IsSpace)
	if w == -1 {
		return &LinkLine{l, ""}
	}
	return &LinkLine{
		l[:w],
		strings.TrimLeftFunc(l[w:], unicode.IsSpace),
	}
}

func parseHeading(l string) Line {
	s := len(l)
	if len(l) > 3 {
		s = 3
	}
	level := strings.Count(l[:s], "#")
	l = l[level:]
	l = strings.TrimSpace(l)
	return &HeadingLine{l, level}
}

func parseUList(l string) Line {
	l = strings.TrimPrefix(l, "*")
	l = strings.TrimSpace(l)
	return (*UnorderedListLine)(&l)
}

func parseQuote(l string) Line {
	l = strings.TrimPrefix(l, ">")
	l = strings.TrimSpace(l)
	return (*QuoteLine)(&l)
}

// For TOC stuff, see toc.go
