package vars

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
)

const Name = "___MMOCK__TEST__ENV__VAR"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	randString := make([]rune, n)
	for i := range randString {
		randString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randString)
}

func TestEnvVar(t *testing.T) {
	var value = randSeq(18)
	os.Setenv(Name, value)

	var tag = fmt.Sprintf("env(%s)", Name)
	holders := []string{tag}
	env := EnvVars{}

	result := env.Fill(holders)
	foundValue, found := result[tag]

	os.Unsetenv(Name)

	if !found {
		t.Errorf("tag %s was not found", tag)
	}

	if !strings.Contains(foundValue[0], value) {
		t.Errorf("EnvVar %s Expected: %s, Actual: %s", Name, value, foundValue[0])
	}

	var varNotPresent = "__NOT_PRESENT___"

	tag = fmt.Sprintf("env(%s)", varNotPresent)
	holders = []string{tag}
	env = EnvVars{}

	result = env.Fill(holders)
	foundValue, found = result[tag]

	os.Unsetenv(varNotPresent)

	if found && foundValue[0] != "" {
		t.Errorf("tag %s was found with value %v", tag, foundValue[0])
	}
}
