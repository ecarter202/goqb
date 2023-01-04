package goqb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestBindRequest(t *testing.T) {
	req := new(http.Request)
	req.Body = ioutil.NopCloser(strings.NewReader(testBody))

	filters, err := BindRequest(req)
	if err != nil {
		t.Fatalf("binding request [ERR: %s]", err)
	}

	fmt.Println(filters.String())
}

const testBody = `{
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
