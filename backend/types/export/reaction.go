package export

import "github.com/mtlynch/whatgotdone/backend/types"

type Reaction struct {
	Username  types.Username `json:"username" yaml:"username"`
	Symbol    string         `json:"symbol" yaml:"symbol"`
	Timestamp string         `json:"timestamp" yaml:"timestamp"`
}
