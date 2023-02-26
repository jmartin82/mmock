package route

import (
	"fmt"
	"testing"
)

func TestMatching(t *testing.T) {

	route := NewRoute("/user/:name")
	equal(t, route.Match("/user/azer").Pattern, "/user/:name")

	route = NewRoute("/fruits/:fruit")
	equal(t, route.Match("/fruits/watermelon").Pattern, "/fruits/:fruit")

	route = NewRoute("/fruits/:fruit/:page")
	equal(t, route.Match("/fruits/cherry/452").Pattern, "/fruits/:fruit/:page")

	route = NewRoute("/api/urn:note:123")
	equal(t, route.Match("/api/urn:note:123").Pattern, "/api/urn:note:123")

	route = NewRoute("/api/:userid.json")
	equal(t, route.Match("/api/99.json").Pattern, "/api/:userid.json")

	route = NewRoute("/")
	equal(t, route.Match("/").Pattern, "/")
}

func TestParams(t *testing.T) {
	route := NewRoute("/user/:name")
	equal(t, route.Match("/user/azer").Params["name"], "azer")

	route = NewRoute("/fruits/:fruit")
	equal(t, route.Match("/fruits/watermelon").Params["fruit"], "watermelon")

	route = NewRoute("/fruits/:fruit/:page")
	equal(t, route.Match("/fruits/cherry/452").Params["fruit"], "cherry")
	equal(t, route.Match("/fruits/cherry/452").Params["page"], "452")
}

func equal(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("%s and %s aren't equal", a, b))
	}
}
