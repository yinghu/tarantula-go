package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"encoding/hex"
	"fmt"
)

type Cipher struct {
	Gcm cipher.AEAD
}


func UseIt(clearText string) {
	key := make([]byte, 32)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	cipherText := gcm.Seal(nonce, nonce, []byte(clearText), nil)
	fmt.Println(hex.EncodeToString(cipherText))
	txt, err := gcm.Open(nil, cipherText[:gcm.NonceSize()], cipherText[gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(txt))
	//encrypt.CryptBlocks([]byte(clearText), clearTextBlock)
}
