package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
)

type Env struct {
	db *sql.DB
}

type NullString string

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("column is not a string")
	}
	*s = NullString(strVal)
	return nil
}

func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}

type Account struct {
	ID          uint64     `json:"id"`
	Root        uint64     `json:"root"`
	Name        string     `json:"name"`
	Name_en     NullString `json:"name_en"`
	Description NullString `json:"description"`
	Is_active   bool       `json:"is_active"`
}
