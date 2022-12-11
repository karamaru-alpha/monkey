package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/code"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/object"
	"github.com/karamaru-alpha/monkey/parser"
)

func TestCompiler_Compile(t *testing.T) {
	type expected struct {
		constants    []interface{}
		instructions []code.Instructions
	}
	for _, tt := range []struct {
		input    string
		expected expected
	}{
		{
			input: "1 + 2",
			expected: expected{
				constants: []interface{}{1, 2},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "1; 2",
			expected: expected{
				constants: []interface{}{1, 2},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "1 / 2",
			expected: expected{
				constants: []interface{}{1, 2},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDiv),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "1 * 2",
			expected: expected{
				constants: []interface{}{1, 2},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpMul),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "true",
			expected: expected{
				constants: []interface{}{},
				instructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "false",
			expected: expected{
				constants: []interface{}{},
				instructions: []code.Instructions{
					code.Make(code.OpFalse),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "1 == 1",
			expected: expected{
				constants: []interface{}{1, 1},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpEqual),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "1 != 1",
			expected: expected{
				constants: []interface{}{1, 1},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpNotEqual),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "2 > 1",
			expected: expected{
				constants: []interface{}{2, 1},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "2 < 1",
			expected: expected{
				constants: []interface{}{1, 2},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "-1",
			expected: expected{
				constants: []interface{}{1},
				instructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpMinus),
					code.Make(code.OpPop),
				},
			},
		},
		{
			input: "!true",
			expected: expected{
				constants: []interface{}{},
				instructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpBang),
					code.Make(code.OpPop),
				},
			},
		},
	} {
		program := parser.New(lexer.New(tt.input)).ParseProgram()

		compiler := New()
		assert.NoError(t, compiler.Compile(program))

		bytecode := compiler.Bytecode()
		assert.Equal(t, concatInstructions(tt.expected.instructions), bytecode.Instructions)

		testConstants(t, tt.expected.constants, bytecode.Constants)
	}
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) {
	t.Helper()
	assert.Equal(t, len(expected), len(actual))
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			result := actual[i].(*object.Integer)
			assert.Equal(t, int64(constant), result.Value)
		}
	}
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}
