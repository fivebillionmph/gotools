package db

type SQLRowInterface interface {
	Scan(dest ...interface{}) error
}
