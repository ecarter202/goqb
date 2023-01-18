package goqb

import (
	"fmt"
	"strings"
)

type (
	Rule struct {
		Field       string      `json:"field,omitempty"`
		Value       interface{} `json:"value,omitempty"`
		Operator    string      `json:"operator,omitempty"`
		Rules       []*Rule     `json:"rules,omitempty"`
		Combinator  string      `json:"combinator,omitempty"`
		ValueSource string      `json:"valueSource,omitempty"`
	}
)

func (r *Rule) String() (s string) {
	switch r.Operator {
	case "=":
		s = r.Equals()
	case "!=":
		s = r.NotEquals()
	case "<":
		s = r.LessThan()
	case ">":
		s = r.GreaterThan()
	case "<=":
		s = r.LessThanEqual()
	case ">=":
		s = r.GreaterThanEqual()
	case "contains":
		s = r.Contains()
	case "beginsWith":
		s = r.BeginsWith()
	case "endsWith":
		s = r.EndsWith()
	case "doesNotContain":
		s = r.DoesNotContain()
	case "doesNotBeginWith":
		s = r.DoesNotBeginWith()
	case "doesNotEndWith":
		s = r.DoesNotEndWith()
	case "null":
		s = fmt.Sprintf("%s IS NULL", r.Field)
	case "notNull":
		s = fmt.Sprintf("%s IS NOT NULL", r.Field)
	case "in":
		s = r.In()
	case "notIn":
		s = r.NotIn()
	case "between":
		s = r.Between()
	case "notBetween":
		s = r.NotBetween()
	}

	return s
}

func (r *Rule) Equals() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s = %s", r.Field, r.Value)
	case "value", "":
		switch r.Value.(type) {
		case string:
			s = fmt.Sprintf("%s = '%s'", r.Field, r.Value)
		default:
			s = fmt.Sprintf("%s = %v", r.Field, r.Value)
		}
	}

	return s
}

func (r *Rule) NotEquals() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s != %s", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s != '%s'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) LessThan() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s < %s", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s < '%s'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) GreaterThan() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s > %s", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s > '%s'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) LessThanEqual() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s <= %s", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s <= '%s'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) GreaterThanEqual() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s >= %s", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s >= '%s'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) Contains() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s LIKE '%%' || %s || '%%'", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s LIKE '%%%s%%'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) BeginsWith() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s LIKE %s || '%%'", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s LIKE '%s%%'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) EndsWith() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s LIKE '%%' || %s", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s LIKE '%%%s'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) DoesNotContain() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s NOT LIKE '%%' || %s || '%%'", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s NOT LIKE '%%%s%%'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) DoesNotBeginWith() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s NOT LIKE %s || '%%'", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s NOT LIKE %s'%%'", r.Field, r.Value)
	}

	return s
}

func (r *Rule) DoesNotEndWith() (s string) {
	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s NOT LIKE '%%' || %s", r.Field, r.Value)
	case "value", "":
		s = fmt.Sprintf("%s NOT LIKE '%%'%s", r.Field, r.Value)
	}

	return s
}

func (r *Rule) In() (s string) {
	var x []string

	fmt.Printf("TYPE: %T\n", r.Value)

	switch r.Value.(type) {
	case string:
		x = strings.Split(r.Value.(string), ",")
	case []string, []int, []float64, []bool, []interface{}:
		x = stringSlice(r.Value)
	}

	switch r.ValueSource {
	case "field":
		for i, v := range x {
			x[i] = fmt.Sprintf("%s", strings.TrimSpace(v))
		}
	case "value", "":
		for i, v := range x {
			x[i] = fmt.Sprintf("'%s'", strings.TrimSpace(v))
		}
	}

	return fmt.Sprintf("%s IN (%s)", r.Field, strings.Join(x, ","))
}

func (r *Rule) NotIn() (s string) {
	x := strings.Split(r.Value.(string), ",")

	switch r.ValueSource {
	case "field":
		for i, v := range x {
			x[i] = fmt.Sprintf("%s", strings.TrimSpace(v))
		}
	case "value", "":
		for i, v := range x {
			x[i] = fmt.Sprintf("'%s'", strings.TrimSpace(v))
		}
	}

	return fmt.Sprintf("%s NOT IN (%s)", r.Field, strings.Join(x, ","))
}

func (r *Rule) Between() (s string) {
	x := strings.Split(r.Value.(string), ",")
	if len(x) != 2 {
		return ""
	}

	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s BETWEEN %s AND %s", r.Field, strings.TrimSpace(x[0]), strings.TrimSpace(x[1]))
	case "value", "":
		s = fmt.Sprintf("%s BETWEEN '%s' AND '%s'", r.Field, strings.TrimSpace(x[0]), strings.TrimSpace(x[1]))
	}

	return s
}

func (r *Rule) NotBetween() (s string) {
	x := strings.Split(r.Value.(string), ",")
	if len(x) != 2 {
		return ""
	}

	switch r.ValueSource {
	case "field":
		s = fmt.Sprintf("%s NOT BETWEEN %s AND %s", r.Field, strings.TrimSpace(x[0]), strings.TrimSpace(x[1]))
	case "value", "":
		s = fmt.Sprintf("%s NOT BETWEEN '%s' AND '%s'", r.Field, strings.TrimSpace(x[0]), strings.TrimSpace(x[1]))
	}

	return s
}
