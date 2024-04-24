package storage

import (
	"database/sql"
)

type Domain interface {
	Commit() error
}

// Handler is a interface for handling work.
type Handler[T Domain] interface {
	Handle(domain Domain, w Work) error
}

type Work interface {
	Type() string
}

// UnitOfWork is a collection of work.
type UnitOfWork[T Domain] struct {
	tx       *sql.Tx
	Works    []Work
	Handlers map[string]Handler[T]
}

func (u *UnitOfWork[T]) Add(w Work) {
	u.Works = append(u.Works, w)
}

func (u *UnitOfWork[T]) Commit() error {
	return nil
}

// Store is a collection of handlers.
type Store[T Domain] struct {
	db         *sql.DB
	domainFunc func(*sql.DB) T
	handlers   map[string]Handler[T]
}

func New[T Domain](db *sql.DB, domain T) *Store[T] {
	return &Store[T]{
		db: db,
	}
}

func (s *Store[T]) RegisterHandler(t string, h Handler[Domain]) {
	s.handlers[t] = h
}

func (s *Store[T]) Start() (*UnitOfWork[Domain], error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	return &UnitOfWork[Domain]{
		tx: tx,
	}, nil
}
