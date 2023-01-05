package goqb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type (
	Filters struct {
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

// Usable will return whether or not the filters obj can be used as a query string.
// Use this instead of checking String() == "" because that will take time to construct the string.
func (filters *Filters) Usable() bool {
	if filters == nil || len(filters.Rules) < 1 {
		return false
	}

	return true
}

// String returns the filters object parsed as a query string.
func (filters *Filters) String() (s string) {
	if !filters.Usable() {
		return ""
	}

	x := recurRules(filters.Rules, filters.Combinator)

	s = strings.Join(x, " ")

	if filters.Not {
		s = "!" + s
	}

	return s
}

func recurRules(rules []*Rule, combinator string) (x []string) {
	if rules == nil {
		return
	}

	x = []string{"("}
	for i, rule := range rules {
		x = append(x, rule.String())
		if rule.Rules != nil {
			x = append(x, recurRules(rule.Rules, rule.Combinator)...)
		}
		if i < len(rules)-1 {
			x = append(x, combinator)
		}
	}
	x = append(x, ")")

	return
}
