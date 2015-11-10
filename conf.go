package fixrupr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type fixrConf struct {
	path    string
	Schemas []struct {
		Name      string   `json:"name"`
		Tables    []string `json:"tables"`
		Functions []string `json:"functions"`
	} `json:"schemas"`
	Data []string `json:"data"`
}

type fixrDef struct {
	schemas []fixrSchemaDef
	data    []fixrDataDef
}

type fixrSchemaDef struct {
	name      string
	tables    []string
	functions []string
}

type fixrDataDef struct {
	schema string
	table  string
	rows   []map[string]fixrCellDef
}

type fixrCellDef struct {
	isParameter bool
	notNil      bool
	value       string
	column      string
}

func (c *fixrConf) load() (def *fixrDef, err error) {
	// make sure the files exist and then load the file content
	def = &fixrDef{}

	var (
		schemaDef   fixrSchemaDef
		tableDef    []byte
		functionDef []byte
		dataDef     fixrDataDef
		rowsDef     []byte
	)

	for _, schema := range c.Schemas {
		schemaDef = fixrSchemaDef{name: schema.Name}
		for _, table := range schema.Tables {
			tableDef, err = ioutil.ReadFile(fmt.Sprintf("%s/schema/%s/tables/%s.sql", c.path, schema.Name, table))
			if err != nil {
				return
			}
			schemaDef.tables = append(schemaDef.tables, string(tableDef))
		}

		for _, function := range schema.Functions {
			functionDef, err = ioutil.ReadFile(fmt.Sprintf("%s/schema/%s/functions/%s.sql", c.path, schema.Name, function))
			if err != nil {
				return
			}
			schemaDef.functions = append(schemaDef.functions, string(functionDef))
		}

		def.schemas = append(def.schemas, schemaDef)
	}

	for _, d := range c.Data {
		rowsDef, err = ioutil.ReadFile(fmt.Sprintf("%s/data/%s.yml", c.path, d))
		if err != nil {
			return
		}

		pieces := strings.Split(d, ".")
		dataDef = fixrDataDef{
			schema: pieces[0],
			table:  pieces[1],
		}

		err = yaml.Unmarshal(rowsDef, &dataDef.rows)
		if err != nil {
			// TODO: wrap in error struct with all the necessary data
			fmt.Println("Error", d)
			return
		}
		def.data = append(def.data, dataDef)
	}

	return
}

func loadConfig(filename string) (*fixrConf, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf := &fixrConf{}

	err = json.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (d *fixrCellDef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	toStr := ""
	toStruct := struct {
		Value       string `yaml:"value"`
		Column      string `yaml:"column"`
		IsParameter *bool  `yaml:"param,omitempty"`
	}{}

	err := unmarshal(&toStr)
	if err == nil {
		d.value = toStr
		d.isParameter = true

	} else {
		err = unmarshal(&toStruct)
		if err != nil {
			return err
		}

		d.column = toStruct.Column
		d.value = toStruct.Value

		if toStruct.IsParameter == nil {
			d.isParameter = true
		} else {
			d.isParameter = *toStruct.IsParameter
		}
	}

	d.notNil = true
	return nil
}
