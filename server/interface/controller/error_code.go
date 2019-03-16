package controller

// ErrCode is error code.
type ErrCode string

// Server error.
const (
	InternalFailure    ErrCode = "InternalFailure"
	InternalDBFailure  ErrCode = "InternalDBFailure"
	InternalSQLFailure ErrCode = "InternalSQLFailure"
	ServerError        ErrCode = "ServerError"
)

// User error.
const (
	InvalidParameterValueFailure  ErrCode = "InvalidParameterValueFailure"
	InvalidParametersValueFailure ErrCode = "InvalidParametersValueFailure"
	NoSuchDataFailure             ErrCode = "NoSuchDataFailure"
	RequiredFailure               ErrCode = "RequiredError"
	AlreadyExistsFailure          ErrCode = "AlreadyExistsFailure"
	AuthenticationFailure         ErrCode = "AuthenticationFailure"
)
