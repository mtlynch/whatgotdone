package parse

import (
	"errors"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/types"
)

var (
	ErrEmptyEntryContent = errors.New("entry must not be empty or whitespace")
)

func EntryContent(content string) (types.EntryContent, error) {
	stripped := strings.TrimSpace(content)

	if len(stripped) == 0 {
		return types.EntryContent(""), ErrEmptyEntryContent
	}
	return types.EntryContent(stripped), nil
}
