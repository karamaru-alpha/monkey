package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/compiler"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/object"
	"github.com/karamaru-alpha/monkey/parser"
)

func TestVM_IntegerArithmetic(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected interface{}
	}{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
	} {
		program := parser.New(lexer.New(tt.input)).ParseProgram()

		c := compiler.New()
		assert.NoError(t, c.Compile(program))

		vm := New(c.Bytecode())
		assert.NoError(t, vm.Run())

		stackElm := vm.StackTop()
		switch expected := tt.expected.(type) {
		case int:
			assert.Equal(t, int64(expected), stackElm.(*object.Integer).Value)
		}
	}
}
