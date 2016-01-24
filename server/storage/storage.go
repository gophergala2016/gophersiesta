package storage

type Storage interface {
	Init()
	GetLabels(appname string) []string
	SetOption(appname  string, label  string, variable  string, value string)
	GetOption(appname  string, label  string, variable string) string
	GetOptions(appname  string, label string) map[string]string
	Close()
}


func CreateSampleData(s Storage){
	s.SetOption("app1", "prod", "datasource.url", "jdbc:mysql://proddatabaseserver:3306/shcema?profileSQL=true")
	s.SetOption("app1", "", "datasource.username", "GOPHER")
	s.SetOption("app1", "dev", "datasource.username", "GOPHER-dev")
	s.SetOption("app1", "prod", "datasource.username", "GOPHER-prod")
	s.SetOption("app1", "", "datasource.password", "FOOBAR")
	s.SetOption("app1", "dev", "datasource.password", "LOREM")
	s.SetOption("app1", "prod", "datasource.password", "IPSUM")

	s.SetOption("app2", "", "datasource.password", "DOCKER-PASS")
	s.SetOption("app2", "dev", "datasource.password", "DEV-PASS")
}