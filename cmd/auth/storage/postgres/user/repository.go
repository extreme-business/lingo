package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/dwethmar/lingo/cmd/auth/storage"
	"github.com/dwethmar/lingo/pkg/database"
	"github.com/google/uuid"
	"github.com/lib/pq"

	_ "embed"
)

const (
	userIDConstraint          = "users_pkey"
	userDisplayNameConstraint = "users_display_name_key"
	userEmailConstraint       = "users_email_key"
)

var _ storage.UserRepository = &Repository{}

type Repository struct {
	dbConn           database.Conn
	listTemplateFunc sync.Once          // compile the list template only once
	listTemplate     *template.Template // compiled list template
}

func New(dbConn database.Conn) *Repository {
	return &Repository{
		dbConn: dbConn,
	}
}

// scan scans a user from a sql.Row or sql.Rows.
func scan(u *storage.User, f func(dest ...any) error) error {
	return f(
		&u.ID,
		&u.OrganizationID,
		&u.DisplayName,
		&u.Email,
		&u.CreateTime,
		&u.UpdateTime,
	)
}

const createQuery = `INSERT INTO users (id, organization_id,  display_name, email, password, create_time, update_time)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, organization_id, display_name, email, create_time, update_time
;`

// Create a new user.
func (r *Repository) Create(ctx context.Context, u *storage.User) (*storage.User, error) {
	row := r.dbConn.QueryRow(
		ctx,
		createQuery,
		u.ID,
		u.OrganizationID,
		u.DisplayName,
		u.Email,
		u.Password,
		u.CreateTime,
		u.UpdateTime,
	)

	var n storage.User
	if err := scan(&n, row.Scan); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				switch pqErr.Constraint {
				case userIDConstraint:
					return nil, storage.ErrConflictUserID
				case userEmailConstraint:
					return nil, storage.ErrConflictUserEmail
				}
			}
		}

		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &n, nil
}

const getByIDQuery = `SELECT id, organization_id, display_name, email, create_time, update_time
FROM users
WHERE id = $1
;`

// GetByID get a user by id.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*storage.User, error) {
	row := r.dbConn.QueryRow(ctx, getByIDQuery, id)

	var u storage.User
	if err := row.Scan(
		&u.ID,
		&u.OrganizationID,
		&u.DisplayName,
		&u.Email,
		&u.CreateTime,
		&u.UpdateTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

const getByEmailQuery = `SELECT id, organization_id, display_name, password, email, create_time, update_time
FROM users
WHERE email = $1
;`

// GetByEmail get a user by email.
// This is used for authentication and so it also returns the password.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*storage.User, error) {
	row := r.dbConn.QueryRow(ctx, getByEmailQuery, email)

	var u storage.User
	if err := row.Scan(
		&u.ID,
		&u.OrganizationID,
		&u.DisplayName,
		&u.Password,
		&u.Email,
		&u.CreateTime,
		&u.UpdateTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

const updateQueryTemplate = `UPDATE users
SET %s
WHERE id = $%d
RETURNING id, organization_id, display_name, email, create_time, update_time;`

func (r *Repository) Update(ctx context.Context, in *storage.User, fields []storage.UserField) (*storage.User, error) {
	if len(fields) == 0 {
		return nil, storage.ErrNoUserFieldsToUpdate
	}

	set := make([]string, 0, len(fields)) // set clauses, e.g. "display_name = $1", "email = $2"
	args := make([]interface{}, 0, len(fields)+1)

	for _, f := range fields {
		switch f {
		case storage.UserDisplayName:
			set = append(set, fmt.Sprintf("display_name = $%d", len(args)+1))
			args = append(args, in.DisplayName)
		case storage.UserEmail:
			set = append(set, fmt.Sprintf("email = $%d", len(args)+1))
			args = append(args, in.Email)
		case storage.UserPassword:
			set = append(set, fmt.Sprintf("password = $%d", len(args)+1))
			args = append(args, in.Password)
		case storage.UserUpdateTime:
			set = append(set, fmt.Sprintf("update_time = $%d", len(args)+1))
			args = append(args, in.UpdateTime)
		case storage.UserOrganizationID:
			set = append(set, fmt.Sprintf("organization_id = $%d", len(args)+1))
			args = append(args, in.OrganizationID)
		case storage.UserID:
			return nil, storage.ErrImmutableUserID
		case storage.UserCreateTime:
			return nil, storage.ErrImmutableUserCreateTime
		default:
			return nil, fmt.Errorf("field %s: %w", f, storage.ErrUserUnknownField)
		}
	}

	// Add the user ID to the end of the args slice
	args = append(args, in.ID)

	query := fmt.Sprintf(
		updateQueryTemplate,
		strings.Join(set, ", "),
		len(args), // the parameter number for the user ID
	)

	row := r.dbConn.QueryRow(ctx, query, args...)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("failed to run update query: %w", err)
	}

	var u storage.User
	if err := scan(&u, row.Scan); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed scan user: %w", err)
	}

	return &u, nil
}

const deleteQuery = `DELETE FROM users WHERE id = $1;`

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.dbConn.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if n == 0 {
		return storage.ErrUserNotFound
	}

	return nil
}

// generatePredicates generates the WHERE clause predicates for the list query.
func generatePredicates(argOffset int, conditions []storage.Condition) ([]string, []interface{}, error) {
	var predicates []string
	var args []interface{}

	for _, c := range conditions {
		switch t := c.(type) {
		case storage.UserByOrganizationIDCondition:
			predicates = append(predicates, fmt.Sprintf("u.organization_id = $%d", len(args)+argOffset+1))
			args = append(args, t.OrganizationID)
		default:
			return nil, nil, fmt.Errorf("unknown or non allowed condition: %T", c)
		}
	}

	return predicates, args, nil
}

//go:embed list.tmpl.sql
var listQueryTemplate []byte

type listQueryTemplateParams struct {
	Predicates  []string
	Sorting     []storage.UserSort
	LimitParam  string
	OffsetParam string
}

// List implements user.Repository.
func (r *Repository) List(ctx context.Context, pagination storage.Pagination, sorting storage.UserOrderBy, conditions ...storage.Condition) ([]*storage.User, error) {
	predicates, args, err := generatePredicates(0, conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var limitParam, offsetParam string
	if pagination.Limit > 0 {
		limitParam = fmt.Sprintf("$%d", len(args)+1)
		args = append(args, pagination.Limit)
	}

	if pagination.Offset > 0 {
		offsetParam = fmt.Sprintf("$%d", len(args)+1)
		args = append(args, pagination.Offset)
	}

	if err = sorting.Validate(); err != nil {
		return nil, fmt.Errorf("sorting validation failed: %w", err)
	}

	// Compile the list template only once
	r.listTemplateFunc.Do(func() {
		t, tErr := template.New("list").Parse(string(listQueryTemplate))
		if err != nil {
			err = fmt.Errorf("failed to parse list query template: %w", tErr)
		}
		r.listTemplate = t
	})

	if err != nil {
		return nil, err
	}

	w := &strings.Builder{}
	err = r.listTemplate.Execute(w, listQueryTemplateParams{
		Predicates:  predicates,
		Sorting:     sorting,
		LimitParam:  limitParam,
		OffsetParam: offsetParam,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute list query template: %w", err)
	}

	query := w.String()
	rows, err := r.dbConn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*storage.User
	for rows.Next() {
		var u storage.User
		if err = scan(&u, rows.Scan); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
