package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/object"
	"github.com/karamaru-alpha/monkey/parser"
)

func TestEval_IntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := Eval(program)
		assert.Equal(t, tt.expected, obj.(*object.Integer).Value)
	}
}

func TestEval_BooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := Eval(program)
		assert.Equal(t, tt.expected, obj.(*object.Boolean).Value)
	}
}

func TestEval_BangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := Eval(program)
		assert.Equal(t, tt.expected, obj.(*object.Boolean).Value)
	}
}

func TestEval_MinusOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-1", -1},
		{"-5", -5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := Eval(program)
		assert.Equal(t, tt.expected, obj.(*object.Integer).Value)
	}
}

func TestEval_InfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1+2", 3},
		{"2*3", 6},
		{"4/2", 2},
		{"4-2", 2},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := Eval(program)
		assert.Equal(t, tt.expected, obj.(*object.Integer).Value)
	}
}
