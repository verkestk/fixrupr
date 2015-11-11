// Package fixrupr provides a utility for setting up and tearing down mysql databases. It can be used to
// create a seed database or for test fixtures.
package fixrupr

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"time"
)

// // Fixr is the interface consuming packages use to interact with this package.
// type Fixr interface {
// 	SetUp() error
// 	TearDown() error
// }

type Fixr struct {
	conn       fixrConn
	def        *fixrDef
	prefix     string
	schemaName string
}

// New gets a new Fixr instance
// conn: db connection
// configPath: path to the directory containing the config file and the schema/data directories
// schemaName (optional): schema to track set-ups/tear-downs
func New(conn *sql.DB, configPath string, schemaName string) (f *Fixr, err error) {
	// parse config file
	var (
		conf = &fixrConf{}
		def  = &fixrDef{}
	)

	// get the file contents & parse
	conf, err = loadConfig(fmt.Sprintf("%s/test.config.json", configPath))
	if err != nil {
		return
	}
	conf.path = configPath

	// validate and load the config data
	def, err = conf.load()
	if err != nil {
		return
	}

	f = &Fixr{
		conn:       conn,
		def:        def,
		prefix:     getPrefix(),
		schemaName: schemaName,
	}

	return
}

// SetUp sets up the database(s) - creates schemas, tables, and functions and
// inserts rows.
func (f *Fixr) SetUp() (err error) {
	// create schema
	err = f.create()
	if err != nil {
		return
	}

	// insert rows
	err = f.insert()
	return
}

// TearDown tears down the database(s) - drops the databases created in SetUp.
func (f *Fixr) TearDown() (err error) {
	// drop schema
	err = f.drop()
	return
}

func (f *Fixr) GetPrefix() string {
	return f.prefix
}

func getPrefix() string {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	// make the hostname safe for unquoted mysql identifier
	re := regexp.MustCompile("[^0-9a-zA-Z$_]")
	host = re.ReplaceAllString(host, "_")

	if len(host) > 19 {
		host = host[0:19]
	}

	return fmt.Sprintf("z_%s_%d", host, time.Now().Unix())
}
