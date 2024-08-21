package utils

import (
	"fmt"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

// GenerateToken return paseto token
func GenerateToken(expire time.Duration, userID, role string) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetIssuer("bb-core")
	token.SetExpiration(time.Now().Add(expire))
	token.SetString("user_id", userID)
	token.SetString("role", role)

	key, err := paseto.V4SymmetricKeyFromHex(os.Getenv("TOKEN_ENCODED"))
	if err != nil {
		return "", fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	encrypted := token.V4Encrypt(key, nil)

	return encrypted, nil
}

// VerifyToken Validate token and return nil if it successes
func VerifyToken(token string) (*paseto.Token, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy("bb-core"))
	parser.AddRule(paseto.NotExpired())

	key, err := paseto.V4SymmetricKeyFromHex(os.Getenv("TOKEN_ENCODED"))
	if err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	parsedToken, err := parser.ParseV4Local(key, token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return parsedToken, nil
}
