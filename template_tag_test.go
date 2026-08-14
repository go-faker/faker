package faker

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"text/template"

	fakerErrors "github.com/go-faker/faker/v4/pkg/errors"
	"github.com/go-faker/faker/v4/pkg/options"
)

// bodyA is the template body shared by most TestDecodeTagsTemplate cases.
const bodyA = "{{.A}}"

// tagOf decodes the faker tag of the last field of the struct type of sample.
func tagOf(t *testing.T, sample any) structTag {
	t.Helper()
	typ := reflect.TypeOf(sample)
	return decodeTags(typ, typ.NumField()-1, "faker")
}

// TestDecodeTagsTemplate guards the tag parser against the ways a template body
// used to be mangled: commas splitting it, a chunk being mistaken for a priority
// tag, and keep/unique blanking the whole tag.
func TestDecodeTagsTemplate(t *testing.T) {
	cases := []struct {
		name         string
		sample       any
		isTemplate   bool
		template     string
		fieldType    string
		unique       bool
		keepOriginal bool
	}{
		{
			name: "plain",
			sample: struct {
				X string `faker:"template:{{.A}}"`
			}{},
			isTemplate: true, template: bodyA,
		},
		{
			name: "case_insensitive_prefix",
			sample: struct {
				X string `faker:"Template:{{.A}}"`
			}{},
			isTemplate: true, template: bodyA,
		},
		{
			name: "bare_without_colon_is_not_a_template",
			sample: struct {
				X string `faker:"template"`
			}{},
			fieldType: "template",
		},
		{
			name: "false_prefix_is_not_a_template",
			sample: struct {
				X string `faker:"templated"`
			}{},
			fieldType: "templated",
		},
		{
			name: "comma_in_body_is_preserved",
			sample: struct {
				X string `faker:"template:{{.A}}, {{.B}}"`
			}{},
			isTemplate: true, template: "{{.A}}, {{.B}}",
		},
		{
			// A chunk after a comma used to be matched against PriorityTags, which
			// replaced the whole template with a bogus "email" tag.
			name: "comma_chunk_is_not_hijacked_by_priority_tags",
			sample: struct {
				X string `faker:"template:{{.A}}, email: {{.B}}"`
			}{},
			isTemplate: true, template: "{{.A}}, email: {{.B}}",
		},
		{
			name: "trailing_unique",
			sample: struct {
				X string `faker:"template:{{.A}},unique"`
			}{},
			isTemplate: true, template: bodyA, unique: true,
		},
		{
			name: "trailing_keep",
			sample: struct {
				X string `faker:"template:{{.A}},keep"`
			}{},
			isTemplate: true, template: bodyA, keepOriginal: true,
		},
		{
			name: "leading_unique",
			sample: struct {
				X string `faker:"unique,template:{{.A}}"`
			}{},
			isTemplate: true, template: bodyA, unique: true,
		},
		{
			name: "empty_body",
			sample: struct {
				X string `faker:"template:"`
			}{},
			isTemplate: true, template: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := tagOf(t, c.sample)
			if got.isTemplate != c.isTemplate {
				t.Fatalf("isTemplate = %v, want %v (%+v)", got.isTemplate, c.isTemplate, got)
			}
			if got.template != c.template {
				t.Errorf("template = %q, want %q", got.template, c.template)
			}
			if got.unique != c.unique {
				t.Errorf("unique = %v, want %v", got.unique, c.unique)
			}
			if got.keepOriginal != c.keepOriginal {
				t.Errorf("keepOriginal = %v, want %v", got.keepOriginal, c.keepOriginal)
			}
			if c.isTemplate {
				if got.fieldType != "" {
					t.Errorf("fieldType = %q, want empty for a template tag", got.fieldType)
				}
			} else if got.fieldType != c.fieldType {
				t.Errorf("fieldType = %q, want %q", got.fieldType, c.fieldType)
			}
		})
	}
}

