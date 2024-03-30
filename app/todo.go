package app

import "time"

type Todo struct {
	ID        int        `db:"id"`
	Task      string     `db:"task"`
	Status    Status     `db:"status"`
	Priority  string     `db:"priority"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
