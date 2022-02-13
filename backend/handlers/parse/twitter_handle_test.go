package parse

import (
	"testing"
)

func TestTwitterHandle(t *testing.T) {
	var tests = []struct {
		explanation string
		handle      string
		errExpected error
	}{
		{
			"well-formed handle address is valid",
			"jack",
			nil,
		},
		{
			"empty string is invalid",
			"",
			ErrInvalidTwitterHandle,
		},
		{
			"undefined value is invalid",
			"undefined",
			ErrInvalidTwitterHandle,
		},
	}

	for _, tt := range tests {
		_, err := TwitterHandle(tt.handle)
		if err != tt.errExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.handle, err, tt.errExpected)
		}
	}
}
