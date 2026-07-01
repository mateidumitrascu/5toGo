package token

import (
	"database/sql"
	"fmt"
)

type TokenSQLiteRepo struct {
	db *sql.DB
}

const TokensTable = "tokens"

func NewTokenSQLiteRepo(db *sql.DB) *TokenSQLiteRepo {
	return &TokenSQLiteRepo{db}
}

func (tr *TokenSQLiteRepo) Create(token *AuthToken) (*AuthToken, error) {
	_, err := tr.db.Exec("INSERT INTO "+TokensTable+"(hash, uid, expiry) VALUES (?, ?, ?)", token.Hash, token.UID, token.Expiry)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}
	return token, nil
}

func (tr *TokenSQLiteRepo) FindToken(token string) (*AuthToken, error) {
	tokenHash := HashToken(token)
	rows, err := tr.db.Query("SELECT hash, uid, expiry FROM "+TokensTable+" WHERE hash=?", tokenHash)
	if err != nil {
		return nil, fmt.Errorf("error in find token: %w", err)
	}

	//nolint:errcheck
	defer rows.Close()

	tokens, err := parseRows(rows)
	if err != nil {
		return nil, fmt.Errorf("error finding token by hash: %w", err)
	}
	if len(tokens) == 0 {
		return nil, nil
	}

	if len(tokens) > 1 {
		return nil, ErrDuplicateTokenFound
	}

	t := tokens[0]
	return &t, nil
}

func (tr *TokenSQLiteRepo) DeleteByHash(hash string) error {
	_, err := tr.db.Exec("DELETE FROM "+TokensTable+" WHERE hash=?", hash)
	if err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}
	return nil
}

func parseRows(rows *sql.Rows) ([]AuthToken, error) {
	var tokens []AuthToken
	for rows.Next() {
		var t AuthToken
		if err := rows.Scan(&t.Hash, &t.UID, &t.Expiry); err != nil {
			return []AuthToken{}, fmt.Errorf("error scanning token rows: %w", err)
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}
