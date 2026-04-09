package faker

import (
	"reflect"
	"testing"
)

// TestTemplateTagPrefixParsing ensures that only tags of the form `template:...`
// are treated as template fields. Tags whose name merely begins with the substring
// "template" (e.g. `templated`) or the bare word `template` without a colon must
// not be hijacked by the template-tag pipeline.
func TestTemplateTagPrefixParsing(t *testing.T) {
	// Bare "template" with no colon: not a template field. Should fall through
	// to the normal tag pipeline, which has no provider named "template" and
	// must therefore return an error rather than silently producing "template".
	type NoColon struct {
		X string `faker:"template"`
	}
	var n NoColon
	if err := FakeData(&n); err == nil && n.X == "template" {
		t.Fatalf("bare 'template' tag was incorrectly treated as a template field, got X=%q", n.X)
	}

	// "templated" shares a prefix with "template" but is a distinct (unknown) tag.
	// It must not be routed through the template engine.
	type FalsePrefix struct {
		X string `faker:"templated"`
	}
	var f FalsePrefix
	if err := FakeData(&f); err == nil && f.X == "templated" {
		t.Fatalf("'templated' tag was incorrectly treated as a template field, got X=%q", f.X)
	}
}

func TestEvaluateTemplateFields(t *testing.T) {
	cases := []struct {
		name        string
		setup       func() (reflect.Type, reflect.Value, []int)
		expectErr   bool
		expectCheck func() error
	}{
		{
			name: "email-lower",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					Name  string
					Email string `faker:"template:{{.Name | lower}}@example.com"`
				}
				var inst T
				inst.Name = "John"
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{1}
			},
			expectErr: false,
			expectCheck: func() error {
				return nil
			},
		},
		{
			name: "multiple",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					Name     string
					Email    string `faker:"template:{{.Name | lower}}@example.com"`
					Username string `faker:"template:{{.Name | lower}}"`
				}
				var inst T
				inst.Name = "John"
				// template fields are at indices 1 and 2
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{1, 2}
			},
			expectErr: false,
		},
		{
			name: "chained",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					Name  string
					Slug  string `faker:"template:{{.Name | slug}}"`
					Email string `faker:"template:{{.Slug}}@example.com"`
				}
				var inst T
				inst.Name = "John"
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{1, 2}
			},
			expectErr: false,
		},
		{
			name: "slug",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					Name string
					Slug string `faker:"template:{{.Name | slug}}"`
				}
				var inst T
				inst.Name = "Hello, World!"
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{1}
			},
			expectErr: false,
		},
		{
			name: "nested",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type Addr struct{ City string }
				type T struct {
					Addr Addr
					Desc string `faker:"template:City is {{.Addr.City}}"`
				}
				var inst T
				inst.Addr.City = "Jakarta"
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{1}
			},
			expectErr: false,
		},
		{
			name: "order",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					Email string `faker:"template:{{.Name | lower}}@example.com"`
					Name  string
				}
				var inst T
				inst.Name = "John"
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{0}
			},
			expectErr: false,
		},
		{
			name: "nonstring-error",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					Name string
					Age  int `faker:"template:{{.Name}}"`
				}
				var inst T
				inst.Name = "John"
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{1}
			},
			expectErr: true,
		},
		{
			name: "circular-reference",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					A string `faker:"template:{{.B}}-x"`
					B string `faker:"template:{{.A}}-y"`
				}
				var inst T
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{0, 1}
			},
			expectErr: true,
		},
		{
			name: "self-reference",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					A string `faker:"template:{{.A}}-x"`
				}
				var inst T
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{0}
			},
			expectErr: true,
		},
		{
			name: "forward-reference",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					A string `faker:"template:{{.B}}"`
					B string `faker:"template:hello"`
				}
				var inst T
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{0, 1}
			},
			expectErr: true,
		},
		{
			name: "missing-field-reference",
			setup: func() (reflect.Type, reflect.Value, []int) {
				type T struct {
					Name string
					Bad  string `faker:"template:{{.DoesNotExist}}"`
				}
				var inst T
				inst.Name = "John"
				return reflect.TypeOf(inst), reflect.ValueOf(&inst).Elem(), []int{1}
			},
			expectErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rt, rv, fields := c.setup()
			err := EvaluateTemplateFields(rt, rv, fields, "faker")
			if c.expectErr {
				if err == nil {
					t.Fatalf("expected error for case %s", c.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for case %s: %v", c.name, err)
			}
			// basic sanity checks per-case
			switch c.name {
			case "email-lower":
				if rv.Field(1).String() != "john@example.com" {
					t.Fatalf("email-lower mismatch: %q", rv.Field(1).String())
				}
			case "slug":
				if rv.Field(1).String() != "hello-world" {
					t.Fatalf("slug mismatch: %q", rv.Field(1).String())
				}
			case "nested":
				if rv.Field(1).String() != "City is Jakarta" {
					t.Fatalf("nested mismatch: %q", rv.Field(1).String())
				}
			case "order":
				if rv.Field(0).String() != "john@example.com" {
					t.Fatalf("order mismatch: %q", rv.Field(0).String())
				}
			}
		})
	}
}
