package TERMINUS

import (
	"gopkg.in/mgo.v2"
	"log"
)

var Host = []string{
"mongodb",
// replica set addrs...
}

const (
	Username   = "YOUR_USERNAME"
	Password   = "YOUR_PASS"
	Database   = "VM_BOOT_RATE_DB"
	Collection = "YOUR_COLLECTION"
	Database_Final   = "VM_BOOT_SHUTDOWN_RATE_DB"
	Database_predicted_boot  = "PREDICTED_VM_BOOT_SHUTDOWN_RATE_DB"

)
var mgoSession   *mgo.Session

// Creates a new session if mgoSession is nil i.e there is no active mongo session.
//If there is an active mongo session it will return a Clone
func GetMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.DialWithInfo(&mgo.DialInfo{
			Addrs: Host,
			 //Username: os.Getenv("MONGODB_USER"),
			 //Password: os.Getenv("MONGODB_PASS"),
			 //Database: os.Getenv("MONGO_INITDB_DATABASE"),
			// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
			// },
		})
		if err != nil {
			log.Fatal("Error: ", err)
			log.Fatal("Error: Failed to start the Mongo session")
		}
	}
	return mgoSession.Clone()
}