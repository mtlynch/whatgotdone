package export

import "github.com/mtlynch/whatgotdone/backend/types"

type JournalEntry struct {
	Date         types.EntryDate `json:"date"`
	Markdown     string          `json:"markdown"`
	LastModified string          `json:"lastModified"`
}
