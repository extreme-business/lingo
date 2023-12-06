package uuidgen

import "github.com/google/uuid"

type Generator struct {
	f func() uuid.UUID
}

func New(f func() uuid.UUID) *Generator {
	return &Generator{f: f}
}

// Default returns a new generator that uses uuid.New.
func Default() *Generator { return New(uuid.New) }

func (g *Generator) New() uuid.UUID { return g.f() }
