package parser

import "testing"

func TestJSONCanParse(t *testing.T) {
	json := JSONReader{}

	var extTest = []struct {
		n        string
		expected bool
	}{
		{"test.json", true},
		{"test.JSON", true},
		{"test", false},
		{"test.json.txt", false},
		{"test.3234.json", true},
	}

	for _, p := range extTest {
		actual := json.CanParse(p.n)
		if actual != p.expected {
			t.Errorf("With value %s expected '%v' actual '%v'", p.n, p.expected, actual)
		}
	}

}

func TestJSONRead(t *testing.T) {
	validDefinition := []byte(`{
		"URI": "name",
		"description": "description",
		"request": {
			"method": "GET",
			"path": "/your/path/:variable",
			"queryStringParameters": {
				"name": ["value"]
			},
			"headers": {
				"name": ["value"]
			},
			"cookies": {
				"name": "value"
			},
			"body": "Expected Body"
		},
		"response": {
			"statusCode": 200,
			"headers": {
				"name": ["value"]
			},
			"cookies": {
				"name": "value"
			},
			"body": "Responsebody"
		},
		"control": {
			"scenario": {
				"name": "string (scenario name)",
				"requiredState": [
					"not_started (default state)",
					"another_state_name"
				],
				"newState": "new_stat_neme"
			},
			"proxyBaseURL": "http://www.jordi.io",
			"delay": 5,
			"crazy": true,
			"priority": 1
		}
	}`)
	invalidDefinition := []byte("sfsdf")

	json := JSONReader{}
	m, err := json.Parse(invalidDefinition)
	if err == nil {
		t.Errorf("Expected error in config")
	}

	m, err = json.Parse(validDefinition)
	if err != nil {
		t.Errorf("Unexpected error in config")
	}

	if m.URI != "name" {
		t.Errorf("Missing name")
	}

	if m.Description != "description" {
		t.Errorf("Missing description")
	}

	//request
	if m.Request.Method != "GET" {
		t.Errorf("Missing description")
	}

	if m.Request.Path != "/your/path/:variable" {
		t.Errorf("Missing description")
	}

	if m, f := m.Request.QueryStringParameters["name"]; f == false || m[0] != "value" {
		t.Errorf("Missing QueryStringParameters")
	}

	if m, f := m.Request.Headers["name"]; f == false || m[0] != "value" {
		t.Errorf("Missing Headers")
	}

	if m, f := m.Request.Cookies["name"]; f == false || m != "value" {
		t.Errorf("Missing Cookies")
	}

	if m.Request.Body != "Expected Body" {
		t.Errorf("Missing Body")
	}

	//response
	if m.Response.StatusCode != 200 {
		t.Errorf("statusCode")
	}
	if m, f := m.Response.Headers["name"]; f == false || m[0] != "value" {
		t.Errorf("Missing Headers")
	}
	if m, f := m.Response.Cookies["name"]; f == false || m != "value" {
		t.Errorf("Missing Cookies")
	}
	if m.Response.Body != "Responsebody" {
		t.Errorf("Missing Body")
	}

	//control

	if m.Control.ProxyBaseURL != "http://www.jordi.io" {
		t.Errorf("Missing ProxyBaseURL")
	}

	if m.Control.Delay != 5 {
		t.Errorf("Missing delay")
	}

	if m.Control.Crazy != true {
		t.Errorf("Missing crazy")
	}

	if m.Control.Priority != 1 {
		t.Errorf("Missing Priority")
	}

	if m.Control.Scenario.Name != "string (scenario name)" {
		t.Errorf("Missing scenario name")
	}

	if m.Control.Scenario.RequiredState[1] != "another_state_name" {
		t.Errorf("Missing scenario RequiredState")
	}

	if m.Control.Scenario.NewState != "new_stat_neme" {
		t.Errorf("Missing scenario NewState")
	}

}
