package faker

import (
	"reflect"
	"strings"
	"testing"
)

func TestUserAgent(t *testing.T) {
	layoutEngines := []string{"Gecko", "KHTML", "Presto", "Trident"}
	uaStringContains := func(elems []string, v string) bool {
		for _, s := range elems {
			if strings.Contains(v, s) {
				return true
			}
		}
		return false
	}

	ua, err := GetUserAgent().UserAgent(reflect.Value{})
	t.Logf("agent: %s", ua.(string))
	if err != nil {
		t.Error("Expected not error, got err", err)
	}
	if strings.Index(ua.(string), "Mozilla/5.0") != 0 {
		t.Error("Expected User Agent string to begin with 'Mozilla/5.0'")
	}
	if !uaStringContains(layoutEngines, ua.(string)) {
		t.Error("Expected User Agent string to contain a valid layout engine")
	}
}
