package faker

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

// fieldRefRegexp matches references to struct fields in a template body, e.g. `.FieldName`.
// Field names in Go are exported (start with uppercase), so we only consider those.
var fieldRefRegexp = regexp.MustCompile(`\.([A-Z]\w*)`)

// EvaluateTemplateFields evaluates the template-tagged fields for a struct value `v` of type `t`.
// tagName is the field tag name to use (typically opts.TagName).
func EvaluateTemplateFields(t reflect.Type, v reflect.Value, templateFields []int, tagName string) error {
	if len(templateFields) == 0 {
		return nil
	}
	// build context map from current struct values
	// this will be updated as we compute each template field so
	// later templates can reference earlier template results.
	ctx := make(map[string]any)
	for i := 0; i < v.NumField(); i++ {
		ctx[t.Field(i).Name] = v.Field(i).Interface()
	}

	// Build the set of all template-field names and track which have been
	// evaluated so we can detect forward / circular references.
	templateFieldNames := make(map[string]bool, len(templateFields))
	for _, idx := range templateFields {
		templateFieldNames[t.Field(idx).Name] = true
	}
	evaluated := make(map[string]bool, len(templateFields))

	for _, idx := range templateFields {
		tags := decodeTags(t, idx, tagName)
		tpl := tags.fieldType
		lower := strings.ToLower(tpl)
		if strings.HasPrefix(lower, TemplateTag+":") {
			tpl = tpl[len(TemplateTag)+1:]
		}
		fieldName := t.Field(idx).Name
		// Reject forward / circular references to other template fields.
		for _, m := range fieldRefRegexp.FindAllStringSubmatch(tpl, -1) {
			ref := m[1]
			if ref == fieldName {
				return fmt.Errorf("template field %q references itself", fieldName)
			}
			if templateFieldNames[ref] && !evaluated[ref] {
				return fmt.Errorf("template field %q references template field %q which has not been evaluated yet (forward or circular reference)", fieldName, ref)
			}
		}
		res, err := evaluateTemplateForField(tpl, ctx)
		if err != nil {
			return err
		}
		if v.Field(idx).Kind() != reflect.String {
			return fmt.Errorf("template tag only supported on string fields: %s", t.Field(idx).Name)
		}

		v.Field(idx).SetString(res)
		//  update context so subsequent template fields can reference it
		ctx[fieldName] = res
		evaluated[fieldName] = true
	}
	return nil
}

// evaluateTemplateForField evaluates a template string tpl with context ctx (map of field name -> value)
// and returns the resulting string. Only a small set of helper funcs are exposed.
func evaluateTemplateForField(tpl string, ctx map[string]any) (string, error) {
	t, err := template.New("faker_template").Option("missingkey=error").Funcs(helpers()).Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func helpers() template.FuncMap {
	return template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"slug":  func(s string) string { return slugify(s) },
		//can add more
	}
}

// slugify is a tiny, permissive slug function: lowercases and replaces non-alphanumeric with '-'
func slugify(s string) string {
	var b strings.Builder
	s = strings.ToLower(s)
	lastDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	res := b.String()
	res = strings.Trim(res, "-")
	return res
}
