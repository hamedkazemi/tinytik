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
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	URL       string    `gorm:"column:url" json:"url"`
	Hash      string    `gorm:"column:hash" json:"hash"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdateAt  time.Time `gorm:"column:update_at" json:"update_at"`
}

// TableName sets the insert table name for this struct type
func (q *Links) TableName() string {
	return "Links"
}