// TestTemplateTagPrefixParsing ensures that only tags of the form `template:...`
// are treated as template fields. Tags whose name merely begins with the substring
// "template" (e.g. `templated`) or the bare word `template` without a colon must
// not be hijacked by the template-tag pipeline.
func TestTemplateTagPrefixParsing(t *testing.T) {
	cases := []struct {
		name string
		tag  string
		fake func() error
	}{
		{
			name: "bare_template_without_colon",
			tag:  "template",
			fake: func() error {
				var v struct {
					X string `faker:"template"`
				}
				return FakeData(&v)
			},
		},
		{
			name: "templated_shares_a_prefix",
			tag:  "templated",
			fake: func() error {
				var v struct {
					X string `faker:"templated"`
				}
				return FakeData(&v)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.fake()
			if err == nil {
				t.Fatalf("expected %q to be rejected as an unknown tag, got nil error", c.tag)
			}
			want := fmt.Sprintf(fakerErrors.ErrTagNotSupported, c.tag)
			if err.Error() != want {
				t.Fatalf("error = %q, want %q", err.Error(), want)
			}
		})
	}
}

// TestTemplateFieldRefs covers dependency extraction. A regular expression over
// the template source cannot get these right: it matches literal text, reports
// sub-selectors as fields, and cannot see that range/with rebind the dot.
func TestTemplateFieldRefs(t *testing.T) {
	cases := []struct {
		name string
		tpl  string
		want []string
	}{
		{"simple", "{{.User}}", []string{"User"}},
		{"literal_text_is_not_a_reference", "{{.User}}@Example.Com", []string{"User"}},
		{"sub_selector_reports_the_root_only", "{{.Addr.City}}", []string{"Addr"}},
		{"pipeline", "{{.A | lower}}", []string{"A"}},
		{"printf_args", `{{printf "%s-%s" .A .B}}`, []string{"A", "B"}},
		{"if_else_do_not_rebind_dot", "{{if .Flag}}{{.A}}{{else}}{{.B}}{{end}}", []string{"Flag", "A", "B"}},
		{"range_body_rebinds_dot", "{{range .Items}}{{.Name}}{{end}}", []string{"Items"}},
		{"with_body_rebinds_dot", "{{with .Addr}}{{.City}}{{else}}{{.Fallback}}{{end}}", []string{"Addr", "Fallback"}},
		{"variable_declaration", "{{$x := .A}}{{$x}}", []string{"A"}},
		{"plain_text", "no fields here .NotAField", nil},
		{"bare_dot", "{{.}}", nil},
		{"deduplicated", "{{.A}}-{{.A}}", []string{"A"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tpl, err := template.New("t").Funcs(templateHelpers).Parse(c.tpl)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			got := templateFieldRefs(tpl.Tree)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("templateFieldRefs(%q) = %v, want %v", c.tpl, got, c.want)
			}
		})
	}
}

