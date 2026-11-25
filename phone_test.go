package faker

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// A number of the North American Numbering Plan: the area code and the exchange start with 2 to 9
// and neither of them is N11, the form of the service codes like 411 and 911
var nanpNumber = regexp.MustCompile(`^[2-9]\d\d-[2-9]\d\d-\d{4}$`)

func followsNumberingPlan(number string) bool {
	return nanpNumber.MatchString(number) && number[1:3] != "11" && number[5:7] != "11"
}

// Phonenumber is "201-886-0269": a number that exists in the numbering plan, with digits drawn one by one
func TestPhonenumberFollowsNumberingPlan(t *testing.T) {
	const samples = 20000
	invalid, example := 0, ""
	for range samples {
		if ph := Phonenumber(); !followsNumberingPlan(ph) {
			invalid++
			example = ph
		}
	}

	if invalid > 0 {
		t.Errorf("%d of %d phone numbers do not follow the numbering plan, e.g. %s", invalid, samples, example)
	}
}

// TollFreePhoneNumber is "(888) 937-7238": a three digit code, the exchange and four digits
func TestTollFreePhoneNumberShape(t *testing.T) {
	shape := regexp.MustCompile(`^\((777|888)\) ([2-9]\d\d)-\d{4}$`)
	const samples = 20000
	invalid, example := 0, ""
	for range samples {
		ph := TollFreePhoneNumber()
		match := shape.FindStringSubmatch(ph)
		if match == nil || match[2][1:] == "11" {
			invalid++
			example = ph
		}
	}

	if invalid > 0 {
		t.Errorf("%d of %d toll free numbers are not of the form (888) 937-7238, e.g. %s", invalid, samples, example)
	}
}

// E164PhoneNumber is a plus and eleven digits, like "+27113456789"
func TestE164PhoneNumberShape(t *testing.T) {
	shape := regexp.MustCompile(`^\+[78]\d{10}$`)
	const samples = 20000
	invalid, example := 0, ""
	for range samples {
		if ph := E164PhoneNumber(); !shape.MatchString(ph) {
			invalid++
			example = ph
		}
	}

	if invalid > 0 {
		t.Errorf("%d of %d E.164 numbers are not a plus and eleven digits, e.g. %s", invalid, samples, example)
	}
}

// The digits are drawn one by one: each digit comes up about as often at the places that are free, and the
// numbers hardly ever come twice (they used to be a permutation of 1 to 10, or of 1 to 9 for the toll free
// ones: 3.6 million and 362,880 numbers). 20,000 draws from the 15.8 million toll free numbers repeat about 13
func TestPhoneNumberDigits(t *testing.T) {
	const samples = 20000
	cases := map[string]struct {
		generate      func() string
		free          []int
		maxDuplicates int
	}{
		"Phonenumber":         {func() string { return Phonenumber() }, []int{8, 9, 10, 11}, 5},
		"TollFreePhoneNumber": {func() string { return TollFreePhoneNumber() }, []int{10, 11, 12, 13}, 40},
		"E164PhoneNumber":     {func() string { return E164PhoneNumber() }, []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 5},
	}

	for name, tc := range cases {
		counts := make(map[int][10]int)
		seen := make(map[string]bool, samples)
		duplicates := 0
		for range samples {
			ph := tc.generate()
			if seen[ph] {
				duplicates++
			}

			seen[ph] = true
			for _, place := range tc.free {
				digits := counts[place]
				digits[ph[place]-'0']++
				counts[place] = digits
			}
		}

		uneven, example := 0, ""
		for _, place := range tc.free {
			for digit, count := range counts[place] {
				if count < samples/10*85/100 || count > samples/10*115/100 {
					uneven++
					example = fmt.Sprintf("the digit %d is at place %d in %d of %d numbers, about %d expected", digit, place, count, samples, samples/10)
				}
			}
		}

		if uneven > 0 {
			t.Errorf("%s: %d digits are far from even at the free places, e.g. %s", name, uneven, example)
		}

		if duplicates > tc.maxDuplicates {
			t.Errorf("%s: %d of %d numbers came up twice", name, duplicates, samples)
		}
	}
}

func TestPhoneNumber(t *testing.T) {
	ph, err := GetPhoner().PhoneNumber(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if strings.Count(ph.(string), "-") != 2 {
		t.Error("Expected no more than two characters '-'")
	}
}

func TestTollFreePhoneNumber(t *testing.T) {
	ph, err := GetPhoner().TollFreePhoneNumber(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !strings.HasPrefix(ph.(string), "(888)") && !strings.HasPrefix(ph.(string), "(777)") {
		t.Error("Expected character '(888)' or (777), in function TollFreePhoneNumber")
	}
}

func TestE164PhoneNumber(t *testing.T) {
	ph, err := GetPhoner().E164PhoneNumber(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !strings.HasPrefix(ph.(string), "+") {
		t.Error("Expected character '(888)', in function TollFreePhoneNumber")
	}
}

func TestFakePhoneNumber(t *testing.T) {
	ph := Phonenumber()
	if strings.Count(ph, "-") != 2 {
		t.Error("Expected no more than two characters '-'")
	}
}

func TestFakeTollFreePhoneNumber(t *testing.T) {
	ph := TollFreePhoneNumber()
	if !strings.HasPrefix(ph, "(888)") && !strings.HasPrefix(ph, "(777)") {
		t.Error("Expected character '(888)' or (777), in function TollFreePhoneNumber")
	}
}

func TestFakeE164PhoneNumber(t *testing.T) {
	ph := E164PhoneNumber()
	if !strings.HasPrefix(ph, "+") {
		t.Error("Expected character '(888)', in function TollFreePhoneNumber")
	}
}

func isValidChineseMobileNumber(number string) bool {
	if len(number) != 11 || number[0] != '1' {
		return false
	}
	for _, r := range number {
		if r < '0' || r > '9' {
			return false
		}
	}
	prefix := number[:3]
	for _, p := range chineseMobilePrefixes {
		if p == prefix {
			return true
		}
	}
	return false
}

func TestChinesePhoneNumber(t *testing.T) {
	ph, err := GetPhoner().ChinesePhoneNumber(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !isValidChineseMobileNumber(ph.(string)) {
		t.Errorf("Expected a valid Chinese mobile number, got: %s", ph.(string))
	}
}

func TestFakeChinesePhoneNumber(t *testing.T) {
	ph := ChinesePhoneNumber()
	if !isValidChineseMobileNumber(ph) {
		t.Errorf("Expected a valid Chinese mobile number, got: %s", ph)
	}
}

type chinesePhoneStruct struct {
	Number string `faker:"chinese_phone_number"`
}

func TestChinesePhoneNumberTag(t *testing.T) {
	s := chinesePhoneStruct{}
	if err := FakeData(&s); err != nil {
		t.Fatal("Expected NoError, but got Err:", err)
	}
	if !isValidChineseMobileNumber(s.Number) {
		t.Errorf("Expected a valid Chinese mobile number, got: %s", s.Number)
	}
}
