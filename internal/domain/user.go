package domain

import "time"

type User struct {
	ID        int64
	Username  string
	Password  string
	Nickname  string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
