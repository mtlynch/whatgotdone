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
			"handle with numbers is valid",
			"jerry123",
			nil,
		},
		{
			"underscore characters are allowed",
			"jack_and_jill",
			nil,
		},
		{
			"single-character handle is invalid",
			"j",
			ErrInvalidTwitterHandle,
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
		{
			"handle with leading @ is invalid",
			"@jack",
			ErrInvalidTwitterHandle,
		},
		{
			"handle with more than 15 characters is invalid",
			"jackandjillwentup",
			ErrInvalidTwitterHandle,
		},
		{
			"handle with illegal characters is invalid",
			"jack.and.jill",
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
