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
)

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
