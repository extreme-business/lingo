package organization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/pkg/database"
	"github.com/google/uuid"
	"github.com/lib/pq"

	_ "embed"
)

const (
	orgIDConstraint         = "organizations_pkey"
	orgLegalNameConstraint  = "organizations_legal_name_key"
	orgSlugUniqueConstraint = "organizations_slug_key"
	orgSlugFormatConstraint = "organizations_chk_slug"
)

var _ storage.OrganizationRepository = &Repository{}

type Repository struct {
	dbConn           database.Conn
	listTemplateFunc sync.Once          // compile the list template only once
	listTemplate     *template.Template // compiled list template
}

func New(db database.Conn) *Repository {
	return &Repository{
		dbConn: db,
	}
}

// scan scans a organization from a sql.Row or sql.Rows.
func scan(f func(dest ...any) error, o *storage.Organization) error {
	if err := f(
		&o.ID,
		&o.LegalName,
		&o.Slug,
		&o.CreateTime,
		&o.UpdateTime,
	); err != nil {
		return fmt.Errorf("failed to scan organization: %w", err)
	}

	return nil
}

const createQuery = `INSERT INTO organizations (id, legal_name, slug, create_time, update_time)
VALUES ($1, $2, $3, $4, $5)
RETURNING id,  legal_name, slug, create_time, update_time
;`

// Create a new organization.
func (r *Repository) Create(ctx context.Context, u *storage.Organization) (*storage.Organization, error) {
	row := r.dbConn.QueryRow(
		ctx,
		createQuery,
		u.ID,
		u.LegalName,
		u.Slug,
		u.CreateTime,
		u.UpdateTime,
	)

	var n storage.Organization
	if err := scan(row.Scan, &n); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				switch pqErr.Constraint {
				case orgIDConstraint:
					return nil, storage.ErrConflictOrganizationID
				case orgLegalNameConstraint:
					return nil, storage.ErrConflictOrganizationLegalName
				case orgSlugUniqueConstraint:
					return nil, storage.ErrConflictOrganizationSlug
				case orgSlugFormatConstraint:
					return nil, storage.ErrInvalidOrganizationSlug
				}
			}
		}

		return nil, err
	}

	return &n, nil
}

const getByIDQuery = `SELECT id, legal_name, slug, create_time, update_time
FROM organizations
WHERE id = $1
;`

// Get.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*storage.Organization, error) {
	row := r.dbConn.QueryRow(ctx, getByIDQuery, id)

	var o storage.Organization
	if err := scan(row.Scan, &o); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrOrganizationNotFound
		}

		return nil, err
	}

	return &o, nil
}

const updateQueryTemplate = `UPDATE organizations
SET %s
WHERE id = $%d
RETURNING id, legal_name, slug, create_time, update_time;`

func (r *Repository) Update(ctx context.Context, in *storage.Organization, fields []storage.OrganizationField) (*storage.Organization, error) {
	if len(fields) == 0 {
		return nil, storage.ErrNoUserFieldsToUpdate
	}

	set := make([]string, 0, len(fields)) // set clauses, e.g. "legal_name = $1"
	args := make([]interface{}, 0, len(fields)+1)

	for _, f := range fields {
		switch f {
		case storage.OrganizationLegalName:
			set = append(set, fmt.Sprintf("legal_name = $%d", len(args)+1))
			args = append(args, in.LegalName)
		case storage.OrganizationSlug:
			set = append(set, fmt.Sprintf("slug = $%d", len(args)+1))
			args = append(args, in.Slug)
		case storage.OrganizationUpdateTime:
			set = append(set, fmt.Sprintf("update_time = $%d", len(args)+1))
			args = append(args, in.UpdateTime)
		case storage.OrganizationID:
			return nil, fmt.Errorf("field %s: %w", f, storage.ErrImmutableOrganizationCreateTime)
		case storage.OrganizationCreateTime:
			return nil, fmt.Errorf("field %s: %w", f, storage.ErrImmutableOrganizationID)
		default:
			return nil, fmt.Errorf("field %s: %w", f, storage.ErrUnknownOrganizationField)
		}
	}

	// Add the organization ID to the end of the args slice
	args = append(args, in.ID)

	query := fmt.Sprintf(
		updateQueryTemplate,     // the update query template
		strings.Join(set, ", "), // set clauses
		len(args),               // the parameter index for the organization ID
	)
	row := r.dbConn.QueryRow(ctx, query, args...)

	var o storage.Organization
	if err := scan(row.Scan, &o); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrOrganizationNotFound
		}

		return nil, err
	}

	return &o, nil
}

const deleteQuery = `DELETE FROM organizations WHERE id = $1;`

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.dbConn.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
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
func generatePredicates(_ int, c []storage.Condition) ([]string, []interface{}, error) {
	var predicates []string
	var args []interface{}

	for _, c := range c {
		switch t := c.(type) {
		case storage.OrganizationByLegalNameCondition:
			arg := t.LegalName
			if t.Wildcard {
				arg = fmt.Sprintf("%%%s%%", arg)
				predicates = append(predicates, fmt.Sprintf("legal_name LIKE $%d", len(args)+1))
			} else {
				predicates = append(predicates, fmt.Sprintf("legal_name = $%d", len(args)+1))
			}
			args = append(args, arg)
		default:
			return nil, nil, fmt.Errorf("unknown condition type: %T", c)
		}
	}

	return predicates, args, nil
}

//go:embed list.tmpl.sql
var listQueryTemplate []byte

type listQueryTemplateParams struct {
	Predicates  []string
	Sorting     []storage.OrganizationSort
	LimitParam  string
	OffsetParam string
}

// List organizations.
func (r *Repository) List(ctx context.Context, pagination storage.Pagination, sorting storage.OrganizationOrderBy, conditions ...storage.Condition) ([]*storage.Organization, error) {
	var predicates []string
	var args []interface{}
	var err error

	if len(conditions) > 0 {
		predicates, args, err = generatePredicates(0, conditions)
		if err != nil {
			return nil, fmt.Errorf("failed to list organizations: %w", err)
		}
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
		panic(err)
	}

	rows, err := r.dbConn.Query(ctx, w.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	defer rows.Close()

	var organizations []*storage.Organization
	for rows.Next() {
		var o storage.Organization
		if err = scan(rows.Scan, &o); err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}

		organizations = append(organizations, &o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	return organizations, nil
}
