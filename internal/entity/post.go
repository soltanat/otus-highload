package entity

import (
	"github.com/google/uuid"
	"time"
)

type RangeFilter struct {
	From *time.Time
	To   *time.Time
}

type PostFilter struct {
	AuthorIDs []uuid.UUID
	//Star      *bool
	Limit *int
	//Offset *int
	//CreatedAt *RangeFilter
}

type Post struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID  uuid.UUID `json:"author_id"`
}

type PostEventType int

const (
	PostCreate PostEventType = iota
	PostUpdate
	PostDelete
)

type PostEvent struct {
	Post  Post
	Event PostEventType
}
