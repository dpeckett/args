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

type ArgMarshaler interface {
	MarshalArg() string
}

// Marshal marshals a struct into a slice of strings suitable for passing to
// a command line program. It's a very naive implementation that only supports
// a limited set of types.
func Marshal(opts any) []string {
	var args []string
	var posArgs = make(map[int]string)

	s := structs.New(opts)
	for _, field := range s.Fields() {
		if field.IsEmbedded() {
			args = append(args, Marshal(field.Value())...)
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

		var argsToAppend []string
		switch v := field.Value().(type) {
		case bool:
			argsToAppend = marshalBoolFlag(tag, v)
		case *bool:
			argsToAppend = marshalBoolFlag(tag, *v)
		case int:
			argsToAppend = marshalIntFlag(tag, v)
		case *int:
			argsToAppend = marshalIntFlag(tag, *v)
		case string:
			argsToAppend = marshalStringFlag(tag, v)
		case *string:
			argsToAppend = marshalStringFlag(tag, *v)
		case []string:
			for _, s := range v {
				argsToAppend = marshalStringFlag(tag, s)
				if isPosArg {
					for i := 0; i < len(argsToAppend); i++ {
						posArgs[pos] = argsToAppend[i]
						pos++
					}
				} else {
					args = append(args, argsToAppend...)
				}
			}

			continue
		default:
			if m, ok := field.Value().(ArgMarshaler); ok {
				argsToAppend = marshalCustomFlag(tag, m)
			} else {
				panic(fmt.Sprintf("unsupported argument type: %s", field.Kind()))
			}
		}

		if len(argsToAppend) > 0 {
			if isPosArg {
				for i := 0; i < len(argsToAppend); i++ {
					posArgs[pos] = argsToAppend[i]
					pos++
				}
			} else {
				args = append(args, argsToAppend...)
			}
		}
	}

	orderedPosArgs := make([]string, len(posArgs))
	for pos, arg := range posArgs {
		orderedPosArgs[pos] = arg
	}

	return append(args, orderedPosArgs...)
}

func marshalBoolFlag(tag string, v bool) []string {
	var isPos bool
	if _, err := strconv.Atoi(tag); err == nil {
		isPos = true
	}

	if isPos {
		if v {
			return []string{"true"}
		}
		return []string{"false"}
	} else if v {
		if len(tag) == 1 {
			return []string{"-" + tag}
		}
		return []string{"--" + tag}
	}

	return nil
}

func marshalIntFlag(tag string, v int) []string {
	var isPos bool
	if _, err := strconv.Atoi(tag); err == nil {
		isPos = true
	}

	if isPos {
		return []string{strconv.Itoa(v)}
	}

	if len(tag) == 1 {
		return []string{"-" + tag, strconv.Itoa(v)}
	}

	return []string{fmt.Sprintf("--%s=%d", tag, v)}
}

func marshalStringFlag(tag, v string) []string {
	var isPos bool
	if _, err := strconv.Atoi(tag); err == nil {
		isPos = true
	}

	if isPos {
		return []string{v}
	}

	if len(tag) == 1 {
		return []string{"-" + tag, v}
	}

	return []string{fmt.Sprintf("--%s=%s", tag, v)}
}

func marshalCustomFlag(tag string, v ArgMarshaler) []string {
	var isPos bool
	if _, err := strconv.Atoi(tag); err == nil {
		isPos = true
	}

	if isPos {
		return []string{v.MarshalArg()}
	}

	if len(tag) == 1 {
		return []string{"-" + tag, v.MarshalArg()}
	}

	return []string{fmt.Sprintf("--%s=%s", tag, v.MarshalArg())}
}
