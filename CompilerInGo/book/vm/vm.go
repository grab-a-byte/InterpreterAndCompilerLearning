package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048
const GlobalsSize uint = 65536
const MaxFrames = 1024

var trueObj = &object.Boolean{Value: true}
var falseObj = &object.Boolean{Value: false}
var nullObj = &object.Null{}

type VM struct {
	constants []object.Object

	stack        []object.Object
	stackPointer int // points to next free spot or last popped
	globals      []object.Object

	frames     []*Frame
	frameIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)
	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack:        make([]object.Object, StackSize),
		stackPointer: 0,

		globals: make([]object.Object, GlobalsSize),

		frames:     frames,
		frameIndex: 1,
	}
}

func NewWithGlobalStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.stackPointer]
}

func (vm *VM) Run() error {

	var ip int
	var ins code.Instructions
	var op code.Opcode

	for vm.currentFrame().instructionPointer < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().instructionPointer++

		ip = vm.currentFrame().instructionPointer
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().instructionPointer += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(trueObj)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(falseObj)
			if err != nil {
				return err
			}
		case code.OpNull:
			err := vm.push(nullObj)
			if err != nil {
				return err
			}

		case code.OpReturnValue:
			returnValue := vm.pop()
			frame := vm.popFrame()
			vm.stackPointer = frame.basePointer - 1
			err := vm.push(returnValue)
			if err != nil {
				return err
			}

		case code.OpReturn:
			frame := vm.popFrame()
			vm.stackPointer = frame.basePointer - 1
			err := vm.push(nullObj)
			if err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().instructionPointer = pos - 1

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().instructionPointer += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().instructionPointer = pos - 1
			}
		case code.OpAdd, code.OpSub, code.OpDiv, code.OpMul:
			vm.executeBinaryOperation(op)
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparrison(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().instructionPointer += 2
			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().instructionPointer += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().instructionPointer += 2
			array := vm.buildArray(vm.stackPointer-numElements, vm.stackPointer)

			vm.stackPointer = vm.stackPointer - numElements
			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().instructionPointer += 2
			hash, err := vm.buildHash(vm.stackPointer-numElements, vm.stackPointer)
			if err != nil {
				return err
			}
			vm.stackPointer = vm.stackPointer - numElements
			err = vm.push(hash)
			if err != nil {
				return err
			}

		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()
			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}

		case code.OpCall:
			fn, ok := vm.stack[vm.stackPointer-1].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("tried to call a non-function")
			}

			frame := NewFrame(fn, vm.stackPointer)
			vm.pushFrame(frame)
			vm.stackPointer = frame.basePointer + fn.NumLocals

		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().instructionPointer += 1
			frame := vm.currentFrame()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()

		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().instructionPointer += 1

			frame := vm.currentFrame()
			localValue := vm.stack[frame.basePointer+int(localIndex)]
			err := vm.push(localValue)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()
		}
	}

	return nil
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

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.frameIndex] = f
	vm.frameIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.frameIndex--
	return vm.frames[vm.frameIndex]
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unable to execute minus operator on type %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()
	if isTruthy(operand) {
		return vm.push(falseObj)
	} else {
		return vm.push(trueObj)
	}
}

func (vm *VM) executeComparrison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	}

	return fmt.Errorf("unable to do comparrison of types %T and %T", left, right)
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal == rightVal))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal != rightVal))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftVal > rightVal))
	default:
		return fmt.Errorf("unknown operator on integers: %d", op)
	}
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeIntegerBinaryOperation(op, left, right)
	}

	if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return vm.executeStringBinaryOperation(op, left, right)
	}

	return fmt.Errorf("unknown operator %d on type %s and %s", op, leftType, rightType)
}

func (vm *VM) executeStringBinaryOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unable to do operation %q on strings", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

func (vm *VM) executeIntegerBinaryOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unable to do operation %d on integers", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("unable to execute index on type %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(left, index object.Object) error {
	arr := left.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arr.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(nullObj)
	}

	return vm.push(arr.Elements[i])
}

func (vm *VM) executeHashIndex(left, index object.Object) error {
	hashObj := left.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unable to use type %s for an index", index.Type())
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return vm.push(nullObj)
	}

	return vm.push(pair.Value)
}

func (vm *VM) push(obj object.Object) error {
	if vm.stackPointer >= StackSize {
		return fmt.Errorf("tried to push but stack is full")
	}

	vm.stack[vm.stackPointer] = obj
	vm.stackPointer++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.stackPointer-1]
	vm.stackPointer--
	return o
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]
		pair := object.HashPair{Key: key, Value: value}

		haskKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unable to has key %s", key.Type())
		}

		hashedPairs[haskKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return trueObj
	}

	return falseObj
}
