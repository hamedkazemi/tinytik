package links

import (
	"database/sql"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// Links Models should only be concerned with database schema, more strict checking should be put in validator.
type Links struct {
	Hash      string    `gorm:"column:hash;index:idx_hash;size:10;uniqueIndex;primary_key" json:"hash"`
	URL       string    `gorm:"column:url" json:"url"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (q *Links) TableName() string {
	return "Links"
}
