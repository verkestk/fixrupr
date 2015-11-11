// Package prefixr offers a utility for inserting prefixed schema names into queries
package prefixr

import (
	"fmt"
	"regexp"
)

// Prefixr keeps track of a single prefix and can apply to to multiple queries.
// If you plan to apply the same prefix to multiple queries, then use one of these.
type Prefixr struct {
	PrefixString string
}

// Prefix takes a query string and applies the prefix to the beginning of schema names
// The schema names must match a specific pattern. For example, if your original query is
//   SELECT *
//   FROM blog
//   JOIN reporing
//
// Then you will want to write that query as
//   SELECT *
//   FROM {{pf:blog}}
//   JOIN {{pf:reporting}}
//
// If you PrefixString is "my-prefix" then, the prefixed query will be:
//   SELECT *
//   FROM `my-prefix_blog`
//   JOIN `my-prefix_reporing`
func (p *Prefixr) Prefix(query string) string {
	return Prefix(p.PrefixString, query)
}

// Prefix takes a query string and applies the prefix to the beginning of schema names
// The schema names must match a specific pattern. For example, if your original query is
//   SELECT *
//   FROM blog
//   JOIN reporing
//
// Then you will want to write that query as
//   SELECT *
//   FROM {{pf:blog}}
//   JOIN {{pf:reporting}}
//
// If you prefix is "my-prefix" then, the prefixed query will be:
//   SELECT *
//   FROM `my-prefix_blog`
//   JOIN `my-prefix_reporing`
func Prefix(prefix, query string) string {
	// if the prefix is empty, then don't include an underscore
	format := "`%s_%s`"
	if prefix == "" {
		format = "`%s%s`"
	}

	r := regexp.MustCompile("{{pf:`.+?`}}")
	prefixed := r.ReplaceAllStringFunc(query, func(match string) string {
		return fmt.Sprintf(format, prefix, match[6:len(match)-3])
	})

	r = regexp.MustCompile("`{{pf:.+?}}`")

	prefixed = r.ReplaceAllStringFunc(prefixed, func(match string) string {
		return fmt.Sprintf(format, prefix, match[6:len(match)-3])
	})

	r = regexp.MustCompile("{{pf:.+?}}")
	return r.ReplaceAllStringFunc(prefixed, func(match string) string {
		return fmt.Sprintf(format, prefix, match[5:len(match)-2])
	})
}
