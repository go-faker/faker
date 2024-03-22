package faker

import (
	"regexp"
	"testing"
)

func TestGetLongitude(t *testing.T) {
	long := Longitude()
	if long > 180 || long < -180 {
		t.Error("function Longitude need return a valid longitude")
	}
}

func TestGetLatitude(t *testing.T) {
	lat := Latitude()
	if lat > 90 || lat < -90 {
		t.Error("function Latitude need return a valid longitude")
	}
}

func TestGetRealAddress(t *testing.T) {
	addr := GetRealAddress()
	if addr.Address == "" || addr.City == "" {
		t.Error("empty address")
	}
	t.Log(addr)
}

func TestGetCountryInfo(t *testing.T) {
	rand.Seed(31)
	countryInfo := GetCountryInfo()

	expectedCountryName := "Morocco"
	if countryInfo.Name != expectedCountryName {
		t.Errorf("Test failed, expected: %s, got: %s", expectedCountryName, countryInfo.Name)
	}

	if len(countryInfo.Abbr) != 2 {
		t.Error("Invalid ISO-3166 abbreviation")
	}

	if len(countryInfo.Continent) != 2 {
		t.Error("Invalid continent abbreviation")
	}

	re := regexp.MustCompile(`^\d{1,3}(,\d{3})*$`)
	if !re.MatchString(countryInfo.Population) {
		t.Error("Invalid population number")
	}
}
