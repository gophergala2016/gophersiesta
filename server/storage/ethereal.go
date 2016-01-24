package storage

// Ethereal is a Volatile store backed up by a map
type Ethereal struct {
	options map[string]map[string]map[string]string
}

// Init initializes the map that stores the values
func (e *Ethereal) Init() {
	e.options = make(map[string]map[string]map[string]string)
}

// GetLabels return all labels for a given appname
func (e *Ethereal) GetLabels(appname string) []string {
	lbls := make([]string, 0)

	if e.options[appname] != nil {
		for k := range e.options[appname] {
			lbls = append(lbls, k)
		}
	}

	return lbls
}

// SetOption stores a placeholders value for a given appname, and label in the storage engine
func (e *Ethereal) SetOption(appname, label, variable, value string) {

	if label == "" {
		label = "default"
	}

	if e.options[appname] == nil {
		e.options[appname] = make(map[string]map[string]string)
	}
	if e.options[appname][label] == nil {
		e.options[appname][label] = make(map[string]string)
	}

	e.options[appname][label][variable] = value

}

// GetOption returns a placeholders value for a given appname, and label in the storage engine
func (e *Ethereal) GetOption(appname, label, variable string) string {

	if label == "" {
		label = "default"
	}

	if e.options[appname] == nil {
		e.options[appname] = make(map[string]map[string]string)
	}
	if e.options[appname][label] == nil {
		e.options[appname][label] = make(map[string]string)
	}

	return e.options[appname][label][variable]

}

// GetOptions returns a map of placeholders value for a given appname, and label in the storage engine
func (e *Ethereal) GetOptions(appname, label string) map[string]string {

	if label == "" {
		label = "default"
	}

	if e.options[appname] == nil {
		e.options[appname] = make(map[string]map[string]string)
	}
	if e.options[appname][label] == nil {
		e.options[appname][label] = make(map[string]string)
	}

	return e.options[appname][label]

}

// Close shutdowns the storage
func (e *Ethereal) Close() {
	// TODO maybe we can write it to disck
}
