package main
import (
"net/http"
"log"

awsHandler "./aws_handler"
)
// @title INSTANCE APIs
// @version 1.0
// @description This is the api page for all APIs in INSTANCE

// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @contact.name Anshul Jindal
// @contact.email anshul.jindal@tum.de


// @host localhost:8080
// @BasePath /
func main() {
	// The router is now formed by calling the `newRouter` constructor function

	//awsHandler.MongoInit()

	router := awsHandler.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
	log.Printf("Server started")
}
