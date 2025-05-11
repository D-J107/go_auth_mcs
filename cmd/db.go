package main

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("environment variable DATABASE URL not set")
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		panic("cant establish connection to remote Postgre db")
	}
	return &DB{pool: pool}, nil
}

func (db *DB) GetRoleID(ctx context.Context, value string) (int, error) {
	var id int
	if err := db.pool.QueryRow(ctx, `SELECT id FROM roles WHERE value=$1`, value).Scan(&id); err != nil {
		return 0, ErrRoleNotFound
	}
	return id, nil
}

func (db *DB) CreateNewUser(ctx context.Context, username, hashedPwd, email string) error {
	tx, _ := db.pool.Begin(ctx)
	defer tx.Rollback(ctx)

	var userId int
	if err := tx.QueryRow(ctx, `INSERT INTO users (username, email, password, balance) VALUES ($1,$2,$3,$4) RETURNING id`,
		username, email, hashedPwd, 3000).Scan(&userId); err != nil {
		return ErrUserExists
	}

	userRoleId, err := db.GetRoleID(ctx, "USER")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES ($1, $2)`,
		userId, userRoleId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (username, hashedPwd string, err error) {
	err = db.pool.QueryRow(ctx, `SELECT username, password FROM users WHERE email=$1`, email).Scan(&username, &hashedPwd)
	if err != nil {
		return "", "", err
	}
	return username, hashedPwd, nil
}
