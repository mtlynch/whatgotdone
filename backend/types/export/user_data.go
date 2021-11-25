package export

import "github.com/mtlynch/whatgotdone/backend/types"

type UserData struct {
	Entries     []JournalEntry   `json:"entries"`
	Drafts      []JournalEntry   `json:"drafts"`
	Following   []types.Username `json:"following"`
	Profile     UserProfile      `json:"profile"`
	Preferences Preferences      `json:"preferences"`
}
