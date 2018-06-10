package DBHelper

// BQuery query seter
type BQuery interface {
	GetCount(string, ...interface{}) (int, error)
	Raw(...interface{}) (interface{}, error)
}
