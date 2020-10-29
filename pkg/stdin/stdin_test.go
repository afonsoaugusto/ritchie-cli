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

package stdin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"
)

const msg = "read stdin test"

type TestStdin struct {
	Test string `json:"test"`
}

func TestReadJson(t *testing.T) {

	// Convert interface to Json for test
	i := TestStdin{Test: msg}
	jsonData, _ := json.Marshal(i)

	// Insert Json inside a new Reader (simulating os.Stdin)
	var stdin bytes.Buffer
	stdin.Write(jsonData)
	reader := bufio.NewReader(&stdin)

	tr := TestStdin{}

	// ReadJson through Reader and convert to chosen interface

	if err := ReadJson(reader, &tr); err != nil {
		t.Errorf("Got error %v", err)
	}

	// Assert the decoder result is the initial message
	if msg != tr.Test {
		t.Errorf("Expected : %v but got %v", msg, tr.Test)
	}
}

// func TestExistsEntry(t *testing.T) {
// 	var tests = []struct {
// 		name           string
// 		expectedResult bool
// 		inputMsg       interface{}
// 	}{
// 		{
// 			name:           "return true when json data inputed",
// 			expectedResult: true,
// 			inputMsg:       TestStdin{Test: msg},
// 		},
// 		{
// 			name:           "return false when json data not inputed",
// 			expectedResult: false,
// 			inputMsg:       nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			jsonData, _ := json.Marshal(tt.inputMsg)

// 			var stdin bytes.Buffer
// 			stdin.Write(jsonData)
// 			reader := bufio.NewReader(&stdin)

// 			if result := ExistsEntry(reader); result != tt.expectedResult {
// 				t.Errorf("Got: %v expected: %v", result, tt.expectedResult)
// 			}
// 		})
// 	}
// }
