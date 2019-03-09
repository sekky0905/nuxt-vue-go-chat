package model

import "fmt"

// RepositoryMethod is method of Repository.
type RepositoryMethod string

// methods of Repository.
const (
	RepositoryMethodREAD   RepositoryMethod = "READ"
	RepositoryMethodInsert RepositoryMethod = "INSERT"
	RepositoryMethodUPDATE RepositoryMethod = "UPDATE"
	RepositoryMethodDELETE RepositoryMethod = "DELETE"
	RepositoryMethodLIST   RepositoryMethod = "LIST"
)

// RequiredError is not existing necessary value error.
type RequiredError struct {
	BaseErr error
	PropertyNameForDeveloper
	PropertyNameForUser
}

// Error returns error message.
func (e *RequiredError) Error() string {
	return fmt.Sprintf("%s is required", e.PropertyNameForDeveloper)
}

// InvalidParamError is inappropriate parameter error。
type InvalidParamError struct {
	BaseErr error
	PropertyNameForDeveloper
	PropertyNameForUser
	PropertyValue             interface{}
	InvalidReasonForDeveloper string
	InvalidReasonForUser      string
}

// Error returns error message.
func (e *InvalidParamError) Error() string {
	return fmt.Sprintf("%s, %v, is invalid, %s", e.PropertyNameForDeveloper, e.PropertyValue, e.InvalidReasonForDeveloper)
}

// NoSuchDataError is not existing specified data error.
type NoSuchDataError struct {
	BaseErr error
	PropertyNameForDeveloper
	PropertyNameForUser
	PropertyValue interface{}
	DomainModelNameForDeveloper
	DomainModelNameForUser
}

// Error returns error message.
func (e *NoSuchDataError) Error() string {
	return fmt.Sprintf("no such data, %s: %v, %s", e.PropertyNameForDeveloper, e.PropertyValue, e.DomainModelNameForDeveloper)
}

// RepositoryError is Repository error.
type RepositoryError struct {
	BaseErr          error
	RepositoryMethod RepositoryMethod
	DomainModelNameForDeveloper
	DomainModelNameForUser
}

// Error returns error message.
func (e *RepositoryError) Error() string {
	return fmt.Sprintf("failed Repository operation, %s, %s", e.RepositoryMethod, e.DomainModelNameForDeveloper)
}

// SQLError is SQL error.
type SQLError struct {
	BaseErr                   error
	InvalidReasonForDeveloper InvalidReasonForDeveloper
}

// Error returns error message.
func (e *SQLError) Error() string {
	return e.InvalidReasonForDeveloper.String()
}
