package faker

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cj/faker/pkg/slice"
)

func TestBloodType(t *testing.T) {
	bloodType, err := GetBlood().BloodType(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(bloodTypes, bloodType.(string)) {
		t.Error("Expected value from variable bloodType in function BloodType")
	}
}

func TestBloodRhFactor(t *testing.T) {
	bloodRhFactor, err := GetBlood().BloodRHFactor(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !slice.Contains(bloodRhFactors, bloodRhFactor.(string)) {
		t.Error("Expected value from variable bloodRhFactor in function BloodType")
	}
}

func TestBloodGroup(t *testing.T) {
	bloodTypes = []string{"O"}
	bloodRhFactors = []string{"+"}
	bloodGroup, err := GetBlood().BloodGroup(reflect.Value{})
	if err != nil {
		t.Error("Expected  not error, got err", err)
	}
	if !strings.Contains(bloodGroup.(string), "O+") {
		t.Error("Expected get url")
	}
}
