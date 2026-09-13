package faker

import (
	"reflect"
	"strings"
	"testing"
)

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
