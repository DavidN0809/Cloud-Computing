// Launch microservice server- main.go
package main
import (
"github.com/DavidN0809/Cloud-Computing/lab11/microservice"
"log"
)
func main() {
s := microservice.NewServer("", "8000")
log.Fatal(s.ListenAndServe())
}
