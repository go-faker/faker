package faker

import (
	"reflect"
	"regexp"
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

func TestNANPPhoneNumber(t *testing.T) {
	// Test single generation with interface method
	ph, err := GetPhoner().NANPPhoneNumber(reflect.Value{})
	if err != nil {
		t.Error("Expected not error, got err", err)
	}
	phoneStr := ph.(string)
	if !isValidNANP(phoneStr) {
		t.Errorf("Generated NANP phone number '%s' is not NANP-compliant", phoneStr)
	}
}

func TestFakeNANPPhoneNumber(t *testing.T) {
	// Test single generation with standalone function
	ph := NANPPhoneNumber()
	if !isValidNANP(ph) {
		t.Errorf("Generated NANP phone number '%s' is not NANP-compliant", ph)
	}
}

func TestNANPPhoneNumberBulk(t *testing.T) {
	// Test bulk generation for comprehensive coverage
	const testCount = 1000

	t.Run("Interface Method", func(t *testing.T) {
		phoner := GetPhoner()
		for i := 0; i < testCount; i++ {
			ph, err := phoner.NANPPhoneNumber(reflect.Value{})
			if err != nil {
				t.Fatalf("Error generating NANP phone number at iteration %d: %v", i, err)
			}
			phoneStr := ph.(string)
			if !isValidNANP(phoneStr) {
				t.Errorf("Generated NANP phone number '%s' is not NANP-compliant at iteration %d", phoneStr, i)
			}
		}
	})

	t.Run("Standalone Function", func(t *testing.T) {
		for i := 0; i < testCount; i++ {
			ph := NANPPhoneNumber()
			if !isValidNANP(ph) {
				t.Errorf("Generated NANP phone number '%s' is not NANP-compliant at iteration %d", ph, i)
			}
		}
	})
}

func TestNANPValidatorHelper(t *testing.T) {
	// Test the validation helper function itself
	testCases := []struct {
		phoneNumber string
		expected    bool
		description string
	}{
		{"234-567-8901", true, "valid NANP number"},
		{"987-654-3210", true, "valid NANP number with high digits"},
		{"212-555-1234", true, "valid NANP number"},
		{"123-456-7890", false, "area code starts with 1"},
		{"034-567-8901", false, "area code starts with 0"},
		{"234-156-7890", false, "central office code starts with 1"},
		{"234-067-8901", false, "central office code starts with 0"},
		{"234-567-890", false, "line number too short"},
		{"234-567-89012", false, "line number too long"},
		{"23-567-8901", false, "area code too short"},
		{"2345-567-8901", false, "area code too long"},
		{"234-56-8901", false, "central office code too short"},
		{"234-5678-901", false, "central office code too long"},
		{"234.567.8901", false, "wrong separator"},
		{"234 567 8901", false, "wrong separator"},
		{"2345678901", false, "no separators"},
		{"", false, "empty string"},
	}

	for _, tc := range testCases {
		result := isValidNANP(tc.phoneNumber)
		if result != tc.expected {
			t.Errorf("isValidNANP('%s') = %v; expected %v (%s)",
				tc.phoneNumber, result, tc.expected, tc.description)
		}
	}
}

// isValidNANP validates if a phone number string is NANP-compliant
// NANP format: NXX-NXX-XXXX where N=2-9, X=0-9
func isValidNANP(phoneNumber string) bool {
	// Use regex for efficient validation
	nanpRegex := regexp.MustCompile(`^[2-9][0-9]{2}-[2-9][0-9]{2}-[0-9]{4}$`)
	return nanpRegex.MatchString(phoneNumber)
}
