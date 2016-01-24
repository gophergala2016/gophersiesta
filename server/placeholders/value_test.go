package placeholders

import (
	"fmt"
	"testing"
)

func TestCreateValuesEmpty(t *testing.T) {
	m := make(map[string]string)

	vs := CreateValues(m)

	if len(vs.Values) != 0 {
		t.Error("For an empty map the values slice must be also empty")
	}
}

func TestCreateValues(t *testing.T) {
	m := make(map[string]string)

	m["data.source"] = "value_data_source"
	m["url"] = "value_url"

	vs := CreateValues(m)

	if len(vs.Values) != 2 {
		t.Fatal("For a map with 2 keys the lenght of the values must be 2")
	}

	checkValue(&vs, "data.source", "value_data_source", t)
	checkValue(&vs, "url", "value_url", t)

}

func checkValue(vs *Values, key string, value string, t *testing.T) {
	v, err := vs.findValue(key)

	if err != nil {
		t.Fatal(err)
	}

	if v.Value != value {
		t.Fatal(fmt.Errorf("The value for key %s is %s and not the expected one %s", key, v.Value, value))
	}
}

func (vs *Values) findValue(k string) (*Value, error) {
	for _, v := range vs.Values {

		if v.Name == k {
			return v, nil
		}
	}
	return nil, fmt.Errorf("No value found for key %s", k)
}
