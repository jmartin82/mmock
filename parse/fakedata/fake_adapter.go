package fakedata

import (
	"github.com/icrowley/fake"
	"strconv"
)

type FakeAdapter struct {
}

func (this FakeAdapter) Brand() string {
	return fake.Brand()
}
func (this FakeAdapter) Character() string {
	return fake.Character()
}
func (this FakeAdapter) Characters() string {
	return fake.Characters()
}
func (this FakeAdapter) City() string {
	return fake.City()
}
func (this FakeAdapter) Color() string {
	return fake.Color()
}
func (this FakeAdapter) Company() string {
	return fake.Company()
}
func (this FakeAdapter) Continent() string {
	return fake.Continent()
}
func (this FakeAdapter) Country() string {
	return fake.Country()
}
func (this FakeAdapter) CreditCardVisa() string {
	return fake.CreditCardNum("Visa")
}
func (this FakeAdapter) CreditCardMasterCard() string {
	return fake.CreditCardNum("MasterCard")
}
func (this FakeAdapter) CreditCardAmericanExpress() string {
	return fake.CreditCardNum("American Express")
}
func (this FakeAdapter) Currency() string {
	return fake.Currency()
}
func (this FakeAdapter) CurrencyCode() string {
	return fake.CurrencyCode()
}
func (this FakeAdapter) Digits() string {
	return fake.Digits()
}
func (this FakeAdapter) EmailAddress() string {
	return fake.EmailAddress()
}
func (this FakeAdapter) FirstName() string {
	return fake.FirstName()
}
func (this FakeAdapter) FullName() string {
	return fake.FullName()
}
func (this FakeAdapter) LastName() string {
	return fake.LastName()
}
func (this FakeAdapter) Gender() string {
	return fake.Gender()
}
func (this FakeAdapter) IPv4() string {
	return fake.IPv4()
}
func (this FakeAdapter) Language() string {
	return fake.Language()
}
func (this FakeAdapter) Model() string {
	return fake.Model()
}
func (this FakeAdapter) Paragraph() string {
	return fake.Paragraph()
}
func (this FakeAdapter) Paragraphs() string {
	return fake.Paragraphs()
}
func (this FakeAdapter) Phone() string {
	return fake.Phone()
}
func (this FakeAdapter) Product() string {
	return fake.Product()
}
func (this FakeAdapter) Sentence() string {
	return fake.Sentence()
}
func (this FakeAdapter) Sentences() string {
	return fake.Sentences()
}
func (this FakeAdapter) SimplePassword() string {
	return fake.SimplePassword()
}
func (this FakeAdapter) State() string {
	return fake.State()
}
func (this FakeAdapter) StateAbbrev() string {
	return fake.StateAbbrev()
}
func (this FakeAdapter) Street() string {
	return fake.Street()
}
func (this FakeAdapter) StreetAddress() string {
	return fake.StreetAddress()
}
func (this FakeAdapter) UserName() string {
	return fake.UserName()
}
func (this FakeAdapter) Day() string {
	return strconv.Itoa(fake.Day())
}
func (this FakeAdapter) Month() string {
	return fake.Month()
}
func (this FakeAdapter) Year() string {
	return strconv.Itoa(fake.Year(1980, 2020))
}
func (this FakeAdapter) MonthShort() string {
	return fake.MonthShort()
}
func (this FakeAdapter) WeekDay() string {
	return fake.WeekDay()
}
func (this FakeAdapter) Word() string {
	return fake.Word()
}
func (this FakeAdapter) Words() string {
	return fake.Words()
}
func (this FakeAdapter) Zip() string {
	return fake.Zip()
}
