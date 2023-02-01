package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type ObjectType string
type BuiltInFunction func(args ...Object) Object

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return fmt.Sprintf("%v", rv.Value.Inspect()) }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJECT }

type ErrorValue struct {
	Message string
}

func (ev *ErrorValue) Inspect() string  { return fmt.Sprintf("Error : %s", ev.Message) }
func (ev *ErrorValue) Type() ObjectType { return ERROR_VALUE_OBJECT }

type FunctionValue struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fv *FunctionValue) Inspect() string {
	var out bytes.Buffer

	params := make([]string, 0)
	for _, p := range fv.Parameters {
		params = append(params, p.Value)
	}

	out.WriteString(fmt.Sprintf("(%s)", strings.Join(params, ",")))
	out.WriteString("{\n")
	out.WriteString(fv.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (fv *FunctionValue) Type() ObjectType { return FUNCTION_VALUE_OBJECT }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return fmt.Sprintf(`"%s"`, s.Value) }
func (s *String) Type() ObjectType { return STRING_OBJECT }

type BuiltIn struct {
	Fn BuiltInFunction
}

func (bi *BuiltIn) Inspect() string  { return "built in function" }
func (bi *BuiltIn) Type() ObjectType { return BUILT_IN_OBJECT }

type Array struct {
	Items []Object
}

func (a *Array) Inspect() string {
	var items []string

	for _, i := range a.Items {
		items = append(items, i.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(items, ","))
}

func (a *Array) Type() ObjectType { return ARRAY_OBJECT }
