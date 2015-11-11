package fixrupr

import (
	"database/sql"
	"os"

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
	conf := s.mock_fixrConf(c)
	def, _ := conf.load()

	c.Assert(def, NotNil)

	hostname, _ := os.Hostname()
	conn := &mockDb{}
	fixr := &fixr{
		conn:       conn,
		def:        def,
		prefix:     "v_test",
		schemaName: "jamila",
	}

	err := fixr.create()
	c.Check(err, IsNil)
	c.Assert(conn.queries, HasLen, 9)
	c.Assert(conn.args, HasLen, 9)

	c.Check(conn.queries[0], Equals, "insert into `jamila`.schemas (name, prefix, hostname) values (?, ?, ?)")
	c.Assert(conn.args[0], HasLen, 3)
	c.Check(conn.args[0][0], Equals, "blog")
	c.Check(conn.args[0][1], Equals, "v_test")
	c.Check(conn.args[0][2], Equals, hostname)

	c.Check(conn.queries[1], Equals, "create schema `v_test_blog`")
	c.Check(conn.args[1], HasLen, 0)

	c.Check(conn.queries[2], Equals, "choo-choo")
	c.Check(conn.args[2], HasLen, 0)

	c.Check(conn.queries[3], Equals, "egyptian")
	c.Check(conn.args[3], HasLen, 0)

	c.Check(conn.queries[4], Equals, "turkish")
	c.Check(conn.args[4], HasLen, 0)

	c.Check(conn.queries[5], Equals, "taqsim")
	c.Check(conn.args[5], HasLen, 0)

	c.Check(conn.queries[6], Equals, "insert into `jamila`.schemas (name, prefix, hostname) values (?, ?, ?)")
	c.Assert(conn.args[6], HasLen, 3)
	c.Check(conn.args[6][0], Equals, "reporting")
	c.Check(conn.args[6][1], Equals, "v_test")
	c.Check(conn.args[6][2], Equals, hostname)

	c.Check(conn.queries[7], Equals, "create schema `v_test_reporting`")
	c.Check(conn.args[7], HasLen, 0)

	c.Check(conn.queries[8], Equals, "samiha")
	c.Check(conn.args[8], HasLen, 0)
}

func (s *MySuite) Test_fixr_drop(c *C) {
	conf := s.mock_fixrConf(c)
	def, _ := conf.load()

	c.Assert(def, NotNil)

	// hostname, _ := os.Hostname()
	conn := &mockDb{}
	fixr := &fixr{
		conn:       conn,
		def:        def,
		prefix:     "v_test",
		schemaName: "jamila",
	}

	err := fixr.drop()
	c.Check(err, IsNil)
	c.Assert(conn.queries, HasLen, 4)
	c.Assert(conn.args, HasLen, 4)

	c.Check(conn.queries[0], Equals, "drop schema `v_test_blog`")
	c.Check(conn.args[0], HasLen, 0)

	c.Check(conn.queries[1], Equals, "update `jamila`.schemas set dropped = now() where name = ? and prefix = ?")
	c.Assert(conn.args[1], HasLen, 2)
	c.Check(conn.args[1][0], Equals, "blog")
	c.Check(conn.args[1][1], Equals, "v_test")

	c.Check(conn.queries[2], Equals, "drop schema `v_test_reporting`")
	c.Check(conn.args[2], HasLen, 0)

	c.Check(conn.queries[3], Equals, "update `jamila`.schemas set dropped = now() where name = ? and prefix = ?")
	c.Assert(conn.args[3], HasLen, 2)
	c.Check(conn.args[3][0], Equals, "reporting")
	c.Check(conn.args[3][1], Equals, "v_test")
}

func (s *MySuite) Test_fixr_insert(c *C) {
	conf := s.mock_fixrConf(c)
	def, _ := conf.load()

	c.Assert(def, NotNil)

	// hostname, _ := os.Hostname()
	conn := &mockDb{}
	fixr := &fixr{
		conn:       conn,
		def:        def,
		prefix:     "v_test",
		schemaName: "jamila",
	}

	err := fixr.insert()
	c.Check(err, IsNil)
	c.Assert(conn.queries, HasLen, 5)
	c.Assert(conn.args, HasLen, 5)

	c.Check(conn.queries[0], Equals, "insert into `v_test_blog`.`users` (`id`,`joined`,`username`) VALUES (?,?,?),(?,?,?)")
	c.Assert(conn.args[0], HasLen, 6)
	c.Check(conn.args[0][0], Equals, "1")
	c.Check(conn.args[0][1], Equals, "2015-05-05")
	c.Check(conn.args[0][2], Equals, "maya")
	c.Check(conn.args[0][3], Equals, "2")
	c.Check(conn.args[0][4], Equals, nil)
	c.Check(conn.args[0][5], Equals, nil)

	c.Check(conn.queries[1], Equals, "insert into `v_test_blog`.`articles` (`article-title`,`id`,`posted`) VALUES (?,?,?)")
	c.Assert(conn.args[1], HasLen, 3)
	c.Check(conn.args[1][0], Equals, "suzyQ")
	c.Check(conn.args[1][1], Equals, "1")
	c.Check(conn.args[1][2], Equals, nil)

	c.Check(conn.queries[2], Equals, "insert into `v_test_blog`.`comments` (`comment`,`id`,`posted`) VALUES (?,?,now())")
	c.Assert(conn.args[2], HasLen, 2)
	c.Check(conn.args[2][0], Equals, "cool!")
	c.Check(conn.args[2][1], Equals, "1")

	c.Check(conn.queries[3], Equals, "insert into `v_test_blog`.`comments` (`comment`,`id`,`posted`) VALUES (?,?,?)")
	c.Assert(conn.args[3], HasLen, 3)
	c.Check(conn.args[3][0], Equals, "now()")
	c.Check(conn.args[3][1], Equals, "2")
	c.Check(conn.args[3][2], Equals, "2015-03-15")

	c.Check(conn.queries[4], Equals, "insert into `v_test_reporting`.`reports` (`id`,`report`) VALUES (?,?)")
	c.Assert(conn.args[4], HasLen, 2)
	c.Check(conn.args[4][0], Equals, "1")
	c.Check(conn.args[4][1], Equals, "now()")
}
