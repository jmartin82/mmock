package fakedata

//DataFaker interface contains the funtions to obtain the fake data to fill the response.
type DataFaker interface {
	Brand() string
	Character() string
	Characters() string
	CharactersN(n int) string
	City() string
	Color() string
	Company() string
	Continent() string
	Country() string
	CreditCardVisa() string
	CreditCardMasterCard() string
	CreditCardAmericanExpress() string
	Currency() string
	CurrencyCode() string
	Day() string
	Digits() string
	DigitsN(n int) string
	EmailAddress() string
	FirstName() string
	FullName() string
	LastName() string
	Gender() string
	Hex(n int) string
	IPv4() string
	Language() string
	Model() string
	Month() string
	Year() string
	MonthShort() string
	Paragraph() string
	Paragraphs() string
	ParagraphsN(n int) string
	Phone() string
	Product() string
	Sentence() string
	Sentences() string
	SentencesN(n int) string
	SimplePassword() string
	State() string
	StateAbbrev() string
	Street() string
	StreetAddress() string
	UserName() string
	WeekDay() string
	Word() string
	Words() string
	Zip() string
	Int(n int) string
	IntMinMax(values ...int) string
	Float(n int) string
	UUID() string
}
