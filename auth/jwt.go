package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type JWTUser struct {
	Name          string `json:"name"`
	Username      string `json:"username"`
	EmailVerified bool   `json:"emailVerified"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
}

func GetUser(encrypted, key string) (JWTUser, error) {
	result := JWTUser{}
	decrypted, err := decodeText(encrypted, key)
	if err != nil {
		err = errors.New(fmt.Sprintf("Decrypt failed: %v", err))
		return result, err
	}

	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(decrypted, jwt.MapClaims{})
	if err != nil {
		err = errors.New(fmt.Sprintf("JWT Token parsing failed: %v", err))
		return result, err
	}

	result.Name = token.Claims.(jwt.MapClaims)["name"].(string)
	result.Username = token.Claims.(jwt.MapClaims)["preferred_username"].(string)
	result.EmailVerified = token.Claims.(jwt.MapClaims)["email_verified"].(bool)
	result.FirstName = token.Claims.(jwt.MapClaims)["given_name"].(string)
	result.LastName = token.Claims.(jwt.MapClaims)["family_name"].(string)
	result.Email = token.Claims.(jwt.MapClaims)["email"].(string)

	return result, nil
}

func decodeText(state, key string) (string, error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(state)
	if err != nil {
		return "", err
	}
	// step: decrypt the cookie back in the expiration|token
	encoded, err := decryptDataBlock(cipherText, []byte(key))
	if err != nil {
		return "", errors.New("invalid encrypted input")
	}

	return string(encoded), nil
}

func decryptDataBlock(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("failed to decrypt the ciphertext, the text is too short")
	}
	nonce, input := cipherText[:nonceSize], cipherText[nonceSize:]

	return gcm.Open(nil, nonce, input, nil)
}
