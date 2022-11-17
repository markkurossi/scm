//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package scm

import (
	"fmt"
	"strings"
)

var (
	_ Value = &Cons{}
	_ Value = &Vector{}
	_ Value = &Identifier{}
	_ Value = &Number{}
	_ Value = &Boolean{}
	_ Value = &String{}
	_ Value = &Character{}
	_ Value = &Lambda{}
)

// ValueType specifies value types.
type ValueType int

// Value types.
const (
	VCons ValueType = iota
	VVector
	VIdentifier
	VNumber
	VBoolean
	VString
	VCharacter
	VLambda
)

// Value implements a Scheme value.
type Value interface {
	Type() ValueType
	Scheme() string
}

// Cons implements cons cell values.
type Cons struct {
	Car Value
	Cdr Value
}

// Type returns the cons cell value type.
func (v *Cons) Type() ValueType {
	return VCons
}

// Scheme returns the value as a Scheme string.
func (v *Cons) Scheme() string {
	return v.String()
}

func (v *Cons) String() string {
	var str strings.Builder
	str.WriteRune('(')

	i := v
	first := true
loop:
	for {
		if first {
			first = false
		} else {
			str.WriteRune(' ')
		}
		if i.Car == nil {
			str.WriteString("nil")
		} else {
			str.WriteString(i.Car.Scheme())
		}
		switch cdr := i.Cdr.(type) {
		case *Cons:
			i = cdr

		case nil:
			break loop

		default:
			str.WriteString(" . ")
			str.WriteString(fmt.Sprintf("%v", i.Cdr))
			break loop
		}
	}
	str.WriteRune(')')

	return str.String()
}

// Vector implements vector values.
type Vector struct {
	Elements []Value
}

// Type returns the vector value type.
func (v *Vector) Type() ValueType {
	return VVector
}

// Scheme returns the value as a Scheme string.
func (v *Vector) Scheme() string {
	return v.String()
}

func (v *Vector) String() string {
	var str strings.Builder
	str.WriteString("#(")

	for idx, el := range v.Elements {
		if idx > 0 {
			str.WriteRune(' ')
		}
		str.WriteString(fmt.Sprintf("%v", el))
	}
	str.WriteRune(')')
	return str.String()
}

// Identifier implements identifier values.
type Identifier struct {
	Name   string
	Global Value
}

// Type returns the identifier value type.
func (v *Identifier) Type() ValueType {
	return VIdentifier
}

// Scheme returns the value as a Scheme string.
func (v *Identifier) Scheme() string {
	return v.String()
}

func (v *Identifier) String() string {
	return v.Name
}

// Boolean implements boolean values.
type Boolean struct {
	Bool bool
}

// Type returns the boolean value type.
func (v *Boolean) Type() ValueType {
	return VBoolean
}

// Scheme returns the value as a Scheme string.
func (v *Boolean) Scheme() string {
	return v.String()
}

func (v *Boolean) String() string {
	return BooleanToScheme(v.Bool)
}

// BooleanToScheme returns the bool as Scheme boolean literal.
func BooleanToScheme(v bool) string {
	var ch rune
	if v {
		ch = 't'
	} else {
		ch = 'f'
	}
	return fmt.Sprintf("#%c", ch)
}

// Lambda implements lambda values.
type Lambda struct {
	MinArgs int
	MaxArgs int
	Native  Native
	Locals  []Value
}

// Type returns the lambda value type.
func (v *Lambda) Type() ValueType {
	return VLambda
}

// Scheme returns the value as a Scheme string.
func (v *Lambda) Scheme() string {
	return v.String()
}

func (v *Lambda) String() string {
	return fmt.Sprintf("(lambda () {native})")
}

// Native implements native functions.
type Native func(vm *VM, args []Value) (Value, error)

// Builtin defines a built-in native function.
type Builtin struct {
	Name    string
	MinArgs int
	MaxArgs int
	Native  Native
}
