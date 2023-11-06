package repository

import (
	"backend-test/internal/api/handler"
	"backend-test/internal/domain"
	"context"
	"database/sql"
	"time"
)

const (
	createUser = `
	INSERT INTO verify_my.user (name, age, email, password, address)
	VALUES (?, ?, ?, ?, ?)`

	getUser = `
	SELECT id, name, age, email, password, address FROM verify_my.user WHERE id = ?
	`

	listUsers = `
	SELECT id, name, age, email, password, address FROM verify_my.user
	`

	deleteUser = `-- name: DeleteUser :exec
	DELETE FROM verify_my.user WHERE id = ?
	`

	updateUser = `
	UPDATE verify_my.user
	SET
		name = IF(?='', name, ?),
		age = IF(?=0, age, ?),
		email = IF(?='', email, ?),
		password = IF(?='', password, ?),
		address = IF(?='', address, ?)
	WHERE id = ?
	`
)

type Repository struct {
	timeout time.Duration
	db      *sql.DB
}

func NewUserRepo(timeout time.Duration, db *sql.DB) handler.Repository {
	return &Repository{timeout, db}
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.ExecContext(timeoutCtx, createUser,
		user.Name,
		user.Age,
		user.Email,
		user.Password,
		user.Address,
	)

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int32(id)

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int32) error {
	cancelCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.db.ExecContext(cancelCtx, deleteUser, id)
	return err
}

func (r *Repository) GetUser(ctx context.Context, id int32) (*domain.User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	row := r.db.QueryRowContext(timeoutCtx, getUser, id)
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Email,
		&user.Password,
		&user.Address,
	)

	return &user, err
}

func (r *Repository) ListUsers(ctx context.Context) ([]domain.User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	rows, err := r.db.QueryContext(timeoutCtx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Age,
			&user.Email,
			&user.Password,
			&user.Address,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.ExecContext(timeoutCtx, updateUser,
		user.Name,
		user.Name,
		user.Age,
		user.Age,
		user.Email,
		user.Email,
		user.Password,
		user.Password,
		user.Address,
		user.Address,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
