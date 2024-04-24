package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestOrderBy_Validate(t *testing.T) {
	tests := []struct {
		name    string
		o       OrderBy
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("OrderBy.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockRepository_Create(t *testing.T) {
	type fields struct {
		CreateFunc     func(context.Context, *User) (*User, error)
		GetFunc        func(context.Context, uuid.UUID) (*User, error)
		ListFunc       func(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
		GetByEmailFunc func(context.Context, string) (*User, error)
		UpdateFunc     func(context.Context, *User, ...Field) (*User, error)
		DeleteFunc     func(context.Context, uuid.UUID) error
	}
	type args struct {
		ctx context.Context
		u   *User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockRepository{
				CreateFunc:     tt.fields.CreateFunc,
				GetFunc:        tt.fields.GetFunc,
				ListFunc:       tt.fields.ListFunc,
				GetByEmailFunc: tt.fields.GetByEmailFunc,
				UpdateFunc:     tt.fields.UpdateFunc,
				DeleteFunc:     tt.fields.DeleteFunc,
			}
			got, err := m.Create(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockRepository_Get(t *testing.T) {
	type fields struct {
		CreateFunc     func(context.Context, *User) (*User, error)
		GetFunc        func(context.Context, uuid.UUID) (*User, error)
		ListFunc       func(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
		GetByEmailFunc func(context.Context, string) (*User, error)
		UpdateFunc     func(context.Context, *User, ...Field) (*User, error)
		DeleteFunc     func(context.Context, uuid.UUID) error
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockRepository{
				CreateFunc:     tt.fields.CreateFunc,
				GetFunc:        tt.fields.GetFunc,
				ListFunc:       tt.fields.ListFunc,
				GetByEmailFunc: tt.fields.GetByEmailFunc,
				UpdateFunc:     tt.fields.UpdateFunc,
				DeleteFunc:     tt.fields.DeleteFunc,
			}
			got, err := m.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockRepository_List(t *testing.T) {
	type fields struct {
		CreateFunc     func(context.Context, *User) (*User, error)
		GetFunc        func(context.Context, uuid.UUID) (*User, error)
		ListFunc       func(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
		GetByEmailFunc func(context.Context, string) (*User, error)
		UpdateFunc     func(context.Context, *User, ...Field) (*User, error)
		DeleteFunc     func(context.Context, uuid.UUID) error
	}
	type args struct {
		ctx context.Context
		p   Pagination
		s   OrderBy
		c   []Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockRepository{
				CreateFunc:     tt.fields.CreateFunc,
				GetFunc:        tt.fields.GetFunc,
				ListFunc:       tt.fields.ListFunc,
				GetByEmailFunc: tt.fields.GetByEmailFunc,
				UpdateFunc:     tt.fields.UpdateFunc,
				DeleteFunc:     tt.fields.DeleteFunc,
			}
			got, err := m.List(tt.args.ctx, tt.args.p, tt.args.s, tt.args.c...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockRepository_GetByEmail(t *testing.T) {
	type fields struct {
		CreateFunc     func(context.Context, *User) (*User, error)
		GetFunc        func(context.Context, uuid.UUID) (*User, error)
		ListFunc       func(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
		GetByEmailFunc func(context.Context, string) (*User, error)
		UpdateFunc     func(context.Context, *User, ...Field) (*User, error)
		DeleteFunc     func(context.Context, uuid.UUID) error
	}
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockRepository{
				CreateFunc:     tt.fields.CreateFunc,
				GetFunc:        tt.fields.GetFunc,
				ListFunc:       tt.fields.ListFunc,
				GetByEmailFunc: tt.fields.GetByEmailFunc,
				UpdateFunc:     tt.fields.UpdateFunc,
				DeleteFunc:     tt.fields.DeleteFunc,
			}
			got, err := m.GetByEmail(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockRepository.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockRepository.GetByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockRepository_Update(t *testing.T) {
	type fields struct {
		CreateFunc     func(context.Context, *User) (*User, error)
		GetFunc        func(context.Context, uuid.UUID) (*User, error)
		ListFunc       func(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
		GetByEmailFunc func(context.Context, string) (*User, error)
		UpdateFunc     func(context.Context, *User, ...Field) (*User, error)
		DeleteFunc     func(context.Context, uuid.UUID) error
	}
	type args struct {
		ctx    context.Context
		u      *User
		fields []Field
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockRepository{
				CreateFunc:     tt.fields.CreateFunc,
				GetFunc:        tt.fields.GetFunc,
				ListFunc:       tt.fields.ListFunc,
				GetByEmailFunc: tt.fields.GetByEmailFunc,
				UpdateFunc:     tt.fields.UpdateFunc,
				DeleteFunc:     tt.fields.DeleteFunc,
			}
			got, err := m.Update(tt.args.ctx, tt.args.u, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockRepository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockRepository_Delete(t *testing.T) {
	type fields struct {
		CreateFunc     func(context.Context, *User) (*User, error)
		GetFunc        func(context.Context, uuid.UUID) (*User, error)
		ListFunc       func(context.Context, Pagination, OrderBy, ...Condition) ([]*User, error)
		GetByEmailFunc func(context.Context, string) (*User, error)
		UpdateFunc     func(context.Context, *User, ...Field) (*User, error)
		DeleteFunc     func(context.Context, uuid.UUID) error
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockRepository{
				CreateFunc:     tt.fields.CreateFunc,
				GetFunc:        tt.fields.GetFunc,
				ListFunc:       tt.fields.ListFunc,
				GetByEmailFunc: tt.fields.GetByEmailFunc,
				UpdateFunc:     tt.fields.UpdateFunc,
				DeleteFunc:     tt.fields.DeleteFunc,
			}
			if err := m.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("MockRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
