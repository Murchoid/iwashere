package errors

import "fmt"

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Domain errors
var (
	ErrNoteNotFound    = &Error{Code: "NOTE_NOT_FOUND", Message: "note not found"}
	ErrSessionNotFound = &Error{Code: "SESSION_NOT_FOUND", Message: "session not found"}
	ErrConfigNotFound  = &Error{Code: "CONFIG_NOT_FOUND", Message: "config not found"}
	ErrGitNotFound     = &Error{Code: "GIT_NOT_FOUND", Message: "not a git repository"}
)
