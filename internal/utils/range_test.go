package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRangeHeader(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput *RangeHeader
		errorExpected  bool
	}{
		{
			name:           "empty header",
			input:          "",
			expectedOutput: nil,
			errorExpected:  true,
		},
		{
			name:           "without range",
			input:          "bytes=",
			expectedOutput: nil,
			errorExpected:  true,
		},
		{
			name:           "string range",
			input:          "bytes=abc-def",
			expectedOutput: nil,
			errorExpected:  true,
		},
		{
			name:           "float range",
			input:          "bytes=1.10-2.20",
			expectedOutput: nil,
			errorExpected:  true,
		},
		{
			name:           "range with 2 separators",
			input:          "bytes=0-1-2",
			expectedOutput: nil,
			errorExpected:  true,
		},
		{
			name:           "range without start",
			input:          "bytes=-100",
			expectedOutput: nil,
			errorExpected:  true,
		},
		{
			name:  "range without end",
			input: "bytes=0-",
			expectedOutput: &RangeHeader{
				Unit: "bytes",
				Ranges: []Range{
					{
						Start: 0,
						End:   -1,
					},
				},
			},
			errorExpected: false,
		},
		{
			name:  "multiple range",
			input: "bytes=0-10, 10-100",
			expectedOutput: &RangeHeader{
				Unit: "bytes",
				Ranges: []Range{
					{
						Start: 0,
						End:   10,
					},
					{
						Start: 10,
						End:   100,
					},
				},
			},
			errorExpected: false,
		},
	}
	for _, test := range tests {
		output, err := ParseRangeHeader(test.input)
		if test.errorExpected {
			assert.Error(t, err, test.name)
			continue
		}
		assert.NoError(t, err, test.name)
		assert.EqualValues(t, test.expectedOutput, output, test.name)
	}
}
