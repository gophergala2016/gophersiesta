package storage


type Ethereal struct {
	options map[string]map[string]map[string]string
}

func (this *Ethereal) Init() {
	this.options = make(map[string]map[string]map[string]string)
}

func (this *Ethereal) SetOption(appname, label, variable, value string) {

	if label=="" {
		label = "default"
	}

	if this.options[appname]==nil {
		this.options[appname] = make(map[string]map[string]string)
	}
	if this.options[appname][label]==nil {
		this.options[appname][label] = make(map[string]string)
	}

	this.options[appname][label][variable] = value

}

func (this *Ethereal) GetOption(appname, label, variable string) string {

	if label=="" {
		label = "default"
	}

	if this.options[appname]==nil {
		this.options[appname] = make(map[string]map[string]string)
	}
	if this.options[appname][label]==nil {
		this.options[appname][label] = make(map[string]string)
	}

	return this.options[appname][label][variable]

}

func (this *Ethereal) GetOptions(appname, label string) map[string]string {

	if label=="" {
		label = "default"
	}

	if this.options[appname]==nil {
		this.options[appname] = make(map[string]map[string]string)
	}
	if this.options[appname][label]==nil {
		this.options[appname][label] = make(map[string]string)
	}

	return this.options[appname][label]

}

func (this *Ethereal) Close(){
	// TODO maybe we can write it to disck
}