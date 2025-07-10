package parse

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/types"
)

var (
	ErrForwardingAddressEmpty         = errors.New("forwarding address cannot be empty")
	ErrForwardingAddressInvalidScheme = errors.New("forwarding address must use http or https")
	ErrForwardingAddressInvalidURL    = errors.New("forwarding address must be a valid URL")
)

// ForwardingAddress parses a raw URL string into a ForwardingAddress object,
// validating that it's a well-formed URL with http or https scheme.
func ForwardingAddress(rawURL string) (types.ForwardingAddress, error) {
	if strings.TrimSpace(rawURL) == "" {
		return types.ForwardingAddress(""), ErrForwardingAddressEmpty
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return types.ForwardingAddress(""), ErrForwardingAddressInvalidURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return types.ForwardingAddress(""), ErrForwardingAddressInvalidScheme
	}

	if parsedURL.Host == "" {
		return types.ForwardingAddress(""), ErrForwardingAddressInvalidURL
	}

	// Strip trailing slashes from the path
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")

	return types.ForwardingAddress(parsedURL.String()), nil
}
