package core

type Authenticator interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password string, hash string) error

}
