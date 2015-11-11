package prefixr

import (
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

func (s *MySuite) Test_Prefixr_Prefix(c *C) {
	query := "SELECT * FROM {{pf:blog}}.users JOIN `{{pf:reporing}}`.reports JOIN {{pf:`schemas`}} JOIN other"
	prefixer := &Prefixr{PrefixString: "this-is-my-prefix"}

	prefixed := prefixer.Prefix(query)

	c.Check(prefixed, Equals, "SELECT * FROM `this-is-my-prefix_blog`.users JOIN `this-is-my-prefix_reporing`.reports JOIN `this-is-my-prefix_schemas` JOIN other")

	prefixer.PrefixString = ""
	prefixed = prefixer.Prefix(query)
	c.Check(prefixed, Equals, "SELECT * FROM `blog`.users JOIN `reporing`.reports JOIN `schemas` JOIN other")
}

func (s *MySuite) Test_Prefix(c *C) {
	query := "SELECT * FROM {{pf:blog}}.users JOIN `{{pf:reporing}}`.reports JOIN {{pf:`schemas`}} JOIN other"
	prefixed := Prefix("this-is-my-prefix", query)

	c.Check(prefixed, Equals, "SELECT * FROM `this-is-my-prefix_blog`.users JOIN `this-is-my-prefix_reporing`.reports JOIN `this-is-my-prefix_schemas` JOIN other")

	prefixed = Prefix("", query)
	c.Check(prefixed, Equals, "SELECT * FROM `blog`.users JOIN `reporing`.reports JOIN `schemas` JOIN other")
}
