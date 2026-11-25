package faker

import (
	"fmt"
	"math"
	mathrand "math/rand"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-faker/faker/v4/pkg/options"
)

func TestCreditCardType(t *testing.T) {
	ccType, err := GetPayment().CreditCardType(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	randCC := strings.ToLower(ccType.(string))

	if _, exist := creditCards[randCC]; !exist {
		t.Errorf("Expected from function creditCardType() : %s", randCC)
	}
}

func TestCreditCardNumber(t *testing.T) {
	_, err := GetPayment().CreditCardNumber(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
}

func TestFakeCreditCardType(t *testing.T) {
	ccType := CCType()
	randCC := strings.ToLower(ccType)

	if _, exist := creditCards[randCC]; !exist {
		t.Errorf("Expected from function creditCardType() : %s", randCC)
	}
}

func TestFakeCreditCardNumber(t *testing.T) {
	ccNumber := CCNumber()

	if ccNumber == "" {
		t.Error("Expected Credit Card Number ")
	}
}

// luhnValid reports whether a number made of digits passes the Luhn check
func luhnValid(number string) bool {
	sum := 0
	for i := range len(number) {
		digit := int(number[len(number)-1-i]) - '0'
		if digit < 0 || digit > 9 {
			return false
		}

		if i%2 == 1 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}

// cardNumberMatches reports whether the number has the length, a prefix and the check digit of the card type
func cardNumberMatches(number, ccType string) bool {
	card := creditCards[strings.ToLower(ccType)]
	if len(number) != card.length || !luhnValid(number) {
		return false
	}

	for _, prefix := range card.prefixes {
		if strings.HasPrefix(number, strconv.Itoa(prefix)) {
			return true
		}
	}

	return false
}

// luhnCheckDigit is the one digit that makes a number pass the Luhn check
func TestLuhnCheckDigit(t *testing.T) {
	// the card numbers that the payment providers publish for tests
	for _, number := range []string{
		"4111111111111111", "5555555555554444", "378282246310005", "6011111111111117", "30569309025904", "3530111333300000",
	} {
		prefix, want := number[:len(number)-1], number[len(number)-1:]
		if got := luhnCheckDigit(prefix); got != want {
			t.Errorf("the check digit of %s is %s, want %s", prefix, got, want)
		}
	}

	// every prefix of 4 and of 5 digits: exactly one check digit makes the number pass, and it is the one
	for _, width := range []int{4, 5} {
		for n := range int(math.Pow10(width)) {
			prefix := fmt.Sprintf("%0*d", width, n)
			passing := ""
			for check := range 10 {
				if luhnValid(prefix + strconv.Itoa(check)) {
					passing += strconv.Itoa(check)
				}
			}

			if got := luhnCheckDigit(prefix); got != passing {
				t.Fatalf("the check digit of %s is %s, the digits that pass the Luhn check are %q", prefix, got, passing)
			}
		}
	}
}

// A card number ends with the Luhn check digit, so the form of a shop that checks it accepts the number
func TestFakeCreditCardNumberLuhn(t *testing.T) {
	const samples = 6000
	invalid, example := 0, ""
	for range samples {
		if number := CCNumber(); !luhnValid(number) {
			invalid++
			example = number
		}
	}

	if invalid > 0 {
		t.Errorf("%d of %d card numbers fail the Luhn check, e.g. %s", invalid, samples, example)
	}
}

// CCType picks a type on every call, it used to return the first one it picked on every later call
func TestFakeCreditCardTypeVaries(t *testing.T) {
	seen := map[string]bool{}
	for range 600 {
		seen[CCType()] = true
	}

	if len(seen) != len(creditCards) {
		t.Errorf("600 calls of CCType gave %d of the %d card types: %v", len(seen), len(creditCards), seen)
	}
}

// Outside a struct the methods of GetPayment give a card of their own on every call
func TestCreditCardProvidersOutsideAStruct(t *testing.T) {
	types := map[string]bool{}
	for range 600 {
		ccType, err := GetPayment().CreditCardType(reflect.Value{})
		if err != nil {
			t.Fatal(err)
		}

		types[ccType.(string)] = true

		number, err := GetPayment().CreditCardNumber(reflect.Value{})
		if err != nil {
			t.Fatal(err)
		}

		if !luhnValid(number.(string)) {
			t.Fatalf("%s fails the Luhn check", number)
		}
	}

	if len(types) != len(creditCards) {
		t.Errorf("600 calls of CreditCardType gave %d of the %d card types: %v", len(types), len(creditCards), types)
	}
}

type cardTypeFirst struct {
	Type   string `faker:"cc_type"`
	Number string `faker:"cc_number"`
}

type cardNumberFirst struct {
	Number string `faker:"cc_number"`
	Type   string `faker:"cc_type"`
}

// The cc_type and the cc_number of a struct are one card, and every struct gets a card of its own
func TestCreditCardTypeAndNumberOfAStruct(t *testing.T) {
	types := map[string]bool{}
	for range 600 {
		var first cardTypeFirst
		if err := FakeData(&first); err != nil {
			t.Fatal(err)
		}

		if !cardNumberMatches(first.Number, first.Type) {
			t.Fatalf("%s is not a %s number", first.Number, first.Type)
		}

		var second cardNumberFirst
		if err := FakeData(&second); err != nil {
			t.Fatal(err)
		}

		if !cardNumberMatches(second.Number, second.Type) {
			t.Fatalf("%s is not a %s number", second.Number, second.Type)
		}

		types[first.Type], types[second.Type] = true, true
	}

	if len(types) != len(creditCards) {
		t.Errorf("1200 structs got %d of the %d card types: %v", len(types), len(creditCards), types)
	}
}

// The elements of a slice are structs of their own, each with a card of its own
func TestCreditCardTypesOfASlice(t *testing.T) {
	var cards []cardTypeFirst
	if err := FakeData(&cards, options.WithRandomMapAndSliceMinSize(600), options.WithRandomMapAndSliceMaxSize(601)); err != nil {
		t.Fatal(err)
	}

	if len(cards) != 600 {
		t.Fatalf("got %d structs, want 600", len(cards))
	}

	types := map[string]bool{}
	for _, card := range cards {
		if !cardNumberMatches(card.Number, card.Type) {
			t.Fatalf("%s is not a %s number", card.Number, card.Type)
		}

		types[card.Type] = true
	}

	if len(types) != len(creditCards) {
		t.Errorf("600 structs of a slice got %d of the %d card types: %v", len(types), len(creditCards), types)
	}
}

type cardPointers struct {
	Number *string `faker:"cc_number"`
	Type   *string `faker:"cc_type"`
}

// The fields of a struct that are pointers to strings are one card too
func TestCreditCardTypeAndNumberOfPointerFields(t *testing.T) {
	for range 600 {
		var card cardPointers
		if err := FakeData(&card); err != nil {
			t.Fatal(err)
		}

		if card.Number == nil || card.Type == nil {
			t.Fatalf("got the number %v and the type %v, want both", card.Number, card.Type)
		}

		if !cardNumberMatches(*card.Number, *card.Type) {
			t.Fatalf("%s is not a %s number", *card.Number, *card.Type)
		}
	}
}

type cardNested struct {
	Type   string `faker:"cc_type"`
	Inner  cardNumberFirst
	Number string `faker:"cc_number"`
}

// A struct inside a struct has a card of its own and leaves the card of the struct around it alone
func TestCreditCardOfNestedStructs(t *testing.T) {
	for range 600 {
		var card cardNested
		if err := FakeData(&card); err != nil {
			t.Fatal(err)
		}

		if !cardNumberMatches(card.Number, card.Type) {
			t.Fatalf("%s is not a %s number", card.Number, card.Type)
		}

		if !cardNumberMatches(card.Inner.Number, card.Inner.Type) {
			t.Fatalf("%s is not a %s number", card.Inner.Number, card.Inner.Type)
		}
	}
}

// The same seed gives the same card types in the same order
func TestFakeCreditCardTypeSeeded(t *testing.T) {
	defer SetRandomSource(NewSafeSource(mathrand.NewSource(time.Now().UnixNano())))

	draw := func() []string {
		SetRandomSource(NewSafeSource(mathrand.NewSource(7)))
		types := make([]string, 30)
		for i := range types {
			types[i] = CCType()
		}

		return types
	}

	first, second := draw(), draw()
	if !slices.Equal(first, second) {
		t.Errorf("the same seed gave %v and then %v", first, second)
	}

	if len(slices.Compact(slices.Sorted(slices.Values(first)))) < 2 {
		t.Errorf("30 card types of one seed are all the same: %v", first)
	}
}

// Structs that are filled from several goroutines at once keep a card of their own: the cc_type and the cc_number of
// every struct agree, whatever the other goroutines do (run with -race as well)
func TestCreditCardConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	var mismatched, total atomic.Int64
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 2000 {
				var first cardTypeFirst
				if err := FakeData(&first); err != nil {
					t.Error(err)
					return
				}

				var second cardNumberFirst
				if err := FakeData(&second); err != nil {
					t.Error(err)
					return
				}

				total.Add(2)
				if !cardNumberMatches(first.Number, first.Type) {
					mismatched.Add(1)
				}

				if !cardNumberMatches(second.Number, second.Type) {
					mismatched.Add(1)
				}

				_, _ = CCType(), CCNumber()
			}
		}()
	}

	wg.Wait()

	if mismatched.Load() > 0 {
		t.Errorf("cc_type and cc_number did not agree in %d of %d structs", mismatched.Load(), total.Load())
	}
}
