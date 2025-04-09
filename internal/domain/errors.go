package domain

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error codes
const (
	CodeNotFound              = "not_found"
	CodeInterestAlreadyExists = "interest_already_exists"
)

var (
	ErrNotFound              = NewError(CodeNotFound, "not found")
	ErrInterestAlreadyExists = NewError(CodeInterestAlreadyExists, "interest already exists")
)
