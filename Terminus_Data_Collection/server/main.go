package main
import (
"net/http"

awsHandler "./aws_handler"
	log "github.com/sirupsen/logrus"
	log2 "log"
	"net"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/logmatic/logmatic-go"
	"os"
)
// @title TERMINUS APIs
// @version 1.0
// @description This is the api page for all APIs in TERMINUS

// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @contact.name Anshul Jindal
// @contact.email anshul.jindal@tum.de


// @host localhost:8081
// @BasePath /
func main() {
	// The router is now formed by calling the `newRouter` constructor function

	//awsHandler.MongoInit()
	// use JSONFormatter
	awsHandler.LogVar = log.New()
	log.SetFormatter(&logmatic.JSONFormatter{})
	conn, err := net.Dial("tcp", "logstash:5000")
	if err != nil {
		awsHandler.LogVar.Fatal(err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(log.Fields{"app": "TERMINUS"}))

	awsHandler.LogVar.Hooks.Add(hook)
	ctx := awsHandler.LogVar.WithFields(log.Fields{
		"functionName": "main",
	})

	f, err := os.OpenFile("/data/logs", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	awsHandler.InitUserConfigurationFirst()

	router := awsHandler.NewRouter()
	// Start processing events
	awsHandler.B.Start()

	go awsHandler.ReadLogsContinously()
	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	router.Handle("/events/", awsHandler.B)

	log2.Fatal(http.ListenAndServe(":8081", router))
	ctx.Info("Server started")
	log.Info("Server started")
}
