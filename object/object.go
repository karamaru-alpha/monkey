package object

import (
	"fmt"
	"strconv"
)

type Type int64

type Object interface {
	Type() Type
	Inspect() string
}

const (
	INTEGER Type = iota + 1
	BOOLEAN
	NULL
	RETURN_VALUE
	ERROR
)

func (typ Type) String() string {
	switch typ {
	case INTEGER:
		return "INTEGER"
	case BOOLEAN:
		return "BOOLEAN"
	case NULL:
		return "NULL"
	case RETURN_VALUE:
		return "RETURN_VALUE"
	case ERROR:
		return "ERROR"
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

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
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
