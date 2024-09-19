package mime

import (
	"encoding/json"
	"mime"
	"sort"
	"strings"
)

const (
	UTF_8 = "utf-8"
)

type Type string

const (
	Text     = Type("text/plain")
	Markdown = Type("text/markdown")
	HTML     = Type("text/html")
	JSON     = Type("application/json")
	CSV      = Type("text/csv")
	XML      = Type("text/xml")
	GZIP     = Type("application/gzip")
	Invalid  = Type("")
)

var MarkdownCompatible = Options{
	Text,
	Markdown,
}

// Parse parses a mimetype string and returns a normalized type and the
// parameters associated with it. If the type has any parameters, they
// are sorted and the canonical type string is rewritten.
func Parse(v string) (Type, map[string]string, error) {
	t, p, err := mime.ParseMediaType(v)
	if err != nil {
		return Invalid, nil, err
	}

	sb := &strings.Builder{}
	sb.WriteString(t)

	if l := len(p); l > 0 {
		keys := make([]string, 0, l)
		for k, _ := range p {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, e := range keys {
			sb.WriteString(";")
			sb.WriteString(e)
			sb.WriteString("=")
			sb.WriteString(p[e])
		}
	}

	return Type(sb.String()), p, nil
}

// Base strips any parameters that may be present off the end of the
// type and returns a new type representing its base.
func (t Type) Base() Type {
	if x := strings.Index(string(t), ";"); x >= 0 {
		return Type(strings.TrimSpace(string(t)[:x]))
	} else {
		return t
	}
}

// Equals compares the provided type to the receiver exactly, including
// all parameters.
func (t Type) Equals(s Type) bool {
	return strings.EqualFold(t.String(), s.String())
}

// Matches compares the provided type's base to the receiver type's base,
// effectively excluding any parameters that either might have.
func (t Type) Matches(s Type) bool {
	return strings.EqualFold(t.Base().String(), s.Base().String())
}

func (t Type) String() string {
	return string(t)
}

// Ext produces a filename extension (including the '.' separator) for
// a variety of known types.
func (t Type) Ext() string {
	switch t.Base() {
	case Invalid:
		return ""
	case Text:
		return ".txt"
	case Markdown:
		return ".md"
	case HTML:
		return ".html"
	case JSON:
		return ".json"
	case CSV:
		return ".csv"
	case XML:
		return ".xml"
	case GZIP:
		return ".gz"
	default:
		return t.firstExt()
	}
}

func (t Type) firstExt() string {
	e, err := mime.ExtensionsByType(string(t))
	if err != nil {
		return ""
	}
	if len(e) < 1 {
		return ""
	}
	return e[0]
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*t = Type(s)
	return nil
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t), nil
}

func (t *Type) UnmarshalText(text []byte) error {
	*t = Type(text)
	return nil
}

type Options []Type

func (o Options) Contains(t Type) bool {
	for _, e := range o {
		if e == t {
			return true
		}
	}
	return false
}

func (o Options) First(d Type) Type {
	if len(o) < 1 {
		return d
	} else {
		return o[0]
	}
}

func (o Options) String() string {
	b := &strings.Builder{}
	for i, e := range o {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e.String())
	}
	return b.String()
}
