package fakedata

import (
	"strconv"

	"github.com/icrowley/fake"
)

type FakeAdapter struct {
}

func (fa FakeAdapter) Brand() string {
	return fake.Brand()
}

func (fa FakeAdapter) Character() string {
	return fake.Character()
}

func (fa FakeAdapter) Characters() string {
	return fake.Characters()
}

func (fa FakeAdapter) City() string {
	return fake.City()
}

func (fa FakeAdapter) Color() string {
	return fake.Color()
}

func (fa FakeAdapter) Company() string {
	return fake.Company()
}

func (fa FakeAdapter) Continent() string {
	return fake.Continent()
}

func (fa FakeAdapter) Country() string {
	return fake.Country()
}

func (fa FakeAdapter) CreditCardVisa() string {
	return fake.CreditCardNum("Visa")
}

func (fa FakeAdapter) CreditCardMasterCard() string {
	return fake.CreditCardNum("MasterCard")
}

func (fa FakeAdapter) CreditCardAmericanExpress() string {
	return fake.CreditCardNum("American Express")
}

func (fa FakeAdapter) Currency() string {
	return fake.Currency()
}

func (fa FakeAdapter) CurrencyCode() string {
	return fake.CurrencyCode()
}

func (fa FakeAdapter) Digits() string {
	return fake.Digits()
}

func (fa FakeAdapter) EmailAddress() string {
	return fake.EmailAddress()
}

func (fa FakeAdapter) FirstName() string {
	return fake.FirstName()
}

func (fa FakeAdapter) FullName() string {
	return fake.FullName()
}

func (fa FakeAdapter) LastName() string {
	return fake.LastName()
}

func (fa FakeAdapter) Gender() string {
	return fake.Gender()
}

func (fa FakeAdapter) IPv4() string {
	return fake.IPv4()
}

func (fa FakeAdapter) Language() string {
	return fake.Language()
}

func (fa FakeAdapter) Model() string {
	return fake.Model()
}

func (fa FakeAdapter) Paragraph() string {
	return fake.Paragraph()
}

func (fa FakeAdapter) Paragraphs() string {
	return fake.Paragraphs()
}

func (fa FakeAdapter) Phone() string {
	return fake.Phone()
}

func (fa FakeAdapter) Product() string {
	return fake.Product()
}

func (fa FakeAdapter) Sentence() string {
	return fake.Sentence()
}

func (fa FakeAdapter) Sentences() string {
	return fake.Sentences()
}

func (fa FakeAdapter) SimplePassword() string {
	return fake.SimplePassword()
}

func (fa FakeAdapter) State() string {
	return fake.State()
}

func (fa FakeAdapter) StateAbbrev() string {
	return fake.StateAbbrev()
}

func (fa FakeAdapter) Street() string {
	return fake.Street()
}

func (fa FakeAdapter) StreetAddress() string {
	return fake.StreetAddress()
}

func (fa FakeAdapter) UserName() string {
	return fake.UserName()
}

func (fa FakeAdapter) Day() string {
	return strconv.Itoa(fake.Day())
}

func (fa FakeAdapter) Month() string {
	return fake.Month()
}

func (fa FakeAdapter) Year() string {
	return strconv.Itoa(fake.Year(1980, 2020))
}

func (fa FakeAdapter) MonthShort() string {
	return fake.MonthShort()
}

func (fa FakeAdapter) WeekDay() string {
	return fake.WeekDay()
}

func (fa FakeAdapter) Word() string {
	return fake.Word()
}

func (fa FakeAdapter) Words() string {
	return fake.Words()
}

func (fa FakeAdapter) Zip() string {
	return fake.Zip()
}