func TestTopoSortTemplateFields(t *testing.T) {
	// build turns a declaration order plus a name -> dependency-names map into the
	// []templateField the sorter expects.
	build := func(order []string, deps map[string][]string) []templateField {
		pos := make(map[string]int, len(order))
		fields := make([]templateField, len(order))
		for i, n := range order {
			pos[n] = i
			fields[i] = templateField{index: i, name: n}
		}
		for _, n := range order {
			for _, d := range deps[n] {
				fields[pos[n]].deps = append(fields[pos[n]].deps, pos[d])
			}
		}
		return fields
	}

	t.Run("linear_chain_declared_in_reverse", func(t *testing.T) {
		fields := build([]string{"C", "B", "A"}, map[string][]string{
			"C": {"B"},
			"B": {"A"},
		})
		order, err := topoSortTemplateFields(fields)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var names []string
		for _, p := range order {
			names = append(names, fields[p].name)
		}
		if want := []string{"A", "B", "C"}; !reflect.DeepEqual(names, want) {
			t.Fatalf("order = %v, want %v", names, want)
		}
	})

	t.Run("diamond", func(t *testing.T) {
		fields := build([]string{"D", "B", "C", "A"}, map[string][]string{
			"D": {"B", "C"},
			"B": {"A"},
			"C": {"A"},
		})
		order, err := topoSortTemplateFields(fields)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order) != 4 {
			t.Fatalf("order = %v, want 4 entries", order)
		}
		seen := make(map[string]int)
		for i, p := range order {
			seen[fields[p].name] = i
		}
		for _, pair := range [][2]string{{"A", "B"}, {"A", "C"}, {"B", "D"}, {"C", "D"}} {
			if seen[pair[0]] > seen[pair[1]] {
				t.Errorf("%s must be evaluated before %s, got %v", pair[0], pair[1], seen)
			}
		}
	})

	t.Run("independent_chains_are_deterministic", func(t *testing.T) {
		fields := build([]string{"A", "B", "C", "D"}, map[string][]string{
			"B": {"A"},
			"D": {"C"},
		})
		first, err := topoSortTemplateFields(fields)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for range 20 {
			again, err := topoSortTemplateFields(fields)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(first, again) {
				t.Fatalf("order is not deterministic: %v then %v", first, again)
			}
		}
	})

	cycles := []struct {
		name    string
		order   []string
		deps    map[string][]string
		members []string
	}{
		{"self", []string{"A"}, map[string][]string{"A": {"A"}}, []string{"A"}},
		{"two_node", []string{"A", "B"}, map[string][]string{"A": {"B"}, "B": {"A"}}, []string{"A", "B"}},
		{"three_node", []string{"A", "B", "C"}, map[string][]string{"A": {"B"}, "B": {"C"}, "C": {"A"}}, []string{"A", "B", "C"}},
	}
	for _, c := range cycles {
		t.Run("cycle_"+c.name, func(t *testing.T) {
			_, err := topoSortTemplateFields(build(c.order, c.deps))
			if err == nil {
				t.Fatal("expected a cycle error")
			}
			for _, m := range c.members {
				if !strings.Contains(err.Error(), m) {
					t.Errorf("error %q does not name cycle member %q", err, m)
				}
			}
		})
	}
}

// TestTemplateTagWithUnexportedField is the repro from the review of
// go-faker/faker#83: building the template context used to call
// reflect.Value.Interface on every field, which panics for unexported ones.
func TestTemplateTagWithUnexportedField(t *testing.T) {
	t.Run("unexported_field_is_ignored", func(t *testing.T) {
		type PrivateField struct {
			password  string
			FirstName string `faker:"first_name"`
			Slug      string `faker:"template:{{.FirstName | slug}}"`
		}

		var u PrivateField
		if err := FakeData(&u); err != nil {
			t.Fatal(err)
		}
		if u.FirstName == "" {
			t.Fatal("expected non-empty FirstName")
		}
		if u.Slug == "" {
			t.Fatal("expected non-empty Slug")
		}
		if want := slugify(u.FirstName); u.Slug != want {
			t.Fatalf("Slug = %q, want %q", u.Slug, want)
		}
		if u.password != "" {
			t.Fatalf("unexported field was written: %q", u.password)
		}
	})

	t.Run("referencing_an_unexported_field_errors", func(t *testing.T) {
		type PrivateRef struct {
			password string
			Bad      string `faker:"template:{{.password}}"`
		}

		var u PrivateRef
		err := FakeData(&u)
		if err == nil {
			t.Fatal("expected an error when a template reads an unexported field")
		}
		if !strings.Contains(err.Error(), "password") {
			t.Fatalf("error %q does not mention the offending field", err)
		}
		_ = u.password
	})
}

