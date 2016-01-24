package placeholders

type Value struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Values struct {
	Values []*Value `json:"values"`
}

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
