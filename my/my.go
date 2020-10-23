package my

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	db *sql.DB
}

func NewConnection(hostname string, database string, username string, password string, port int) (*Connection, error) {
	db_str := username + ":" + password + "@" + hostname + "/" + database
	db, err := sql.Open("mysql", db_str)
	if err != nil { return nil, err }

	cxn := Connection {
		db: db,
	}
	return &cxn, nil
}
