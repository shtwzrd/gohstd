package common

type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	PublicKey string `json:"publickey"`
}
