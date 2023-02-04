package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

type ObjectType string
type BuiltInFunction func(args ...Object) Object

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{
		Type:  b.Type(),
		Value: value,
	}
}

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
func (s *String) HashKey() HashKey {

	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}

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

type HashPair struct {
	Key   Object
	Value Object
}

type HashMap struct {
	Pairs map[HashKey]HashPair
}

func (hm *HashMap) Type() ObjectType { return HASH_OBJECT }
func (hm *HashMap) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, val := range hm.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s : %s", val.Key.Inspect(), val.Value.Inspect()))
	}

	out.WriteRune('{')
	out.WriteString(strings.Join(pairs, ","))
	out.WriteRune('}')

	return out.String()
}
