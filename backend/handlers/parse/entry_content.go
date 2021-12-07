package parse

import (
	"github.com/mtlynch/whatgotdone/backend/types"
)

func EntryContent(content string) (types.EntryContent, error) {
	// TODO: Apply validation on the entry content.
	return types.EntryContent(content), nil
}
