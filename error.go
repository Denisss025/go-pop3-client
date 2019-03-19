package pop3

import "fmt"

var (
	// ErrAlreadyQuit is an error that is returned when an attempt to
	// write is made after the "QUIT" command is sent.
	ErrAlreadyQuit = fmt.Errorf("pop3: already quit from server")
	// ErrWriteAfterClose is an error that is returned when an attempt
	// to write is made after the connection is closed.
	ErrWriteAfterClose = fmt.Errorf("pop3: write after close")
)

// Error is a POP3 error.
type Error struct {
	// Context is context where an error happend.
	Context string
	// Err is an error.
	Err error
}

func newError(context string, errx error) (err Error) {
	var ok bool
	err, ok = errx.(Error)

	if ok {
		err.Context = fmt.Sprintf("%s: %s", err.Context, context)
	} else {
		err.Context = context
		err.Err = errx
	}
	return
}

// Error returns POP3 error message.
func (e Error) Error() (s string) {
	return fmt.Sprintf("pop3: %s: %v", e.Context, e.Err)
}
