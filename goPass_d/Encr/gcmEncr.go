package encr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"log"

	"golang.org/x/crypto/pbkdf2"
)

type GCM_Encr struct {
	gcmCipher cipher.AEAD
}

func (g *GCM_Encr) GenNonce() ([]byte, error) {
	s := g.gcmCipher.NonceSize()
	nonce := make([]byte, s)
	_, err := rand.Read(nonce)
	return nonce, err
}

func (g *GCM_Encr) Init(key []byte, salt []byte) error {
	paddedKey := pbkdf2.Key(key, salt, 32000, 32, sha512.New)
	log.Println("PADDED KEY: ", paddedKey)
	if block, err := aes.NewCipher(paddedKey); err != nil {
		return err
	} else {
		if g.gcmCipher, err = cipher.NewGCM(block); err != nil {
			return err
		}
	}
	return nil
}

func (g *GCM_Encr) Encrypt(plainText []byte) ([]byte, []byte, error) {
	nonce, e := g.GenNonce()
	cipherText := g.gcmCipher.Seal(nil, nonce, plainText, nil)
	log.Println("NONCE", nonce)
	log.Println("CIPHER", cipherText)
	return cipherText, nonce, e
}

func (g *GCM_Encr) Decrypt(cipher []byte, nonce []byte) ([]byte, error) {
	return g.gcmCipher.Open(nil, nonce, cipher, nil)
}
