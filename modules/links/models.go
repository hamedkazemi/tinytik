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
	MD5       string    `gorm:"column:md5;index:idx_md5;size:32;"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (q *Links) TableName() string {
	return "Links"
}

type CreateResponse struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

//func insertTilExist(md5 string) gorm.DB {
//	for {
//
//		if !condition {
//			break
//		}
//	}
//}
