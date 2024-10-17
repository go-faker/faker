package faker_test

import (
	"fmt"

	"github.com/cj/faker"
	"github.com/cj/faker/pkg/options"
)

// SomeStructWithTags ...
type SomeStructWithACustomTagName struct {
	Latitude           float32           `custom:"lat"`
	Longitude          float32           `custom:"long"`
	RealAddress        faker.RealAddress `custom:"real_address"`
	CreditCardNumber   string            `custom:"cc_number"`
	CreditCardType     string            `custom:"cc_type"`
	Email              string            `custom:"email"`
	DomainName         string            `custom:"domain_name"`
	IPV4               string            `custom:"ipv4"`
	IPV6               string            `custom:"ipv6"`
	Password           string            `custom:"password"`
	Jwt                string            `custom:"jwt"`
	PhoneNumber        string            `custom:"phone_number"`
	MacAddress         string            `custom:"mac_address"`
	URL                string            `custom:"url"`
	UserName           string            `custom:"username"`
	TollFreeNumber     string            `custom:"toll_free_number"`
	E164PhoneNumber    string            `custom:"e_164_phone_number"`
	TitleMale          string            `custom:"title_male"`
	TitleFemale        string            `custom:"title_female"`
	FirstName          string            `custom:"first_name"`
	FirstNameMale      string            `custom:"first_name_male"`
	FirstNameFemale    string            `custom:"first_name_female"`
	LastName           string            `custom:"last_name"`
	Name               string            `custom:"name"`
	UnixTime           int64             `custom:"unix_time"`
	Date               string            `custom:"date"`
	Time               string            `custom:"time"`
	MonthName          string            `custom:"month_name"`
	Year               string            `custom:"year"`
	DayOfWeek          string            `custom:"day_of_week"`
	DayOfMonth         string            `custom:"day_of_month"`
	Timestamp          string            `custom:"timestamp"`
	Century            string            `custom:"century"`
	TimeZone           string            `custom:"timezone"`
	TimePeriod         string            `custom:"time_period"`
	Word               string            `custom:"word"`
	Sentence           string            `custom:"sentence"`
	Paragraph          string            `custom:"paragraph"`
	Currency           string            `custom:"currency"`
	Amount             float64           `custom:"amount"`
	AmountWithCurrency string            `custom:"amount_with_currency"`
	UUIDHypenated      string            `custom:"uuid_hyphenated"`
	UUID               string            `custom:"uuid_digit"`
	Skip               string            `custom:"-"`
	PaymentMethod      string            `custom:"oneof: cc, paypal, check, money order"` // oneof will randomly pick one of the comma-separated values supplied in the tag
	AccountID          int               `custom:"oneof: 15, 27, 61"`                     // use commas to separate the values for now. Future support for other separator characters may be added
	Price32            float32           `custom:"oneof: 4.95, 9.99, 31997.97"`
	Price64            float64           `custom:"oneof: 47463.9463525, 993747.95662529, 11131997.978767990"`
	NumS64             int64             `custom:"oneof: 1, 2"`
	NumS32             int32             `custom:"oneof: -3, 4"`
	NumS16             int16             `custom:"oneof: -5, 6"`
	NumS8              int8              `custom:"oneof: 7, -8"`
	NumU64             uint64            `custom:"oneof: 9, 10"`
	NumU32             uint32            `custom:"oneof: 11, 12"`
	NumU16             uint16            `custom:"oneof: 13, 14"`
	NumU8              uint8             `custom:"oneof: 15, 16"`
	NumU               uint              `custom:"oneof: 17, 18"`
	PtrNumU            *uint             `custom:"oneof: 19, 20"`
}

func Example_withTagsAndCustomTagName() {

	a := SomeStructWithTags{}
	// just set the custom tag name option
	err := faker.FakeData(&a, options.WithTagName("custom"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", a)
	/*
		Result:
		{
			Latitude: 81.12195
			Longitude: -84.38158
			RealAddress: {Address:107 Guaymas Place City:Davis State:CA PostalCode:95616 Coordinates:{Latitude:38.567048 Longitude:-121.746046}}
			CreditCardType: American Express
			CreditCardNumber: 373641309057568
			Email: mJBJtbv@OSAaT.ru
			DomainName: FWZcaRE.ru,
			IPV4: 99.23.42.63
			IPV6: 975c:fb2c:2133:fbdd:beda:282e:1e0a:ec7d
			Password: dfJdyHGuVkHBgnHLQQgpINApynzexnRpgIKBpiIjpTPOmNyMFb
			Jwt: HDMNSOKhEIYkPIuHcVjfCtHlKkaqLGrUEqjKVkgR.HDMNSOKhEIYkPIuHcVjfCtHlKkaqLGrUEqjKVkgR.HDMNSOKhEIYkPIuHcVjfCtHlKkaqLGrUEqjKVkgR
			PhoneNumber: 792-153-4861
			MacAddress: cd:65:e1:d4:76:c6
			URL: https://www.oEuqqAY.org/QgqfOhd
			UserName: lVxELHS
			TollFreeNumber: (777) 831-964572
			E164PhoneNumber: +724891571063
			TitleMale: Mr.
			TitleFemale: Queen
			FirstName: Whitney
			FirstNameMale: Kenny
			FirstNameFemale: Jana
			LastName: Rohan
			Name: Miss Casandra Kiehn
			UnixTime: 1197930901
			Date: 1982-02-27
			Time: 03:10:25
			MonthName: February
			Year: 1996
			DayOfWeek: Sunday
			DayOfMonth: 20
			Timestamp: 1973-06-21 14:50:46
			Century: IV
			TimeZone: Canada/Eastern
			TimePeriod: AM
			Word: nesciunt
			Sentence: Consequatur perferendis aut sit voluptatem accusantium.
			Paragraph: Aut consequatur sit perferendis accusantium voluptatem. Accusantium perferendis consequatur voluptatem sit aut. Aut sit accusantium consequatur voluptatem perferendis. Perferendis voluptatem aut accusantium consequatur sit.
			Currency: IRR,
			Amount: 88.990000,
			AmountWithCurrency: XBB 49257.100000,
			UUIDHypenated: 8f8e4463-9560-4a38-9b0c-ef24481e4e27,
			UUID: 90ea6479fd0e4940af741f0a87596b73,
			PaymentMethod: paypal,
			AccountID: 61,
			Price32: 4.95,
			Price64: 993747.95662529
			NumS64:	1
			NumS32:	-3
			NumS16:	5
			NumS8:	-8
			NumU64:	9
			NumU32:	11
			NumU16:	13
			NumU8:	15
			NumU:	17
			PtrNumU: 19
			Skip:
		}
	*/

}
