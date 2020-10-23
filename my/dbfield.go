package my

import (
	"errors"
)

var ALLOWED_FIELD_TYPES = [...]string{"int", "float", "text", "blob"}

type DBField struct {
	field_type string
	field_name string
	db_name string
}

func NewDBField(field_type string, field_name string, db_name string) (*DBField, error) {
	found := false
	for _, ft := range ALLOWED_FIELD_TYPES {
		if field_type == ft {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("Invalid field type " + field_type)
	}

	db_field := DBField {
		field_type: field_type,
		field_name: field_name,
		db_name: db_name,
	}

	return &db_field, nil
}

type DBFieldConfig struct {
	Is_primary_key bool
	Map_db_to_go func(interface{}) interface{}
	Map_go_to_db func(interface{}) interface{}
	Validator func(interface{}) string
}
