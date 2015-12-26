package common

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Secret indicates an encrypted user password
type Secret []byte
