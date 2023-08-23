package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("bytes: didn't read enough random bytes")
	}
	return b, nil
}

func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("String: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// 32 bytes = 256 (caracteres existentes) * 32 (quantidade de caracteres) possibilidades
// 115792089237316195423570985008687907853269984665640564039457584007913129639936
// OWASP: tokens de sessão devem ter pelo menos 16 bytes
const SessionTokenBytes = 32

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
