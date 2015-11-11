package fixrupr

import (
	"errors"

	. "gopkg.in/check.v1"
)

func (s *MySuite) Test_newDbError(c *C) {
	err := newDbError(errors.New("this-is-my-error"), "this-is-my-query", []interface{}{"these", "are", "my", "parameters"})
	c.Assert(err, FitsTypeOf, &dbError{})

	dbErr, ok := err.(*dbError)
	c.Assert(dbErr, NotNil)
	c.Assert(ok, Equals, true)
	c.Assert(dbErr.err, NotNil)
	c.Check(dbErr.err.Error(), Equals, "this-is-my-error")
	c.Check(dbErr.query, Equals, "this-is-my-query")
	c.Assert(dbErr.parameters, HasLen, 4)
	c.Check(dbErr.parameters[0], Equals, "these")
	c.Check(dbErr.parameters[1], Equals, "are")
	c.Check(dbErr.parameters[2], Equals, "my")
	c.Check(dbErr.parameters[3], Equals, "parameters")
}

func (s *MySuite) Test_dbError_Error(c *C) {
	err := newDbError(errors.New("this-is-my-error"), "this-is-my-query", []interface{}{"this", "are", "my", "parameters"})
	c.Check(err.Error(), Equals, "this-is-my-error")
}
