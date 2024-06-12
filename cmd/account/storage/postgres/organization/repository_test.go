package organization_test

import (
	"context"
	_ "embed"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/extreme-business/lingo/cmd/account/storage"
	"github.com/extreme-business/lingo/cmd/account/storage/postgres/organization"
	"github.com/extreme-business/lingo/cmd/account/storage/postgres/seed"
	"github.com/extreme-business/lingo/pkg/database"
	"github.com/extreme-business/lingo/pkg/database/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Run("should return a new repository", func(t *testing.T) {
		if got := organization.New(nil); got == nil {
			t.Errorf("New() = %v, want %v", got, nil)
		}
	})
}

func TestRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtest.SetupPostgres(t, "account", seed.RunMigrations)

	t.Run("should create a new organization", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		r := organization.New(database.NewDB(db))

		o := seed.NewOrganization(
			"268d2306-030a-43d5-9269-8afe666a8cf8",
			"Test Organization",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		expected := seed.NewOrganization(
			"268d2306-030a-43d5-9269-8afe666a8cf8",
			"Test Organization",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		got, err := r.Create(context.Background(), o)
		if err != nil {
			t.Errorf("Create() error = %v, want %v", err, nil)
		}

		if diff := cmp.Diff(expected, got); diff != "" {
			t.Errorf("Create() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the organization id already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		r := organization.New(database.NewDB(db))

		o := seed.NewOrganization(
			"1a707dff-65a3-4cb7-83ee-1b4e7a0ae29e",
			"Test 765434",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		if _, err := r.Create(ctx, o); err != nil {
			t.Errorf("Create() error = %v, want %v", err, nil)
		}

		o = seed.NewOrganization(
			"1a707dff-65a3-4cb7-83ee-1b4e7a0ae29e", // Same ID as the previous organization
			"Test 765437",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		rr, err := r.Create(ctx, o)
		if err == nil {
			t.Errorf("Create() error = %v, want %v", nil, err)
		}

		if !errors.Is(err, storage.ErrConflictOrganizationID) {
			t.Errorf("Create() error = %v, want %v", err, storage.ErrConflictOrganizationID)
		}

		if rr != nil {
			t.Errorf("Create() error = %v, want %v", r, nil)
		}
	})

	t.Run("should return an error if the organization legal name already exists", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		r := organization.New(database.NewDB(db))

		o := seed.NewOrganization(
			"2f5c8650-2913-41f1-a196-343c4a27ed75",
			"Test 9855475",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		_, err := r.Create(ctx, o)
		if err != nil {
			t.Errorf("Create() error = %v, want %v", err, nil)
		}

		o = seed.NewOrganization(
			"2f5c8650-2913-41f1-a196-343c4a27ed76", // Different ID
			"Test 9855475",                         // Same legal name as the previous organization
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		rr, err := r.Create(ctx, o)
		if err == nil {
			t.Errorf("Create() error = %v, want %v", nil, err)
		}

		if !errors.Is(err, storage.ErrConflictOrganizationLegalName) {
			t.Errorf("Create() error = %v, want %v", err, storage.ErrConflictOrganizationLegalName)
		}

		if rr != nil {
			t.Errorf("Create() error = %v, want %v", r, nil)
		}
	})
}

func TestRepository_Get(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtest.SetupPostgres(t, "account", seed.RunMigrations)

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"95d2f153-bbb9-4104-8f82-d619f0df5ca9",
				"Test Organization",
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
	})

	t.Run("should return an organization", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		r := organization.New(database.NewDB(db))

		expected := seed.NewOrganization(
			"95d2f153-bbb9-4104-8f82-d619f0df5ca9",
			"Test Organization",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)

		got, err := r.Get(ctx, uuid.MustParse("95d2f153-bbb9-4104-8f82-d619f0df5ca9"))
		if err != nil {
			t.Errorf("Get() error = %v, want %v", err, nil)
		}

		if diff := cmp.Diff(expected, got); diff != "" {
			t.Errorf("Get() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the organization does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		r := organization.New(database.NewDB(db))

		_, err := r.Get(ctx, uuid.MustParse("b12428a7-f5b0-49a1-b1a5-2f6a6cf3baf5"))
		if err == nil {
			t.Errorf("Get() error = %v, want %v", nil, err)
		}

		if !errors.Is(err, storage.ErrOrganizationNotFound) {
			t.Errorf("Get() error = %v, want %v", err, storage.ErrOrganizationNotFound)
		}
	})
}

func TestRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtest.SetupPostgres(t, "account", seed.RunMigrations)

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"c75823e3-ddc7-4170-ade8-9e7a8152f274",
				"Test Organization",
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
	})

	t.Run("should update an organization", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(database.NewDB(db))
		r := organization.New(recorder)

		o := seed.NewOrganization(
			"c75823e3-ddc7-4170-ade8-9e7a8152f274",
			"Test Organization 3",
			time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC), // create time is also changed, but it should not be updated
			time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		)

		expected := seed.NewOrganization(
			"c75823e3-ddc7-4170-ade8-9e7a8152f274",
			"Test Organization 3",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		)

		got, err := r.Update(ctx, o, []storage.OrganizationField{
			storage.OrganizationLegalName,
			storage.OrganizationUpdateTime,
		})
		if err != nil {
			t.Errorf("Update() error = %v, want %v", err, nil)
		}

		if diff := cmp.Diff(expected, got); diff != "" {
			t.Errorf("Update() mismatch (-want +got):\n%s", diff)
		}

		query := recorder.RowQueries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "UPDATE organizations SET legal_name = $1, update_time = $2 WHERE id = $3 RETURNING id, legal_name, create_time, update_time;"

		if query != expectedQuery {
			t.Errorf("Update() query = %v, want %v", query, expectedQuery)
		}
	})

	t.Run("should return an error if the organization does not exist", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		r := organization.New(database.NewDB(db))

		o := seed.NewOrganization(
			"a1f67eb8-f2f6-4321-85d5-34690ce9ec5d",
			"Test Organization 3",
			time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		)

		rr, err := r.Update(ctx, o, []storage.OrganizationField{
			storage.OrganizationLegalName,
			storage.OrganizationUpdateTime,
		})
		if err == nil {
			t.Errorf("Update() error = %v, want %v", nil, err)
		}

		if !errors.Is(err, storage.ErrOrganizationNotFound) {
			t.Errorf("Update() error = %v, want %v", err, storage.ErrOrganizationNotFound)
		}

		if rr != nil {
			t.Errorf("Update() error = %v, want %v", rr, nil)
		}
	})
}

func TestRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := dbtest.SetupPostgres(t, "account", seed.RunMigrations)

	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"debf1bcc-55c6-4816-9ff4-bf53a00084be",
				"Test Organization",
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		},
	})

	t.Run("should delete an organization", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.Connect(ctx, t, dbc.ConnectionString)
		recorder := dbtest.NewRecorder(database.NewDB(db))
		r := organization.New(recorder)

		err := r.Delete(ctx, uuid.MustParse("debf1bcc-55c6-4816-9ff4-bf53a00084be"))
		if err != nil {
			t.Errorf("Delete() error = %v, want %v", err, nil)
		}

		query := recorder.ExecQueries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		expectedQuery := "DELETE FROM organizations WHERE id = $1;"

		if query != expectedQuery {
			t.Errorf("Delete() query = %v, want %v", query, expectedQuery)
		}
	})
}

func setupTestDatabaseForList(t *testing.T) *dbtest.PostgresContainer {
	dbc := dbtest.SetupPostgres(t, "account", seed.RunMigrations)
	seed.Run(t, dbc.ConnectionString, seed.State{
		Organizations: []*storage.Organization{
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205264",
				"test",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			),
			seed.NewOrganization(
				"7bb443e5-8974-44c2-8b7c-b95124205265",
				"test2",
				time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			),
		},
		Users: []*storage.User{},
	})
	return dbc
}

func connectAndCreateRepo(t *testing.T, db database.Conn) (*organization.Repository, *dbtest.Recorder) {
	t.Helper()
	recorder := dbtest.NewRecorder(db)
	return organization.New(recorder), recorder
}

func assertOrganizations(t *testing.T, expect, actual []*storage.Organization) {
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("Create() mismatch (-want +got):\n%s", diff)
	}
}

