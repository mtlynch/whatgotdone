package export

import "github.com/mtlynch/whatgotdone/backend/types"

type JournalEntry struct {
	Date         types.EntryDate `json:"date" yaml:"date"`
	Markdown     string          `json:"markdown" yaml:"markdown"`
	LastModified string          `json:"lastModified" yaml:"last_modified"`
}
