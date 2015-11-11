package fixrupr

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
)

type fixrConn interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// creates all the schemas and tables and functions
func (f *fixr) create() (err error) {
	for _, schema := range f.def.schemas {
		err = f.schema(schema.name)
		if err != nil {
			return
		}

		for _, table := range schema.tables {
			err = f.table(schema.name, table)
			if err != nil {
				return
			}
		}

		for _, function := range schema.functions {
			err = f.function(schema.name, function)
			if err != nil {
				return
			}
		}
	}

	return
}

// inserts all the rows
func (f *fixr) insert() (err error) {
	for _, d := range f.def.data {
		err = f.load(f.prefix, d)
		if err != nil {
			return
		}
	}
	return
}

// drops all the schemas
func (f *fixr) drop() (err error) {
	for _, schema := range f.def.schemas {
		query := fmt.Sprintf("drop schema `%s_%s`", f.prefix, schema.name)
		_, e := f.conn.Exec(query)
		if e != nil {
			err = newDbError(e, query, []interface{}{})
		} else {
			query = "update zombie.schemas set dropped = now() where name = ? and prefix = ?"
			_, e := f.conn.Exec(query, schema.name, f.prefix)
			if e != nil {
				err = newDbError(e, query, []interface{}{})
			}
		}
		// don't return right away this time - even if there was and error, want to still clean up the rest
		// log it though - user needs to know what db stuck around
	}

	return
}

// creates a schema
func (f *fixr) schema(name string) (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "<unknown>"
	}
	query := "insert into zombie.schemas (name, prefix, hostname) values (?, ?, ?)"

	_, err = f.conn.Exec(query, name, f.prefix, hostname)
	if err != nil {
		err = newDbError(err, query, []interface{}{name, f.prefix, hostname})
		return
	}

	query = fmt.Sprintf("create schema `%s_%s`", f.prefix, name)

	_, err = f.conn.Exec(query)
	if err != nil {
		err = newDbError(err, query, []interface{}{})
	}
	return
}

// creates a table
func (f *fixr) table(schema string, ddl string) (err error) {
	return f.exec(schema, ddl)
}

// creates a function
func (f *fixr) function(schema string, ddl string) error {
	return f.exec(schema, ddl)
}

// executes ddl
func (f *fixr) exec(schema string, ddl string) (err error) {
	query := strings.Replace(string(ddl), "{{schema}}", fmt.Sprintf("%s_%s", f.prefix, schema), -1)
	_, err = f.conn.Exec(query)
	if err != nil {
		err = newDbError(err, query, []interface{}{})
	}
	return
}

// inserts a group of rows
func (f *fixr) load(prefix string, data fixrDataDef) (err error) {
	if len(data.rows) == 0 {
		return
	}

	fields := getInsertFields(data.rows)
	rows := []string{}
	parameters := []interface{}{}
	for _, row := range data.rows {
		rowInsert, rowParams := generateInsert(fields, row)
		rows = append(rows, fmt.Sprintf("(%s)", strings.Join(rowInsert, ",")))
		parameters = append(parameters, rowParams...)
	}

	query := fmt.Sprintf(
		"insert into `%s_%s`.`%s` (%s) VALUES %s",
		prefix,
		data.schema,
		data.table,
		strings.Join(fields, ","),
		strings.Join(rows, ","),
	)

	_, err = f.conn.Exec(query, parameters...)
	if err != nil {
		err = newDbError(err, query, parameters)
	}
	return
}

// gets all the fields to use in an insert statement based on the row data
func getInsertFields(rows []map[string]fixrCellDef) []string {
	var (
		fieldMap = map[string]bool{}
		fields   = []string{}
	)

	for _, row := range rows {
		for field, def := range row {
			if def.column == "" {
				fieldMap[field] = true
			} else {
				fieldMap[def.column] = true
			}
		}
	}

	for field := range fieldMap {
		fields = append(fields, fmt.Sprintf("`%s`", field))
	}

	// this sorting step is for testability only
	// otherwise, the order of the fields is non-deterministic
	// the insert will still work, but it's more difficult to test
	sort.Strings(fields)
	return fields
}

// generates the insert values and params for a single row
func generateInsert(fields []string, row map[string]fixrCellDef) (values []string, params []interface{}) {
	params = []interface{}{}
	values = []string{}

	columns := map[string]fixrCellDef{}
	for field, cellDef := range row {
		if cellDef.column == "" {
			columns[fmt.Sprintf("`%s`", field)] = cellDef
		} else {
			columns[fmt.Sprintf("`%s`", cellDef.column)] = cellDef
		}
	}

	for _, field := range fields {
		if cellDef, ok := columns[field]; ok {
			if !cellDef.notNil {
				values = append(values, "?")
				params = append(params, nil)
			} else if cellDef.isParameter {
				values = append(values, "?")
				params = append(params, cellDef.value)
			} else {
				values = append(values, cellDef.value)
			}
		} else {
			// treat as nil
			values = append(values, "?")
			params = append(params, nil)
		}
	}

	return
}
