package parse

import (
	"testing"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestForwardingAddress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    types.ForwardingAddress
		expectError bool
	}{
		{
			name:        "valid https URL",
			input:       "https://example.com",
			expected:    types.ForwardingAddress("https://example.com"),
			expectError: false,
		},
		{
			name:        "valid http URL",
			input:       "http://example.com",
			expected:    types.ForwardingAddress("http://example.com"),
			expectError: false,
		},
		{
			name:        "strip trailing slashes",
			input:       "http://example.com///",
			expected:    types.ForwardingAddress("http://example.com"),
			expectError: false,
		},
		{
			name:        "valid URL with path",
			input:       "https://example.com/path/to/page",
			expected:    types.ForwardingAddress("https://example.com/path/to/page"),
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    types.ForwardingAddress(""),
			expectError: true,
		},
		{
			name:        "invalid scheme",
			input:       "ftp://example.com",
			expected:    types.ForwardingAddress(""),
			expectError: true,
		},
		{
			name:        "no scheme",
			input:       "example.com",
			expected:    types.ForwardingAddress(""),
			expectError: true,
		},
		{
			name:        "invalid URL",
			input:       "not-a-url",
			expected:    types.ForwardingAddress(""),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ForwardingAddress(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %q but got %q", tt.expected, result)
				}
			}
		})
	}
}
