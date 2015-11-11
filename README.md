# fixrupr&nbsp;[![Build Status](https://travis-ci.org/verkestk/fixrupr.svg?branch=master)](https://travis-ci.org/verkestk/fixrupr)&nbsp;[![godoc reference](https://godoc.org/github.com/verkestk/fixrupr?status.png)](https://godoc.org/github.com/verkestk/fixrupr)

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

#### Setting Up Your Database

This package will be creating and destroying schemas, tables, and functions. It will also be inserting. All schemas created will be prefixed with "z_". Make sure the user your code will connect with has permissions to do so. We recommend full permissions on

- ```z_%``` schemas: full permissions (minus grant)
- and these administrative roles:
  - create
  - create routine
  - drop

###### (Optional) Tracking Schema Set-Ups/Tear-Downs

fixrupr gives you the option of tracking schema set-ups and tear-downs in a database table. If you don't want to use this feature, use an empty string for the ```schemaName``` parameter to the ```New``` function.

```
f, err := fixrupr.New(conn, './test-data', "")
```

If you do want to use this feature, you'll need to create this table in the desired schema:

```
CREATE TABLE `schemas` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `prefix` varchar(32) NOT NULL,
  `hostname` varchar(64) NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `dropped` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index2` (`name`,`prefix`)
)
```

The user your code connects will will need to have insert/update privileges on this schema.

#### Go Code


Once you've got all your configuration and database ready, you are ready to use the fixrupr go package:

```
import "github.com/verkest/fixrupr"

// ...

conn, _ := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?timeout=30s", dbuser, dbpass, dbserver, dbschema))
f, _ := fixrupr.New(conn, './test-data', 'test_schemas')
f.SetUp()

// ... run tests

f.TearDown()
```

#### Keeping Your DB Code Testable

This package is designed to support concurrent creations of the same configured fixures. In order to do that, the schemas created are prefixed uniquely (based on the hostname of the client and the unix time). That means that when your code connects to a database and makes queries, it cannot hardcode schema names.

You could write your queries like this:

```
query := fmt.Sprintf("select * from %s.users", schemaName)
```

Or, you could use fixrupr's prefix service:

```
import "github.com/verkest/fixrupr/prefixr"

// ...

pf := &prefixr.Prefixr{PrefixString: "my-prefix"}
query := pf.Prefix("SELECT * FROM {{pf:blog}} JOIN {{pf:reporing}}")

// ...
```

The service is pretty forgiving and supports prefixes that required surrounding with back ticks and also handles when backticks are already in the query.  The following all produce the same result:

- ``` {{pf:blog}}```
- ``` {{pf:`blog`}}```
- ``` `{{pf:blog}}` ```

If the prefix is "my-prefix" those all resolve to ``` `my-prefix_blog` ```
