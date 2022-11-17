//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package scm

import (
	"strings"
)

// String implements string values.
type String struct {
	Data []byte
}

// Type returns the string value type.
func (v *String) Type() ValueType {
	return VString
}

// Scheme returns the value as a Scheme string.
func (v *String) Scheme() string {
	return StringToScheme(string(v.Data))
}

func (v *String) String() string {
	return string(v.Data)
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
