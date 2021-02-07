package gmitxt

// LineType describes the type of line in Gemini formatted text.
type LineType int

const (
	// Head1 is a heading level 1 line.  It is a line that starts with # and is
	// optionally followed by whitespace.
	Head1 LineType = iota + 1
	// Head2 is a heading level 2 line.  It is a line that starts with ## and is
	// optionally followed by whitespace.
	Head2
	// Head3 is a heading level 3 line.  It is a line that starts with ### and
	// is optionally followed by whitspace.
	Head3
	// Text is a normal text line.  This is the default line type.
	Text
	// Link is a line containing a link.  It is a line that starts with => and
	// is followed by a URL and optional text to describe the link, each
	// separated by whitespace.
	Link
	// PreStart is a line that starts preformatted text.  It is a line that
	// starts with ``` and may have additional alternative text.
	PreStart
	// PreBody is a line of text that should be rendered as preformatted text.
	// It's a line in-between lines that start with ```.
	PreBody
	// PreEnd is a line that ends preformatted text.  It's a line that starts
	// with ``` that comes after a previous line that starts with ```.  Any text
	// after the ``` should be ignored.
	PreEnd
	// List is an unordered list line.  It's a line that starts with * .
	List
	// Quote is a line containing a quoted text.  It's a line that starts with
	// >.
	Quote
)

// String returns the string representation of the line type.  For example, for
// Head3 it will return the string "Head3".
func (typ LineType) String() string {
	switch typ {
	case Head1:
		return "Head1"
	case Head2:
		return "Head2"
	case Head3:
		return "Head3"
	case Text:
		return "Text"
	case Link:
		return "Link"
	case PreStart:
		return "PreStart"
	case PreBody:
		return "PreBody"
	case PreEnd:
		return "PreEnd"
	case List:
		return "List"
	case Quote:
		return "Quote"
	default:
		return "UNKNOWN"
	}
}
