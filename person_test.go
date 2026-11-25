package faker

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4/pkg/slice"
)

func TestTitleMale(t *testing.T) {
	male, err := GetPerson().TitleMale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(titlesMale, male.(string)) {
		t.Error("Expected value from variable titleMales in function TitleMale")
	}
}

func TestTitleFemale(t *testing.T) {
	female, err := GetPerson().TitleFeMale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(titlesFemale, female.(string)) {
		t.Error("Expected value from variable titleFemales in function TitleFeMale")
	}
}

func TestFirstNameMale(t *testing.T) {
	firstName, err := GetPerson().FirstNameMale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(firstNamesMale, firstName.(string)) {
		t.Error("Expected value from variable firstNamesMale in function FirstNameMale")
	}
}

func TestFirstNameFemale(t *testing.T) {
	firstName, err := GetPerson().FirstNameFemale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(firstNamesFemale, firstName.(string)) {
		t.Error("Expected value from variable firstNamesFemale in function FirstNameFemale")
	}
}

func TestFirstName(t *testing.T) {
	firstname, err := GetPerson().FirstName(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(firstNames, firstname.(string)) {
		t.Error("Expected value from either firstNamesMale or firstNamesFemale in function FirstName")
	}
}

func TestLastName(t *testing.T) {
	lastname, err := GetPerson().LastName(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(lastNames, lastname.(string)) {
		t.Error("Expected value from variable lastNames in function LastName")
	}
}

func TestNameMale(t *testing.T) {
	name, err := GetPerson().Name(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if name.(string) == "" {
		t.Error("Expected from function name string get empty string")
	}
}
func TestNameFemale(t *testing.T) {
	name, err := GetPerson().Name(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if name.(string) == "" {
		t.Error("Expected from function name string get empty string")
	}
}

// Name is a male or a female name by the call. It used to be one of them for the whole run of the program,
// because the choice was drawn once when the package started
func TestNameGenderVaries(t *testing.T) {
	const calls = 2000
	male, female := 0, 0
	for range calls {
		title, _, _ := strings.Cut(Name(), " ")
		inMale, inFemale := slice.Contains(titlesMale, title), slice.Contains(titlesFemale, title)
		switch {
		case inMale && !inFemale:
			male++
		case inFemale && !inMale:
			female++
		}
	}

	// the titles in both lists (Dr. and Prof.) are left out, 4 of 6 male and 6 of 8 female titles stay
	if male < 500 || female < 500 {
		t.Errorf("%d names have a title that only a man has and %d a title that only a woman has, about 667 and 750 expected", male, female)
	}
}

// The title and the first name of a Name come from the same list
func TestNameIsOneGender(t *testing.T) {
	fits := func(title, first, last string, titles, firsts []string) bool {
		return slice.Contains(titles, title) && slice.Contains(firsts, first) && slice.Contains(lastNames, last)
	}

	for range 2000 {
		name := Name()
		title, rest, _ := strings.Cut(name, " ")
		ok := false
		for i, r := range rest {
			if r != ' ' {
				continue
			}

			first, last := rest[:i], rest[i+1:]
			if fits(title, first, last, titlesMale, firstNamesMale) || fits(title, first, last, titlesFemale, firstNamesFemale) {
				ok = true
				break
			}
		}

		if !ok {
			t.Fatalf("%q is not a title, a first name and a last name of one gender", name)
		}
	}
}

func TestFakeTitleMale(t *testing.T) {
	male := TitleMale()
	if !slice.Contains(titlesMale, male) {
		t.Error("Expected value from variable titleMales in function TitleMale")
	}
}

func TestFakeTitleFemale(t *testing.T) {
	female := TitleFemale()
	if !slice.Contains(titlesFemale, female) {
		t.Error("Expected value from variable titleFemales in function TitleFeMale")
	}
}

func TestFakeFirstNameMale(t *testing.T) {
	firstName := FirstNameMale()
	if !slice.Contains(firstNamesMale, firstName) {
		t.Error("Expected value from variable firstNamesMale in function FirstNameMale")
	}
}

func TestFakeFirstNameFemale(t *testing.T) {
	firstName := FirstNameFemale()
	if !slice.Contains(firstNamesFemale, firstName) {
		t.Error("Expected value from variable firstNamesFemale in function FirstNameFemale")
	}
}

func TestFakeFirstName(t *testing.T) {
	firstname := FirstName()
	if !slice.Contains(firstNames, firstname) {
		t.Error("Expected value from either firstNamesMale or firstNamesFemale in function FirstName")
	}
}

func TestFakeLastName(t *testing.T) {
	lastname := LastName()
	if !slice.Contains(lastNames, lastname) {
		t.Error("Expected value from variable lastNames in function LastName")
	}
}

func TestFakeNameMale(t *testing.T) {
	name := Name()
	if name == "" {
		t.Error("Expected from function name string get empty string")
	}
}
func TestFakeNameFemale(t *testing.T) {
	name := Name()
	if name == "" {
		t.Error("Expected from function name string get empty string")
	}
	t.Log(name)
}

func TestFakeGender(t *testing.T) {
	gender, err := GetPerson().Gender(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(genders, gender.(string)) {
		t.Error("Expected value from variable genders in function Gender")
	}
}

func TestFakeGenderPublicFunction(t *testing.T) {
	gender := Gender()
	if !slice.Contains(genders, gender) {
		t.Error("Expected value from variable genders in function Gender")
	}
}

func TestChineseFirstName(t *testing.T) {
	firstname, err := GetPerson().ChineseFirstName(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	t.Log(firstname)
	if !slice.Contains(chineseFirstNames, firstname.(string)) {
		t.Error("Expected value from either chineseFirstNames in function ChineseFirstName")
	}
}

func TestChineseLastName(t *testing.T) {
	firstname, err := GetPerson().ChineseLastName(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	t.Log(firstname)
	if !slice.Contains(chineseLastNames, firstname.(string)) {
		t.Error("Expected value from either chineseLastNames in function ChineseLastName")
	}
}

func TestChineseName(t *testing.T) {
	firstname, err := GetPerson().ChineseName(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	t.Log(firstname)
}

func TestRussianFirstNameMale(t *testing.T) {
	firstName, err := GetPerson().RussianFirstNameMale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(russianFirstNamesMale, firstName.(string)) {
		t.Error("Expected value from variable russianFirstNamesMale in function RussianFirstNameMale")
	}
}

func TestRussianMiddleNameMale(t *testing.T) {
	middleName, err := GetPerson().RussianMiddleNameMale(reflect.Value{})
	if err != nil {
		t.Error("Expected not error, got err", err)
	}
	if !slice.Contains(russianMiddleNamesMale, middleName.(string)) {
		t.Error("Expected value from variable russianMiddleNamesMale in function RussianMiddleNameMale")
	}
}

func TestRussianLastNameMale(t *testing.T) {
	lastName, err := GetPerson().RussianLastNameMale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(russianLastNamesMale, lastName.(string)) {
		t.Error("Expected value from variable russianLastNamesMale in function RussianLastNameMale")
	}
}

func TestRussianFirstNameFemale(t *testing.T) {
	firstName, err := GetPerson().RussianFirstNameFemale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(russianFirstNamesFemale, firstName.(string)) {
		t.Error("Expected value from variable russianFirstNamesFemale in function RussianFirstNameFemale")
	}
}

func TestRussianMiddleNameFemale(t *testing.T) {
	middleName, err := GetPerson().RussianMiddleNameFemale(reflect.Value{})
	if err != nil {
		t.Error("Expected not error, got err", err)
	}
	if !slice.Contains(russianMiddleNamesFemale, middleName.(string)) {
		t.Error("Expected value from variable russianMiddleNamesFemale in function RussianMiddleNameFemale")
	}
}

func TestRussianLastNameFemale(t *testing.T) {
	lastName, err := GetPerson().RussianLastNameFemale(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(russianLastNamesFemale, lastName.(string)) {
		t.Error("Expected value from variable russianLastNamesFemale in function RussianLastNameFemale")
	}
}

// TestCheckRussianRunesInRussianNames - homoglyph test
func TestCheckRussianRunesInRussianNames(t *testing.T) {
	var russianNames []string
	russianNames = append(russianNames, russianFirstNamesMale...)
	russianNames = append(russianNames, russianFirstNamesFemale...)
	russianNames = append(russianNames, russianLastNamesMale...)
	russianNames = append(russianNames, russianLastNamesFemale...)
	russianNames = append(russianNames, russianMiddleNamesMale...)
	russianNames = append(russianNames, russianMiddleNamesFemale...)

	for _, v := range russianNames {
		for _, r := range v {
			if (r < 'а' || r > 'я') && (r < 'А' || r > 'Я') && r != 'ё' && r != 'Ё' {
				t.Errorf("Expected russian rune, but got [%c] in word [%s]", r, v)
			}
		}
	}

}

// TestCheckLatinRunesInLatinNames - homoglyph test
func TestCheckLatinRunesInLatinNames(t *testing.T) {
	var latinNames []string
	latinNames = append(latinNames, firstNames...)
	latinNames = append(latinNames, lastNames...)

	for _, v := range latinNames {
		for _, r := range v {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && r != '\'' && r != '-' {
				t.Errorf("Expected latin rune, but got [%c] in word [%s]", r, v)
			}
		}
	}
}
