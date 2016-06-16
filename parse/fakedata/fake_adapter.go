package fakedata

import (
	"strconv"

	"github.com/icrowley/fake"
)

type FakeAdapter struct {
}

//Brand return a random Brand
func (fa FakeAdapter) Brand() string {
	return fake.Brand()
}

//Character return a random Character
func (fa FakeAdapter) Character() string {
	return fake.Character()
}

//Characters return a random Characters
func (fa FakeAdapter) Characters() string {
	return fake.Characters()
}

//City return a random City
func (fa FakeAdapter) City() string {
	return fake.City()
}

//Color return a random Color
func (fa FakeAdapter) Color() string {
	return fake.Color()
}

//Company return a random Company
func (fa FakeAdapter) Company() string {
	return fake.Company()
}

//Continent return a random Continent
func (fa FakeAdapter) Continent() string {
	return fake.Continent()
}

//Country return a random Country
func (fa FakeAdapter) Country() string {
	return fake.Country()
}

//CreditCardVisa return a random CreditCardVisa
func (fa FakeAdapter) CreditCardVisa() string {
	return fake.CreditCardNum("Visa")
}

//CreditCardMasterCard return a random CreditCardMasterCard
func (fa FakeAdapter) CreditCardMasterCard() string {
	return fake.CreditCardNum("MasterCard")
}

//CreditCardAmericanExpress return a random CreditCardAmericanExpress
func (fa FakeAdapter) CreditCardAmericanExpress() string {
	return fake.CreditCardNum("American Express")
}

//Currency return a random Currency
func (fa FakeAdapter) Currency() string {
	return fake.Currency()
}

//CurrencyCode return a random CurrencyCode
func (fa FakeAdapter) CurrencyCode() string {
	return fake.CurrencyCode()
}

//Digits return a random Digits
func (fa FakeAdapter) Digits() string {
	return fake.Digits()
}

//EmailAddress return a random EmailAddress
func (fa FakeAdapter) EmailAddress() string {
	return fake.EmailAddress()
}

//FirstName return a random FirstName
func (fa FakeAdapter) FirstName() string {
	return fake.FirstName()
}

//FullName return a random FullName
func (fa FakeAdapter) FullName() string {
	return fake.FullName()
}

//LastName return a random LastName
func (fa FakeAdapter) LastName() string {
	return fake.LastName()
}

//Gender return a random Gender
func (fa FakeAdapter) Gender() string {
	return fake.Gender()
}

//IPv4 return a random IPv4
func (fa FakeAdapter) IPv4() string {
	return fake.IPv4()
}

//Language return a random Language
func (fa FakeAdapter) Language() string {
	return fake.Language()
}

//Model return a random Model
func (fa FakeAdapter) Model() string {
	return fake.Model()
}

//Paragraph return a random Paragraph
func (fa FakeAdapter) Paragraph() string {
	return fake.Paragraph()
}

//Paragraphs return a random Paragraphs
func (fa FakeAdapter) Paragraphs() string {
	return fake.Paragraphs()
}

//Phone return a random Phone
func (fa FakeAdapter) Phone() string {
	return fake.Phone()
}

//Product return a random Product
func (fa FakeAdapter) Product() string {
	return fake.Product()
}

//Sentence return a random Sentence
func (fa FakeAdapter) Sentence() string {
	return fake.Sentence()
}

//Sentences return a random Sentences
func (fa FakeAdapter) Sentences() string {
	return fake.Sentences()
}

//SimplePassword return a random SimplePassword
func (fa FakeAdapter) SimplePassword() string {
	return fake.SimplePassword()
}

//State return a random State
func (fa FakeAdapter) State() string {
	return fake.State()
}

//StateAbbrev return a random StateAbbrev
func (fa FakeAdapter) StateAbbrev() string {
	return fake.StateAbbrev()
}

//Street return a random Street
func (fa FakeAdapter) Street() string {
	return fake.Street()
}

//StreetAddress return a random StreetAddress
func (fa FakeAdapter) StreetAddress() string {
	return fake.StreetAddress()
}

//UserName return a random UserName
func (fa FakeAdapter) UserName() string {
	return fake.UserName()
}

//Day return a random Day
func (fa FakeAdapter) Day() string {
	return strconv.Itoa(fake.Day())
}

//Month return a random Month
func (fa FakeAdapter) Month() string {
	return fake.Month()
}

//Year return a random Year
func (fa FakeAdapter) Year() string {
	return strconv.Itoa(fake.Year(1980, 2020))
}

//MonthShort return a random MonthShort
func (fa FakeAdapter) MonthShort() string {
	return fake.MonthShort()
}

//WeekDay return a random WeekDay
func (fa FakeAdapter) WeekDay() string {
	return fake.WeekDay()
}

//Word return a random Word
func (fa FakeAdapter) Word() string {
	return fake.Word()
}

//Words return a random Words
func (fa FakeAdapter) Words() string {
	return fake.Words()
}

//Zip return a random Zip
func (fa FakeAdapter) Zip() string {
	return fake.Zip()
}
