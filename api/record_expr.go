package api

import (
	"fmt"
	forcequery "github.com/ForceCLI/force/lib/query"
	"github.com/expr-lang/expr"
)

// Record wraps a forcequery.Record for convenient field access.
type Record struct {
	forcequery.Record
}

// StringValue returns the string value at the given expression path within the record.
func (r Record) StringValue(path string) (string, error) {
	v, err := r.Value(path)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("field %q is not a string (got %T)", path, v)
	}
	return s, nil
}

// Value evaluates the given expression path against the record and returns the value.
func (r Record) Value(path string) (interface{}, error) {
	env := r.toEnv()
	v, err := expr.Eval(path, env)
	if err != nil {
		return nil, fmt.Errorf("evaluating %q: %w", path, err)
	}
	return v, nil
}

// toEnv converts the record's fields into a nested map for expression evaluation.
func (r Record) toEnv() map[string]interface{} {
	env := make(map[string]interface{}, len(r.Fields))
	for k, v := range r.Fields {
		switch val := v.(type) {
		case forcequery.Record:
			env[k] = Record{val}.toEnv()
		case []forcequery.Record:
			arr := make([]interface{}, len(val))
			for i, rec := range val {
				arr[i] = Record{rec}.toEnv()
			}
			env[k] = arr
		default:
			env[k] = val
		}
	}
	return env
}
