package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/google/uuid"
)

var _ user.Repository = &Repository{}

type Repository struct {
	db database.DB
}

func NewRepository(db database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

const createQuery = `INSERT INTO users (id, username, email, password, create_time, update_time)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, username, email, create_time, update_time
`

// Create a new user
func (r *Repository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		createQuery,
		u.ID,
		u.Username,
		u.Email,
		u.Password,
		u.CreateTime,
		u.UpdateTime,
	)

	var n user.User
	if err := row.Scan(
		&n.ID,
		&n.Username,
		&n.Email,
		&n.CreateTime,
		&n.UpdateTime,
	); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &n, nil
}

const getByIDQuery = `SELECT id, username, email, create_time, update_time
FROM users
WHERE id = $1
`

// GetByID get a user by id
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, getByIDQuery, id)

	var user user.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreateTime, &user.UpdateTime); err != nil {
		return nil, err
	}

	return &user, nil
}

const getByUsernameQuery = `SELECT id, username, email, create_time, update_time
FROM users
WHERE username = $1`

func (r *Repository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, getByUsernameQuery, username)

	var user user.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreateTime, &user.UpdateTime); err != nil {
		return nil, err
	}

	return &user, nil
}

const updateQueryTemplate = `UPDATE users u
SET %s
WHERE u.id = $%d
RETURNING u.id, u.username, u.email, u.create_time, u.update_time`

func (r *Repository) Update(ctx context.Context, u *user.User, fields ...user.Field) (*user.User, error) {
	set := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields)+1)

	for _, f := range fields {
		switch f {
		case user.Username:
			set = append(set, fmt.Sprintf("u.username = $%d", len(values)+1))
			values = append(values, u.Username)
		case user.Email:
			set = append(set, fmt.Sprintf("u.email = $%d", len(values)+1))
			values = append(values, u.Email)
		case user.Password:
			set = append(set, fmt.Sprintf("u.password = $%d", len(values)+1))
			values = append(values, u.Password)
		default:
			return nil, fmt.Errorf("unknown field %q", f)
		}
	}

	if len(set) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add the user ID to the end of the values slice
	values = append(values, u.ID)

	query := fmt.Sprintf(updateQueryTemplate, strings.Join(set, ", "), len(values))
	row := r.db.QueryRowContext(ctx, query, values...)

	var user user.User
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreateTime,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

const deleteQuery = `DELETE FROM users WHERE id = $1`

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteQuery, id)
	return err
}
