package db

import (
	"encoding/json"
	"errors"
	"strings"
	"database/sql"
)

type Server struct {

}

type DBObj interface {
	To_json() ([]byte, error)
	Update_from_json([]byte) error
}

type DBType interface {
	Name() string
	Table_name() string
	Db_fields() []string
	New_from_json([]byte) (*DBObj, error)
	Delete(*DBObj) error
	Load_obj_from_row(SQLRowInterface) (*DBObj, error)
}

func Load_by_json(db *sql.DB, db_type *DBType, data []byte, limit int, offset int) (res []DBObj, err error) {
	var args []interface{}
	var predicate string
	args, predicate, err = stmt_predicate(db, db_type, data, limit, offset)
	if err != nil { return }

	rows, err := db.Query("SELECT * " + predicate, args...)
	if err != nil { return }

	res = make([]DBObj, 0, limit)
	var single_result *DBObj
	for rows.Next() {
		single_result, err = (*db_type).Load_obj_from_row(rows)
		if err != nil { return }
		res = append(res, *single_result)
	}

	return
}

func Count_by_json(db *sql.DB, db_type *DBType, data[]byte) (count int, err error) {
	var args []interface{}
	var predicate string
	args, predicate, err = stmt_predicate(db, db_type, data, 1, 0)
	if err != nil { return }
	row := db.QueryRow("SELECT COUNT(*) " + predicate, args...)
	err = row.Scan(&count)
	return
}

func stmt_predicate(db *sql.DB, db_type *DBType, data []byte, limit int, offset int) (args []interface{}, predicate string, err error) {
	var data_map map[string]string
	err = json.Unmarshal(data, &data_map)
	if err != nil { return }

	seen_fields := make(map[string]struct{})
	where_predicate := make([]string, 0, len(data_map))
	args = make([]interface{}, 0, len(data_map) + 2)
	for k, v := range data_map {
		if _, ok := seen_fields[k]; ok {
			err = errors.New("Repeated field: " + k)
			return
		}
		if Is_valid_field(db_type, k) {
			err = errors.New("Field does not exist: " + k)
			return
		}
		seen_fields[k] = struct{}{}
		where_predicate = append(where_predicate, k + " = ?")
		args = append(args, v)
	}
	where_str := ""
	if len(where_predicate) > 0 {
		where_str = "WHERE " + strings.Join(where_predicate, " AND ")
	}
	args = append(args, offset)
	args = append(args, limit)

	predicate = `
		FROM ` + (*db_type).Table_name() + `
		` + where_str + `
		ORDER BY id ASC
		LIMIT ?, ?
	`

	return
}

func Is_valid_field(db_type *DBType, name string) bool {
	fields := (*db_type).Db_fields()
	for _, f := range fields {
		if f == name {
			return true
		}
	}
	return false
}
