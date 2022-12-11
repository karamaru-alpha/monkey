package vm

import (
	"encoding/binary"
	"errors"

	"github.com/karamaru-alpha/monkey/code"
	"github.com/karamaru-alpha/monkey/compiler"
	"github.com/karamaru-alpha/monkey/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int // Always points to the next value. top of stack is stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (v *VM) StackTop() object.Object {
	if v.sp == 0 {
		return nil
	}
	return v.stack[v.sp-1]
}

func (v *VM) Run() error {
	for i := 0; i < len(v.instructions); i++ {
		op := code.Opcode(v.instructions[i])
		switch op {
		case code.OpConstant:
			constIndex := int(binary.BigEndian.Uint16(v.instructions[i+1:]))
			i += 2

			if err := v.push(v.constants[constIndex]); err != nil {
				return err
			}
		case code.OpAdd:
			left := v.pop()
			right := v.pop()
			v.push(&object.Integer{Value: left.(*object.Integer).Value + right.(*object.Integer).Value})
		}
	}
	return nil
}

func (v *VM) push(obj object.Object) error {
	if v.sp >= StackSize {
		return errors.New("stack overflow")
	}

	v.stack[v.sp] = obj
	v.sp++

	return nil
}

func (v *VM) pop() object.Object {
	obj := v.stack[v.sp-1]
	v.sp--
	return obj
}
