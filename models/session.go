package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/vitoraalmeida/lenslocked/rand"
)

// 32 bytes = 256 (caracteres existentes) * 32 (quantidade de caracteres) possibilidades
// 115792089237316195423570985008687907853269984665640564039457584007913129639936
// OWASP: tokens de sessão devem ter pelo menos 16 bytes
const MinBytesPerToken = 32

type Session struct {
	ID     int
	UserID int
	// Token só será atribuido quando criarmos uma nova sessão. Caso olhe os atributos de
	// uma instância de Session, não tera o token, pois só mantemos o token hash que não
	// pode ser revertido para o token original
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// Valor que determina qual será o tamanho em bytes do token de sessão
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	// tenta primeiro atualizar uma sessão existente com um novo token,
	// se não existir uma sessão, cria uma nova
	row := ss.DB.QueryRow(
		`UPDATE sessions SET token_hash = $2 WHERE user_id = $1 RETURNING id;`,
		session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	// quando não há nenhuma linha retornada, o pacote sql do go gera o erro ErrNoRows
	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(
			`INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id;`,
			session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}

	// checa por outros erros
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	row := ss.DB.QueryRow(
		`SELECT user_id FROM sessions WHERE token_hash=$1;`,
		tokenHash)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	var user User

	row = ss.DB.QueryRow(
		`SELECT email FROM users WHERE id=$1`, userID)
	err = row.Scan(&user.Email)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) hash(token string) string {
	// não utiliza bcrypt pois ele adiciona um salt em cada geração de hash,
	// de forma que seria necessário adicionar uma lógica para definir qual
	// é o salt utilizado (user o userId por exemplo) para que possamos
	// gerar novamente e comparar se é o mesmo token de sessão que geramos
	// originalmente. O bcrypt é um pouco demorado (pois faz um trabalho maior) e
	// no nosso caso utilizamos um payload grande (32 bytes). Para o caso em questão
	// uma função de hash sem salt, mas com entropia suficiente para tornar improvavel
	// gerar o mesmo hash é suficiente
	tokenHash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(tokenHash[:])
}
