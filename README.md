# fixrupr&nbsp;[![Build Status](https://travis-ci.org/verkestk/fixrupr.svg?branch=master)](https://travis-ci.org/verkestk/fixupr)&nbsp;[![godoc reference](https://godoc.org/github.com/verkestk/fixrupr?status.png)][godoc]

Fixrupr is a golang mysql database seeder. Its primary use case is for setting up and tearing down test fixture data.

## Usage

The go code you'll write is very simple - most of what you need to know is how to organize your fixture data.

#### Config File

First things's first: the config file. It will look something like this:

**file** ```test.config.json```

```
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
    "blog.comments.artcile2"
  ]
}

```

This file describes 2 schemas - ```blog``` and ```reporting```.  The ```blog``` schema has 3 tables and 1 function. The ```reporting``` schema has a single table.

After the schema definition is the data loading definition. The names in this list are important. If you split the name by the ```.``` character, the first part is the schema name, and the second part is the table name. Always. So, ```blog.users``` being the first in the list means that the rows for the users table in the blog schema are inserted first.

#### Directory Structure

In the above example, the files in the ```tables```, ```functions```, and ```data``` directories map to files in the following directory structure:

for example

- 📁 **data**
  - 📄 **blog.users.yml** _(rows for blog.users table)_
  - 📄 **blog.articles.yml** _(rows for blog.articles table)_
  - 📄 **blog.comments.article1.yml** _(rows for blog.comments table)_
  - 📄 **blog.comments.article2.yml** _(more rows for blog.comments table)_
  - 📄 **reporting.reports.yml** _(rows for reporting.reports table)_
- 📁 **schema**
  - 📁 **blog**
    - 📁 **functions**
      - 📄 **copy_article.sql** _(ddl for creating copy_article function)_
    - 📁 **tables**
      - 📄 **articles.sql** _(ddl for creating articles table)_
      - 📄 **comments.sql** _(ddl for creating comments table)_
      - 📄 **users.sql** _(ddl for creating users table)_
  - 📁 **reporting**
    - 📁 **tables**
      - 📄 **reports.sql** _(ddl for creating reports table)_
- 📄 **test.config.json** _(the config file)_

#### Data Files

Above there are yaml files containing row data to insert into the tables. Here's what those look like:

**file** ```blog.users.yml```

```
# in most cases, you can just defined a simple list of objects
# whose properties map to table columns
- id: 1
  username: babyBuggy
  password: 35dcddd1057c32cce2b5ac5b5060101c
  joined: "2015-01-05"

# sometimes column names don't make good object properties. in
# these cases, you can specify the column name like this and use
# a different object property. in this example, the value of
# "join_date" will be inserted into the "joined" column.
- id: 2
  username: stinkBug
  password: 40cdad7da4b1f2373e26499acf00bf7e
  join_date:
    column: joined
    value: "2015-02-10"

# by default all values are inserted as parameters ("?" in the
# insert statement). you can disable this by setting "param" to
# false, letting you do things like use database functions.
- id: 3
  username: prettyPretzel
  password:
    param: false
    value: "md5(\"puffinPop\")"
```

#### Example


**go code**:

```
conn, _ := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?timeout=30s", dbuser, dbpass, dbserver, dbschema))
f, _ := fixrup.New(conn, './test-data')
f.SetUp()

// ... run tests

f.TearDown()
```