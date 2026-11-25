package faker

import (
	"maps"
	"reflect"
	"slices"
	"strconv"

	"github.com/go-faker/faker/v4/pkg/options"
)

const (
	numberBytes = "0123456789"
)

// creditCard struct
type creditCard struct {
	ccType   string
	length   int
	prefixes []int
}

var creditCards = map[string]creditCard{
	"visa":             {"VISA", 16, []int{4539, 4556, 4916, 4532, 4929, 40240071, 4485, 4716, 4}},
	"mastercard":       {"MasterCard", 16, []int{51, 52, 53, 54, 55}},
	"american express": {"American Express", 15, []int{34, 37}},
	"discover":         {"Discover", 16, []int{6011}},
	"jcb":              {"JCB", 16, []int{3528, 3538, 3548, 3558, 3568, 3578, 3588}},
	"diners club":      {"Diners Club", 14, []int{36, 38, 39}},
}

// creditCardKeys lists the keys of creditCards in a fixed order: a map is read in a random order,
// so the same seed would not pick the same card type
var creditCardKeys = slices.Sorted(maps.Keys(creditCards))

// GetPayment returns a new Render interface of Payment struct
func GetPayment() Render {
	pay := &Payment{}
	return pay
}

// Render contains Whole Random Credit Card Generators with their types
type Render interface {
	CreditCardType(v reflect.Value) (any, error)
	CreditCardNumber(v reflect.Value) (any, error)
}

// Payment struct
type Payment struct{}

// randomCreditCard picks a card type
func randomCreditCard() creditCard {
	return creditCards[creditCardKeys[rand.Intn(len(creditCardKeys))]]
}

// structCreditCard is the card of a struct that is being filled, so that the cc_type and the cc_number of the struct are
// one card. The first of the two that is filled picks it. It belongs to the code that fills the struct, so structs that
// are filled at the same time, for example from several goroutines, do not share one
type structCreditCard struct {
	card   creditCard
	picked bool
}

func (s *structCreditCard) get() creditCard {
	if !s.picked {
		s.card = randomCreditCard()
		s.picked = true
	}

	return s.card
}

// setField fills a field of the struct with the type or the number of the card if its tag is cc_type or cc_number, and
// reports whether it did. The field is a string or a pointer to one
func (s *structCreditCard) setField(field reflect.Value, tag string) bool {
	if tag != CreditCardType && tag != CreditCardNumber {
		return false
	}

	var target reflect.Value

	switch {
	case field.Kind() == reflect.String:
		target = field
	case field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.String:
		target = reflect.New(field.Type().Elem()).Elem()
	default:
		return false
	}

	if tag == CreditCardType {
		target.SetString(s.get().ccType)
	} else {
		target.SetString(creditCardNumber(s.get()))
	}

	if field.Kind() == reflect.Pointer {
		field.Set(target.Addr())
	}

	return true
}

// cctype returns the type of a card of its own. The cc_type field of a struct gets the card of the struct instead
func (p Payment) cctype() string {
	return randomCreditCard().ccType
}

// CreditCardType returns one of the following credit values:
// VISA, MasterCard, American Express, Discover, JCB and Diners Club
func (p Payment) CreditCardType(v reflect.Value) (any, error) {
	return p.cctype(), nil
}

// CCType get a credit card type randomly in string (VISA, MasterCard, etc)
func CCType(opts ...options.OptionFunc) string {
	return singleFakeData(CreditCardType, func() any {
		return randomCreditCard().ccType
	}, opts...).(string)
}

// luhnCheckDigit returns the digit that makes the number followed by it pass the Luhn check
func luhnCheckDigit(number string) string {
	sum := 0
	for i := range len(number) {
		digit := int(number[len(number)-1-i]) - '0'
		// the digit next to the check digit is doubled, and every second one after it
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return strconv.Itoa((10 - sum%10) % 10)
}

// creditCardNumber returns a number of the card: one of its prefixes, random digits and the Luhn check digit
func creditCardNumber(card creditCard) string {
	prefix := strconv.Itoa(card.prefixes[rand.Intn(len(card.prefixes))])
	num := prefix + randomStringNumber(card.length-len(prefix)-1)

	return num + luhnCheckDigit(num)
}

// ccnumber returns a number of a card of its own. The cc_number field of a struct gets the card of the struct instead
func (p Payment) ccnumber() string {
	return creditCardNumber(randomCreditCard())
}

// CreditCardNumber generated credit card number according to the card number rules:
// the length and a prefix of its type, and the Luhn check digit as the last digit
func (p Payment) CreditCardNumber(v reflect.Value) (any, error) {
	return p.ccnumber(), nil
}

// CCNumber get a credit card number randomly in string (VISA, MasterCard, etc)
func CCNumber(opts ...options.OptionFunc) string {
	return singleFakeData(CreditCardNumber, func() any {
		return creditCardNumber(randomCreditCard())
	}, opts...).(string)
}