// TestTemplateTagDependencyOrder covers references to template fields declared
// later in the struct, which used to be rejected outright.
func TestTemplateTagDependencyOrder(t *testing.T) {
	t.Run("forward_reference", func(t *testing.T) {
		type U struct {
			FirstName string `faker:"first_name"`
			Domain    string `faker:"domain_name"`
			Email     string `faker:"template:{{.Username}}@{{.Domain}}"`
			Username  string `faker:"template:{{.FirstName | lower}}"`
		}

		var u U
		if err := FakeData(&u); err != nil {
			t.Fatal(err)
		}
		if want := strings.ToLower(u.FirstName); u.Username != want {
			t.Fatalf("Username = %q, want %q", u.Username, want)
		}
		if want := u.Username + "@" + u.Domain; u.Email != want {
			t.Fatalf("Email = %q, want %q", u.Email, want)
		}
	})

	t.Run("three_level_chain_declared_in_reverse", func(t *testing.T) {
		type U struct {
			C    string `faker:"template:{{.B}}-c"`
			B    string `faker:"template:{{.A}}-b"`
			A    string `faker:"template:{{.Seed | lower}}"`
			Seed string `faker:"first_name"`
		}

		var u U
		if err := FakeData(&u); err != nil {
			t.Fatal(err)
		}
		wantA := strings.ToLower(u.Seed)
		if u.A != wantA || u.B != wantA+"-b" || u.C != wantA+"-b-c" {
			t.Fatalf("chain resolved incorrectly: A=%q B=%q C=%q (seed %q)", u.A, u.B, u.C, u.Seed)
		}
	})
}

func TestTemplateTagCycleDetection(t *testing.T) {
	t.Run("self", func(t *testing.T) {
		var u struct {
			A string `faker:"template:{{.A}}-x"`
		}
		err := FakeData(&u)
		if err == nil || !strings.Contains(err.Error(), "cycle") {
			t.Fatalf("expected a cycle error, got %v", err)
		}
	})

	t.Run("two_node", func(t *testing.T) {
		var u struct {
			A string `faker:"template:{{.B}}-x"`
			B string `faker:"template:{{.A}}-y"`
		}
		err := FakeData(&u)
		if err == nil {
			t.Fatal("expected a cycle error")
		}
		for _, name := range []string{"A", "B"} {
			if !strings.Contains(err.Error(), name) {
				t.Errorf("error %q does not name cycle member %q", err, name)
			}
		}
	})

	t.Run("three_node", func(t *testing.T) {
		var u struct {
			A string `faker:"template:{{.B}}"`
			B string `faker:"template:{{.C}}"`
			C string `faker:"template:{{.A}}"`
		}
		if err := FakeData(&u); err == nil {
			t.Fatal("expected a cycle error")
		}
	})
}

func TestTemplateTagNonStringField(t *testing.T) {
	type U struct {
		Name string `faker:"first_name"`
		Age  int    `faker:"template:{{.Name}}"`
		Ok   string `faker:"template:{{.Name | upper}}"`
	}

	var u U
	err := FakeData(&u)
	if err == nil {
		t.Fatal("expected an error for a template tag on an int field")
	}
	if !strings.Contains(err.Error(), "Age") {
		t.Fatalf("error %q does not name the offending field", err)
	}
	// Validation must happen before anything is evaluated, so the sibling template
	// field must not have been written.
	if u.Ok != "" {
		t.Fatalf("sibling template field was evaluated despite the error: %q", u.Ok)
	}
}

func TestTemplateTagPointerAndNamedStringField(t *testing.T) {
	type Slug string
	type U struct {
		Name    string  `faker:"first_name"`
		Slug    Slug    `faker:"template:{{.Name | slug}}"`
		Pointer *string `faker:"template:{{.Name | upper}}"`
	}

	var u U
	if err := FakeData(&u); err != nil {
		t.Fatal(err)
	}
	if want := Slug(slugify(u.Name)); u.Slug != want {
		t.Fatalf("Slug = %q, want %q", u.Slug, want)
	}
	if u.Pointer == nil {
		t.Fatal("expected Pointer to be allocated")
	}
	if want := strings.ToUpper(u.Name); *u.Pointer != want {
		t.Fatalf("*Pointer = %q, want %q", *u.Pointer, want)
	}
}

