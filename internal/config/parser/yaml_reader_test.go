package parser

import (
	"testing"
	"time"
)

func TestYamlCanParse(t *testing.T) {
	yaml := YAMLReader{}

	var extTest = []struct {
		n        string
		expected bool
	}{
		{"test.yaml", true},
		{"test.YAML", true},
		{"test.yml", true},
		{"test", false},
		{"test.json.txt", false},
		{"test.3234.yaml", true},
	}

	for _, p := range extTest {
		actual := yaml.CanParse(p.n)
		if actual != p.expected {
			t.Errorf("With value %s expected '%v' actual '%v'", p.n, p.expected, actual)
		}
	}

}

func TestYamlRead(t *testing.T) {
	validDefinition := []byte(`URI: name
description: description
request: 
 method: GET
 path: "/your/path/:variable"
 queryStringParameters: 
  name: 
   - value
 headers: 
  name: 
   - value
 cookies: 
  name: value
 body: "Expected Body"
response: 
 statusCode: 200
 headers: 
  name: 
   - value
 cookies: 
  name: value
 body: Responsebody
control: 
 scenario: 
  name: "string (scenario name)"
  requiredState: 
   - "not_started (default state)"
   - another_state_name
  newState: new_stat_neme
 proxyBaseURL: "http://www.jordi.io"
 delay: 5s
 crazy: true
 priority: 1`)
	invalidDefinition := []byte("sfsdf")

	yaml := YAMLReader{}
	m, err := yaml.Parse(invalidDefinition)
	if err == nil {
		t.Error("Expected error in config")
	}

	m, err = yaml.Parse(validDefinition)
	if err != nil {
		t.Errorf("Unexpected error in config: %s", err)
	}

	if m.URI != "name" {
		t.Error("Missing name")
	}

	if m.Description != "description" {
		t.Error("Missing description")
	}

	//request
	if m.Request.Method != "GET" {
		t.Error("Missing description")
	}

	if m.Request.Path != "/your/path/:variable" {
		t.Error("Missing description")
	}

	if m, f := m.Request.QueryStringParameters["name"]; f == false || m[0] != "value" {
		t.Error("Missing QueryStringParameters")
	}

	if m, f := m.Request.Headers["name"]; f == false || m[0] != "value" {
		t.Error("Missing Headers")
	}

	if m, f := m.Request.Cookies["name"]; f == false || m != "value" {
		t.Error("Missing Cookies")
	}

	if m.Request.Body != "Expected Body" {
		t.Error("Missing Body")
	}

	//response
	if m.Response.StatusCode != 200 {
		t.Error("statusCode")
	}
	if m, f := m.Response.Headers["name"]; f == false || m[0] != "value" {
		t.Error("Missing Headers")
	}
	if m, f := m.Response.Cookies["name"]; f == false || m != "value" {
		t.Error("Missing Cookies")
	}
	if m.Response.Body != "Responsebody" {
		t.Error("Missing Body")
	}

	//control

	if m.Control.ProxyBaseURL != "http://www.jordi.io" {
		t.Error("Missing ProxyBaseURL")
	}

	if m.Control.Delay.Duration != 5*time.Second {
		t.Errorf("Missing delay")
	}

	if m.Control.Crazy != true {
		t.Error("Missing crazy")
	}

	if m.Control.Priority != 1 {
		t.Error("Missing Priority")
	}

	if m.Control.Scenario.Name != "string (scenario name)" {
		t.Error("Missing scenario name")
	}

	if m.Control.Scenario.RequiredState[1] != "another_state_name" {
		t.Error("Missing scenario RequiredState")
	}

	if m.Control.Scenario.NewState != "new_stat_neme" {
		t.Error("Missing scenario NewState")
	}

}
