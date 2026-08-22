package user

import "errors"


var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound 	  = errors.New("user not found")
	ErrFailedToDeleteUser = errors.New("failed to delete user")
)
