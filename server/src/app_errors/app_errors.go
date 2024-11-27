package app_errors

import "errors"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrRegistryNotFound        = errors.New("registry not found")
	ErrUserIsNotAdmin          = errors.New("the user is not an admin")
	ErrVerificationPinNotFound = errors.New("verification pin not found")
)

const (
	IncorrectUsernameOrPassword = "Incorrect username or password"
	InternalServerError         = "Internal server error"
	EmailAlreadyExists          = "Email already exists"
	UsernameAlreadyExists       = "Username already exists"
	UsernameNotFound            = "User not found"
	RegistryNotFound            = "Registry not found"
	InvalidRegistryStep         = "Invalid registry step"
	VerificationPinNotFound     = "Verification pin not found"
	InvalidInterest             = "Invalid interest"
	CantFollowYourself          = "Can't follow yourself"
	AlreadyFollowing            = "The user already follows this user"
	NotFollowing                = "The user is not following this user"
	UserShouldModifyItself      = "The user should modify its own profile"
	UserIsNotAdmin              = "The user is not an admin"
	UserBlocked                 = "User is blocked"
)
