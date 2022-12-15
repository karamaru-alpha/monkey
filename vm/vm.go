package vm

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/karamaru-alpha/monkey/code"
	"github.com/karamaru-alpha/monkey/compiler"
	"github.com/karamaru-alpha/monkey/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int // Always points to the next value. top of stack is stack[sp-1]
	globals      []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      make([]object.Object, GlobalsSize),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (v *VM) StackTop() object.Object {
	if v.sp == 0 {
		return nil
	}
	return v.stack[v.sp-1]
}

func (v *VM) LastPoppedStackElem() object.Object {
	return v.stack[v.sp]
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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := v.executeBinaryOperation(op); err != nil {
				return err
			}
		case code.OpTrue:
			if err := v.push(True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := v.push(False); err != nil {
				return err
			}
		case code.OpNull:
			if err := v.push(Null); err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			if err := v.push(True); err != nil {
				return err
			}
		case code.OpBang:
			if err := v.executeBangOperator(); err != nil {
				return err
			}
		case code.OpMinus:
			if err := v.executeMinusOperator(); err != nil {
				return err
			}
		case code.OpJump:
			position := int(binary.BigEndian.Uint16(v.instructions[i+1:]))
			i = position - 1
		case code.OpJumpNotTruthy:
			position := int(binary.BigEndian.Uint16(v.instructions[i+1:]))
			i += 2

			condition := v.pop()
			if !isTruthy(condition) {
				i = position - 1
			}
		case code.OpSetGlobal:
			globalIndex := int(binary.BigEndian.Uint16(v.instructions[i+1:]))
			i += 2
			v.globals[globalIndex] = v.pop()
		case code.OpGetGlobal:
			globalIndex := int(binary.BigEndian.Uint16(v.instructions[i+1:]))
			i += 2
			if err := v.push(v.globals[globalIndex]); err != nil {
				return err
			}
		case code.OpPop:
			v.pop()
		}
	}
	return nil
}

func (v *VM) executeBinaryOperation(op code.Opcode) error {
	right := v.pop()
	left := v.pop()

	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return v.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", left.Type(), right.Type())
}

func (v *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result int64
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	case code.OpMul:
		result = leftValue * rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return v.push(&object.Integer{Value: result})
}

func (v *VM) executeComparison(op code.Opcode) error {
	right := v.pop()
	left := v.pop()

	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return v.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return v.push(v.nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return v.push(v.nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (v *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch op {
	case code.OpEqual:
		return v.push(v.nativeBoolToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return v.push(v.nativeBoolToBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return v.push(v.nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (v *VM) nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (v *VM) executeBangOperator() error {
	operand := v.pop()
	switch operand {
	case True:
		return v.push(False)
	case False:
		return v.push(True)
	case Null:
		return v.push(True)
	default:
		return v.push(False)
	}
}

func (v *VM) executeMinusOperator() error {
	operand := v.pop()
	if operand.Type() != object.INTEGER {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}
	return v.push(&object.Integer{Value: -operand.(*object.Integer).Value})
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
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
