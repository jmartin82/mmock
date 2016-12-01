package fakedata

import "strconv"

//DummyDataFaker is used in tests
type DummyDataFaker struct {
	Dummy string
}

func NewDummyDataFaker(dummyString string) DummyDataFaker {
	result := DummyDataFaker{Dummy: dummyString}
	return result
}

func (ddf DummyDataFaker) Brand() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Character() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Characters() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CharactersN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) City() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Color() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Company() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Continent() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Country() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CreditCardVisa() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CreditCardMasterCard() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CreditCardAmericanExpress() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Currency() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CurrencyCode() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Day() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Digits() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) DigitsN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) EmailAddress() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) FirstName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) FullName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) LastName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Gender() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) IPv4() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Language() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Model() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Month() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Year() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) MonthShort() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Paragraph() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Paragraphs() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) ParagraphsN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) Phone() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Product() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Sentence() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Sentences() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) SentencesN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) SimplePassword() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) State() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) StateAbbrev() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Street() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) StreetAddress() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) UserName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) WeekDay() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Word() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Words() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) WordsN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) Zip() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Int(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) Float(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) UUID() string {
	return "00000000-0000-0000-0000-000000000000"
}
