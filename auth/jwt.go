package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type JWTUser struct {
	Name          string   `json:"name"`
	Username      string   `json:"username"`
	EmailVerified bool     `json:"emailVerified"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Email         string   `json:"email"`
	Roles         []string `json:"roles"`
}

// from https://github.com/keycloak/keycloak-gatekeeper/blob/d87453446b6dbe6aea36d069dac7aef7b42e6c5e/doc.go
const (
	claimRealmAccess    = "realm_access"
	claimResourceAccess = "resource_access"
	claimResourceRoles  = "roles"
	claimGroups         = "groups"
)

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

	log.Info("Claim list")

	for key, val := range token.Claims.(jwt.MapClaims) {
		log.Info(fmt.Sprintf("Claim Key: %v, value: %v", key, val))
	}

	result.Roles = getRoles(token.Claims.(jwt.MapClaims))

	return result, nil
}

// from https://github.com/keycloak/keycloak-gatekeeper/blob/1b7ee69ed9ef1b471be86a18b37db82bc950a4f6/user_context.go
func getRoles(claims jwt.MapClaims) []string {
	var roleList []string

	if realmRoles, found := claims[claimRealmAccess].(map[string]interface{}); found {
		if roles, found := realmRoles[claimResourceRoles]; found {
			for _, r := range roles.([]interface{}) {
				roleList = append(roleList, fmt.Sprintf("%s", r))
			}
		}
	}

	if accesses, found := claims[claimResourceAccess].(map[string]interface{}); found {
		for name, list := range accesses {
			scopes := list.(map[string]interface{})
			if roles, found := scopes[claimResourceRoles]; found {
				for _, r := range roles.([]interface{}) {
					roleList = append(roleList, fmt.Sprintf("%s:%s", name, r))
				}
			}
		}
	}

	return roleList
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
