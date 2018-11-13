package main
import (
"net/http"
	log "github.com/sirupsen/logrus"
	log2 "log"
	"os"
	"./sandbox_microservices"
)
// @title sandbox_microservices APIs
// @version 1.0
// @description This is the api page for all APIs in sandbox_microservices

// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @contact.name Anshul Jindal
// @contact.email anshul.jindal@tum.de


// @host localhost:8083
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

	sandbox_microservices.InitUserConfigurationFirst()

	router := sandbox_microservices.NewRouter()
	// Start processing events
	sandbox_microservices.B.Start()

	go sandbox_microservices.ReadLogsContinously()
	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	router.Handle("/events/", sandbox_microservices.B)

	log2.Fatal(http.ListenAndServe(":8083", router))
	log.Info("Server started")
}
