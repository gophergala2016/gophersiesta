package storage


type Ethereal struct {
	options map[string]map[string]map[string]string
}

func (e *Ethereal) Init() {
	e.options = make(map[string]map[string]map[string]string)
}

func (e *Ethereal) GetLabels(appname string) []string{
	lbls := make([]string, 0)

	if e.options[appname] !=nil {
		for k, _ := range e.options[appname] {
			lbls = append(lbls, k)
		}
	}

	return lbls
}

func (e *Ethereal) SetOption(appname, label, variable, value string) {

	if label=="" {
		label = "default"
	}

	if e.options[appname]==nil {
		e.options[appname] = make(map[string]map[string]string)
	}
	if e.options[appname][label]==nil {
		e.options[appname][label] = make(map[string]string)
	}

	e.options[appname][label][variable] = value

}

func (e *Ethereal) GetOption(appname, label, variable string) string {

	if label=="" {
		label = "default"
	}

	if e.options[appname]==nil {
		e.options[appname] = make(map[string]map[string]string)
	}
	if e.options[appname][label]==nil {
		e.options[appname][label] = make(map[string]string)
	}

	return e.options[appname][label][variable]

}

func (e *Ethereal) GetOptions(appname, label string) map[string]string {

	if label=="" {
		label = "default"
	}

	if e.options[appname]==nil {
		e.options[appname] = make(map[string]map[string]string)
	}
	if e.options[appname][label]==nil {
		e.options[appname][label] = make(map[string]string)
	}

	return e.options[appname][label]

}

func (e *Ethereal) Close(){
	// TODO maybe we can write it to disck
}