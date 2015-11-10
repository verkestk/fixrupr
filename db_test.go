package fixrupr

import (
	"database/sql"

	. "gopkg.in/check.v1"
)

type mockDb struct {
	queries []string
	args    [][]interface{}
}

func (m *mockDb) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.queries = append(m.queries, query)
	m.args = append(m.args, args)
	return nil, nil
}

func (m *mockDb) clear() {
	m.queries = []string{}
	m.args = [][]interface{}{}
}

func (s *MySuite) Test_fixr_create(c *C) {
}

func (s *MySuite) Test_fixr_insert(c *C) {
}

func (s *MySuite) Test_fixr_drop(c *C) {
}

func (s *MySuite) Test_fixr_schema(c *C) {
}

func (s *MySuite) Test_fixr_table(c *C) {
}

func (s *MySuite) Test_fixr_function(c *C) {
}

func (s *MySuite) Test_fixr_exec(c *C) {
}

func (s *MySuite) Test_fixr_load(c *C) {
}

func (s *MySuite) Test_getInsertFields(c *C) {
}

func (s *MySuite) Test_generateInsert(c *C) {
}
