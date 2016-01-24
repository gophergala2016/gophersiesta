package main
import "github.com/gophergala2016/gophersiesta/server"

func main()  {

	s := server.StartServer()

	defer s.Storage.Close()
}