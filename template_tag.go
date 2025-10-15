package faker

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

// EvaluateTemplateFields evaluates the template-tagged fields for a struct value `v` of type `t`.
// tagName is the field tag name to use (typically opts.TagName).
func EvaluateTemplateFields(t reflect.Type, v reflect.Value, templateFields []int, tagName string) error {
	if len(templateFields) == 0 {
		return nil
	}
	// build context map from current struct values
	ctx := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		ctx[t.Field(i).Name] = v.Field(i).Interface()
	}

	for _, idx := range templateFields {
		tags := decodeTags(t, idx, tagName)
		tpl := tags.fieldType
		lower := strings.ToLower(tpl)
		if strings.HasPrefix(lower, TemplateTag+":") {
			tpl = tpl[len(TemplateTag)+1:]
		}
		res, err := evaluateTemplateForField(tpl, ctx)
		if err != nil {
			return err
		}
		if v.Field(idx).Kind() != reflect.String {
			return fmt.Errorf("template tag only supported on string fields: %s", t.Field(idx).Name)
		}
		v.Field(idx).SetString(res)
	}
	return nil
}

// evaluateTemplateForField evaluates a template string tpl with context ctx (map of field name -> value)
// and returns the resulting string. Only a small set of helper funcs are exposed.
func evaluateTemplateForField(tpl string, ctx map[string]interface{}) (string, error) {
	t, err := template.New("faker_template").Funcs(helpers()).Parse(tpl)
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
