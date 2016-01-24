package placeholders

import (
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"strings"
)


// Placeholder groups the basic information to work with placeholders
type Placeholder struct { // ${DATASOURCE_URL:jdbc:mysql://localhost:3306/shcema?profileSQL=true}
	PropertyName  string `json:"property_name"`  // the full path to the property datasource.url
	PropertyValue string `json:"property_value"` // jdbc:mysql://localhost:3306/shcema?profileSQL=true
	PlaceHolder   string `json:"placeholder"`    // DATASOURCE_URL
}

// Placeholders is a collection of Placeholder
type Placeholders struct {
	Placeholders []*Placeholder `json:"placeholders"`
}

// GetPlaceHolders uses the provided viper configuration to extract properties that have placeholders in is values
func GetPlaceHolders(conf *viper.Viper) Placeholders {
	list := parseMap(conf.AllSettings())

	properties := CreateProperties(list)

	return Placeholders{properties}
}

// CreateProperties transform the propsMap into a Property struct slice
func CreateProperties(propsMap map[string]string) []*Placeholder {
	count := len(propsMap)

	ps := make([]*Placeholder, count)
	i := 0
	for k, v := range propsMap {
		p, d, err := extractPlaceholder(v)
		if err == nil {
			p := &Placeholder{k, d, p}
			ps[i] = p
		}

		i++
	}

	return ps
}

func extractPlaceholder(s string) (string, string, error) {
	if s[:2] != "${" {
		return "", "", fmt.Errorf("%s does not contain any placeholder with format ${PLACEHOLER_VARIABLE[:defaultvalue]}", s)
	}

	if s[len(s)-1:len(s)] != "}" {
		return "", "", fmt.Errorf("%s does not contain any placeholder with format ${PLACEHOLER_VARIABLE[:defaultvalue]}", s)
	}

	s = s[2:]
	s = s[0 : len(s)-1]

	defaultValue := ""
	if strings.Contains(s, ":") {
		dv := strings.Split(s, ":")
		defaultValue = strings.Join(dv[1:], ":")
	}

	return strings.Split(s, ":")[0], defaultValue, nil
}

func parseMap(props map[string]interface{}) map[string]string {
	list := make(map[string]string)
	for key, value := range props {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			l := parseMapInterface(v, key, list)
			for pkey, pvalue := range l {
				list[pkey] = pvalue
			}
		case string:
			if v[:2] == "${" {
				list[key] = v
			}
		default:
		}
	}
	return list
}

func parseMapInterface(props map[interface{}]interface{}, key string, list map[string]string) map[string]string {
	for k, value := range props {
		actKey := key + "." + fmt.Sprint(k)

		switch v := value.(type) {
		case map[interface{}]interface{}:
			list = parseMapInterface(v, actKey, list)
		case string:
			if v[:2] == "${" {
				keystr := fmt.Sprint(actKey) // <-- HACK
				list[keystr] = v
			}
		default:
		}
	}
	return list
}
