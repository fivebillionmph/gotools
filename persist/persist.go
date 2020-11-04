package persist

type SQLRowInterface interface {
	Scan(dest ...interface{}) error
}
