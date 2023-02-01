package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanExpression:
		return nativeToBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Expression, env)
		if isError(right) {
			return right
		}
		return processPrefixOperator(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(left, right, node.Operator)
	case *ast.IfExpression:
		return processIfExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node.Value, env)
	case *ast.FunctionLiteral:
		return &object.FunctionValue{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return evalFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Array:
		var objects []object.Object
		for _, item := range node.Items {
			obj := Eval(item, env)
			objects = append(objects, obj)
		}
		return &object.Array{Items: objects}
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	default:
		return newError("unknown type: %T", node)
	}

	return nil
}

func evalIndexExpression(exp *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(exp.Left, env)
	if left == nil {
		return object.NULL
	}
	idx := Eval(exp.Index, env)
	if idx == nil {
		return object.NULL
	}

	if left.Type() != object.ARRAY_OBJECT || idx.Type() != object.INTEGER_OBJ {
		return object.NULL
	}

	arr := left.(*object.Array)
	index := idx.(*object.Integer)

	if index.Value < 0 || index.Value > int64(len(arr.Items))-1 {
		return object.NULL
	}

	return arr.Items[index.Value]
}

func evalFunction(function object.Object, params []object.Object) object.Object {
	switch fn := function.(type) {
	case *object.FunctionValue:
		extendedEnv := createExtendedEnvironment(fn, params)
		evaluated := Eval(fn.Body, extendedEnv)
		if rv, ok := evaluated.(*object.ReturnValue); ok {
			return rv.Value
		}
		return evaluated

	case *object.BuiltIn:
		return fn.Fn(params...)
	}

	f, ok := function.(*object.FunctionValue)
	if !ok {
		return newError("Object is not a function!!!!")
	}

	extendedEnv := createExtendedEnvironment(f, params)
	evaluated := Eval(f.Body, extendedEnv)
	if rv, ok := evaluated.(*object.ReturnValue); ok {
		return rv.Value
	}
	return evaluated
}

func createExtendedEnvironment(function *object.FunctionValue, params []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Env)
	for idx, p := range function.Parameters {
		env.Set(p.Value, params[idx])
	}
	return env
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range expressions {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIdentifier(name string, env *object.Environment) object.Object {
	val, ok := env.Get(name)
	if ok {
		return val
	}

	fn, ok := builtIns[name]
	if ok {
		return fn
	}

	return newError("identifier not found: '%s'", name)
}

func evalInfixExpression(left, right object.Object, operator string) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return processIntegerInfixExpresion(left.(*object.Integer), operator, right.(*object.Integer))
	} else if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return processBooleanInfixOperator(left.(*object.Boolean), operator, right.(*object.Boolean))
	} else if left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT {
		return processStringInfixOperator(left.(*object.String), operator, right.(*object.String))
	}

	return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
}

func processStringInfixOperator(left *object.String, operator string, right *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	}
	return newError("Unknown string operand %q", operator)
}

func processBooleanInfixOperator(left *object.Boolean, operand string, right *object.Boolean) object.Object {
	switch operand {
	case "==":
		return &object.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &object.Boolean{Value: left.Value != right.Value}

	}
	return newError("unknown operator: %s %s %s", left.Type(), operand, right.Type())
}

func processIfExpression(expression *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(expression.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(expression.Consequence, env)
	} else if expression.Alternative != nil {
		return Eval(expression.Alternative, env)
	} else {
		return object.NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch node := obj.(type) {
	case *object.Boolean:
		return node.Value
	default:
		return false
	}
}

func processIntegerInfixExpresion(left *object.Integer, operand string, right *object.Integer) object.Object {
	rightValue := right.Value
	leftValue := left.Value

	switch operand {
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "<":
		return nativeToBoolean(leftValue < rightValue)
	case ">":
		return nativeToBoolean(leftValue > rightValue)
	}

	return newError("unknown operator: %s %s %s", left.Type(), operand, right.Type())
}

func processPrefixOperator(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return processBangOperator(right)
	case "-":
		return processMinusOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func processMinusOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: %s%s", "-", right.Type())
	}
	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func processBangOperator(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.FALSE
	default:
		return object.FALSE
	}
}

func nativeToBoolean(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range statements {
		result = Eval(s, env)
		switch node := result.(type) {
		case *object.ReturnValue:
			return node.Value
		case *object.ErrorValue:
			return node
		}
	}

	return result
}

func evalBlockExpression(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range block.Statements {
		result = Eval(s, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJECT || rt == object.ERROR_VALUE_OBJECT {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) object.Object {
	return &object.ErrorValue{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_VALUE_OBJECT
	}
	return false
}
