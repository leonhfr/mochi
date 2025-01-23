package ast

import "github.com/yuin/goldmark/ast"

// CalloutType represents the callout type.
type CalloutType int

// nolint:revive
const (
	Note     CalloutType = iota // note (default)
	Abstract                    // abstract, summary, tldr
	Info                        // info
	Todo                        // todo
	Tip                         // tip, hint, important
	Success                     // success, check, done
	Question                    // question, help, faq
	Warning                     // warning, caution, attention
	Failure                     // failure, fail, missing
	Danger                      // danger, error
	Bug                         // bug
	Example                     // example
	Quote                       // quote, cite
)

// ParseCalloutType parses a callout type from a byte slice.
func ParseCalloutType(calloutType []byte) CalloutType {
	if t, ok := calloutAliases[string(calloutType)]; ok {
		return t
	}
	return Note
}

var calloutAliases = map[string]CalloutType{
	"note":      Note,
	"abstract":  Abstract,
	"summary":   Abstract,
	"tldr":      Abstract,
	"info":      Info,
	"todo":      Todo,
	"tip":       Tip,
	"hint":      Tip,
	"important": Tip,
	"success":   Success,
	"check":     Success,
	"done":      Success,
	"question":  Question,
	"help":      Question,
	"faq":       Question,
	"warning":   Warning,
	"caution":   Warning,
	"attention": Warning,
	"failure":   Failure,
	"fail":      Failure,
	"missing":   Failure,
	"danger":    Danger,
	"error":     Danger,
	"bug":       Bug,
	"example":   Example,
	"quote":     Quote,
	"cite":      Quote,
}

// Fold represents whether the callout is foldable and its default state.
type Fold int

// nolint:revive
const (
	None      Fold = iota // <empty> (default)
	Expanded              // +
	Collapsed             // -
)

// ParseFold parses a fold from a byte slice.
func ParseFold(fold []byte) Fold {
	switch string(fold) {
	case "+":
		return Expanded
	case "-":
		return Collapsed
	default:
		return None
	}
}

// Callout represents a callout node.
type Callout struct {
	ast.BaseBlock
	Title       []byte
	CalloutType CalloutType
	Fold        Fold
}

var _ ast.Node = (*Callout)(nil)

// NewCallout returns a new Callout node.
func NewCallout(title []byte, calloutType CalloutType, fold Fold) *Callout {
	return &Callout{
		Title:       title,
		CalloutType: calloutType,
		Fold:        fold,
	}
}

// KindCallout is a NodeKind of the Callout node.
var KindCallout = ast.NewNodeKind("Callout")

// Kind implements Node.Kind.
func (n *Callout) Kind() ast.NodeKind {
	return KindCallout
}

// Dump implements Node.Dump.
func (n *Callout) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
