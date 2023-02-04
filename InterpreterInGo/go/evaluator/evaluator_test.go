package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"-10", -10},
		{"5 + 5", 10},
		{"-5 + 5", 0},
		{"-5 + 10 + 5 - 6", 4},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 /3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"true != true", false},
		{"false != false", false},
		{"true != false", true},
		{"false != true", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"true == (1 < 2)", true},
		{"true == (1 > 2)", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBangOpration(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if(10 < 20) { 10 } else { 20 }", 10},
		{"if(10 > 20) { 10 } else { 20 }", 20},
		{"if(10 > 20) { 10 }", nil},
		{"if(10 > 1){if(10 > 1){return 10} return 1}", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		expected, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(expected))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEvalReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 42", 42},
		{"15; return 75", 75},
		{"18; return 24; 3+4", 24},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 42; a", 42},
		{"let a = 5 * 5; a", 25},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.FunctionValue)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x }; identity(42)", 42},
		{"let identity = fn(x) {return x;}; identity(42)", 42},
		{"let double = fn(x) {return x + x;}; double(5)", 10},
		{"let add = fn(x,y) {return x + y;}; add(5,6)", 11},
		{"let add = fn(x,y) {return x + y;}; add(5 + 5, add(6, 6))", 22},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionClosures(t *testing.T) {
	input := `
	let sample = fn(x) {
		return fn(y) {
			return x + y;
		}
	}

	let test = sample(1);
	test(3)
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)

}

func TestBuiltInFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello world")`, 11},
		{`len("")`, 0},
		{`len("one", "two)`, "too many args for builtin function 'len'"},
		{`len(42)`, "Unsupported Type for builtin function 'len'"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if str, ok := tt.expected.(string); ok {
			err, ok := evaluated.(*object.ErrorValue)
			if !ok {
				t.Fatalf("Not an error object")
			}

			if err.Message != str {
				t.Fatalf("Built in error incorrect, expected %q found %q", str, err.Message)
			}

		} else if i, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(i))
		}
	}
}

func TestIndexingArrays(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let myArray = [1,2,3]; myArray[0]", 1},
		{"let myArray = [1,2,3]; myArray[1 + 1]", 3},
		{"let myArray = [1,2,3]; myArray[5]", nil},
		{"let myArray = [1,2,3]; myArray[3]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if tt.expected == nil && evaluated.Type() != object.NULL_OBJ {
			t.Fatalf("Expected NUll from out of range but got %q", evaluated.Type())
		} else {
			continue
		}

		expected := tt.expected.(int)
		obj, ok := evaluated.(*object.Integer)
		if !ok {
			t.Fatalf("Expected Integer, found %T", evaluated)
		}

		if obj.Value != int64(expected) {
			t.Fatalf("Expected %d got %d", expected, obj.Value)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input                string
		expectedErrorMessage string
	}{
		{"-true", "unknown operator: -BOOLEAN"},
		{"1 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"1 + true; 5", "type mismatch: INTEGER + BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if(10 > 2) {if(10 > 2){return true + true} else {1}}", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: 'foobar'"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errorObj, ok := evaluated.(*object.ErrorValue)
		if !ok {
			t.Fatalf("Expected errorMessage %s, found %T", tt.expectedErrorMessage, evaluated)
		}

		if errorObj.Message != tt.expectedErrorMessage {
			t.Fatalf("Incorrct error message, expected: %s, found: %s", tt.expectedErrorMessage, errorObj.Message)
		}
	}
}

func testNullObject(t *testing.T, evaluated object.Object) bool {
	if evaluated.Type() != object.NULL_OBJ {
		t.Errorf("Expected NULL object, got %T", evaluated)
		return false
	}

	return true
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) bool {
	integer, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("Expected Integer Object, found %T", evaluated)
		return false
	}

	if integer.Value != expected {
		t.Errorf("Wrong value parsed, expeced %d, found %d", expected, integer.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
	boolean, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("Expected Boolean Object, found %T", evaluated)
		return false
	}

	if boolean.Value != expected {
		t.Errorf("Wrong value parsed, expeced %v, found %v", expected, boolean.Value)
		return false
	}

	return true
}

func TestHashLiterals(t *testing.T) {
	input := `
		let two = "two";
		{
			"one": 10 - 9,
			"two": 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
			true: 5,
			false: 6,
		}
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("Eval didn't return %T. got=%T (%+v)", object.HashMap{}, evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		object.TRUE.HashKey():                      5,
		object.FALSE.HashKey():                     6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"hello": 1}["hello"]`, 1},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testEval(s string) object.Object {
	l := lexer.New(s)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}