func TestTemplateTagWithUnique(t *testing.T) {
	type U struct {
		Name string `faker:"first_name"`
		Slug string `faker:"template:{{.Name | lower}},unique"`
	}

	var u U
	err := FakeData(&u)
	if err == nil {
		t.Fatal("expected unique+template to be rejected")
	}
	if want := fmt.Sprintf(fakerErrors.ErrTemplateWithUnique, "Slug"); err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
	// The old behavior silently dropped the tag and filled the field with a random
	// string; make sure that has not come back.
	if u.Slug != "" {
		t.Fatalf("Slug was written despite the error: %q", u.Slug)
	}
}

func TestTemplateTagWithKeep(t *testing.T) {
	type U struct {
		Name string `faker:"first_name"`
		Slug string `faker:"template:{{.Name | lower}},keep"`
	}

	t.Run("non_zero_value_is_kept", func(t *testing.T) {
		u := U{Slug: "preset"}
		if err := FakeData(&u); err != nil {
			t.Fatal(err)
		}
		if u.Slug != "preset" {
			t.Fatalf("Slug = %q, want the preset value to be kept", u.Slug)
		}
	})

	t.Run("zero_value_is_templated", func(t *testing.T) {
		var u U
		if err := FakeData(&u); err != nil {
			t.Fatal(err)
		}
		if want := strings.ToLower(u.Name); u.Slug != want {
			t.Fatalf("Slug = %q, want %q", u.Slug, want)
		}
	})
}

func TestTemplateTagWithCommaInBody(t *testing.T) {
	type U struct {
		First string `faker:"first_name"`
		Last  string `faker:"last_name"`
		Full  string `faker:"template:{{.Last}}, {{.First}}"`
	}

	var u U
	if err := FakeData(&u); err != nil {
		t.Fatal(err)
	}
	if want := u.Last + ", " + u.First; u.Full != want {
		t.Fatalf("Full = %q, want %q", u.Full, want)
	}
}

func TestTemplateTagWithOptionOverrides(t *testing.T) {
	type U struct {
		Name    string `faker:"first_name"`
		Slug    string `faker:"template:{{.Name | lower}}"`
		Derived string `faker:"template:{{.Slug}}-x"`
	}

	t.Run("ignored_field_wins_and_dependents_see_it", func(t *testing.T) {
		var u U
		if err := FakeData(&u, options.WithFieldsToIgnore("Slug")); err != nil {
			t.Fatal(err)
		}
		if u.Slug != "" {
			t.Fatalf("Slug = %q, want it left at the zero value", u.Slug)
		}
		if u.Derived != "-x" {
			t.Fatalf("Derived = %q, want %q", u.Derived, "-x")
		}
	})

	t.Run("field_provider_wins_and_dependents_see_it", func(t *testing.T) {
		var u U
		err := FakeData(&u, options.WithCustomFieldProvider("Slug", func() (any, error) {
			return "provided", nil
		}))
		if err != nil {
			t.Fatal(err)
		}
		if u.Slug != "provided" {
			t.Fatalf("Slug = %q, want %q", u.Slug, "provided")
		}
		if u.Derived != "provided-x" {
			t.Fatalf("Derived = %q, want %q", u.Derived, "provided-x")
		}
	})

	t.Run("only_zero_fields_keeps_a_preset_template_field", func(t *testing.T) {
		u := U{Slug: "preset"}
		if err := FakeData(&u, options.WithOnlyZeroFields()); err != nil {
			t.Fatal(err)
		}
		if u.Slug != "preset" {
			t.Fatalf("Slug = %q, want the preset value", u.Slug)
		}
		if u.Derived != "preset-x" {
			t.Fatalf("Derived = %q, want %q", u.Derived, "preset-x")
		}
	})
}

