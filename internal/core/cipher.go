package core

const CIPHER_KEY_NAME string = "cipherkey"

type Cipher interface {
	Encrypt(clearText string) string
	Decrypt(cipherText string) (string, error)
}
