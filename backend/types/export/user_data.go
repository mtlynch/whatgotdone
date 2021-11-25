package export

import "github.com/mtlynch/whatgotdone/backend/types"

type UserData struct {
	Entries     []JournalEntry   `json:"entries" yaml:"entries"`
	Drafts      []JournalEntry   `json:"drafts" yaml:"drafts"`
	Following   []types.Username `json:"following" yaml:"following"`
	Profile     UserProfile      `json:"profile" yaml:"profile"`
	Preferences Preferences      `json:"preferences" yaml:"preferences"`
}
