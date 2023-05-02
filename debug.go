//
// Copyright (c) 2022-2023 Markku Rossi
//
// All rights reserved.
//

package scheme

import (
	"fmt"
	"sort"
)

const (
	fLambda int = 1 << iota
	fNative
	fScheme
)

var debugBuiltins = []Builtin{
	{
		Name: "print-env",
		Args: []string{"sym..."},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {

			var flags int

			if len(args) == 0 {
				flags = 0xffff
			} else {
				for _, arg := range args {
					switch flag := arg.(type) {
					case *Identifier:
						switch flag.Name {
						case "native":
							flags |= fNative
						case "scheme":
							flags |= fScheme

						default:
							return nil, l.Errorf("unknown flag: %v", arg)
						}

					case Keyword:
						switch flag {
						case KwLambda:
							flags |= fLambda

						default:
							return nil, l.Errorf("unknown flag: %v", arg)
						}

					default:
						return nil, l.Errorf("invalid flag: %v(%T)", arg, arg)
					}
				}
			}
			var names []string
			for k, v := range scm.symbols {
				if v.Global != nil {
					switch val := v.Global.(type) {
					case *Lambda:
						if flags&fLambda != 0 {
							if (flags&(fNative|fScheme)) == 0 ||
								((val.Impl.Native != nil) &&
									(flags&fNative != 0)) ||
								((val.Impl.Native == nil) &&
									(flags&fScheme != 0)) {
								names = append(names, k)
							}
						}

					default:
						if flags&fLambda == 0 {
							names = append(names, k)
						}
					}
				}
			}
			sort.Strings(names)

			var max int
			for _, name := range names {
				if len(name) > max {
					max = len(name)
				}
			}

			fmt.Printf("Global symbols:\n")

			for _, name := range names {
				for i := 0; i+len(name) < max; i++ {
					fmt.Print(" ")
				}
				fmt.Printf("%s : %s\n", name, scm.symbols[name].Global)
			}
			fmt.Printf("%d symbols matched\n", len(names))
			return nil, nil
		},
	},
	{
		Name: "disassemble",
		Args: []string{"obj"},
		Native: func(scm *Scheme, l *Lambda, args []Value) (Value, error) {
			switch arg := args[0].(type) {
			case *Lambda:
				scm.Stdout.Printf("lambda: %v\n", arg)
				if arg.Impl.Native == nil {
					for _, c := range arg.Impl.Code {
						scm.Stdout.Printf("%s\n", c)
					}
				}

			default:
				scm.Stdout.Printf("value: %v\n", arg)
			}
			return nil, nil
		},
	},
}
