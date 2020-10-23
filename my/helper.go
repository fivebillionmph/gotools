package my

type SQLRowInterface interface {
	Scan(dest ...interface{}) error
}
