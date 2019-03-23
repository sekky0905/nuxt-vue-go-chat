package model

import (
	"time"
)

// Thread is thread model.
type Thread struct {
	ID        uint32 `json:"id"`
	Title     string `json:"title"`
	User      `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ThreadList is list of thread.
type ThreadList struct {
	Threads []*Thread `json:"threads"`
	HasNext bool      `json:"hasNext"`
	Cursor  uint32    `json:"cursor"`
}
