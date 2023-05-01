//
// Copyright (c) 2023 Markku Rossi
//
// All rights reserved.
//

package scheme

import (
	"bytes"
	"fmt"
	"strings"
)

// ByteVector implements bytevector values.
type ByteVector []byte

// Scheme returns the value as a Scheme string.
func (v ByteVector) Scheme() string {
	return v.String()
}

// Eq tests if the argument value is eq? to this value.
func (v ByteVector) Eq(o Value) bool {
	ov, ok := o.(ByteVector)
	if !ok {
		return false
	}
	if len(v) == 0 && len(ov) == 0 {
		return true
	}
	return false
}

// Equal tests if the argument value is equal to this value.
func (v ByteVector) Equal(o Value) bool {
	ov, ok := o.(ByteVector)
	if !ok || len(v) != len(ov) {
		return false
	}
	for idx, vv := range v {
		if vv != ov[idx] {
			return false
		}
	}
	return true
}

func (v ByteVector) String() string {
	var str strings.Builder
	str.WriteString("#vu8(")

	for idx, el := range v {
		if idx > 0 {
			str.WriteRune(' ')
		}
		str.WriteString(fmt.Sprintf("%d", el))
	}
	str.WriteRune(')')
	return str.String()
}

var rnrsBytevectorBuiltins = []Builtin{
	{
		Name: "bytevector?",
		Args: []string{"obj"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			_, ok := args[0].(ByteVector)
			return Boolean(ok), nil
		},
	},
	{
		Name: "make-bytevector",
		Args: []string{"k", "[k:fill]"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			length, err := Int64(args[0])
			if err != nil || length < 0 {
				return nil, l.Errorf("negative length: %v", args[0])
			}

			var fill byte
			if len(args) == 2 {
				f, err := Int64(args[1])
				if err != nil || f < -128 || f > 255 {
					return nil, l.Errorf("invalid fill: %v", args[1])
				}
				fill = byte(f)
			}
			elements := make([]byte, length, length)
			for i := 0; i < int(length); i++ {
				elements[i] = fill
			}

			return ByteVector(elements), nil
		},
	},
	{
		Name: "bytevector-length",
		Args: []string{"bytevector"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			v, ok := args[0].(ByteVector)
			if !ok {
				return nil, l.Errorf("not a bytevector: %v", args[0])
			}
			return NewNumber(0, len(v)), nil
		},
	},
	{
		Name: "bytevector=?",
		Args: []string{"bytevector1", "bytevector2"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			v1, ok := args[0].(ByteVector)
			if !ok {
				return nil, l.Errorf("not a bytevector: %v", args[0])
			}
			v2, ok := args[1].(ByteVector)
			if !ok {
				return nil, l.Errorf("not a bytevector: %v", args[1])
			}
			return Boolean(bytes.Equal(v1, v2)), nil
		},
	},
	{
		Name: "bytevector-fill",
		Args: []string{"bytevector", "k:fill"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			v, ok := args[0].(ByteVector)
			if !ok {
				return nil, l.Errorf("invalid bytevector: %v", args[0])
			}
			f, err := Int64(args[1])
			if err != nil || f < -128 || f > 255 {
				return nil, l.Errorf("invalid fill: %v", args[1])
			}
			fill := byte(f)

			for i := 0; i < len(v); i++ {
				v[i] = fill
			}

			return v, nil
		},
	},
	{
		Name: "bytevector-copy!",
		Args: []string{
			"bytevector:source", "k:source-start",
			"bytevector:target", "k:target-start", "k",
		},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			source, ok := args[0].(ByteVector)
			if !ok {
				return nil, l.Errorf("invalid source: %v", args[0])
			}
			target, ok := args[2].(ByteVector)
			if !ok {
				return nil, l.Errorf("invalid target: %v", args[2])
			}
			sourceStart, err := Int64(args[1])
			if err != nil || sourceStart < 0 {
				return nil, l.Errorf("invalid source-start: %v", args[1])
			}

			targetStart, err := Int64(args[3])
			if err != nil || targetStart < 0 {
				return nil, l.Errorf("invalid target-start: %v", args[3])
			}

			k, err := Int64(args[4])
			if err != nil || k < 0 {
				return nil, l.Errorf("invalid k: %v", args[4])
			}

			if sourceStart+k > int64(len(source)) {
				return nil, l.Errorf("invalid source range: %v+%v>%v",
					sourceStart, k, len(source))
			}
			if targetStart+k > int64(len(target)) {
				return nil, l.Errorf("invalid target range: %v+%v>%v",
					targetStart, k, len(target))
			}

			copy(target[targetStart:targetStart+k],
				source[sourceStart:sourceStart+k])

			return nil, nil
		},
	},
	{
		Name: "bytevector-copy",
		Args: []string{"bytevector"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			v, ok := args[0].(ByteVector)
			if !ok {
				return nil, l.Errorf("not a bytevector: %v", args[0])
			}
			arr := make([]byte, len(v))
			copy(arr, v)

			return ByteVector(arr), nil
		},
	},
	{
		Name: "bytevector-u8-ref",
		Args: []string{"bytevector", "k"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			v, ok := args[0].(ByteVector)
			if !ok {
				return nil, l.Errorf("not a bytevector: %v", args[0])
			}
			k, err := Int64(args[1])
			if err != nil {
				return nil, l.Errorf("invalid index: %v", args[1])
			}
			if k < 0 || k >= int64(len(v)) {
				return nil, l.Errorf("invalid index: 0 <= %v < %v", k, len(v))
			}

			return NewNumber(0, int(uint8(v[k]))), nil
		},
	},
	{
		Name: "bytevector-s8-ref",
		Args: []string{"bytevector", "k"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			v, ok := args[0].(ByteVector)
			if !ok {
				return nil, l.Errorf("not a bytevector: %v", args[0])
			}
			k, err := Int64(args[1])
			if err != nil {
				return nil, l.Errorf("invalid index: %v", args[1])
			}
			if k < 0 || k >= int64(len(v)) {
				return nil, l.Errorf("invalid index: 0 <= %v < %v", k, len(v))
			}

			return NewNumber(0, int(int8(v[k]))), nil
		},
	},
}
