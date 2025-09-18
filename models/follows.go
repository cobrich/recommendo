package models

import "time"

type Follow struct {
	ID          int       `db:"follows_id"`
	FolllowerID int       `db:"follower_id"`
	FollowingID int       `db:"following_id"`
	CreatedAt   time.Time `db:"created_at"`
}
