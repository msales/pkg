package stats

// Tags is a map of strings that adds some convenience functionality over the vanilla go map.
// Deprecated.
// Warning! There is a known bug where Tags passed to any stats func (like Inc, Gauge, etc.) will
// not be merge correctly into the global tags. This might result in the stats loss! Do not pass
// Tags object to those functions, use the variadic list instead..
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

func (t Tags) toArray() []interface{} {
	arr := make([]interface{}, len(t)*2)

	i := 0
	for k, v := range t {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}
