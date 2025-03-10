#!/bin/bash
# Script to generate Go test files with unit test and benchmark templates
# Usage: ./goaddtest.sh CreateFile create

# Check if required parameters are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <testName> <fileName> e.g.: CreateFile create"
    exit 1
fi

# Test function name
testName=$1

# File name without extension
file=$2"_test.go"

# Current folder name
packageName=$(basename "$(pwd)")

# Go unit test file template
template=$(cat <<EOF
package ${packageName}_test

import (
	"bytes"
	"reflect"
	"testing"
)

var testData${testName} = []struct {
	Comment       string
	DataIN        any
	DataExpected  any
	ErrorExpected error
}{
	{
		Comment:       "test",
		DataIN:        []any{},
		DataExpected:  map[string]string{"price": "2000"},
		ErrorExpected: nil,
	},
}

func Test${testName}(t *testing.T) {
	message := func(comment string, expected, response any) {
		t.Fatalf("\n❌=> in %v the expectation is:\n[%v]\n=> but got:\n[%v]\n", comment, expected, response)
	}
	
	compare := func(comment string, expected, response any) {
		if !reflect.DeepEqual(expected, response) {
			message(comment, expected, response)
		}
	}

	for _, test := range testData${testName} {
		t.Run(("\n" + test.Comment), func(t *testing.T) {
			var buf = &bytes.Buffer{}
			response, err := ${packageName}.${testName}(buf, test.DataIN...)
			compare("ErrorExpected", test.ErrorExpected, err)
			if err == "" {
				compare(test.Comment, test.DataExpected, response)
			}
		})
	}
}

func Benchmark${testName}(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++) {
		for _, test := testData${testName} {
			var buf = &bytes.Buffer{}
			${packageName}.${testName}(buf, test.DataIN...)
		}
	}
}
EOF
)

# Create file in current directory
echo "$template" > "$file"
echo "File $file created."
