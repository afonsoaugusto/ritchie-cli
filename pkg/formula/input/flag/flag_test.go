/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package flag

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
)

func TestInputs(t *testing.T) {
	setup := formula.Setup{
		FormulaPath: os.TempDir(),
	}

	type in struct {
		creResolver      env.Resolvers
		defaultFlagValue string
		valueForList     string
		operator         string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success flags",
			in: in{
				creResolver:      env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				defaultFlagValue: "text",
				valueForList:     "in_list2",
				operator:         "==",
			},
			want: nil,
		},
		{
			name: "success with input omitted",
			in: in{
				creResolver:      env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				defaultFlagValue: "text",
				valueForList:     "in_list2",
				operator:         "!=",
			},
			want: nil,
		},
		{
			name: "error flags empty",
			in: in{
				creResolver:  env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				valueForList: "in_list2",
				operator:     "==",
			},
			want: errors.New("these flags cannot be empty [--sample_text_cache, --sample_text_2, --sample_password]"),
		},
		{
			name: "error env resolver",
			in: in{
				creResolver:      env.Resolvers{"CREDENTIAL": envResolverMock{in: "test", err: errors.New("credential not found")}},
				defaultFlagValue: "text",
				valueForList:     "in_list2",
				operator:         "==",
			},
			want: errors.New("credential not found"),
		},
		{
			name: "invalid value for item",
			in: in{
				creResolver:      env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				defaultFlagValue: "text",
				valueForList:     "invalid",
				operator:         "==",
			},
			want: errors.New("only these input items [in_list1, in_list2, in_list3, in_listN] are accepted in the \"--sample_list\" flag"),
		},
		{
			name: "invalid operator",
			in: in{
				creResolver:      env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				defaultFlagValue: "text",
				valueForList:     "in_list2",
				operator:         "eq",
			},
			want: errors.New("config.json: conditional operator eq not valid. Use any of (==, !=, >, >=, <, <=)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(fmt.Sprintf(inputJson, tt.in.operator)), &inputs)

			inputManager := NewInputManager(tt.in.creResolver)

			cmd := &exec.Cmd{}
			flags := pflag.NewFlagSet("test", 0)

			for _, in := range inputs {
				switch in.Type {
				case input.TextType, input.PassType:
					if len(in.Items) > 0 {
						flags.String(in.Name, tt.in.valueForList, in.Tutorial)
					} else {
						flags.String(in.Name, tt.in.defaultFlagValue, in.Tutorial)
					}
				case input.BoolType:
					flags.Bool(in.Name, false, in.Tutorial)
				}
			}

			setup.Config = formula.Config{
				Inputs: inputs,
			}

			got := inputManager.Inputs(cmd, setup, flags)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Inputs(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

const inputJson = `[
    {
        "name": "sample_text_cache",
        "type": "text",
        "label": "Type : ",
        "cache": {
            "active": true,
            "qty": 6,
            "newLabel": "Type new value. "
        },
		"tutorial": "Add a text for this field."
    },
 	{
        "name": "sample_text",
        "type": "text",
        "label": "Type : ",
		"default": "test"
    },
	{
        "name": "sample_text_2",
        "type": "text",
        "label": "Type : ",
		"required": true
    },
    {
        "name": "sample_list",
        "type": "text",
        "default": "in1",
        "items": [
            "in_list1",
            "in_list2",
            "in_list3",
            "in_listN"
        ],
 		"cache": {
            "active": true,
            "qty": 3,
            "newLabel": "Type new value?"
        },
        "label": "Pick your : ",
		"tutorial": "Select an item for this field."
    },
	{
        "name": "sample_text_condition",
        "type": "text",
        "label": "Type : ",
		"default": "test",
		"condition": {
			"variable": "sample_list",
			"operator": "%s",
			"value":    "in_list2"
		}
    },
    {
        "name": "sample_bool",
        "type": "bool",
        "default": "false",
        "items": [
            "false",
            "true"
        ],
        "label": "Pick: ",
		"tutorial": "Select true or false for this field."
    },
    {
        "name": "sample_password",
        "type": "password",
        "label": "Pick: ",
		"tutorial": "Add a secret password for this field."
    },
    {
        "name": "test_resolver",
        "type": "CREDENTIAL_TEST"
    }
]`
