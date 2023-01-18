package goqb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestBindRequest(t *testing.T) {
	req := new(http.Request)
	req.Body = ioutil.NopCloser(strings.NewReader(testBody1))

	filters, err := BindRequest(req)
	if filters == nil || err != nil {
		t.Fatalf("binding request [ERR: %s]", err)
	}
}

func TestString(t *testing.T) {
	var filters *Filters

	err := json.NewDecoder(strings.NewReader(testBody1)).Decode(&filters)
	if err != nil {
		t.Fatalf("unmarshaling body to *Filters [ERR: %s]", err)
	}

	expected := `( firstName LIKE 'Stev%' and lastName IN ('Vai','Vaughan') and age > '28' and  ( isMusician = true or instrument = 'Guitar' or firstName = 'howdy' or  ( job = 'handyman' ) ) and groupedField1 = groupedField4 and birthdate < '1969-06-01' and firstName = 'bob' )`
	got := filters.String()

	if got != expected {
		t.Fatalf("invalid query string constructed\n\nwanted: %q\n\ngot: %q", expected, got)
	}
}

func TestAllowedFields(t *testing.T) {
	var filters *Filters

	err := json.NewDecoder(strings.NewReader(testBody1)).Decode(&filters)
	if err != nil {
		t.Fatalf("unmarshaling body to *Filters [ERR: %s]", err)
	}

	allowedFields := []string{
		"lastName",
		"age",
		"isMusician",
		"birthdate",
		"instrument",
		"job",
	}

	filters = filters.AllowFields(allowedFields)

	expected := `( lastName IN ('Vai','Vaughan') and age > '28' and  ( isMusician = true or instrument = 'Guitar' or  ( job = 'handyman' ) ) and groupedField1 = groupedField4 and birthdate < '1969-06-01' )`
	got := filters.String()

	if got != expected {
		t.Fatalf("invalid query string constructed\n\nwanted: %q\n\ngot: %q", expected, got)
	}
}

func TestAllowedFields2(t *testing.T) {
	var filters *Filters

	err := json.NewDecoder(strings.NewReader(testBody2)).Decode(&filters)
	if err != nil {
		t.Fatalf("unmarshaling body to *Filters [ERR: %s]", err)
	}

	allowedFields := []string{
		"lastName",
		"age",
		"isMusician",
		"birthdate",
		"instrument",
		"job",
	}

	filters = filters.AllowFields(allowedFields)

	expected := `( lastName IN ('Vai','Vaughan') and age > '28' and  ( isMusician = true or  ( job = 'handyman' ) ) and groupedField1 = groupedField4 )`
	got := filters.String()

	if got != expected {
		t.Fatalf("invalid query string constructed\n\nwanted: %q\n\ngot: %q", expected, got)
	}
}

const testBody1 = `{
	"rules": [
	  {
		"field": "firstName",
		"value": "Stev",
		"operator": "beginsWith"
	  },
	  {
		"field": "lastName",
		"value": ["Vai", "Vaughan"],
		"operator": "in"
	  },
	  {
		"field": "age",
		"value": "28",
		"operator": ">"
	  },
	  {
		"rules": [
		  {
			"field": "isMusician",
			"value": true,
			"operator": "="
		  },
		  {
			"field": "instrument",
			"value": "Guitar",
			"operator": "="
		  },
		  {
			"field": "firstName",
			"value": "howdy",
			"operator": "=",
			"valueSource": "value"
		  },
		  {
			"rules": [
			  {
				"field": "job",
				"value": "handyman",
				"operator": "=",
				"valueSource": "value"
			  }
			],
			"combinator": "and",
			"not": false
		  }
		],
		"combinator": "or"
	  },
	  {
		"field": "groupedField1",
		"value": "groupedField4",
		"operator": "=",
		"valueSource": "field"
	  },
	  {
		"field": "birthdate",
		"value": "1969-06-01",
		"operator": "<"
	  },
	  {
		"field": "firstName",
		"value": "bob",
		"operator": "=",
		"valueSource": "value"
	  }
	],
	"combinator": "and",
	"not": false
  }`

const testBody2 = `{
	"rules": [
	  {
		"field": "firstName",
		"value": "Stev",
		"operator": "beginsWith"
	  },
	  {
		"field": "lastName",
		"value": "Vai, Vaughan",
		"operator": "in"
	  },
	  {
		"field": "age",
		"value": "28",
		"operator": ">"
	  },
	  {
		"rules": [
		  {
			"field": "isMusician",
			"value": true,
			"operator": "="
		  },
		  {
			"rules": [
			  {
				"field": "job",
				"value": "handyman",
				"operator": "=",
				"valueSource": "value"
			  }
			],
			"combinator": "and",
			"not": false
		  }
		],
		"combinator": "or"
	  },
	  {
		"field": "groupedField1",
		"value": "groupedField4",
		"operator": "=",
		"valueSource": "field"
	  },
	  {
		"field": "firstName",
		"value": "bob",
		"operator": "=",
		"valueSource": "value"
	  }
	],
	"combinator": "and",
	"not": false
  }`
