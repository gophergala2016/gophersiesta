package placeholders

// Value is a pair of name and his value
type Value struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Values is a collection of Value
type Values struct {
	Values []*Value `json:"values"`
}

// CreateValues transforms a map of string to Values struct
func CreateValues(m map[string]string) Values {
	values := make([]*Value, len(m))
	i := 0
	for k, v := range m {
		value := &Value{}
		value.Name = k
		value.Value = v
		values[i] = value
		i++
	}

	return Values{values}
}
