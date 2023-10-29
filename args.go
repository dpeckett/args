/* SPDX-License-Identifier: Apache-2.0
 *
 * Copyright 2023 Damian Peckett <damian@pecke.tt>.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF AintNY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package args

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"
)

type TextMarshaler interface {
	MarshalText() string
}

// MarshalArgs marshals a struct into a slice of strings suitable for passing to
// a command line program. It's a very naive implementation that only supports
// a limited set of types.
func MarshalArgs(opts any) []string {
	var args []string
	var posArgs = make(map[int]string)

	s := structs.New(opts)
	for _, field := range s.Fields() {
		if field.IsEmbedded() {
			args = append(args, MarshalArgs(field.Value())...)
			continue
		}

		tag := field.Tag("arg")
		if tag == "" {
			continue
		}

		var isPosArg bool
		pos, err := strconv.Atoi(tag)
		if err == nil {
			isPosArg = true
		}

		if !field.IsExported() || field.IsZero() {
			continue
		}

		var arg string
		switch v := field.Value().(type) {
		case bool:
			if isPosArg {
				if v {
					arg = "true"
				} else {
					arg = "false"
				}
			} else {
				if v {
					arg = tag
				}
			}
		case *bool:
			if isPosArg {
				if *v {
					arg = "true"
				} else {
					arg = "false"
				}
			} else {
				if *v {
					arg = tag
				}
			}
		case int:
			if isPosArg {
				arg = strconv.Itoa(v)
			} else {
				arg = fmt.Sprintf("%s=%d", tag, v)
			}
		case *int:
			if isPosArg {
				arg = strconv.Itoa(*v)
			} else {
				arg = fmt.Sprintf("%s=%d", tag, *v)
			}
		case string:
			if isPosArg {
				arg = v
			} else {
				arg = fmt.Sprintf("%s=%s", tag, v)
			}
		case *string:
			if isPosArg {
				arg = *v
			} else {
				arg = fmt.Sprintf("%s=%s", tag, *v)
			}
		case []string:
			if isPosArg {
				for _, s := range v {
					posArgs[pos] = s
					pos++
				}
			} else {
				for _, s := range v {
					args = append(args, fmt.Sprintf("%s=%s", tag, s))
				}
			}

			continue
		default:
			if m, ok := field.Value().(TextMarshaler); ok {
				if isPosArg {
					arg = m.MarshalText()
				} else {
					arg = fmt.Sprintf("%s=%s", tag, m.MarshalText())
				}
			} else {
				panic(fmt.Sprintf("unsupported argument type: %s", field.Kind()))
			}
		}

		if arg != "" {
			if isPosArg {
				posArgs[pos] = arg
			} else {
				args = append(args, arg)
			}
		}
	}

	orderedPosArgs := make([]string, len(posArgs))
	for pos, arg := range posArgs {
		orderedPosArgs[pos] = arg
	}

	return append(args, orderedPosArgs...)
}

// YesNo is a boolean type that marshals to "y" or "n".
type YesNo bool

var Yes = PtrTo(YesNo(true))
var No = PtrTo(YesNo(false))

func (yn *YesNo) MarshalText() string {
	if *yn {
		return "y"
	}

	return "n"
}

func PtrTo[T any](v T) *T {
	return &v
}
