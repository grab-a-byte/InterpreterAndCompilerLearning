package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	function           *object.CompiledFunction
	instructionPointer int
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{
		function:           fn,
		instructionPointer: -1,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.function.Instructions
}
