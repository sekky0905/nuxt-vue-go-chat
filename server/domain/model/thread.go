package model

import (
	"time"

	"go.uber.org/zap"
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
	if err := enc.AddObject("user", t.User); err != nil {
		return err
	}
	enc.AddTime("createdAt", t.CreatedAt)
	enc.AddTime("updatedAt", t.CreatedAt)
	return nil
}

// ThreadList is list of thread.
type ThreadList struct {
	Threads []*Thread `json:"threads"`
	HasNext bool      `json:"hasNext"`
	Cursor  uint32    `json:"cursor"`
}

// MarshalLogObject for zap logger.
func (t ThreadList) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	zap.Array("threads", zapcore.ArrayMarshalerFunc(func(inner zapcore.ArrayEncoder) error {
		for _, t := range t.Threads {
			if err := enc.AddObject("thread", t); err != nil {
				return err
			}
		}
		return nil
	}))

	enc.AddBool("hasNext", t.HasNext)
	enc.AddInt32("cursor", int32(t.Cursor))
	return nil
}
