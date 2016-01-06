package common

import (
	"errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Secret indicates an encrypted user password
type Secret []byte

// Validate returns an error if a User struct does not satisfy the following
// conditions:
//  * Username is 3 or more characters
//  * Username is composed solely of alphanumeric characters
//  * Email is a valid email address
//  * Password is 8 or more characters
func (u User) Validate() (err error) {
	if !valid.IsEmail(u.Email) {
		err = errors.New(fmt.Sprintf("'%s' is not a valid email address", u.Email))
	}
	if len(u.Username) < 3 {
		err = errors.New(fmt.Sprintf("'%s' is less than 3 characters", u.Username))
	}
	if len(u.Password) < 8 {
		err = errors.New("Password must be at least 8 characters long")
	}
	if !valid.IsAlphanumeric(u.Username) {
		err = errors.New(fmt.Sprintf("'%s' is less than 3 characters", u.Username))
	}

	return err
}
