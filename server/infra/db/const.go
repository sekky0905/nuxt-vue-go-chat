package db

// RepositoryMethod is method of Repository.
type RepositoryMethod string

// methods of Repository.
const (
	RepositoryMethodREAD   = "READ"
	RepositoryMethodInsert = "INSERT"
	RepositoryMethodUPDATE = "UPDATE"
	RepositoryMethodDELETE = "DELETE"
	RepositoryMethodLIST   = "LIST"
)
