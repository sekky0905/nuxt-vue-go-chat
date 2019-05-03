package model

// const
const (
	InvalidID         = 0
	SessionIDAtCookie = "SESSION_ID"
)

// InvalidReason is InvalidReason message for developer.
type InvalidReason string

// String return as string.
func (p InvalidReason) String() string {
	return string(p)
}

// DomainModelName is Model name for developer.
type DomainModelName string

// String return as string.
func (p DomainModelName) String() string {
	return string(p)
}

// Model name.
const (
	DomainModelNameUser    DomainModelName = "User"
	DomainModelNameSession DomainModelName = "Session"
	DomainModelNameThread  DomainModelName = "Thread"
	DomainModelNameComment DomainModelName = "Comment"
)

// PropertyName is property name for developer.
type PropertyName string

// String return as string.
func (p PropertyName) String() string {
	return string(p)
}

// Property name for developer.
const (
	IDProperty       PropertyName = "ID"
	NameProperty     PropertyName = "Name"
	TitleProperty    PropertyName = "Title"
	PassWordProperty PropertyName = "Password"
	ThreadIDProperty PropertyName = "ThreadID"
)

// FailedToBeginTx is error of tx begin.
const FailedToBeginTx InvalidReason = "failed to begin tx"

// User
const (
	UserNameForTest             = "testUserName"
	PasswordForTest             = "testPasswor"
	UserValidIDForTest   uint32 = 1
	UserInValidIDForTest uint32 = 2
)

// Session
const (
	SessionValidIDForTest   = "testValidSessionID12345678"
	SessionInValidIDForTest = "testInvalidSessionID12345678"
	TitleForTest            = "TitleForTest"
)

// Thread
const (
	ThreadValidIDForTest   uint32 = 1
	ThreadInValidIDForTest uint32 = 2
)

// Comment
const (
	CommentValidIDForTest   uint32 = 1
	CommentInValidIDForTest uint32 = 2
	CommentContentForTest          = "ContentForTest"
)

// error message for test
const (
	ErrorMessageForTest = "some error has occurred"
)
