package faker

// The template tag (`faker:"template:{{.Other}}"`) derives a field's value from
// other fields of the same struct. Because a template reads values that faker
// itself generates, template fields are filled in a second pass that runs after
// every other field of the struct is done — see getFakedValueForStruct.
//
// Within that second pass, template fields are evaluated in dependency order
// rather than declaration order, so a template may reference another template
// field declared later. Only a true cycle is an error.
//
// TemplateTag is deliberately absent from PriorityTags: it is not a provider, and
// listing it there would route the tag through setDataWithTag.

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"text/template/parse"
	"unicode"

	fakerErrors "github.com/go-faker/faker/v4/pkg/errors"
)

// templateField is one parsed, validated template-tagged field of a struct.
type templateField struct {
	index int    // index of the field within the struct
	name  string // field name, used in error messages
	tpl   *template.Template
	deps  []int // positions within templatePlan.fields this field depends on
}

// templatePlan is the validated evaluation plan for the template fields of one
// struct type.
type templatePlan struct {
	fields []templateField // declaration order
	order  []int           // positions within fields, in dependency order
}

// evaluateTemplateFields fills the template-tagged fields of v listed in fields.
// v must be an addressable struct value of type t. fields holds struct field
// indices collected by the first pass; a template field missing from it was
// handled by a higher-precedence option (IgnoreFields, OnlyZeroFields,
// FieldProviders, keep) and is left alone.
func evaluateTemplateFields(t reflect.Type, v reflect.Value, fields []int, tagName string) error {
	if len(fields) == 0 {
		return nil
	}
	if !v.CanAddr() {
		return errors.New(fakerErrors.ErrValueNotPtr)
	}
	plan, err := buildTemplatePlan(t, tagName)
	if err != nil {
		return err
	}
	// Templates read the struct itself rather than a map built by reflection. That
	// is what makes unexported fields safe: text/template reports them as an error
	// instead of us panicking on reflect.Value.Interface. Taking the address means
	// each template also sees the results written by earlier ones.
	data := v.Addr().Interface()
	var buf bytes.Buffer
	for _, p := range plan.order {
		tf := &plan.fields[p]
		if !containsFieldIndex(fields, tf.index) {
			continue
		}
		buf.Reset()
		if err := tf.tpl.Execute(&buf, data); err != nil {
			return fmt.Errorf(fakerErrors.ErrTemplateEvaluate, tf.name, err)
		}
		setTemplateResult(v.Field(tf.index), buf.String())
	}
	return nil
}

func containsFieldIndex(fields []int, idx int) bool {
	for _, f := range fields {
		if f == idx {
			return true
		}
	}
	return false
}

// buildTemplatePlan parses and validates every template-tagged field of t and
// works out the order they must be evaluated in.
func buildTemplatePlan(t reflect.Type, tagName string) (*templatePlan, error) {
	plan := &templatePlan{}
	posByName := make(map[string]int)

	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		tags := decodeTags(t, i, tagName)
		if !tags.isTemplate {
			continue
		}
		// Validate the destination before anything is evaluated, so a bad tag can
		// never leave some fields written and others not.
		if !templateAssignable(f.Type) {
			return nil, fmt.Errorf(fakerErrors.ErrTemplateNotStringField, f.Name, f.Type)
		}
		if strings.TrimSpace(tags.template) == "" {
			return nil, fmt.Errorf(fakerErrors.ErrTemplateEmptyBody, f.Name)
		}
		tpl, err := template.New(f.Name).Funcs(templateHelpers).Parse(tags.template)
		if err != nil {
			return nil, fmt.Errorf(fakerErrors.ErrTemplateParse, f.Name, err)
		}
		posByName[f.Name] = len(plan.fields)
		plan.fields = append(plan.fields, templateField{index: i, name: f.Name, tpl: tpl})
	}
	if len(plan.fields) == 0 {
		return plan, nil
	}

	// Only references to other template fields become edges. References to ordinary
	// fields need no ordering: the first pass already filled them.
	for p := range plan.fields {
		for _, ref := range templateFieldRefs(plan.fields[p].tpl.Tree) {
			if q, ok := posByName[ref]; ok {
				plan.fields[p].deps = append(plan.fields[p].deps, q)
			}
		}
	}

	order, err := topoSortTemplateFields(plan.fields)
	if err != nil {
		return nil, err
	}
	plan.order = order
	return plan, nil
}

// templateAssignable reports whether the template result can be stored in a field
// of type t. Named string types work through SetString; *string is allocated.
func templateAssignable(t reflect.Type) bool {
	return t.Kind() == reflect.String ||
		(t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.String)
}

func setTemplateResult(f reflect.Value, s string) {
	if f.Kind() == reflect.Pointer {
		p := reflect.New(f.Type().Elem())
		p.Elem().SetString(s)
		f.Set(p)
		return
	}
	f.SetString(s)
}

