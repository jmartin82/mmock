package fakedata

//DataFacker contains all available functions to create random data in the mock response.
type DataFaker interface {
	//Brand returns a random brand
	Brand() string
	//Character returns a random character
	Character() string
	//Characters returns a random characters
	Characters() string
	//City returns a random city
	City() string
	//Color returns a random color
	Color() string
	//Company returns a random company
	Company() string
	//Continent returns a random continent
	Continent() string
	//Country returns a random country
	Country() string
	//CreditCardVisa returns a random creditCardVisa
	CreditCardVisa() string
	//CreditCardMasterCard returns a random creditCardMasterCard
	CreditCardMasterCard() string
	//CreditCardAmericanExpress returns a random creditCardAmericanExpress
	CreditCardAmericanExpress() string
	//Currency returns a random currency
	Currency() string
	//CurrencyCode returns a random currencyCode
	CurrencyCode() string
	//Day returns a random day
	Day() string
	//Digits returns a random digits
	Digits() string
	//EmailAddress returns a random emailAddress
	EmailAddress() string
	//FirstName returns a random firstName
	FirstName() string
	//FullName returns a random fullName
	FullName() string
	//LastName returns a random lastName
	LastName() string
	//Gender returns a random gender
	Gender() string
	//IPv4 returns a random iPv4
	IPv4() string
	//Language returns a random language
	Language() string
	//Model returns a random model
	Model() string
	//Month returns a random month
	Month() string
	//Year returns a random year
	Year() string
	//MonthShort returns a random monthShort
	MonthShort() string
	//Paragraph returns a random paragraph
	Paragraph() string
	//Paragraphs returns a random paragraphs
	Paragraphs() string
	//Phone returns a random phone
	Phone() string
	//Product returns a random product
	Product() string
	//Sentence returns a random sentence
	Sentence() string
	//Sentences returns a random sentences
	Sentences() string
	//SimplePassword returns a random simplePassword
	SimplePassword() string
	//State returns a random state
	State() string
	//StateAbbrev returns a random stateAbbrev
	StateAbbrev() string
	//Street returns a random street
	Street() string
	//StreetAddress returns a random streetAddress
	StreetAddress() string
	//UserName returns a random userName
	UserName() string
	//WeekDay returns a random weekDay
	WeekDay() string
	//Word returns a random word
	Word() string
	//Words returns a random words
	Words() string
	//Zip returns a random zip
	Zip() string
}