func TestTemplateTagInNestedStruct(t *testing.T) {
	type Inner struct {
		Name string `faker:"first_name"`
		Slug string `faker:"template:{{.Name | slug}}"`
	}

	t.Run("nested_struct_and_pointer", func(t *testing.T) {
		type Outer struct {
			A Inner
			B *Inner
		}
		var o Outer
		if err := FakeData(&o); err != nil {
			t.Fatal(err)
		}
		if o.A.Slug != slugify(o.A.Name) {
			t.Errorf("A: Slug = %q, want %q", o.A.Slug, slugify(o.A.Name))
		}
		if o.B == nil {
			t.Fatal("expected B to be allocated")
		}
		if o.B.Slug != slugify(o.B.Name) {
			t.Errorf("B: Slug = %q, want %q", o.B.Slug, slugify(o.B.Name))
		}
	})

	t.Run("slice_elements_are_independent", func(t *testing.T) {
		type Outer struct {
			Items []Inner `faker:"slice_len=5"`
		}
		var o Outer
		if err := FakeData(&o); err != nil {
			t.Fatal(err)
		}
		if len(o.Items) != 5 {
			t.Fatalf("len(Items) = %d, want 5", len(o.Items))
		}
		for i, it := range o.Items {
			if it.Slug != slugify(it.Name) {
				t.Errorf("Items[%d]: Slug = %q, want %q derived from its own Name %q",
					i, it.Slug, slugify(it.Name), it.Name)
			}
		}
	})
}

func TestTemplateTagWithEmbeddedStruct(t *testing.T) {
	type Address struct {
		City string `faker:"word"`
	}
	type U struct {
		Address
		Promoted string `faker:"template:{{.City}}"`
		Explicit string `faker:"template:{{.Address.City}}"`
	}

	var u U
	if err := FakeData(&u); err != nil {
		t.Fatal(err)
	}
	if u.City == "" {
		t.Fatal("expected the embedded City to be generated")
	}
	if u.Promoted != u.City {
		t.Errorf("Promoted = %q, want the promoted City %q", u.Promoted, u.City)
	}
	if u.Explicit != u.City {
		t.Errorf("Explicit = %q, want %q", u.Explicit, u.City)
	}
}

func TestTemplateTagErrors(t *testing.T) {
	cases := []struct {
		name  string
		field string
		fake  func() error
	}{
		{
			name:  "missing_field",
			field: "Bad",
			fake: func() error {
				var v struct {
					Name string `faker:"first_name"`
					Bad  string `faker:"template:{{.DoesNotExist}}"`
				}
				return FakeData(&v)
			},
		},
		{
			name:  "invalid_syntax",
			field: "Bad",
			fake: func() error {
				var v struct {
					Bad string `faker:"template:{{.Unclosed"`
				}
				return FakeData(&v)
			},
		},
		{
			name:  "empty_body",
			field: "Bad",
			fake: func() error {
				var v struct {
					Bad string `faker:"template:"`
				}
				return FakeData(&v)
			},
		},
		{
			name:  "unknown_helper",
			field: "Bad",
			fake: func() error {
				var v struct {
					Name string `faker:"first_name"`
					Bad  string `faker:"template:{{.Name | nosuchhelper}}"`
				}
				return FakeData(&v)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.fake()
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Error(), c.field) {
				t.Fatalf("error %q does not name the offending field %q", err, c.field)
			}
		})
	}
}

func TestTemplateTagWithTagName(t *testing.T) {
	type U struct {
		Name string `fake:"first_name"`
		Slug string `fake:"template:{{.Name | lower}}"`
	}

	var u U
	if err := FakeData(&u, options.WithTagName("fake")); err != nil {
		t.Fatal(err)
	}
	if u.Name == "" {
		t.Fatal("expected Name to be generated")
	}
	if want := strings.ToLower(u.Name); u.Slug != want {
		t.Fatalf("Slug = %q, want %q", u.Slug, want)
	}
}

func TestTemplateTagNonStringValueInTemplate(t *testing.T) {
	type U struct {
		Age   int    `faker:"boundary_start=30, boundary_end=31"`
		Label string `faker:"template:age-{{.Age}}"`
		Slug  string `faker:"template:{{.Age | slug}}"`
	}

	var u U
	if err := FakeData(&u); err != nil {
		t.Fatal(err)
	}
	if want := fmt.Sprintf("age-%d", u.Age); u.Label != want {
		t.Fatalf("Label = %q, want %q", u.Label, want)
	}
	if want := fmt.Sprintf("%d", u.Age); u.Slug != want {
		t.Fatalf("Slug = %q, want %q", u.Slug, want)
	}
}

