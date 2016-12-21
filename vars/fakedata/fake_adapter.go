package fakedata

import (
	"math/rand"
	"strconv"

	"github.com/icrowley/fake"
	"github.com/twinj/uuid"
)

//FakeAdapter contains all available functions to create random data in the mock response.
type FakeAdapter struct {
}

//Brand returns a random Brand
func (fa FakeAdapter) Brand() string {
	return fake.Brand()
}

//Character returns a random Character
func (fa FakeAdapter) Character() string {
	return fake.Character()
}

//Characters returns from 1 to 5 random Characters
func (fa FakeAdapter) Characters() string {
	return fake.Characters()
}

//CharactersN returns n random Characters
func (fa FakeAdapter) CharactersN(n int) string {
	return fake.CharactersN(n)
}

//City returns a random City
func (fa FakeAdapter) City() string {
	return fake.City()
}

//Color returns a random Color
func (fa FakeAdapter) Color() string {
	return fake.Color()
}

//Company returns a random Company
func (fa FakeAdapter) Company() string {
	return fake.Company()
}

//Continent returns a random Continent
func (fa FakeAdapter) Continent() string {
	return fake.Continent()
}

//Country returns a random Country
func (fa FakeAdapter) Country() string {
	return fake.Country()
}

//CreditCardVisa returns a random CreditCardVisa
func (fa FakeAdapter) CreditCardVisa() string {
	return fake.CreditCardNum("Visa")
}

//CreditCardMasterCard returns a random CreditCardMasterCard
func (fa FakeAdapter) CreditCardMasterCard() string {
	return fake.CreditCardNum("MasterCard")
}

//CreditCardAmericanExpress returns a random CreditCardAmericanExpress
func (fa FakeAdapter) CreditCardAmericanExpress() string {
	return fake.CreditCardNum("American Express")
}

//Currency returns a random Currency
func (fa FakeAdapter) Currency() string {
	return fake.Currency()
}

//CurrencyCode returns a random CurrencyCode
func (fa FakeAdapter) CurrencyCode() string {
	return fake.CurrencyCode()
}

//Digits returns from 1 to 5 random Digits
func (fa FakeAdapter) Digits() string {
	return fake.Digits()
}

//DigitsN returns n random Digits
func (fa FakeAdapter) DigitsN(n int) string {
	return fake.DigitsN(n)
}

//EmailAddress returns a random EmailAddress
func (fa FakeAdapter) EmailAddress() string {
	return fake.EmailAddress()
}

//FirstName returns a random FirstName
func (fa FakeAdapter) FirstName() string {
	return fake.FirstName()
}

//FullName returns a random FullName
func (fa FakeAdapter) FullName() string {
	return fake.FullName()
}

//LastName returns a random LastName
func (fa FakeAdapter) LastName() string {
	return fake.LastName()
}

//Gender returns a random Gender
func (fa FakeAdapter) Gender() string {
	return fake.Gender()
}

//IPv4 returns a random IPv4
func (fa FakeAdapter) IPv4() string {
	return fake.IPv4()
}

//Language returns a random Language
func (fa FakeAdapter) Language() string {
	return fake.Language()
}

//Model returns a random Model
func (fa FakeAdapter) Model() string {
	return fake.Model()
}

//Paragraph returns a random Paragraph
func (fa FakeAdapter) Paragraph() string {
	return fake.Paragraph()
}

//Paragraphs returns from 1 to 5 random Paragraphs
func (fa FakeAdapter) Paragraphs() string {
	return fake.Paragraphs()
}

//ParagraphsN returns n random Paragraphs
func (fa FakeAdapter) ParagraphsN(n int) string {
	return fake.ParagraphsN(n)
}

//Phone returns a random Phone
func (fa FakeAdapter) Phone() string {
	return fake.Phone()
}

//Product returns a random Product
func (fa FakeAdapter) Product() string {
	return fake.Product()
}

//Sentence returns a random sentence
func (fa FakeAdapter) Sentence() string {
	return fake.Sentence()
}

//Sentences returns from 1 to 5 random sentences
func (fa FakeAdapter) Sentences() string {
	return fake.Sentences()
}

//SentencesN returns n random sentences
func (fa FakeAdapter) SentencesN(n int) string {
	return fake.SentencesN(n)
}

//SimplePassword returns a random simple password
func (fa FakeAdapter) SimplePassword() string {
	return fake.SimplePassword()
}

//State returns a random state
func (fa FakeAdapter) State() string {
	return fake.State()
}

//StateAbbrev returns a random state abbrev
func (fa FakeAdapter) StateAbbrev() string {
	return fake.StateAbbrev()
}

//Street returns a random street
func (fa FakeAdapter) Street() string {
	return fake.Street()
}

//StreetAddress returns a random street address
func (fa FakeAdapter) StreetAddress() string {
	return fake.StreetAddress()
}

//UserName returns a random username
func (fa FakeAdapter) UserName() string {
	return fake.UserName()
}

//Day returns a random day
func (fa FakeAdapter) Day() string {
	return strconv.Itoa(fake.Day())
}

//Month returns a random month
func (fa FakeAdapter) Month() string {
	return fake.Month()
}

//Year returns a random year between (1980,2020)
func (fa FakeAdapter) Year() string {
	return strconv.Itoa(fake.Year(1980, 2020))
}

//MonthShort returns a random month (Short Version)
func (fa FakeAdapter) MonthShort() string {
	return fake.MonthShort()
}

//WeekDay returns a random day of week
func (fa FakeAdapter) WeekDay() string {
	return fake.WeekDay()
}

//Word returns a random word
func (fa FakeAdapter) Word() string {
	return fake.Word()
}

//Words returns from 1 to 5 random words
func (fa FakeAdapter) Words() string {
	return fake.Words()
}

//WordsN returns n random words
func (fa FakeAdapter) WordsN(n int) string {
	return fake.WordsN(n)
}

//Zip returns a random zip
func (fa FakeAdapter) Zip() string {
	return fake.Zip()
}

//Number returns a random positive number less than or equal to n
func (fa FakeAdapter) Int(n int) string {
	return strconv.Itoa(rand.Intn(n + 1))
}

//Float returns a random positive floating point number less than n
func (fa FakeAdapter) Float(n int) string {
	f := float64(n)
	value := rand.Float64() * f
	return strconv.FormatFloat(value, 'f', 4, 64)
}

//UUID generates a unique id
func (fa FakeAdapter) UUID() string {
	u := uuid.NewV4()
	return u.String()
}
