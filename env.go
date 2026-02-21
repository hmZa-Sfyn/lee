// env.go
package main

type Environment struct {
	values    map[string]Value
	functions map[string]*FunctionDecl
	outer     *Environment
	mutables  map[string]bool // track mutables
}

func NewEnvironment() *Environment {
	return &Environment{
		values:    make(map[string]Value),
		functions: make(map[string]*FunctionDecl),
		mutables:  make(map[string]bool),
	}
}

func (e *Environment) Define(name string, val Value, mutable bool) *EsoError {
	if _, ok := e.values[name]; ok {
		return NewEsoErrorf(0, 0, "redefinition of %q", name)
	}
	e.values[name] = val
	e.mutables[name] = mutable
	return nil
}

func (e *Environment) Assign(name string, val Value) *EsoError {
	if _, ok := e.values[name]; !ok {
		if e.outer != nil {
			return e.outer.Assign(name, val)
		}
		return NewEsoErrorf(0, 0, "undefined variable %q", name)
	}
	if !e.mutables[name] {
		return NewEsoErrorf(0, 0, "cannot assign to immutable %q", name)
	}
	e.values[name] = val
	return nil
}

func (e *Environment) DefineFn(name string, fn *FunctionDecl) {
	e.functions[name] = fn
}

func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.values[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) GetFn(name string) (*FunctionDecl, bool) {
	fn, ok := e.functions[name]
	if !ok && e.outer != nil {
		return e.outer.GetFn(name)
	}
	return fn, ok
}