func TestTemplateTagWithMaxFieldDepth(t *testing.T) {
	type U struct {
		Name string `faker:"first_name"`
		Slug string `faker:"template:{{.Name | lower}}"`
	}

	var u U
	if err := FakeData(&u, options.WithMaxFieldDepthOption(0)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Slug != "" {
		t.Fatalf("Slug = %q, want empty at depth 0", u.Slug)
	}
}

func TestTemplateTagConcurrent(t *testing.T) {
	type U struct {
		First    string `faker:"first_name"`
		Last     string `faker:"last_name"`
		Domain   string `faker:"domain_name"`
		Email    string `faker:"template:{{.Username}}@{{.Domain}}"`
		Username string `faker:"template:{{.First | lower}}.{{.Last | lower}}"`
	}

	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				var u U
				if err := FakeData(&u); err != nil {
					t.Error(err)
					return
				}
				want := strings.ToLower(u.First) + "." + strings.ToLower(u.Last)
				if u.Username != want {
					t.Errorf("Username = %q, want %q", u.Username, want)
					return
				}
				if u.Email != u.Username+"@"+u.Domain {
					t.Errorf("Email = %q, want %q", u.Email, u.Username+"@"+u.Domain)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestSlugify(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Hello, World!", "hello-world"},
		{"  --Hello--  ", "hello"},
		{"already-a-slug", "already-a-slug"},
		{"José", "josé"},
		{"Иван", "иван"},
		{"李雷", "李雷"},
		{"a1b2", "a1b2"},
		{"!!!", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := slugify(c.in); got != c.want {
			t.Errorf("slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// abc is the sample value shared by the TestTemplateString cases.
const abc = "abc"

// testStringer exercises the fmt.Stringer branch of templateString.
type testStringer struct{ s string }

func (t testStringer) String() string { return t.s }

func TestTemplateString(t *testing.T) {
	type named string
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"nil", nil, ""},
		{"string", abc, abc},
		{"named_string", named(abc), abc},
		{"int", 42, "42"},
		{"float", 1.5, "1.5"},
		{"bool", true, "true"},
		{"stringer", testStringer{abc}, abc},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := templateString(c.in); got != c.want {
				t.Fatalf("templateString(%v) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestTemplateHelpers(t *testing.T) {
	type named string
	type U struct {
		S string `faker:"-"`
		N named  `faker:"-"`
		I int    `faker:"-"`

		Lower string `faker:"template:{{.S | lower}}"`
		Upper string `faker:"template:{{.S | upper}}"`
		Trim  string `faker:"template:{{.S | trim}}"`
		Named string `faker:"template:{{.N | upper}}"`
		Num   string `faker:"template:{{.I | lower}}"`
	}

	u := U{S: "  MiXeD  ", N: named("ab"), I: 7}
	if err := FakeData(&u); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ name, got, want string }{
		{"lower", u.Lower, "  mixed  "},
		{"upper", u.Upper, "  MIXED  "},
		{"trim", u.Trim, "MiXeD"},
		{"upper on a named string", u.Named, "AB"},
		{"lower on an int", u.Num, "7"},
	} {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func BenchmarkFakerDataTemplateTag(b *testing.B) {
	type U struct {
		First    string `faker:"first_name"`
		Last     string `faker:"last_name"`
		Domain   string `faker:"domain_name"`
		Username string `faker:"template:{{.First | lower}}.{{.Last | lower}}"`
		Slug     string `faker:"template:{{.First | slug}}-{{.Last | slug}}"`
		Email    string `faker:"template:{{.Username}}@{{.Domain}}"`
	}
	for b.Loop() {
		var u U
		if err := FakeData(&u); err != nil {
			b.Fatal(err)
		}
	}
}
