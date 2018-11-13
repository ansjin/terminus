package main
import (
"net/http"

awsHandler "./aws_handler"
	log "github.com/sirupsen/logrus"
	log2 "log"
	"os"
)
// @title TERMINUSINTF APIs
// @version 1.0
// @description This is the api page for all APIs in TERMINUS

// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @contact.name Anshul Jindal
// @contact.email anshul.jindal@tum.de


// @host localhost:8082
// @BasePath /
func main() {
	// The router is now formed by calling the `newRouter` constructor function

	//awsHandler.MongoInit()
	// use JSONFormatter
	//log.SetFormatter(&logmatic.JSONFormatter{})
	f, err := os.OpenFile("/data/logs", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	awsHandler.InitUserConfigurationFirst()

	router := awsHandler.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(awsHandler.NotFoundHandler)
	// Start processing events
	awsHandler.B.Start()

	go awsHandler.ReadLogsContinously()
	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	router.Handle("/events/", awsHandler.B)

	log2.Fatal(http.ListenAndServe(":8082", router))
	log.Info("Server started")
}
