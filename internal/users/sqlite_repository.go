package users

import "database/sql"

type UserSQLiteRepo struct {
	db *sql.DB
}

func (r *UserSQLiteRepo) Create() error {
	return nil
}
