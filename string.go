//
// Copyright (c) 2022-2023 Markku Rossi
//
// All rights reserved.
//

package scheme

import (
	"strings"
)

// String implements string values.
type String string

// IsString tests if the value is string.
func IsString(value Value) (v string, ok bool) {
	var str String
	str, ok = value.(String)
	if !ok {
		return
	}
	return string(str), true
}

// Scheme returns the value as a Scheme string.
func (v String) Scheme() string {
	return StringToScheme(string(v))
}

// Eq tests if the argument value is eq? to this value.
func (v String) Eq(o Value) bool {
	return v.Equal(o)
}

// Equal tests if the argument value is equal to this value.
func (v String) Equal(o Value) bool {
	ov, ok := o.(String)
	return ok && v == ov
}

func (v String) String() string {
	return string(v)
}

// StringToScheme returns the string as Scheme string literal.
func StringToScheme(s string) string {
	var str strings.Builder
	str.WriteRune('"')
	for _, r := range s {
		switch r {
		case '\\', '"', '|', '(':
			str.WriteRune('\\')
			str.WriteRune(r)
		case '\a':
			str.WriteRune('\\')
			str.WriteRune('a')
		case '\f':
			str.WriteRune('\\')
			str.WriteRune('f')
		case '\n':
			str.WriteRune('\\')
			str.WriteRune('n')
		case '\r':
			str.WriteRune('\\')
			str.WriteRune('r')
		case '\t':
			str.WriteRune('\\')
			str.WriteRune('t')
		case '\v':
			str.WriteRune('\\')
			str.WriteRune('v')
		case '\b':
			str.WriteRune('\\')
			str.WriteRune('b')
		case 0:
			str.WriteRune('\\')
			str.WriteRune('0')
		default:
			str.WriteRune(r)
		}
	}
	str.WriteRune('"')
	return str.String()
}

var stringBuiltins = []Builtin{
	{
		Name: "string?",
		Args: []string{"obj"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			_, ok := args[0].(String)
			return Boolean(ok), nil
		},
	},
	{
		Name: "make-string",
		Args: []string{"k", "[char]"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			length, err := Int64(args[0])
			if err != nil || length < 0 {
				return nil, l.Errorf("invalid length: %v", args[0])
			}

			fill := ' '
			if len(args) == 2 {
				ch, ok := args[1].(Character)
				if !ok {
					return nil, l.Errorf("invalid char: %v", args[1])
				}
				fill = rune(ch)
			}
			str := make([]rune, length, length)
			for i := 0; i < int(length); i++ {
				str[i] = fill
			}

			return String(string(str)), nil
		},
	},
	{
		Name: "string",
		Args: []string{"char..."},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			length := len(args)
			str := make([]byte, length, length)

			for i := 0; i < length; i++ {
				ch, ok := args[i].(Character)
				if !ok {
					return nil, l.Errorf("invalid character: %v", args[i])
				}
				str[i] = byte(ch)
			}
			return String(str), nil
		},
	},
	{
		Name: "string-length",
		Args: []string{"string"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			switch v := args[0].(type) {
			case String:
				return NewNumber(0, len([]rune(string(v)))), nil

			default:
				return nil, l.Errorf("invalid argument: %v", args[0])
			}
		},
	},
	{
		Name: "string-ref",
		Args: []string{"string", "k"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			str, ok := args[0].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[0])
			}
			chars := []rune(string(str))

			k, err := Int64(args[1])
			if err != nil || k < 0 || k >= int64(len(chars)) {
				return nil, l.Errorf("invalid index: %v", args[1])
			}
			return Character(chars[k]), nil
		},
	},
	{
		Name: "scheme::string=?",
		Args: []string{"string1", "string2"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			str1, ok := args[0].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[0])
			}
			str2, ok := args[1].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[1])
			}
			return Boolean(string(str1) == string(str2)), nil
		},
	},
	{
		Name: "scheme::string<?",
		Args: []string{"string1", "string2"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			str1, ok := args[0].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[0])
			}
			str2, ok := args[1].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[1])
			}
			return Boolean(string(str1) < string(str2)), nil
		},
	},
	{
		Name: "scheme::string>?",
		Args: []string{"string1", "string2"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			str1, ok := args[0].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[0])
			}
			str2, ok := args[1].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[1])
			}
			return Boolean(string(str1) > string(str2)), nil
		},
	},
	{
		Name: "substring",
		Args: []string{"string", "start", "end"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			strv, ok := args[0].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[0])
			}
			str := []rune(string(strv))

			start, err := Int64(args[1])
			if err != nil {
				return nil, l.Errorf("invalid start index: %v", args[1])
			}

			end, err := Int64(args[2])
			if err != nil {
				return nil, l.Errorf("invalid end index: %v", args[2])
			}

			if start < 0 || end < start || end > int64(len(strv)) {
				return nil, l.Errorf("invalid indices: 0 <= %v <= %v <= %v",
					start, end, len(strv))
			}

			return String(string(str[start:end])), nil
		},
	},
	{
		Name: "string-append",
		Args: []string{"string..."},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			var result string

			for _, arg := range args {
				str, ok := arg.(String)
				if !ok {
					return nil, l.Errorf("invalid string: %v", arg)
				}
				result += string(str)
			}
			return String(result), nil
		},
	},
	{
		Name: "string->list",
		Args: []string{"string"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			str, ok := args[0].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[0])
			}
			runes := []rune(string(str))

			var head, tail Pair
			for i := 0; i < len(runes); i++ {
				item := NewPair(Character(runes[i]), nil)
				if head == nil {
					head = item
				} else {
					tail.SetCdr(item)
				}
				tail = item
			}
			return head, nil
		},
	},
	{
		Name: "list->string",
		Args: []string{"chars"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			var str []rune
			err := Map(func(idx int, v Value) error {
				ch, ok := v.(Character)
				if !ok {
					return l.Errorf("invalid character: %v", v)
				}
				str = append(str, rune(ch))
				return nil
			}, args[0])
			if err != nil {
				return nil, err
			}
			return String(string(str)), nil
		},
	},
	{
		Name: "string-copy",
		Args: []string{"string"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			str, ok := args[0].(String)
			if !ok {
				return nil, l.Errorf("invalid string: %v", args[0])
			}
			runes := []rune(string(str))
			new := make([]rune, len(runes), len(runes))
			copy(new, runes)
			return String(string(new)), nil
		},
	},
}
