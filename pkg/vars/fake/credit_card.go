package fake

import (
	"math/rand"
	"time"

	luhn "github.com/joeljunstrom/go-luhn"
)

type creditCard struct {
	length   int
	prefixes []string
}

var creditCards = map[string]creditCard{
	"visa":         {16, []string{"4"}},
	"visaelectron": {16, []string{"4539", "4556", "4916", "4532", "4929", "40240071", "4485", "4716"}},
	"mastercard":   {16, []string{"51", "52", "53", "54", "55"}},
	"amex":         {15, []string{"34", "37"}},
	"discover":     {16, []string{"6011"}},
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type CreditCardGenerator struct {
}

func NewCreditCardGenerator() *CreditCardGenerator {
	return &CreditCardGenerator{}
}

func (c *CreditCardGenerator) getRndPrefix(prefixes []string) string {
	return prefixes[r.Intn(len(prefixes))]
}

func (c *CreditCardGenerator) getNumber(name string) string {
	cardInfo := creditCards[name]
	prefix := c.getRndPrefix(cardInfo.prefixes)
	return luhn.GenerateWithPrefix(cardInfo.length, prefix)
}

func (c *CreditCardGenerator) CreditCardVisa() string {
	return c.getNumber("visa")
}
func (c *CreditCardGenerator) CreditCardVisaElectron() string {
	return c.getNumber("visaelectron")
}

func (c *CreditCardGenerator) CreditCardMasterCard() string {
	return c.getNumber("mastercard")
}

func (c *CreditCardGenerator) CreditCardAmericanExpress() string {
	return c.getNumber("amex")
}

func (c *CreditCardGenerator) CreditCardDiscover() string {
	return c.getNumber("discover")
}
