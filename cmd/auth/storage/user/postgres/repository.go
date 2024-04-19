package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	userIDConstraint       = "users_pkey"
	userUsernameConstraint = "users_username_key"
	userEmailConstraint    = "users_email_key"
)

var _ user.Repository = &Repository{}

type Repository struct {
	db database.DB
}

func New(db database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

const createQuery = `INSERT INTO users (id, username, email, password, create_time, update_time)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, username, email, create_time, update_time
;`

// Create a new user.
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				switch pqErr.Constraint {
				case userIDConstraint:
					return nil, user.ErrUniqueIDConflict
				case userUsernameConstraint:
					return nil, user.ErrUniqueUsernameConflict
				case userEmailConstraint:
					return nil, user.ErrUniqueEmailConflict
				}
			}
		}

		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &n, nil
}

const getByIDQuery = `SELECT id, username, email, create_time, update_time
FROM users
WHERE id = $1
;`

// GetByID get a user by id.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, getByIDQuery, id)

	var u user.User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.CreateTime,
		&u.UpdateTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		return nil, err
	}

	return &u, nil
}

const getByEmailQuery = `SELECT id, username, password, email, create_time, update_time
FROM users
WHERE email = $1
;`

// GetByEmail get a user by email.
// This is used for authentication and so it also returns the password.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, getByEmailQuery, email)

	var u user.User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Email,
		&u.CreateTime,
		&u.UpdateTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		return nil, err
	}

	return &u, nil
}

const updateQueryTemplate = `UPDATE users
SET %s
WHERE id = $%d
RETURNING id, username, email, create_time, update_time;`

func (r *Repository) Update(ctx context.Context, in *user.User, fields ...user.Field) (*user.User, error) {
	if len(fields) == 0 {
		return nil, user.ErrNoFieldsToUpdate
	}

	set := make([]string, 0, len(fields)) // set clauses, e.g. "username = $1", "email = $2"
	args := make([]interface{}, 0, len(fields)+1)

	for _, f := range fields {
		switch f {
		case user.Username:
			set = append(set, fmt.Sprintf("username = $%d", len(args)+1))
			args = append(args, in.Username)
		case user.Email:
			set = append(set, fmt.Sprintf("email = $%d", len(args)+1))
			args = append(args, in.Email)
		case user.Password:
			set = append(set, fmt.Sprintf("password = $%d", len(args)+1))
			args = append(args, in.Password)
		case user.UpdateTime:
			set = append(set, fmt.Sprintf("update_time = $%d", len(args)+1))
			args = append(args, in.UpdateTime)
		case user.CreateTime:
			fallthrough
		default:
			return nil, fmt.Errorf("unknown or non allowed field: %q", f)
		}
	}

	if len(set) == 0 {
		return nil, user.ErrNoFieldsToUpdate
	}

	// Add the user ID to the end of the args slice
	args = append(args, in.ID)

	query := fmt.Sprintf(
		updateQueryTemplate,
		strings.Join(set, ", "),
		len(args), // the parameter number for the user ID
	)
	row := r.db.QueryRowContext(ctx, query, args...)

	var u user.User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.CreateTime,
		&u.UpdateTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		return nil, err
	}

	return &u, nil
}

const deleteQuery = `DELETE FROM users WHERE id = $1;`

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if n == 0 {
		return user.ErrNotFound
	}

	return nil
}

func generatePredicates(argOffset int, conditions []user.Condition) (string, []interface{}, error) {
	sw := strings.Builder{}
	var args []interface{}

	for _, c := range conditions {
		switch t := c.(type) {
		case *user.SearchEmail:
			sw.WriteString(fmt.Sprintf("email LIKE $%d", argOffset+len(args)+1))
			args = append(args, "%"+t.Email+"%")
		default:
			return "", nil, fmt.Errorf("unknown condition type: %T", t)
		}
	}

	return sw.String(), args, nil
}

func generateSorting(sorting []user.Sort) (string, error) {
	if len(sorting) == 0 {
		return "", nil
	}

	sortFields := make([]string, 0, len(sorting))
	for _, s := range sorting {
		var dir string
		switch s.Direction {
		case user.ASC:
			dir = "ASC"
		case user.DESC:
			dir = "DESC"
		default:
			return "", fmt.Errorf("unknown direction: %q", s.Direction)
		}

		switch s.Field {
		case user.Username:
			sortFields = append(sortFields, fmt.Sprintf("username %s", dir))
		case user.Email:
			sortFields = append(sortFields, fmt.Sprintf("email %s", dir))
		case user.CreateTime:
			sortFields = append(sortFields, fmt.Sprintf("create_time %s", dir))
		case user.UpdateTime:
			sortFields = append(sortFields, fmt.Sprintf("update_time %s", dir))
		case user.Password:
			fallthrough
		default:
			return "", fmt.Errorf("unknown or non allowed field: %q", s.Field)
		}
	}

	return strings.Join(sortFields, ", "), nil
}

const listQueryTemplate = `SELECT 
    id, 
    username, 
    email, 
    create_time, 
    update_time
FROM users
%s
ORDER BY %s
LIMIT $%d
OFFSET $%d;`

// List implements user.Repository.
func (r *Repository) List(ctx context.Context, pagination user.Pagination, sorting []user.Sort, conditions ...user.Condition) ([]*user.User, error) {
	whereClause := "WHERE "
	var args []interface{}

	predicates, args, err := generatePredicates(0, conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if len(predicates) > 0 {
		whereClause += predicates
	} else {
		whereClause = ""
	}

	sortClause, err := generateSorting(sorting)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	limitArgIndex := len(args) + 1
	offsetArgIndex := limitArgIndex + 1

	query := fmt.Sprintf(
		listQueryTemplate,
		whereClause,
		sortClause,
		limitArgIndex,  // Update for LIMIT
		offsetArgIndex, // Update for OFFSET
	)

	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		if err = rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.CreateTime,
			&u.UpdateTime,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &u)
	}

	return users, nil
}
