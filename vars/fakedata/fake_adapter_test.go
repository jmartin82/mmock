package fakedata

import (
	"regexp"
	"strconv"
	"testing"
)

func TestInt(t *testing.T) {
	faker := FakeAdapter{}

	for i := 0; i < 10000; i++ {
		result := faker.Int(1000)
		randomInt, err := strconv.Atoi(result)
		if err != nil {
			t.Error("The result should be an integer in string format", result)
		}
		if randomInt < 0 || randomInt > 1000 {
			t.Error("The random number should be between 0 and 1000", randomInt)
		}
	}
}

func TestFloat(t *testing.T) {
	faker := FakeAdapter{}

	for i := 0; i < 10000; i++ {
		result := faker.Float(1000)
		randomFloat, err := strconv.ParseFloat(result, 64)
		if err != nil {
			t.Error("The result should be a float64 in string format", result)
		}
		if randomFloat < 0 || randomFloat > 1000 {
			t.Error("The random number should be between 0 and 1000", randomFloat)
		}
	}
}

func TestUUID(t *testing.T) {
	faker := FakeAdapter{}
	r := regexp.MustCompile(`[0-9A-Fa-f]{32}|[0-9A-Fa-f\\-]{36}`)

	for i := 0; i < 10000; i++ {
		result := faker.UUID()
		if !r.MatchString(result) {
			t.Error("The generated unique id is not a valid UUID", result)
		}
	}
}

func TestBasicFakeVars(t *testing.T) {
	faker := FakeAdapter{}
	if faker.Brand() == "" {
		t.Error("Brand fake doesn't work")
	}
	if faker.Character() == "" && len(faker.Character()) != 1 {
		t.Error("Character fake doesn't work")
	}

	if faker.CharactersN(5) == "" && len(faker.Character()) != 5 {
		t.Error("CharactersN fake doesn't work")
	}

}
