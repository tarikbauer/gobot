package domain

type Code int

const (
	Ok Code = iota
	Cancelled
	Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated
)

func (c Code) String() string {
	return [...]string{
		"ok",
		"cancelled",
		"unknown",
		"invalid argument",
		"deadline exceeded",
		"not found",
		"already exists",
		"permission denied",
		"resource exhausted",
		"failed precondition",
		"aborted",
		"out of range",
		"unimplemented",
		"internal",
		"unavailable",
		"data loss",
		"unauthenticated",
	}[c]
}
