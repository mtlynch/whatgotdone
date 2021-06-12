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
			InvalidTwitterHandleError,
		},
		{
			"empty string is invalid",
			"",
			InvalidTwitterHandleError,
		},
		{
			"undefined value is invalid",
			"undefined",
			InvalidTwitterHandleError,
		},
		{
			"handle with leading @ is invalid",
			"@jack",
			InvalidTwitterHandleError,
		},
		{
			"handle with more than 15 characters is invalid",
			"jackandjillwentup",
			InvalidTwitterHandleError,
		},
		{
			"handle with illegal characters is invalid",
			"jack.and.jill",
			InvalidTwitterHandleError,
		},
	}

	for _, tt := range tests {
		_, err := TwitterHandle(tt.handle)
		if err != tt.errExpected {
			t.Errorf("%s: input [%s], got %v, want %v", tt.explanation, tt.handle, err, tt.errExpected)
		}
	}
}
