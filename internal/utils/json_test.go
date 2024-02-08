package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestStringOrStringSlice(t *testing.T) {
	testCases := []struct {
		name           string
		input          json.RawMessage
		expectedOutput []string
		expectedError  bool
	}{
		{
			name:           "string test 1",
			input:          json.RawMessage(`"abc"`),
			expectedOutput: []string{"abc"},
			expectedError:  false,
		},
		{
			name:           "list test 1",
			input:          json.RawMessage(`["abc"]`),
			expectedOutput: []string{"abc"},
			expectedError:  false,
		},
		{
			name:           "list test 2",
			input:          json.RawMessage(`["abc", "def", "ghi"]`),
			expectedOutput: []string{"abc", "def", "ghi"},
			expectedError:  false,
		},
		{
			name:           "empty input test",
			input:          json.RawMessage(``),
			expectedOutput: []string{},
			expectedError:  true,
		},
		{
			name:           "int input test",
			input:          json.RawMessage(`1`),
			expectedOutput: []string{},
			expectedError:  true,
		},
	}
	for _, testCase := range testCases {
		var v utils.StringOrStringSlice
		err := json.Unmarshal(testCase.input, &v)
		if testCase.expectedError {
			if !assert.NotNil(t, err, testCase.name) {
				t.Fail()
				return
			}
			continue
		}
		if !assert.Equal(t, testCase.expectedOutput, v.List(), testCase.name) {
			t.Fail()
			return
		}
	}
}
