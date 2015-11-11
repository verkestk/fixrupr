// fixrupr provides a utility for setting up and tearing down mysql databases. It can be used to
// create a seed database or for test fixtures.
package fixrupr

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C)     {}
func (s *MySuite) TearDownTest(c *C)  {}
func (s *MySuite) SetUpSuite(c *C)    {}
func (s *MySuite) TearDownSuite(c *C) {}

func (s *MySuite) Test_New(c *C) {
	configPath := s.help_mockFiles(c)
	var conn *sql.DB

	f, err := New(conn, configPath, "jamila")
	c.Assert(f, NotNil)
	c.Assert(err, IsNil)

	fConn, ok := f.conn.(*sql.DB)
	c.Assert(ok, Equals, true)
	c.Check(fConn, Equals, conn)
	c.Check(f.def, NotNil)

	c.Check(f.prefix[0:2], Equals, "z_")
}

func (s *MySuite) Test_fixr_SetUp(c *C) {
	configPath := s.help_mockFiles(c)
	f, err := New(nil, configPath, "jamila")
	c.Assert(f, NotNil)
	c.Assert(err, IsNil)
	f.conn = &mockDb{}

	err = f.SetUp()
	c.Check(err, IsNil)
}

func (s *MySuite) Test_fixr_TearDown(c *C) {
	configPath := s.help_mockFiles(c)
	f, err := New(nil, configPath, "jamila")
	c.Assert(f, NotNil)
	c.Assert(err, IsNil)
	f.conn = &mockDb{}

	f.SetUp()
	err = f.TearDown()
	c.Check(err, IsNil)
}

func (s *MySuite) Test_fixr_GetPrefix(c *C) {
	f := &Fixr{prefix: "this-is-my-prefix"}
	c.Check(f.GetPrefix(), Equals, "this-is-my-prefix")
}

func (s *MySuite) help_mockFiles(c *C) (dir string) {
	// make directories
	dir = c.MkDir()
	os.MkdirAll(fmt.Sprintf("%s/schema/blog/tables", dir), 0755)
	os.MkdirAll(fmt.Sprintf("%s/schema/blog/functions", dir), 0755)
	os.MkdirAll(fmt.Sprintf("%s/schema/reporting/tables", dir), 0755)
	os.MkdirAll(fmt.Sprintf("%s/data", dir), 0755)

	// write files
	configFile := fmt.Sprintf("%s/test.config.json", dir)
	usersTableFile := fmt.Sprintf("%s/schema/blog/tables/users.sql", dir)
	articlesTableFile := fmt.Sprintf("%s/schema/blog/tables/articles.sql", dir)
	commentsTableFile := fmt.Sprintf("%s/schema/blog/tables/comments.sql", dir)
	copyArticleFunctionFile := fmt.Sprintf("%s/schema/blog/functions/copy_article.sql", dir)
	reportsTableFile := fmt.Sprintf("%s/schema/reporting/tables/reports.sql", dir)
	usersDataFile := fmt.Sprintf("%s/data/blog.users.yml", dir)
	articlesDataFile := fmt.Sprintf("%s/data/blog.articles.yml", dir)
	comments1DataFile := fmt.Sprintf("%s/data/blog.comments.article1.yml", dir)
	comments2DataFile := fmt.Sprintf("%s/data/blog.comments.article2.yml", dir)
	reportsDataFile := fmt.Sprintf("%s/data/reporting.reports.yml", dir)

	ioutil.WriteFile(configFile, []byte(configJSON), 0755)
	ioutil.WriteFile(usersTableFile, []byte("choo-choo"), 0755)
	ioutil.WriteFile(articlesTableFile, []byte("egyptian"), 0755)
	ioutil.WriteFile(commentsTableFile, []byte("turkish"), 0755)
	ioutil.WriteFile(copyArticleFunctionFile, []byte("taqsim"), 0755)
	ioutil.WriteFile(reportsTableFile, []byte("samiha"), 0755)
	ioutil.WriteFile(usersDataFile, []byte(usersYAML), 0755)
	ioutil.WriteFile(articlesDataFile, []byte(articlesYAML), 0755)
	ioutil.WriteFile(comments1DataFile, []byte(comments1YAML), 0755)
	ioutil.WriteFile(comments2DataFile, []byte(comments2YAML), 0755)
	ioutil.WriteFile(reportsDataFile, []byte(reportsYAML), 0755)

	return
}

func (s *MySuite) mock_fixrConf(c *C) *fixrConf {

	dir := s.help_mockFiles(c)

	conf, err := loadConfig(fmt.Sprintf("%s/test.config.json", dir))
	if conf != nil {
		conf.path = dir
	}

	c.Check(err, IsNil)
	c.Assert(conf, NotNil)

	return conf
}
