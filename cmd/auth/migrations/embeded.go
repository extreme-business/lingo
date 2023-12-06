package migrations

import "embed"

// FS contains the SQL migrations and atlas sum
//
//go:embed *.sql atlas.sum
var FS embed.FS