// templateFieldRefs returns the names of the struct fields a template reads, in
// source order and deduplicated.
//
// It walks the parse tree rather than matching the template source with a regular
// expression: only the tree can tell a field reference from literal text such as
// "@Example.Com", can report `{{.Addr.City}}` as a reference to Addr alone, and
// can tell that the dot inside a `range` or `with` body is the loop value rather
// than the struct.
func templateFieldRefs(tree *parse.Tree) []string {
	if tree == nil || tree.Root == nil {
		return nil
	}
	var out []string
	seen := make(map[string]struct{})
	add := func(name string) {
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			out = append(out, name)
		}
	}
	var walk func(parse.Node)
	// walkBranch handles if/range/with. The controlling pipe is always evaluated
	// with the outer dot; range and with rebind dot for their body, so the body
	// cannot reference our struct's fields. An else body runs with the outer dot.
	walkBranch := func(b *parse.BranchNode, rebindsDot bool) {
		walk(b.Pipe)
		if !rebindsDot && b.List != nil {
			walk(b.List)
		}
		if b.ElseList != nil {
			walk(b.ElseList)
		}
	}
	walk = func(n parse.Node) {
		switch n := n.(type) {
		case *parse.ListNode:
			if n == nil {
				return
			}
			for _, c := range n.Nodes {
				walk(c)
			}
		case *parse.ActionNode:
			walk(n.Pipe)
		case *parse.PipeNode:
			if n == nil {
				return
			}
			for _, d := range n.Decl {
				walk(d)
			}
			for _, c := range n.Cmds {
				walk(c)
			}
		case *parse.CommandNode:
			for _, a := range n.Args {
				walk(a)
			}
		case *parse.FieldNode:
			if len(n.Ident) > 0 {
				add(n.Ident[0])
			}
		case *parse.ChainNode:
			if _, isDot := n.Node.(*parse.DotNode); isDot && len(n.Field) > 0 {
				add(n.Field[0])
			}
			walk(n.Node)
		case *parse.IfNode:
			walkBranch(&n.BranchNode, false)
		case *parse.RangeNode:
			walkBranch(&n.BranchNode, true)
		case *parse.WithNode:
			walkBranch(&n.BranchNode, true)
		case *parse.TemplateNode:
			walk(n.Pipe)
		}
	}
	walk(tree.Root)
	return out
}

// topoSortTemplateFields returns the positions of fields in dependency order:
// every field appears after the template fields it references. Roots are visited
// in declaration order so the result is deterministic. A cycle, including a field
// that references itself, is an error.
func topoSortTemplateFields(fields []templateField) ([]int, error) {
	const (
		white = iota
		grey
		black
	)
	state := make([]int8, len(fields))
	order := make([]int, 0, len(fields))
	stack := make([]int, 0, len(fields))

	var visit func(p int) error
	visit = func(p int) error {
		switch state[p] {
		case black:
			return nil
		case grey:
			return fmt.Errorf(fakerErrors.ErrTemplateCycle, cyclePath(fields, stack, p))
		}
		state[p] = grey
		stack = append(stack, p)
		for _, d := range fields[p].deps {
			if err := visit(d); err != nil {
				return err
			}
		}
		stack = stack[:len(stack)-1]
		state[p] = black
		order = append(order, p)
		return nil
	}
	for p := range fields {
		if err := visit(p); err != nil {
			return nil, err
		}
	}
	return order, nil
}

// cyclePath renders a cycle as "A -> B -> A", starting from where p first appears
// on the visit stack.
func cyclePath(fields []templateField, stack []int, p int) string {
	start := 0
	for i, s := range stack {
		if s == p {
			start = i
			break
		}
	}
	names := make([]string, 0, len(stack)-start+1)
	for _, s := range stack[start:] {
		names = append(names, fields[s].name)
	}
	names = append(names, fields[p].name)
	return strings.Join(names, " -> ")
}

// templateHelpers is read-only after initialization and shared by every template.
var templateHelpers = template.FuncMap{
	"lower": func(v any) string { return strings.ToLower(templateString(v)) },
	"upper": func(v any) string { return strings.ToUpper(templateString(v)) },
	"trim":  func(v any) string { return strings.TrimSpace(templateString(v)) },
	"slug":  func(v any) string { return slugify(templateString(v)) },
}

// templateString coerces a template argument to a string so the helpers also work
// on named string types, numbers and fmt.Stringer values.
func templateString(v any) string {
	switch s := v.(type) {
	case nil:
		return ""
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.String {
		return rv.String()
	}
	return fmt.Sprint(v)
}

// slugify lowercases s and replaces each run of non-alphanumeric characters with
// a single dash, trimming dashes from both ends. Letters and digits of any script
// are kept, so names generated by the Chinese and Russian providers survive.
func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	b.Grow(len(s))
	lastDash := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case !lastDash:
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
