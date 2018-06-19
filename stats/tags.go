package stats

// Tags is a map of strings that adds some convenience functionality over the vanilla go map.
type Tags map[string]string

// With adds the value to the key.
func (t Tags) With(key, value string) Tags {
	t[key] = value

	return t
}

// Merge adds the values from t2 to t, overriding the conflicting keys.
func (t Tags) Merge(t2 Tags) Tags {
	for k, v := range t2 {
		t[k] = v
	}

	return t
}
