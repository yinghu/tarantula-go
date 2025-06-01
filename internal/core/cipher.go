package core

type Cipher interface {
	Encrypt(clearText string) string
	Decrypt(cipherText string) (string, error)
}
