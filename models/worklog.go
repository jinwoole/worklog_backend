package models

import "time"

type WorkLog struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
