package domains

import "time"

type User struct {
	Id         int64      `db:"id,omitempty"`
	Username   string     `db:"username"`
	Password   string     `db:"password" json:"-"`
	Nickname   string     `db:"nickname"`
	CreateTime *time.Time `db:"create_time,omitempty"`
	UpdateTime *time.Time `db:"update_time,omitempty"`
}
