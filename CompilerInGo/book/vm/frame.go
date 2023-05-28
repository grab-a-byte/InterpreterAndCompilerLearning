package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	function           *object.CompiledFunction
	instructionPointer int
	basePointer        int
}

func NewFrame(fn *object.CompiledFunction, baseInstructionPointer int) *Frame {
	return &Frame{
		function:           fn,
		instructionPointer: -1,
		basePointer:        baseInstructionPointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.function.Instructions
}
