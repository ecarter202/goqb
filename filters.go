package goqb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type (
	Filters struct {
		query         string // if called more than once
		allowedFields map[string]struct{}

		Rules      []*Rule `json:"rules"`
		Combinator string  `json:"combinator"`
		Not        bool    `json:"not"`
	}
)

// BindRequest will bind the request body to a filters struct.
func BindRequest(req *http.Request) (filters *Filters, err error) {
	filters = new(Filters)

	err = json.NewDecoder(req.Body).Decode(&filters)
	if err != nil {
		return nil, fmt.Errorf("binding filters [ERR: %s]", err)
	}

	return filters, nil
}

// Exists will return whether or not the filters obj can be used as a query string.
// Use this instead of checking String() == "" because that will take time to construct the string.
func (filters *Filters) Exists() bool {
	if filters == nil || len(filters.Rules) < 1 {
		return false
	}

	return true
}

// AllowFields will restrict the query from using unsupplied fields.
func (filters *Filters) AllowFields(fields []string) *Filters {
	if filters.allowedFields == nil {
		filters.allowedFields = make(map[string]struct{})
	}

	for _, field := range fields {
		filters.allowedFields[field] = struct{}{}
	}

	return filters
}

// String returns the filters object parsed as a query string.
func (filters *Filters) String() (s string) {
	if !filters.Exists() {
		return ""
	} else if filters.query != "" {
		return filters.query
	}

	x := recurRules(filters.allowedFields, filters.Rules, filters.Combinator)

	s = strings.Join(x, " ")

	if filters.Not {
		s = "!" + s
	}

	filters.query = s

	return s
}

func recurRules(allowedFields map[string]struct{}, rules []*Rule, combinator string) (x []string) {
	if rules == nil {
		return
	}

	allowedrules := allowedRules(allowedFields, rules)

	for i, rule := range allowedrules {
		if len(x) == 0 {
			x = append(x, "(")
		}

		x = append(x, rule.String())
		if rule.Rules != nil {
			x = append(x, recurRules(allowedFields, rule.Rules, rule.Combinator)...)
		}
		if len(allowedrules)-i >= 2 {
			x = append(x, combinator)
		}
	}

	if len(x) > 0 {
		x = append(x, ")")
	}

	return
}

func allowedRules(allowedFields map[string]struct{}, rules []*Rule) (allowed []*Rule) {
	if len(allowedFields) == 0 {
		return rules
	}

	for _, rule := range rules {

		_, allowedField := allowedFields[rule.Field]

		if !allowedField && len(rule.Rules) == 0 && rule.ValueSource != "field" {
			continue
		}

		allowed = append(allowed, rule)
	}

	return allowed
}
