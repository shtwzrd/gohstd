package common

// CommandRepo [sitory] is an interface wrapping the functions for working with
// command and invocation data
type CommandRepo interface {
	// InsertInvocations writes Invocations to storage
	InsertInvocations(user string, invocs Invocations) (err error)

	// GetInvocations returns the [n] most recent Invocations for the given
	// user
	GetInvocations(user string, n int) (result Invocations, err error)

	// GetCommands returns the [n] most recent Commands for the given user
	GetCommands(user string, n int) (result Commands, err error)
}
