package export

import "github.com/mtlynch/whatgotdone/backend/types"

type Reaction struct {
	Username  types.Username `json:"username"`
	Symbol    string         `json:"symbol"`
	Timestamp string         `json:"timestamp"`
}
