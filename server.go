package main
 
import (
  "fmt"
	"net/http"
	"os"
  "github.com/rmonjo/instant/web_sockets"
  "github.com/rmonjo/instant/container"
)
 
func main() {

  //first of all create the baseCN
  container.CreateBaseCn()
	
  http.HandleFunc("/", root)

  // Normal resources
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "9999"
	}

  go web_sockets.StartServer()

	fmt.Println("listening on port: ", port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
	  panic(err)
	}
}
 
//routes
func root(res http.ResponseWriter, req *http.Request) {
  http.ServeFile(res, req, "./public/terminal.html")
}

