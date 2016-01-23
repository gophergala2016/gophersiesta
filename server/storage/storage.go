package server

type Storage interface {
	Init()
	SetOption(appname, label, variable, value string)
	GetOption(appname, label, variable string) string
	GetOptions(appname, label string) map[string]string
}
