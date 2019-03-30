package model

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Comment is comment model.
type Comment struct {
	ID        uint32 `json:"id"`
	Content   string `json:"content"`
	ThreadID  uint32 `json:"threadId"`
	User      `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MarshalLogObject for zap logger.
func (c Comment) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt32("id", int32(c.ID))
	enc.AddString("content", c.Content)
	enc.AddInt32("threadID)", int32(c.ThreadID))
	if err := enc.AddObject("user", c.User); err != nil {
		return err
	}
	enc.AddTime("createdAt", c.CreatedAt)
	enc.AddTime("updatedAt", c.UpdatedAt)
	return nil
}

// CommentLList is list of comment.
type CommentList struct {
	Comments []*Comment `json:"comments"`
	HasNext  bool       `json:"hasNext"`
	Cursor   uint32     `json:"cursor"`
}

// MarshalLogObject for zap logger.
func (cl CommentList) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	zap.Array("comments", zapcore.ArrayMarshalerFunc(func(inner zapcore.ArrayEncoder) error {
		for _, c := range cl.Comments {
			if err := enc.AddObject("comment", c); err != nil {
				return err
			}
		}
		return nil
	}))

	enc.AddBool("hasNext", cl.HasNext)
	enc.AddInt32("cursor", int32(cl.Cursor))
	return nil
}
