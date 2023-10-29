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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package args_test

import (
	"testing"

	"github.com/dpeckett/args"
	"github.com/stretchr/testify/require"
)

type Options struct {
	CredentialOptions
	Verbose     *YesNo   `arg:"verbose"`
	File        string   `arg:"1"`
	Destination string   `arg:"0"`
	Repeated    []string `arg:"repeated"`
}

type CredentialOptions struct {
	User     string `arg:"user"`
	Password string `arg:"password"`
}

type YesNo bool

func (yn *YesNo) MarshalArg() string {
	if *yn {
		return "y"
	}

	return "n"
}

func TestMarshal(t *testing.T) {
	yes := YesNo(true)

	args := args.Marshal(Options{
		CredentialOptions: CredentialOptions{
			User: "root",
		},
		Verbose:     &yes,
		File:        "/tmp/foo",
		Destination: "/tmp/bar",
		Repeated:    []string{"foo", "bar"},
	})

	expected := []string{
		"--user=root",
		"--verbose=y",
		"--repeated=foo",
		"--repeated=bar",
		"/tmp/bar",
		"/tmp/foo",
	}

	require.Len(t, args, len(expected))
	require.Equal(t, expected, args)
}
