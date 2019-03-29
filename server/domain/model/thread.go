package model

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Thread is thread model.
type Thread struct {
	ID        uint32 `json:"id"`
	Title     string `json:"title"`
	*User     `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MarshalLogObject for zap logger.
func (t Thread) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt32("id", int32(t.ID))
	enc.AddString("title", t.Title)
	if err := enc.AddObject("sessionID", t.User); err != nil {
		return err
	}
	enc.AddTime("created_at", t.CreatedAt)
	enc.AddTime("updated_at", t.CreatedAt)
	return nil
}

// ThreadList is list of thread.
type ThreadList struct {
	Threads []*Thread `json:"threads"`
	HasNext bool      `json:"hasNext"`
	Cursor  uint32    `json:"cursor"`
}
