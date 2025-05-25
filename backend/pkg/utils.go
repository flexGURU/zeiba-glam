package pkg

import (
	"log"
	"strconv"

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

func Float64ToPgTypeNumeric(f float64) pgtype.Numeric {
	var amount pgtype.Numeric
	if err := amount.Scan(strconv.FormatFloat(f, 'f', -1, 64)); err != nil {
		log.Println("not float")
		return pgtype.Numeric{
			Valid: false,
		}
	}

	return amount
}

func PgTypeArrayToString(a pgtype.Array[string]) []string {
	return a.Elements
}

func StringToUint32(s string) (uint32, error) {
	if s == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, Errorf(INVALID_ERROR, "invalid id/page: %s", err.Error())
	}

	return uint32(id), nil
}

func StringToBool(s string) (bool, error) {
	if s == "" {
		return false, nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, Errorf(INVALID_ERROR, "invalid is_admin: %s", err.Error())
	}

	return b, nil
}

func StringToFloat64(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, Errorf(INVALID_ERROR, "invalid price: %s", err.Error())
	}

	return f, nil
}
