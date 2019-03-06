package model

// InvalidReasonForDeveloper is InvalidReason message for developer.
type InvalidReasonForDeveloper string

// String return as string.
func (p InvalidReasonForDeveloper) String() string {
	return string(p)
}
