// value.go
package main

import "fmt"

type ValueType int

const (
	Void ValueType = iota
	type_Int
	type_Float
	type_Bool
	type_String
	type_Vec
	type_Func
)

type Value interface {
	Type() ValueType
	String() string
}

type VoidVal struct{}

func (v VoidVal) Type() ValueType { return Void }
func (v VoidVal) String() string  { return "void" }

type IntVal struct{ Val int64 }

func (v IntVal) Type() ValueType { return type_Int }
func (v IntVal) String() string  { return fmt.Sprintf("%d", v.Val) }

type BoolVal struct{ Val bool }

func (v BoolVal) Type() ValueType { return type_Bool }
func (v BoolVal) String() string  { return fmt.Sprintf("%t", v.Val) }

type StringVal struct{ Val string }

func (v StringVal) Type() ValueType { return type_String }
func (v StringVal) String() string  { return v.Val }
