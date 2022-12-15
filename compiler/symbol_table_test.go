package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolTable_Define(t *testing.T) {
	global := NewSymbolTable()

	assert.Equal(t, Symbol{
		Name:  "a",
		Scope: GlobalScope,
		Index: 0,
	}, global.Define("a"))
	assert.Equal(t, Symbol{
		Name:  "b",
		Scope: GlobalScope,
		Index: 1,
	}, global.Define("b"))
}

func TestSymbolTable_Resolve(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	type expected struct {
		symbol Symbol
		ok     bool
	}
	for _, tt := range []struct {
		name     string
		expected expected
	}{
		{
			name: "a",
			expected: expected{
				symbol: Symbol{
					Name:  "a",
					Scope: GlobalScope,
					Index: 0,
				},
				ok: true,
			},
		},
		{
			name: "b",
			expected: expected{
				symbol: Symbol{},
				ok:     false,
			},
		},
	} {
		symbol, ok := global.Resolve(tt.name)
		assert.Equal(t, tt.expected.ok, ok)
		assert.Equal(t, tt.expected.symbol, symbol)
	}
}
