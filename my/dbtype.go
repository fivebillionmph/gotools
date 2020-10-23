package my

import (
	"errors"
)

type DBType struct {
	table_name string
	db_fields []DBField
	locked bool
}

func NewDBType(table_name string) (*DBType, error) {
	if table_name == "" {
		return nil, errors.New("Table name cannot be blank")
	}

	db_type := DBType {
		table_name: table_name,
		db_fields: make([]DBField, 0),
		locked: false,
	}

	return &db_type, nil
}

func (self *DBType) AddDBField(field DBField) error {
	if self.locked {
		return errors.New("Table is locked")
	}

	for _, existing_field := range self.db_fields {
		if existing_field.db_name == field.db_name {
			return errors.New("Duplicated database field name: " + field.db_name)
		}
		if existing_field.field_name == field.db_name {
			return errors.New("Duplicated field name: " + field.field_name)
		}
	}

	self.db_fields = append(self.db_fields, field)

	return nil
}

func (self *DBType) Lock() {
	self.locked = true
}
