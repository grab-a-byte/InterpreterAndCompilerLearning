package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestParseProgram(t *testing.T) {
	input := `let x = 5;
	let y =10;
	let foobar=  09876;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("Program was not parsed and returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Expecetd 3 statements, found %d", len(program.Statements))
	}

	testcases := []struct{ expectedIdentifier string }{
		{"x"}, {"y"}, {"foobar"},
	}

	for i, tt := range testcases {
		if !testLetStatement(t, program.Statements[i], tt.expectedIdentifier) {
			t.Errorf("Failed Parsing %v", program.Statements[i])
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteral(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("Not a let token : %v", stmt)
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("Not a let statement")
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("Invalid name.Value, expected %s got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Invalid name.TokenLiteral(), expected %s got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteral(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpressions(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Program failed to parse and returned nil")
	}
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatal("Not enough Statements")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected to find expression statement, found type %T", program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected to find identifier but found %T", exp)
		return false
	}

	if ident.Value != value {
		t.Fatalf("Expected to find value %q, found %q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Fatalf("Expected to find TokenLiteral %q, found %q", value, ident.Value)
		return false
	}

	return true
}

func TestIntegerLterals(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Program failed to parse and returned nil")
	}
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatal("Not enough Statements")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected to find expression statement, found type %T", program.Statements[0])
	}

	testIntegerLiterals(t, stmt.Expression, 5)
}

func testIntegerLiterals(t *testing.T, exp ast.Expression, value int64) bool {
	intLit, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected to find Integer Literal but found %T", exp)
		return false
	}

	if intLit.Value != value {
		t.Fatalf("Expected to find value %q, found %q", value, intLit.Value)
		return false
	}

	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("Expected to find TokenLiteral %q, found %q", fmt.Sprintf("%d", value), intLit.Value)
		return false
	}
	return true
}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		prog := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := prog.Statements[0]

		exp, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Not an expression statement")
		}

		testBooleanLiteral(t, exp.Expression, tt.expected)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, expected bool) bool {
	boolLit, ok := exp.(*ast.BooleanExpression)
	if !ok {
		t.Fatalf("Expected to find Boolean Literal but found %T", exp)
		return false
	}

	if boolLit.Value != expected {
		t.Fatalf("Expected to find value %t, found %t", expected, boolLit.Value)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) == 0 {
		return
	}

	for _, err := range p.errors {
		t.Errorf("parser error: %q", err)
	}

	t.FailNow()
}

func TestErrors(t *testing.T) {
	input := `let x 5`
	l := lexer.New(input)
	p := New(l)

	p.ParseProgram() //need to parse to check error however expecting error so program shoudl not be successfully parsed

	if len(p.errors) == 0 {
		t.Fatalf("No error generated for %q", input)
	}

	expectedError := `expected to find token = but found INT instead`

	if p.errors[0] != expectedError {
		t.Fatalf("Incorrect error generated: expected %s got %s", expectedError, p.errors[0])
	}
}

func TestPrefixOperators(t *testing.T) {
	tests := []struct {
		input         string
		prefixLiteral string
		value         int64
	}{
		{"!5", "!", 5},
		{"-10", "-", 10},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		prog := p.ParseProgram()
		checkParserErrors(t, p)
		if prog == nil {
			t.Fatalf("Program failed to parse")
		}

		if len(prog.Statements) != 1 {
			t.Fatalf("Incorrect Number of statements")
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Not a expression statement, got %T", prog.Statements[0])
		}

		pe, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Not a prefix expression, got %T", stmt.Expression)
		}

		if pe.TokenLiteral() != tt.prefixLiteral {
			t.Fatalf("Operator %q not equal to expected %q", pe.Operator, tt.prefixLiteral)
		}

		testLiteral(t, pe.Expression, tt.value)
	}
}

func TestInfixOperators(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"false != true", false, "!=", true},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infix, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("Expected a infix expression but got %T", exp)
		return false
	}

	testLiteral(t, infix.Left, left)
	testLiteral(t, infix.Right, right)

	if infix.Operator != operator {
		t.Fatalf("Expected %q but found %q", operator, infix.Operator)
		return false
	}

	return true
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3 ) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for i, tt := range tests {
		t.Logf("Running test number %d", i)
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testLiteral(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiterals(t, exp, int64(v))
	case int64:
		return testIntegerLiterals(t, exp, v)
	case string:
		return testIdentifier(t, exp, string(v))
	case bool:
		return testBooleanLiteral(t, exp, bool(v))
	}

	t.Errorf("Type not handled in testLiteral %T", expected)
	return false
}

func TestIfStatements(t *testing.T) {
	input := "if(x < y) {x}"
	l := lexer.New(input)
	p := New(l)

	prog := p.ParseProgram()
	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("Unable to parse program")
		return
	}

	es, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected Expression Statement Got Type %T", prog.Statements[0])
	}

	ie, ok := es.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected If Expression Got Type %T", es.Expression)
	}

	testInfixExpression(t, ie.Condition, "x", "<", "y")
	if len(ie.Consequence.Statements) != 1 {
		t.Fatalf("Expected 1 expression but got %d", len(ie.Consequence.Statements))
	}

	consequence, ok := ie.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expeted expression as consequcne but got %T", ie.Consequence.Statements[0])
	}

	testIdentifier(t, consequence.Expression, "x")
}

func TestIfElseStatements(t *testing.T) {
	input := "if(x < y) {x} else {y}"
	l := lexer.New(input)
	p := New(l)

	prog := p.ParseProgram()
	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("Unable to parse program")
		return
	}

	es, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected Expression Statement Got Type %T", prog.Statements[0])
	}

	ie, ok := es.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected If Expression Got Type %T", es.Expression)
	}

	testInfixExpression(t, ie.Condition, "x", "<", "y")
	if len(ie.Consequence.Statements) != 1 {
		t.Fatalf("Expected 1 expression but got %d", len(ie.Consequence.Statements))
	}

	consequence, ok := ie.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expeted expression as consequcne but got %T", ie.Consequence.Statements[0])
	}

	testLiteral(t, consequence.Expression, "x")

	alternative, ok := ie.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expeted expression as alternative but got %T", ie.Alternative.Statements[0])
	}

	testLiteral(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteral(t, function.Parameters[0], "x")
	testLiteral(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteral(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteral(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foobar"`, "foobar"},
		{`"Hello world"`, "Hello world"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("Unable to parse program of just strings")
		}

		expr, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Expected ast.ExpressionStatement, got %T", program.Statements[0])
		}

		str, ok := expr.Expression.(*ast.StringLiteral)

		if !ok {
			t.Fatalf("Expected ast.StringLiteral, got %T", expr.Expression)
		}

		if str.Value != tt.expected {
			t.Errorf("Expected string to be %q but found %q", tt.expected, str.Value)
		}
	}
}

func TestArrayliterals(t *testing.T) {
	tests := []struct {
		input         string
		numItems      int
		expectedItems []interface{}
	}{
		{"[]", 0, nil},
		{"[1]", 1, []interface{}{1}},
		{`[true, 42]`, 2, []interface{}{true, 42}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("Unable to parse program of just arrays")
		}

		if len(program.Statements) != 1 {
			t.Fatal("too many arguments in program")
		}

		es := program.Statements[0].(*ast.ExpressionStatement)
		arr, ok := es.Expression.(*ast.Array)
		if !ok {
			t.Fatal("expression not an array")
		}

		for i, item := range arr.Items {
			testLiteral(t, item, tt.expectedItems[i])
		}
	}
}

func TestIndexStatements(t *testing.T) {
	input := "myArray[1+1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}

}

func TestHasLiterals(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("expected HasLiteral but got type %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("Unexpected number of pairs expected 3 got %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("expected string literal as key, got %T", k)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiterals(t, v, expectedValue)
	}

}
