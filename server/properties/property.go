package properties
import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/viper"
	"fmt"
	"strings"
)

type Property struct {
	PropertyName string `json:"property_name"` // the full path to the property datasource.url
	PropertyValue string `json:"property_value"` // ${DATASOURCE_URL:jdbc:mysql://localhost:3306/shcema?profileSQL=true}
	PlaceHolder string `json:"placeholder"`// DATASOURCE_URL
}

type Properties struct {
	Properties []*Property `json:"properties"`
}

// GetPlaceHolders uses the provided viper configuration to extract properties that have placeholders in is values
func GetPlaceHolders(conf *viper.Viper) Properties {
	list := parseMap(conf.AllSettings())

	properties := CreateProperties(list)

	return Properties{properties}
}

// CreatePropertiesGiven transform the propsMap into a Property struct slice
func CreateProperties(propsMap map[string]string) []*Property{
	count := len(propsMap)

	ps := make([]*Property, count)
	i := 0
	for k, v := range propsMap {
		p, err := extractPlaceholder(v)
		if (err == nil){
			p := &Property{k, v, p}
			ps[i] = p
		}

		i++
	}

	return ps
}

func extractPlaceholder(s string) (string, error){
	if s[:2] != "${" {
		return "", fmt.Errorf("%s does not contain any placeholder with format ${PLACEHOLER_VARIABLE[:defaultvalue]}", s)
	}

	if s[len(s)-1:len(s)] != "}" {
		return "", fmt.Errorf("%s does not contain any placeholder with format ${PLACEHOLER_VARIABLE[:defaultvalue]}", s)
	}

	s = s[2:]
	s = s[0:len(s)-1]

	return strings.Split(s, ":")[0], nil
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