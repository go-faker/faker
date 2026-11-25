package faker

import (
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4/pkg/options"
)

var uuidDigitRe = regexp.MustCompile(`^[0-9a-f]{32}$`)

// Issue #37: a built-in tag on a slice field must fill each element with that
// tag's value instead of panicking on reflect.Set.
func TestSliceWithBuiltInTag(t *testing.T) {
	type devices struct {
		IDs    []string `faker:"uuid_digit"`
		URLs   []string `faker:"url"`
		Emails []string `faker:"email"`
	}

	for i := 0; i < 20; i++ {
		d := devices{}
		if err := FakeData(&d, options.WithRandomMapAndSliceMinSize(4), options.WithRandomMapAndSliceMaxSize(8)); err != nil {
			t.Fatalf("FakeData: %v", err)
		}
		for name, got := range map[string][]string{"IDs": d.IDs, "URLs": d.URLs, "Emails": d.Emails} {
			if len(got) < 4 || len(got) >= 8 {
				t.Fatalf("%s: want length in [4,8), got %d", name, len(got))
			}
		}
		for _, id := range d.IDs {
			if !uuidDigitRe.MatchString(id) {
				t.Fatalf("not a uuid_digit: %q", id)
			}
		}
		for _, u := range d.URLs {
			if p, err := url.Parse(u); err != nil || p.Scheme == "" || p.Host == "" {
				t.Fatalf("not a url: %q", u)
			}
		}
		for _, e := range d.Emails {
			if !strings.Contains(e, "@") {
				t.Fatalf("not an email: %q", e)
			}
		}
	}
}

func TestSliceWithBuiltInTagAndSliceLen(t *testing.T) {
	type s struct {
		A []string `faker:"slice_len=3,url"`
		B []string `faker:"url,slice_len=3"`
		C []string `faker:"slice_len=3, uuid_digit"`
		D []string `faker:"uuid_digit,slice_len=0"`
		E []string `faker:"slice_len=3, url"` // exact form from issue #37
	}
	v := s{}
	if err := FakeData(&v); err != nil {
		t.Fatalf("FakeData: %v", err)
	}
	if len(v.A) != 3 || len(v.B) != 3 || len(v.C) != 3 || len(v.D) != 0 || len(v.E) != 3 {
		t.Fatalf("wrong lengths: A=%d B=%d C=%d D=%d E=%d", len(v.A), len(v.B), len(v.C), len(v.D), len(v.E))
	}
	for _, u := range append(v.A, v.E...) {
		if p, err := url.Parse(u); err != nil || p.Scheme == "" || p.Host == "" {
			t.Fatalf("not a url: %q", u)
		}
	}
	for _, id := range v.C {
		if !uuidDigitRe.MatchString(id) {
			t.Fatalf("not a uuid_digit: %q", id)
		}
	}
}

func TestArrayWithBuiltInTag(t *testing.T) {
	type s struct {
		IDs [3]string `faker:"uuid_digit"`
	}
	v := s{}
	if err := FakeData(&v); err != nil {
		t.Fatalf("FakeData: %v", err)
	}
	for _, id := range v.IDs {
		if !uuidDigitRe.MatchString(id) {
			t.Fatalf("not a uuid_digit: %q", id)
		}
	}
}

func TestSliceWithBuiltInTagNamedElemType(t *testing.T) {
	type Email string
	type s struct {
		Emails []Email `faker:"email"`
	}
	v := s{}
	if err := FakeData(&v, options.WithRandomMapAndSliceMinSize(1)); err != nil {
		t.Fatalf("FakeData: %v", err)
	}
	for _, e := range v.Emails {
		if !strings.Contains(string(e), "@") {
			t.Fatalf("not an email: %q", e)
		}
	}
}

func TestSliceWithCustomElementProvider(t *testing.T) {
	if err := AddProvider("slice-elem-provider", func(v reflect.Value) (any, error) {
		return "x", nil
	}); err != nil {
		t.Fatal(err)
	}
	type s struct {
		Vals []string `faker:"slice-elem-provider,slice_len=4"`
	}
	v := s{}
	if err := FakeData(&v); err != nil {
		t.Fatalf("FakeData: %v", err)
	}
	if !reflect.DeepEqual(v.Vals, []string{"x", "x", "x", "x"}) {
		t.Fatalf("got %v", v.Vals)
	}
}

// A tag whose value fits neither the slice nor its element must return an
// error, not panic.
func TestSliceWithIncompatibleTag(t *testing.T) {
	type s struct {
		Nums []int `faker:"uuid_digit"`
	}
	v := s{}
	if err := FakeData(&v); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

// A custom provider that builds the whole slice keeps working when slice_len is
// also set; its slice is used as returned.
func TestSliceWithCustomWholeSliceProviderAndSliceLen(t *testing.T) {
	if err := AddProvider("whole-slice-provider", func(v reflect.Value) (any, error) {
		return []string{"a", "b"}, nil
	}); err != nil {
		t.Fatal(err)
	}
	type s struct {
		Vals []string `faker:"whole-slice-provider,slice_len=3"`
	}
	v := s{}
	if err := FakeData(&v); err != nil {
		t.Fatalf("FakeData: %v", err)
	}
	if !reflect.DeepEqual(v.Vals, []string{"a", "b"}) {
		t.Fatalf("got %v", v.Vals)
	}
}
