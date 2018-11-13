package sandbox_microservices

import (
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"os"
)

var Host = []string{
"mongodb",
// replica set addrs...
}

const (
	Database   = "sandbox_microservices"
	Collection_REQUEST_RESPONSE_TESTS = "REQUEST_RESPONSE_TESTS"
	Collection_API_REQUEST_RESPONSE_APIS = "REQUEST_RESPONSE_TESTS_APIS"
	Collection_Individual_Modified_Services = "INDIVIDUAL_MODIFIED_SERVICES"
)
var mgoSession   *mgo.Session
type TestInformation struct {
	DockerComposeFilePath 			string `json:"DockerComposeFilePath"`
	StartTimestamp     				int64 `json:"StartTimestamp"`
	MainServiceName         		string `json:"MainServiceName"`
	EndpointAPI 					string `json:"EndpointAPI"`
	Phase          					string `json:"Phase"`
	EndTimestamp     				int64 `json:"EndTimestamp"`
	ModifiedDockerCompose 			interface{} `json:"ModifiedDockerCompose"`
	ExternalPort					string `json:"ExternalPort"`
	Testvmip						string `json:"Testvmip"`
}

type NewDockerComposeService struct {
	DockerComposeFilePath 			string `json:"DockerComposeFilePath,omitempty"`
	MainServiceName         		string `json:"MainServiceName,omitempty"`
	EndpointAPI 					string `json:"EndpointAPI,omitempty"`
	ModifiedDockerCompose 			interface{} `json:"ModifiedDockerCompose,omitempty"`
	ExternalPort					string `json:"ExternalPort,omitempty"`
	Response						interface{} `json:"Response,omitempty"`
	Original						bool `json:"Original,omitempty"`
}


// Creates a new session if mgoSession is nil i.e there is no active mongo session.
//If there is an active mongo session it will return a Clone
func GetMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.DialWithInfo(&mgo.DialInfo{
			Addrs: Host,
			 Username: os.Getenv("MONGODB_USER"),
			 Password: os.Getenv("MONGODB_PASS"),
			 //Database: os.Getenv("MONGO_INITDB_DATABASE"),
			// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
			// },
		})
		if err != nil {
			log.Info("Error: ", err)
			log.Info("Error: Failed to start the Mongo session")
		}
	}
	return mgoSession.Clone()
}