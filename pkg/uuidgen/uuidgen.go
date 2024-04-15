package uuidgen

import "github.com/google/uuid"

type Generator struct {
	genFunc func() uuid.UUID
}

func New(genFunc func() uuid.UUID) *Generator {
	return &Generator{
		genFunc: genFunc,
	}
}

// Default returns a new generator that uses uuid.New
func Default() *Generator {
	return New(uuid.New)
}

func (g *Generator) New() uuid.UUID { return g.genFunc() }
