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
	MaxFrames   = 1024
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

type VM struct {
	constants []object.Object
	stack     []object.Object
	sp        int // Always points to the next value. top of stack is stack[sp-1]
	globals   []object.Object

	frames     []*Frame
	frameIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame
	return &VM{
		constants:  bytecode.Constants,
		stack:      make([]object.Object, StackSize),
		sp:         0,
		globals:    make([]object.Object, GlobalsSize),
		frames:     frames,
		frameIndex: 1,
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
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for v.currentFrame().ip < len(v.currentFrame().Instructions())-1 {
		v.currentFrame().ip++

		ip = v.currentFrame().ip
		ins = v.currentFrame().Instructions()
		op = code.Opcode(ins[ip])
		switch op {
		case code.OpConstant:
			constIndex := int(binary.BigEndian.Uint16(ins[ip+1:]))
			v.currentFrame().ip += 2

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
			if err := v.executeComparison(op); err != nil {
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
			position := int(binary.BigEndian.Uint16(ins[ip+1:]))
			v.currentFrame().ip = position - 1
		case code.OpJumpNotTruthy:
			position := int(binary.BigEndian.Uint16(ins[ip+1:]))
			v.currentFrame().ip += 2

			condition := v.pop()
			if !isTruthy(condition) {
				v.currentFrame().ip = position - 1
			}
		case code.OpSetGlobal:
			globalIndex := int(binary.BigEndian.Uint16(ins[ip+1:]))
			v.currentFrame().ip += 2
			v.globals[globalIndex] = v.pop()
		case code.OpGetGlobal:
			globalIndex := int(binary.BigEndian.Uint16(ins[ip+1:]))
			v.currentFrame().ip += 2
			if err := v.push(v.globals[globalIndex]); err != nil {
				return err
			}
		case code.OpArray:
			numElements := int(binary.BigEndian.Uint16(ins[ip+1:]))
			v.currentFrame().ip += 2

			array := v.buildArray(v.sp-numElements, v.sp)
			if err := v.push(array); err != nil {
				return err
			}
		case code.OpHash:
			numElements := int(binary.BigEndian.Uint16(ins[ip+1:]))
			v.currentFrame().ip += 2

			hash, err := v.buildHash(v.sp-numElements, v.sp)
			if err != nil {
				return err
			}
			if err := v.push(hash); err != nil {
				return err
			}
		case code.OpIndex:
			index := v.pop()
			left := v.pop()
			if err := v.executeIndexExpression(left, index); err != nil {
				return err
			}
		case code.OpCall:
			fn, ok := v.stack[v.sp-1].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling non-function")
			}
			frame := NewFrame(fn)
			v.pushFrame(frame)
		case code.OpReturn:
			returnValue := v.pop()
			v.popFrame()
			v.pop()
			if err := v.push(returnValue); err != nil {
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
	if left.Type() == object.STRING && right.Type() == object.STRING {
		return v.executeBinaryStringOperation(op, left, right)
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

func (v *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return v.push(&object.String{Value: leftValue + rightValue})
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

func (v *VM) buildArray(startIdx, endIdx int) object.Object {
	elements := make([]object.Object, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		elements[i-startIdx] = v.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (v *VM) buildHash(startIdx, endIdx int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)
	for i := startIdx; i < endIdx; i += 2 {
		key := v.stack[i]
		value := v.stack[i+1]
		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}
		hashedPairs[hashKey.HashKey()] = pair
	}
	return &object.Hash{Pairs: hashedPairs}, nil
}

func (v *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		if err := v.executeArrayIndex(left, index); err != nil {
			return err
		}
		return nil
	case left.Type() == object.HASH:
		if err := v.executeHashIndex(left, index); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("invalid index. left: %s, index: %s", left.Type(), index.Type())
}

func (v *VM) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if i < 0 || i > max {
		return v.push(Null)
	}
	return v.push(arrayObject.Elements[i])
}

func (v *VM) executeHashIndex(array, index object.Object) error {
	hashObject := array.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return v.push(Null)
	}
	return v.push(pair.Value)
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

func (v *VM) currentFrame() *Frame {
	return v.frames[v.frameIndex-1]
}

func (v *VM) pushFrame(f *Frame) {
	v.frames[v.frameIndex] = f
	v.frameIndex++
}

func (v *VM) popFrame() *Frame {
	v.frameIndex--
	return v.frames[v.frameIndex]
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
