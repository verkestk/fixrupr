package fixrupr

import (
	. "gopkg.in/check.v1"
)

var (
	configJSON = `
{
  "schemas": [{
    "name": "blog",
    "tables": [
      "users",
      "articles",
      "comments"
    ],
    "functions": [
      "copy_article"
    ]
  }, {
    "name": "reporting",
    "tables": [
      "reports"
    ]
  }],
  "data": [
      "blog.users",
      "blog.articles",
      "blog.comments.article1",
      "blog.comments.article2",
      "reporting.reports"
  ]
}
	`

	usersYAML = `
- id: 1
  username: maya
  joined: "2015-05-05"

- id: 2

`

	articlesYAML = `
- id: 1
  title:
    column: article-title
    value: "suzyQ"
  posted: null

`

	comments1YAML = `
- id: 1
  comment: cool!
  posted:
    param: false
    value: "now()"

`

	comments2YAML = `
- id: 2
  comment:
    param: true
    value: "now()"
  posted: "2015-03-15"

`

	reportsYAML = `
- id: 1
  report: "now()"

`
)

func (s *MySuite) Test_fixrConf_load(c *C) {
	conf := s.mock_fixrConf(c)

	c.Assert(conf.Schemas, HasLen, 2)
	c.Check(conf.Schemas[0].Name, Equals, "blog")
	c.Assert(conf.Schemas[0].Tables, HasLen, 3)
	c.Check(conf.Schemas[0].Tables[0], Equals, "users")
	c.Check(conf.Schemas[0].Tables[1], Equals, "articles")
	c.Check(conf.Schemas[0].Tables[2], Equals, "comments")
	c.Assert(conf.Schemas[0].Functions, HasLen, 1)
	c.Check(conf.Schemas[0].Functions[0], Equals, "copy_article")
	c.Check(conf.Schemas[1].Name, Equals, "reporting")
	c.Assert(conf.Schemas[1].Tables, HasLen, 1)
	c.Check(conf.Schemas[1].Tables[0], Equals, "reports")
	c.Check(conf.Schemas[1].Functions, HasLen, 0)
	c.Assert(conf.Data, HasLen, 5)
	c.Check(conf.Data[0], Equals, "blog.users")
	c.Check(conf.Data[1], Equals, "blog.articles")
	c.Check(conf.Data[2], Equals, "blog.comments.article1")
	c.Check(conf.Data[3], Equals, "blog.comments.article2")
	c.Check(conf.Data[4], Equals, "reporting.reports")

	def, err := conf.load()
	c.Check(err, IsNil)
	c.Assert(def, NotNil)

	c.Assert(def.schemas, HasLen, 2)
	c.Check(def.schemas[0].name, Equals, "blog")
	c.Assert(def.schemas[0].tables, HasLen, 3)
	c.Check(def.schemas[0].tables[0], Equals, "choo-choo")
	c.Check(def.schemas[0].tables[1], Equals, "egyptian")
	c.Check(def.schemas[0].tables[2], Equals, "turkish")
	c.Assert(def.schemas[0].functions, HasLen, 1)
	c.Check(def.schemas[0].functions[0], Equals, "taqsim")
	c.Check(def.schemas[1].name, Equals, "reporting")
	c.Assert(def.schemas[1].tables, HasLen, 1)
	c.Check(def.schemas[1].tables[0], Equals, "samiha")
	c.Check(def.schemas[1].functions, HasLen, 0)

	c.Assert(def.data, HasLen, 5)
	c.Check(def.data[0].schema, Equals, "blog")
	c.Check(def.data[0].table, Equals, "users")
	c.Assert(def.data[0].rows, HasLen, 2)
	c.Assert(def.data[0].rows[0], HasLen, 3)
	c.Check(def.data[0].rows[0]["id"].column, Equals, "")
	c.Check(def.data[0].rows[0]["id"].isParameter, Equals, true)
	c.Check(def.data[0].rows[0]["id"].notNil, Equals, true)
	c.Check(def.data[0].rows[0]["id"].value, Equals, "1")
	c.Check(def.data[0].rows[0]["username"].column, Equals, "")
	c.Check(def.data[0].rows[0]["username"].isParameter, Equals, true)
	c.Check(def.data[0].rows[0]["username"].notNil, Equals, true)
	c.Check(def.data[0].rows[0]["username"].value, Equals, "maya")
	c.Check(def.data[0].rows[0]["joined"].column, Equals, "")
	c.Check(def.data[0].rows[0]["joined"].isParameter, Equals, true)
	c.Check(def.data[0].rows[0]["joined"].notNil, Equals, true)
	c.Check(def.data[0].rows[0]["joined"].value, Equals, "2015-05-05")
	c.Check(def.data[0].rows[1], HasLen, 1)
	c.Check(def.data[0].rows[1]["id"].column, Equals, "")
	c.Check(def.data[0].rows[1]["id"].isParameter, Equals, true)
	c.Check(def.data[0].rows[1]["id"].notNil, Equals, true)
	c.Check(def.data[0].rows[1]["id"].value, Equals, "2")

	c.Check(def.data[1].schema, Equals, "blog")
	c.Check(def.data[1].table, Equals, "articles")
	c.Assert(def.data[1].rows, HasLen, 1)
	c.Assert(def.data[1].rows[0], HasLen, 3)
	c.Assert(def.data[1].rows[0]["id"].column, Equals, "")
	c.Assert(def.data[1].rows[0]["id"].isParameter, Equals, true)
	c.Assert(def.data[1].rows[0]["id"].notNil, Equals, true)
	c.Assert(def.data[1].rows[0]["id"].value, Equals, "1")
	c.Assert(def.data[1].rows[0]["title"].column, Equals, "article-title")
	c.Assert(def.data[1].rows[0]["title"].isParameter, Equals, true)
	c.Assert(def.data[1].rows[0]["title"].notNil, Equals, true)
	c.Assert(def.data[1].rows[0]["title"].value, Equals, "suzyQ")
	c.Assert(def.data[1].rows[0]["posted"].column, Equals, "")
	c.Assert(def.data[1].rows[0]["posted"].isParameter, Equals, false)
	c.Assert(def.data[1].rows[0]["posted"].notNil, Equals, false)
	c.Assert(def.data[1].rows[0]["posted"].value, Equals, "")

	c.Check(def.data[2].schema, Equals, "blog")
	c.Check(def.data[2].table, Equals, "comments")
	c.Assert(def.data[2].rows, HasLen, 1)
	c.Assert(def.data[2].rows[0], HasLen, 3)
	c.Assert(def.data[2].rows[0]["id"].column, Equals, "")
	c.Assert(def.data[2].rows[0]["id"].isParameter, Equals, true)
	c.Assert(def.data[2].rows[0]["id"].notNil, Equals, true)
	c.Assert(def.data[2].rows[0]["id"].value, Equals, "1")
	c.Assert(def.data[2].rows[0]["comment"].column, Equals, "")
	c.Assert(def.data[2].rows[0]["comment"].isParameter, Equals, true)
	c.Assert(def.data[2].rows[0]["comment"].notNil, Equals, true)
	c.Assert(def.data[2].rows[0]["comment"].value, Equals, "cool!")
	c.Assert(def.data[2].rows[0]["posted"].column, Equals, "")
	c.Assert(def.data[2].rows[0]["posted"].isParameter, Equals, false)
	c.Assert(def.data[2].rows[0]["posted"].notNil, Equals, true)
	c.Assert(def.data[2].rows[0]["posted"].value, Equals, "now()")

	c.Check(def.data[3].schema, Equals, "blog")
	c.Check(def.data[3].table, Equals, "comments")
	c.Assert(def.data[3].rows, HasLen, 1)
	c.Assert(def.data[3].rows[0], HasLen, 3)
	c.Assert(def.data[3].rows[0]["id"].column, Equals, "")
	c.Assert(def.data[3].rows[0]["id"].isParameter, Equals, true)
	c.Assert(def.data[3].rows[0]["id"].notNil, Equals, true)
	c.Assert(def.data[3].rows[0]["id"].value, Equals, "2")
	c.Assert(def.data[3].rows[0]["comment"].column, Equals, "")
	c.Assert(def.data[3].rows[0]["comment"].isParameter, Equals, true)
	c.Assert(def.data[3].rows[0]["comment"].notNil, Equals, true)
	c.Assert(def.data[3].rows[0]["comment"].value, Equals, "now()")
	c.Assert(def.data[3].rows[0]["posted"].column, Equals, "")
	c.Assert(def.data[3].rows[0]["posted"].isParameter, Equals, true)
	c.Assert(def.data[3].rows[0]["posted"].notNil, Equals, true)
	c.Assert(def.data[3].rows[0]["posted"].value, Equals, "2015-03-15")

	c.Check(def.data[4].schema, Equals, "reporting")
	c.Check(def.data[4].table, Equals, "reports")
	c.Assert(def.data[4].rows, HasLen, 1)
	c.Assert(def.data[4].rows[0], HasLen, 2)
	c.Assert(def.data[4].rows[0]["id"].column, Equals, "")
	c.Assert(def.data[4].rows[0]["id"].isParameter, Equals, true)
	c.Assert(def.data[4].rows[0]["id"].notNil, Equals, true)
	c.Assert(def.data[4].rows[0]["id"].value, Equals, "1")
	c.Assert(def.data[4].rows[0]["report"].column, Equals, "")
	c.Assert(def.data[4].rows[0]["report"].isParameter, Equals, true)
	c.Assert(def.data[4].rows[0]["report"].notNil, Equals, true)
	c.Assert(def.data[4].rows[0]["report"].value, Equals, "now()")
}