func assertQuery(t *testing.T, recorder *dbtest.Recorder, expectedQuery string) {
	if expectedQuery != "" {
		if len(recorder.Queries) != 1 {
			t.Errorf("expected 1 query, got %d", len(recorder.Queries))
			return // no need to check the query
		}

		query := recorder.Queries[0].Query
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
		if query != expectedQuery {
			t.Errorf("expected %q, got %q", expectedQuery, query)
		}
	}
}

func assertError(t *testing.T, expected, actual error) {
	if expected == nil && actual != nil {
		t.Errorf("unexpected error: %v", actual)
	}
	if expected != nil && actual == nil {
		t.Errorf("expected error: %v, got nil", expected)
	}
	if expected != nil && !errors.Is(actual, expected) {
		t.Errorf("expected error: %v, got: %v", expected, actual)
	}
}

func TestRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dbc := setupTestDatabaseForList(t)

	testCases := []struct {
		name                  string
		ctx                   context.Context
		db                    database.Conn
		pagination            storage.Pagination
		orderBy               storage.OrganizationOrderBy
		conditions            []storage.Condition
		expectedOrganizations []*storage.Organization
		expectedQuery         string
		expectedError         error
	}{
		{
			name:       "should return all organizations",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{},
			orderBy:    storage.OrganizationOrderBy{},
			conditions: []storage.Condition{},
			expectedOrganizations: []*storage.Organization{
				seed.NewOrganization(
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test",
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				),
				seed.NewOrganization(
					"7bb443e5-8974-44c2-8b7c-b95124205265",
					"test2",
					time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				),
			},
			expectedQuery: "SELECT id, legal_name, create_time, update_time FROM organizations;",
			expectedError: nil,
		},
		{
			name:       "should return all organizations with pagination",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{Limit: 1, Offset: 1},
			orderBy:    storage.OrganizationOrderBy{},
			conditions: []storage.Condition{},
			expectedOrganizations: []*storage.Organization{
				// 7bb443e5-8974-44c2-8b7c-b95124205264 is skipped because of the offset
				seed.NewOrganization(
					"7bb443e5-8974-44c2-8b7c-b95124205265",
					"test2",
					time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				),
			},
			expectedQuery: "SELECT id, legal_name, create_time, update_time FROM organizations LIMIT $1 OFFSET $2;",
			expectedError: nil,
		},
		{
			name:       "should return all organizations with order by",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{},
			orderBy: storage.OrganizationOrderBy{
				{Field: storage.OrganizationUpdateTime, Direction: storage.DESC},
			},
			conditions: []storage.Condition{},
			expectedOrganizations: []*storage.Organization{
				seed.NewOrganization(
					"7bb443e5-8974-44c2-8b7c-b95124205265",
					"test2",
					time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				),
				seed.NewOrganization(
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test",
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				),
			},
			expectedQuery: "SELECT id, legal_name, create_time, update_time FROM organizations ORDER BY update_time DESC;",
			expectedError: nil,
		},
		{
			name:       "should return all organizations with conditions",
			ctx:        context.Background(),
			db:         database.NewDB(dbtest.Connect(context.Background(), t, dbc.ConnectionString)),
			pagination: storage.Pagination{},
			orderBy:    storage.OrganizationOrderBy{},
			conditions: []storage.Condition{
				storage.OrganizationByLegalNameCondition{
					Wildcard:  true,
					LegalName: "test",
				},
				storage.OrganizationByLegalNameCondition{
					Wildcard:  false,
					LegalName: "test",
				},
			},
			expectedOrganizations: []*storage.Organization{
				seed.NewOrganization(
					"7bb443e5-8974-44c2-8b7c-b95124205264",
					"test",
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				),
			},
			expectedQuery: "SELECT id, legal_name, create_time, update_time FROM organizations WHERE legal_name LIKE $1 AND legal_name = $2;",
			expectedError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, recorder := connectAndCreateRepo(t, tc.db)
			organizations, err := repo.List(tc.ctx, tc.pagination, tc.orderBy, tc.conditions...)
			assertOrganizations(t, tc.expectedOrganizations, organizations)
			assertQuery(t, recorder, tc.expectedQuery)
			assertError(t, tc.expectedError, err)
		})
	}
}
