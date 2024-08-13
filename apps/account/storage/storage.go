package storage

type Pagination struct {
	Limit  int
	Offset int
}

type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)
