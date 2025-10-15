package faker

import (
	"reflect"
	"testing"
)

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
