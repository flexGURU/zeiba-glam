package pkg

import (
	"log"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgTypeNumericToFloat64(n pgtype.Numeric) float64 {
	f, err := n.Float64Value()
	if err != nil {
		log.Println("not float")
		return 0
	}

	return f.Float64
}
