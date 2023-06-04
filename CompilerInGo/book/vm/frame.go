package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	closure            *object.Closure
	instructionPointer int
	basePointer        int
}

func NewFrame(closure *object.Closure, baseInstructionPointer int) *Frame {
	return &Frame{
		closure:            closure,
		instructionPointer: -1,
		basePointer:        baseInstructionPointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.closure.Fn.Instructions
}
