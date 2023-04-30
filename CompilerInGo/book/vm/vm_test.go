package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("Object is not an integer. Got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("Integer incorrect value, expected %d, got %d", expected, result.Value)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("Object is not a boolean. Got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("Boolean incorrect value, expected %t, got %t", expected, result.Value)
	}

	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)

		if err != nil {
			t.Fatalf("Error while compiling program: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()

		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)

	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObjectFailed: %s", err)
		}
	case object.Null:
		if actual != nullObj {
			t.Errorf("expected object to be null, got (%+v)", actual)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"2 - 1", 1},
		{"1 - 2", -1},
		{"4 / 2", 2},
		{"2 * 2", 4},
		{"2 * 2 * 2", 8},
		{"(2 * 2) / 4", 1},
		{"5 + 2 * 10", 25},
		{"-5", -5},
		{"-10", -10},
		{"-5 + 10", 5},
		{"-50+100+-50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestBooleanLogic(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 2", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"false == true", false},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"(1 > 2) != true", true},
		{"(1 < 2) != false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!( if(false) { 5 } )", true},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if(true) {5}", 5},
		{"if(false) { 5 } else { 10 }", 10},
		{"if(1) { 5 }", 5},
		{"if(1 < 2) { 5 }", 5},
		{"if(1 < 2) { 10 } else { 20 }", 10},
		{"if(1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", nullObj},
		{"if(false) {5}", nullObj},
		{"if((if(false){ 10 })) { 10 } else { 20 }", 20},
	}

	runVmTests(t, tests)
}
