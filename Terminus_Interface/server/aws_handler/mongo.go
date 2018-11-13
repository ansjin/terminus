package TERMINUS

import (
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"os"
	"gopkg.in/mgo.v2/bson"
)

var Host = []string{
"mongodb",
// replica set addrs...
}

const (
	Database   = "TERMINUS"
	Collection_MSC = "ALL_TESTS_MSC_PREDICTION"
	Collection_MSC_EXPERIMENTAL = "ALL_EXPERIMENTAL_MSC_CALCULATED"
	Collection_MSC_REGRESSION = "ALL_REGRESSION_MSC_CALCULATED"
	Collection_All_Test_Names = "ALL_BRUTE_FORCE_CONDUCTED_TEST_NAMES"
	Collection_All_TrainingsError = "ALL_TRAINING_RMS_ERROR"
	Collection_All_Replicas_TrainingsError = "ALL_TRAINING_RMS_ERROR_REPLICAS_PREDICTION"
	Collection_All_Smart_Trainings_Error = "ALL_SMART_TRAINING_RMS_ERROR"
	Collection_Pod_Events = "ALL_POD_EVENTS"
	Collection_Pod_Events_Analyzed = "ALL_POD_EVENTS_ANALYZED"
	Collection_Pod_Bootingtime = "ALL_POD_BOOTINGTIME"

)
var mgoSession   *mgo.Session
type TestInformation struct {
	ContainerName 			string `json:"ContainerName"`
	StartTimestamp     		int64 `json:"StartTimestamp"`
	InfluxDbK6Name         	string `json:"InfluxDbK6Name"`
	InfluxDbK8sName 		string `json:"InfluxDbK8sName"`
	NodeCount          		string `json:"NodeCount"`
	NodeInstanceType   		string `json:"NodeInstanceType"`
	MasterInstanceType 		string `json:"MasterInstanceType"`
	NumMicroSvs        		string `json:"NumMicroSvs"`
	TypeApp            		string `json:"TypeApp"`
	AppName            		string `json:"AppName"`
	Limits             		string `json:"Limits"`
	LimitingInstanceType 	string `json:"LimitingInstanceType"`
	Replicas           		string `json:"Replicas"`
	TestIter           		string `json:"TestIter"`
	TestAPI            		string `json:"TestAPI"`
	TestVMIP            	string `json:"TestVMIP"`
	Grafana            		string `json:"Grafana"`
	Logs            		string `json:"Logs"`
	Kibana            		string `json:"Kibana"`
	Phase              		string `json:"Phase"`
	EndTimestamp     		int64 `json:"EndTimestamp"`
}

type MetadataObjectStruct struct {
	Name string `json:"Name"`
	Namespace string `json:"Namespace"`
	SelfLink string `json:"SelfLink"`
	Uid string `json:"Uid"`
	ResourceVersion string `json:"ResourceVersion"`
	CreationTimestamp string `json:"CreationTimestamp"`
}
type InvolvedObjectStruct struct {
	Kind string `json:"Kind"`
	Namespace string `json:"Namespace"`
	Name string `json:"Name"`
	Uid string `json:"Uid"`
	ApiVersion string `json:"ApiVersion"`
	ResourceVersion string `json:"ResourceVersion"`
}
type SourceStruct struct {
	Component string `json:"Component"`
}
type ItemsObject struct {
	Metadata MetadataObjectStruct `json:"Metadata"`
	InvolvedObject InvolvedObjectStruct `json:"InvolvedObject"`
	Reason string `json:"Reason"`
	Message string `json:"Message"`
	Source SourceStruct `json:"Source"`
	FirstTimestamp string  `json:"FirstTimestamp"`
	LastTimestamp string `json:"LastTimestamp"`
	Count int64 `json:"Count"`
	Type string `json:"Type"`
	EventTime string `json:"EventTime"`
	ReportingComponent string `json:"ReportingComponent"`
	ReportingInstance string `json:"ReportingInstance"`
}
type MetadataObject2 struct {
	SelfLink string `json:"SelfLink"`
	ResourceVersion string `json:"ResourceVersion"`

}
type KubeEventsObject struct {
	Kind       string          `json:"Kind"`
	ApiVersion string          `json:"ApiVersion"`
	Metadata   MetadataObject2 `json:"Metadata"`
	Items      []ItemsObject   `json:"Items"`
}
type KubeMongoEventsObject struct {
	Replicas string `json:"Replicas"`
	AppName string `json:"AppName"`
	AppType string `json:"AppType"`
	KubeEvents KubeEventsObject `json:"KubeEvents"`
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

func GetTestFolderName(hostInstance,comparingInstance, numReplicas, appType, appName string, limits bool) string{
	var resultDbNameInfo InfluDbTabaseNamesInfo
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_All_Test_Names)

	log.Info("hostinstancetype:",  hostInstance,", servicetype:", appType,	", servicename:", appName,", replicas:", numReplicas,", limitinginstancetype:", comparingInstance)
	if limits{
		err :=  collection.Find(bson.M{
			"hostinstancetype" : hostInstance,
			"servicetype" : appType,
			"servicename" : appName,
			"replicas" : numReplicas,
			"limitinginstancetype" : comparingInstance}).One(&resultDbNameInfo)
		if err != nil {
			log.Info("Db Error : ", err)
		}
	}else{
		err :=  collection.Find(bson.M{
			"hostinstancetype" : hostInstance,
			"servicetype" : appType,
			"servicename" : appName,
			"replicas" : numReplicas,
			"limitinginstancetype" : ""}).One(&resultDbNameInfo)
		if err != nil {
			log.Info("Db Error : ", err)
		}
	}
	defer mongoSession.Close()
	return  resultDbNameInfo.FolderName
}