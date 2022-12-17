package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/compiler"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/object"
	"github.com/karamaru-alpha/monkey/parser"
)

func TestVM(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected interface{}
	}{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"2 - 1", 1},
		{"4 / 2", 2},
		{"2 * 2", 4},
		{"true", true},
		{"false", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"2 > 1", true},
		{"2 < 1", false},
		{"true == true", true},
		{"!true", false},
		{"-1", -1},
		{"if (true) {1}", 1},
		{"if (false) {1} else {2}", 2},
		{"if (false) {1}", nil},
		{"let a = 1; a", 1},
		{`"kara"+"maru"`, "karamaru"},
		{`[1, 2]`, []int{1, 2}},
		{`[1, 2][0]`, 1},
		{`{1: 2}[1]`, 2},
	} {
		program := parser.New(lexer.New(tt.input)).ParseProgram()

		c := compiler.New()
		assert.NoError(t, c.Compile(program))

		vm := New(c.Bytecode())
		assert.NoError(t, vm.Run())

		stackElem := vm.LastPoppedStackElem()
		switch expected := tt.expected.(type) {
		case int:
			assert.Equal(t, int64(expected), stackElem.(*object.Integer).Value)
		case string:
			assert.Equal(t, expected, stackElem.(*object.String).Value)
		case bool:
			assert.Equal(t, expected, stackElem.(*object.Boolean).Value)
		case nil:
			_, ok := stackElem.(*object.Null)
			assert.True(t, ok)
		case []int:
			for i, e := range stackElem.(*object.Array).Elements {
				assert.Equal(t, int64(expected[i]), e.(*object.Integer).Value)
			}
		case map[int]int:
			for _, pair := range stackElem.(*object.Hash).Pairs {
				assert.Equal(t, expected[int(pair.Key.(*object.Integer).Value)], int(pair.Value.(*object.Integer).Value))
			}
		}
	}
}
