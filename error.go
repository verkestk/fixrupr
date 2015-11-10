package fixrupr

type dbError struct {
	query      string
	parameters []interface{}
	err        error
}

func newDbError(err error, query string, parameters []interface{}) error {
	return &dbError{query: query, parameters: parameters, err: err}
}

func (e dbError) Error() string {
	// verbose option:
	// return fmt.Sprintf("%s\n%s\n%v", e.err.Error(), e.query, e.parameters)
	return e.err.Error()
}
