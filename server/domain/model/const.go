package model

const InvalidID = 0

// InvalidReasonForDeveloper is InvalidReason message for developer.
type InvalidReasonForDeveloper string

// String return as string.
func (p InvalidReasonForDeveloper) String() string {
	return string(p)
}

// DomainModelNameForDeveloper is Model name for developer.
type DomainModelNameForDeveloper string

// String return as string.
func (p DomainModelNameForDeveloper) String() string {
	return string(p)
}

// Model name for developer.
const (
	DomainModelNameUserForDeveloper    DomainModelNameForDeveloper = "User"
	DomainModelNameSessionForDeveloper DomainModelNameForDeveloper = "Session"
)

// DomainModelNameForUser is Model name for user.
type DomainModelNameForUser string

// String return as string.
func (p DomainModelNameForUser) String() string {
	return string(p)
}

// Model name for user.
const (
	DomainModelNameUserForUser    DomainModelNameForUser = "ユーザー"
	DomainModelNameSessionForUser DomainModelNameForUser = "セッション"
)

// PropertyNameForDeveloper is property name for developer.
type PropertyNameForDeveloper string

// String return as string.
func (p PropertyNameForDeveloper) String() string {
	return string(p)
}

// Property name for developer.
const (
	IDPropertyForDeveloper       PropertyNameForDeveloper = "id"
	NamePropertyForDeveloper     PropertyNameForDeveloper = "name"
	PassWordPropertyForDeveloper PropertyNameForDeveloper = "password"
)

// PropertyNameForUser is Property name for user.
type PropertyNameForUser string

// String return as string.
func (p PropertyNameForUser) String() string {
	return string(p)
}

// Property name for user.
const (
	IDPropertyForUser       PropertyNameForUser = "ID"
	NamePropertyForUser     PropertyNameForUser = "名前"
	PassWordPropertyForUser PropertyNameForUser = "パスワード"
)

// PropertyNameKV is the Key/Value of PropertyNameForDeveloper and PropertyNameForUser,
var PropertyNameKV = map[PropertyNameForDeveloper]PropertyNameForUser{
	IDPropertyForDeveloper:       IDPropertyForUser,
	NamePropertyForDeveloper:     NamePropertyForUser,
	PassWordPropertyForDeveloper: PassWordPropertyForUser,
}

// == for test ==
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
)

// error message for test
const (
	ErrorMessageForTest = "some error has occurred"
)
