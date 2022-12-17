package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/code"
)

type Type int64

type Object interface {
	Type() Type
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

const (
	INTEGER Type = iota + 1
	STRING
	BOOLEAN
	NULL
	RETURN_VALUE
	ERROR
	FUNCTION
	ARRAY
	HASH
	BUILTIN
	COMPILED_FUNCTION
)

func (typ Type) String() string {
	switch typ {
	case INTEGER:
		return "INTEGER"
	case STRING:
		return "STRING"
	case BOOLEAN:
		return "BOOLEAN"
	case NULL:
		return "NULL"
	case RETURN_VALUE:
		return "RETURN_VALUE"
	case ERROR:
		return "ERROR"
	case FUNCTION:
		return "FUNCTION"
	case ARRAY:
		return "ARRAY"
	case HASH:
		return "HASH"
	case BUILTIN:
		return "BUILTIN"
	case COMPILED_FUNCTION:
		return "COMPILED_FUNCTION"
	}
	return "UNKNOWN"
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: INTEGER, Value: uint64(i.Value)}
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: STRING, Value: h.Sum64()}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	hashKey := HashKey{Type: BOOLEAN}
	if b.Value {
		hashKey.Value = 1
	} else {
		hashKey.Value = 0
	}
	return hashKey
}

type Null struct{}

func (n *Null) Type() Type {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() Type {
	return RETURN_VALUE
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s", e.Message)
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := make([]string, 0)
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() Type {
	return ARRAY
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := make([]string, 0)
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashKey struct {
	Type  Type
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() Type {
	return HASH
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type {
	return BUILTIN
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type CompiledFunction struct {
	Instructions code.Instructions
}

func (c *CompiledFunction) Type() Type {
	return COMPILED_FUNCTION
}

func (c *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompoledFunction[%p]", c)
}
