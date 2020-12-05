package fake

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"

	"github.com/icrowley/fake"
	"github.com/twinj/uuid"
)

//Provider contains all available functions to create random data in the mock response.
type Provider struct {
	ccg *CreditCardGenerator
}

func NewFakeDataProvider(ccg *CreditCardGenerator) *Provider {
	return &Provider{ccg: ccg}
}

//Brand returns a random Brand
func (p Provider) Brand() string {
	return fake.Brand()
}

//Character returns a random Character
func (p Provider) Character() string {
	return fake.Character()
}

//Characters returns from 1 to 5 random Characters
func (p Provider) Characters() string {
	return fake.Characters()
}

//CharactersN returns n random Characters
func (p Provider) CharactersN(n int) string {
	return fake.CharactersN(n)
}

//City returns a random City
func (p Provider) City() string {
	return fake.City()
}

//Color returns a random Color
func (p Provider) Color() string {
	return fake.Color()
}

//Company returns a random Company
func (p Provider) Company() string {
	return fake.Company()
}

//Continent returns a random Continent
func (p Provider) Continent() string {
	return fake.Continent()
}

//Country returns a random Country
func (p Provider) Country() string {
	return fake.Country()
}

//CreditCardVisa returns a random CreditCardVisa
func (p Provider) CreditCardVisa() string {
	return p.ccg.CreditCardVisa()
}

//CreditCardVisaElectron returns a random CreditCardVisaElectron
func (p Provider) CreditCardVisaElectron() string {
	return p.ccg.CreditCardVisaElectron()
}

//CreditCardMasterCard returns a random CreditCardMasterCard
func (p Provider) CreditCardMasterCard() string {
	return p.ccg.CreditCardMasterCard()
}

//CreditCardAmericanExpress returns a random CreditCardAmericanExpress
func (p Provider) CreditCardAmericanExpress() string {
	return p.ccg.CreditCardAmericanExpress()
}

//Currency returns a random Currency
func (p Provider) Currency() string {
	return fake.Currency()
}

//CurrencyCode returns a random CurrencyCode
func (p Provider) CurrencyCode() string {
	return fake.CurrencyCode()
}

//Digits returns from 1 to 5 random Digits
func (p Provider) Digits() string {
	return fake.Digits()
}

//DigitsN returns n random Digits
func (p Provider) DigitsN(n int) string {
	return fake.DigitsN(n)
}

//EmailAddress returns a random EmailAddress
func (p Provider) EmailAddress() string {
	return fake.EmailAddress()
}

//FirstName returns a random FirstName
func (p Provider) FirstName() string {
	return fake.FirstName()
}

//FullName returns a random FullName
func (p Provider) FullName() string {
	return fake.FullName()
}

//LastName returns a random LastName
func (p Provider) LastName() string {
	return fake.LastName()
}

//Gender returns a random Gender
func (p Provider) Gender() string {
	return fake.Gender()
}

//Hex returns a random hexidecimal string of length n
func (p Provider) Hex(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)

	return strings.ToLower(hex.EncodeToString(bytes))
}

//IPv4 returns a random IPv4
func (p Provider) IPv4() string {
	return fake.IPv4()
}

//Language returns a random Language
func (p Provider) Language() string {
	return fake.Language()
}

//Model returns a random Model
func (p Provider) Model() string {
	return fake.Model()
}

//Paragraph returns a random Paragraph
func (p Provider) Paragraph() string {
	return fake.Paragraph()
}

//Paragraphs returns from 1 to 5 random Paragraphs
func (p Provider) Paragraphs() string {
	return fake.Paragraphs()
}

//ParagraphsN returns n random Paragraphs
func (p Provider) ParagraphsN(n int) string {
	return fake.ParagraphsN(n)
}

//Phone returns a random Phone
func (p Provider) Phone() string {
	return fake.Phone()
}

//Product returns a random Product
func (p Provider) Product() string {
	return fake.Product()
}

//Sentence returns a random sentence
func (p Provider) Sentence() string {
	return fake.Sentence()
}

//Sentences returns from 1 to 5 random sentences
func (p Provider) Sentences() string {
	return fake.Sentences()
}

//SentencesN returns n random sentences
func (p Provider) SentencesN(n int) string {
	return fake.SentencesN(n)
}

//SimplePassword returns a random simple password
func (p Provider) SimplePassword() string {
	return fake.SimplePassword()
}

//State returns a random state
func (p Provider) State() string {
	return fake.State()
}

//StateAbbrev returns a random state abbrev
func (p Provider) StateAbbrev() string {
	return fake.StateAbbrev()
}

//Street returns a random street
func (p Provider) Street() string {
	return fake.Street()
}

//StreetAddress returns a random street address
func (p Provider) StreetAddress() string {
	return fake.StreetAddress()
}

//UserName returns a random username
func (p Provider) UserName() string {
	return fake.UserName()
}

//Day returns a random day
func (p Provider) Day() string {
	return strconv.Itoa(fake.Day())
}

//Month returns a random month
func (p Provider) Month() string {
	return fake.Month()
}

//Year returns a random year between (1980,2020)
func (p Provider) Year() string {
	return strconv.Itoa(fake.Year(1980, 2020))
}

//MonthShort returns a random month (Short Version)
func (p Provider) MonthShort() string {
	return fake.MonthShort()
}

//MonthNum returns a random month (Numeric Version)
func (p Provider) MonthNum() string {
	month := fake.MonthNum()
	retval := strconv.Itoa(month)

	if month < 10 {
		retval = "0" + retval
	}

	return retval
}

//WeekDay returns a random day of week
func (p Provider) WeekDay() string {
	return fake.WeekDay()
}

//Word returns a random word
func (p Provider) Word() string {
	return fake.Word()
}

//Words returns from 1 to 5 random words
func (p Provider) Words() string {
	return fake.Words()
}

//WordsN returns n random words
func (p Provider) WordsN(n int) string {
	return fake.WordsN(n)
}

//Zip returns a random zip
func (p Provider) Zip() string {
	return fake.Zip()
}

//Int returns a random positive number less than or equal to n
func (p Provider) Int(n int) string {
	return strconv.Itoa(rand.Intn(n + 1))
}

//IntMinMax returns a random positive number greater than min and lower than max
func (p Provider) IntMinMax(values ...int) string {
	return strconv.Itoa(rand.Intn(values[1]-values[0]) + values[0])
}

//Float returns a random positive floating point number less than n
func (p Provider) Float(n int) string {
	f := float64(n)
	value := rand.Float64() * f
	return strconv.FormatFloat(value, 'f', 4, 64)
}

//UUID generates a unique id
func (p Provider) UUID() string {
	u := uuid.NewV4()
	return u.String()
}
