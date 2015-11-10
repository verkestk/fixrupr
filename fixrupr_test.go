// fixrupr provides a utility for setting up and tearing down mysql databases. It can be used to
// create a seed database or for test fixtures.
package fixrupr

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {

}

func (s *MySuite) Test_New(c *C) {
	// conn := &mockDb{}

	// dir := c.MkDir()
}

func (s *MySuite) Test_fixr_SetUp(c *C) {

}

func (s *MySuite) Test_fixr_TearDown(c *C) {

}
