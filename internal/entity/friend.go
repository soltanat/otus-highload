package entity

import (
	"github.com/google/uuid"
)

type Friend struct {
	UserID   uuid.UUID
	FriendID uuid.UUID
}

type FriendEventType int

const (
	FriendAdd FriendEventType = iota
	FriendDelete
)

type FriendEvent struct {
	FriendID uuid.UUID
	UserID   uuid.UUID
	Event    FriendEventType
}
