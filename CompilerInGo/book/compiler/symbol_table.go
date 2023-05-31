package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltInScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: st.numDefinitions}
	if st.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	st.store[name] = symbol
	st.numDefinitions++
	return symbol
}

func (st *SymbolTable) DefineBuiltIn(index int, name string) Symbol {
	symbol := Symbol{Name: name, Scope: BuiltInScope, Index: index}
	st.store[name] = symbol
	return symbol
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := st.store[name]

	if !ok && st.Outer != nil {
		obj, ok := st.Outer.Resolve(name)
		return obj, ok
	}

	return obj, ok
}

func NewEnclosedSymbolTable(symTable *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = symTable
	return s
}
