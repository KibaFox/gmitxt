package gmitxt_test

import (
	"testing"

	"git.sr.ht/~kiba/gmitxt"
)

func TestLineTypeString(t *testing.T) {
	if gmitxt.LineType(0).String() != "UNKNOWN" {
		t.Errorf("Expected `UNKNOWN` for line type 0, got: `%s`",
			gmitxt.LineType(0))
	}

	if gmitxt.Head1.String() != "Head1" {
		t.Errorf("Expected `Head1` for line type, got: `%s`", gmitxt.Head1)
	}

	if gmitxt.Head2.String() != "Head2" {
		t.Errorf("Expected `Head2` for line type, got: `%s`", gmitxt.Head2)
	}

	if gmitxt.Head3.String() != "Head3" {
		t.Errorf("Expected `Head3` for line type, got: `%s`", gmitxt.Head3)
	}

	if gmitxt.Text.String() != "Text" {
		t.Errorf("Expected `Text` for line type, got: `%s`", gmitxt.Text)
	}

	if gmitxt.Link.String() != "Link" {
		t.Errorf("Expected `Link` for line type, got: `%s`", gmitxt.Link)
	}

	if gmitxt.PreStart.String() != "PreStart" {
		t.Errorf("Expected `PreStart` for line type, got: `%s`",
			gmitxt.PreStart)
	}

	if gmitxt.PreBody.String() != "PreBody" {
		t.Errorf("Expected `Head1` for line type, got: `%s`", gmitxt.PreBody)
	}

	if gmitxt.PreEnd.String() != "PreEnd" {
		t.Errorf("Expected `Head1` for line type, got: `%s`", gmitxt.PreEnd)
	}

	if gmitxt.List.String() != "List" {
		t.Errorf("Expected `List` for line type, got: `%s`", gmitxt.List)
	}

	if gmitxt.Quote.String() != "Quote" {
		t.Errorf("Expected `Quote` for line type, got: `%s`", gmitxt.Quote)
	}
}
