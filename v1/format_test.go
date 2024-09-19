package mime

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		In         string
		Type, Base Type
		Params     map[string]string
		Err        func(string, error) error
	}{
		{
			In:     "text/plain",
			Type:   Type("text/plain"),
			Base:   Type("text/plain"),
			Params: map[string]string{},
		},
		{
			In:     "text/plain+json",
			Type:   Type("text/plain+json"),
			Base:   Type("text/plain+json"),
			Params: map[string]string{},
		},
		{
			In:   "text/plain+json;charset=utf8",
			Type: Type("text/plain+json;charset=utf8"),
			Base: Type("text/plain+json"),
			Params: map[string]string{
				"charset": "utf8",
			},
		},
		{
			In:   "text/plain+json; charset=utf8",
			Type: Type("text/plain+json;charset=utf8"),
			Base: Type("text/plain+json"),
			Params: map[string]string{
				"charset": "utf8",
			},
		},
		{
			In:   "text/plain+json; charset=utf8; alabama=state",
			Type: Type("text/plain+json;alabama=state;charset=utf8"),
			Base: Type("text/plain+json"),
			Params: map[string]string{
				"charset": "utf8",
				"alabama": "state",
			},
		},
		{
			In: "text?plain+json; charset=utf8; alabama=state",
			Err: func(str string, err error) error {
				if err != nil {
					return nil
				} else {
					return errors.New("Expected an error: " + str)
				}
			},
		},
	}
	for i, e := range tests {
		mt, p, err := Parse(e.In)
		fmt.Printf("%s → %v %v; %v\n", e.In, mt, p, err)
		if e.Err != nil {
			assert.NoError(t, e.Err(e.In, err))
		} else if assert.NoError(t, err) {
			assert.Equal(t, e.Type, mt, "#%d", i)
			assert.Equal(t, e.Base, mt.Base(), "#%d", i)
			assert.Equal(t, e.Params, p, "#%d", i)
		}
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		A, B         Type
		Equal, Match bool
	}{
		{
			A:     Type("text/plain"),
			B:     Type("text/plain"),
			Equal: true,
			Match: true,
		},
		{
			A:     Type("text/plain;charset=utf8"),
			B:     Type("text/plain"),
			Equal: false,
			Match: true,
		},
		{
			A:     Type("text/plain;charset=utf8"),
			B:     Type("text/plain;charset=utf8"),
			Equal: true,
			Match: true,
		},
		{
			A:     Type("text/markdown"),
			B:     Type("text/plain"),
			Equal: false,
			Match: false,
		},
	}
	for i, e := range tests {
		assert.Equal(t, e.Equal, e.A.Equals(e.B), "#%d", i)
		assert.Equal(t, e.Match, e.A.Matches(e.B), "#%d", i)
	}
}
