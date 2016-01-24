package main
import (
	"github.com/gophergala2016/gophersiesta/server"
)

// Another entry point to start the server
func main()  {

	s := server.StartServer()

	defer s.Storage.Close()
}
