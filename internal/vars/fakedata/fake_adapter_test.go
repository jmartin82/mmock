package fakedata

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
	"time"
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

func TestIntMinMax(t *testing.T) {
	faker := FakeAdapter{}
	rand.Seed(time.Now().Unix())

	min := 0
	max := 10000
	for i := min; i < max; i++ {
		params := []int{min, max}
		result := faker.IntMinMax(params...)

		randomInt, err := strconv.Atoi(result)
		if err != nil {
			t.Error("The result should be an integer in string format", result)
		}
		if randomInt < min || randomInt > max {
			errorMsg := fmt.Sprintf("The random number should be between %d and %d", min, max)
			t.Error(errorMsg, randomInt)
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

func TestHex(t *testing.T) {
	faker := FakeAdapter{}

	for i := 1; i < 65; i++ {
		r := regexp.MustCompile(fmt.Sprintf("[0-9a-f]{%d}", i))

		result := faker.Hex(i)
		if !r.MatchString(result) {
			t.Error("The generated string is not a valid lower case hexidecimal string", result)
		}
	}
}
