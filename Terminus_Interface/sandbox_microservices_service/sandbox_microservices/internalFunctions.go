package sandbox_microservices

import (
	"encoding/json"
	"strings"
	log "github.com/sirupsen/logrus"
	"time"
	"gopkg.in/mgo.v2/bson"
	"os"
	"net/http"
	"io/ioutil"
)



func appendHoxyEnvironmentVar(parseddockerComposeObj *DockerComposeParsedStruct){

	log.Info("Appending HOXY_APP as the environment variable")
	for key,service:= range parseddockerComposeObj.Services{

		if(len(service.Environment)==0){
			service.Environment = make([]string, 0)
		}
		if(len(service.Depends_on)==0){
			service.Depends_on = make([]string, 0)
		}
		if(len(service.Links)==0){
			service.Links = make([]string, 0)
		}
		service.Environment = append(service.Environment, "HTTP_PROXY=http://hoxy_app:8085")
		service.Depends_on = append(service.Depends_on, "hoxy_app")
		service.Links = append(service.Links, "hoxy_app")
		parseddockerComposeObj.Services[key]=service
	}
	parseddockerComposeObj.Services["hoxy_app"]=ServicesStruct{Container_name:"hoxy_app", Image:"terminusimages/hoxy_app:latest",
	Environment:[]string{"MONGODB_HOST="+GetOutboundIP(), "MONGO_PORT=27017", "MONGO_USER="+os.Getenv("MONGODB_USER"),
							"MONGO_PASS="+os.Getenv("MONGODB_PASS"), "MONGO_DB="+Database,
							"MONGO_COLLECTION="+Collection_API_REQUEST_RESPONSE_APIS},
	Ports:[]string{"8085:8085"},
	}
}
func createDummyResponseApp(serviceName, port, response, apiEndpoint string ) ServicesStruct{

	var Environment []string
	log.Info("Appending dummy response service ")
	Environment = append(Environment, "TEST_API="+apiEndpoint)
	Environment = append(Environment, "DUMMY_RESPONSE="+response)
	Environment = append(Environment, "PORT="+port)



	serviceInfo:=ServicesStruct{Container_name:serviceName, Image:"terminusimages/dummy_response_app:latest",
		Environment:Environment,
		Ports:[]string{port+":"+port},
	}
	return serviceInfo
}

func getServicePort(mainServiceName string,parseddockerComposeObj DockerComposeParsedStruct ) string{

	if parseddockerComposeObj.Services[mainServiceName].Ports[0]!=""{
		portStr:=parseddockerComposeObj.Services[mainServiceName].Ports[0]
		portArr:=strings.Split(portStr, ":")
		if(len(portArr)>0){
			return portArr[0]
		}else{
			return ""
		}
	}else{
		return ""
	}
}


func deployserviceToRecordAPIAndResponse(endpointAPI,mainServiceName,  dockerComposeFilePath string) bool{

	parseddockerComposeObj:=parseyamlToJson(dockerComposeFilePath)
	if parseddockerComposeObj.Version== "" {
		log.Info("not able to parse version is empty")
		return false
	}
	AllData := NewDockerComposeService{
		DockerComposeFilePath: dockerComposeFilePath,
		MainServiceName:       mainServiceName,
		EndpointAPI:           strings.Replace(endpointAPI, "_", "/", -1),
		ModifiedDockerCompose: parseddockerComposeObj,
		ExternalPort:          getServicePort(mainServiceName, parseddockerComposeObj),
		Original:true,

	}
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_Individual_Modified_Services)

	if err := collection.Insert(AllData); err != nil {
		log.Info("error ", err)
	} else {
		log.Info("#inserted into ", Collection_Individual_Modified_Services)
	}
	defer mongoSession.Close()

	appendHoxyEnvironmentVar(&parseddockerComposeObj)
	log.Info(parseddockerComposeObj)
	endpointAPI = strings.Replace(endpointAPI, "_", "/", -1)
	externalPortofService:=getServicePort(mainServiceName, parseddockerComposeObj)
	if(externalPortofService==""){
		return false
	} else {

		mongoSession := GetMongoSession()
		collection := mongoSession.DB(Database).C(Collection_REQUEST_RESPONSE_TESTS)

		AllData := TestInformation{
			DockerComposeFilePath: dockerComposeFilePath,
			StartTimestamp:        time.Now().Unix(),
			MainServiceName:       mainServiceName,
			EndpointAPI:           endpointAPI,
			Phase:                 "Deployment",
			ModifiedDockerCompose: parseddockerComposeObj,
			ExternalPort:          externalPortofService,
		}
		if err := collection.Insert(AllData); err != nil {
			log.Info("error ", err)
		} else {
			log.Info("#inserted into ", Collection_REQUEST_RESPONSE_TESTS)
		}

		parseddockerComposeObjBytes, err := json.Marshal(parseddockerComposeObj)

		if err != nil {
			log.Error(err)
			return false
		} else {

			log.Info("Starting a test VM t2.large and deploying data collector")
			startedInstanceId := startTestVM(string(parseddockerComposeObjBytes))
			time.Sleep(20 * time.Second)
			stopDataCollectionDeployment := Schedule(func() {
				log.Info("waiting for the deployment of the VM and application to get finish(waiting for 10 minutes).....")
				getVMPublicIP(startedInstanceId)
			}, 30*time.Second)
			time.Sleep(5 * time.Minute)

			publicAddress := getVMPublicIP(startedInstanceId)
			if publicAddress != "" {
				stopDataCollectionDeployment <- true
				log.Info("Stoping the recurring wait as public address is found")
			} else {
				log.Info("5 minutes more wait as VM has not started")
				time.Sleep(5 * time.Minute)
				stopDataCollectionDeployment <- true
				log.Info("VM has not started")
				log.Info("VM has not started... Terminating Test")
				errMongoU := collection.Update(bson.M{"dockercomposefilepath": dockerComposeFilePath}, bson.M{"$set": bson.M{"phase": "Terminated"}})
				if errMongoU != nil {
					log.Info("Error : %s", errMongoU)
				}
				return false
			}

			errMongoU := collection.Update(bson.M{"dockercomposefilepath": dockerComposeFilePath}, bson.M{"$set": bson.M{"testvmip": publicAddress,
				"phase": "Deployed"}})
			if errMongoU != nil {
				log.Info("Error : %s", errMongoU)
			}

			testURL:="http://"+publicAddress+":"+externalPortofService+endpointAPI
			log.Info("IP for the service to test should be",testURL )

			responseExp, err := http.Get(testURL)
			if err != nil {
				log.Info(err)
			} else {
				defer responseExp.Body.Close()
				contentsExp, err := ioutil.ReadAll(responseExp.Body)
				if err != nil {
					log.Info(err)
				}
					log.Info(string(contentsExp))

				myobj :=APIResponseObj{Headers:responseExp.Header, Protocol:"http:",Method:"GET", Hostname:mainServiceName, Port:externalPortofService,
						Api_endpoint: endpointAPI, Response:responseExp.Body}
				mongoSession := GetMongoSession()
				collection := mongoSession.DB(Database).C(Collection_API_REQUEST_RESPONSE_APIS)

				if err := collection.Insert(myobj); err != nil {
					log.Info("error ", err)
				} else {
					log.Info("#inserted into ", Collection_API_REQUEST_RESPONSE_APIS)
				}
				defer mongoSession.Close()
			}
			time.Sleep(5 * time.Minute)
			log.Info(" Terminating Test VM")
			terminateTestVM(startedInstanceId)

			errMonFin := collection.Update(bson.M{"dockercomposefilepath": dockerComposeFilePath}, bson.M{"$set": bson.M{"endtimestamp": time.Now().Unix(),
				"phase": "Completed"}})
			if errMonFin != nil {
				log.Info("Error::%s", errMonFin)
			}

			defer mongoSession.Close()
		}
		return true
	}
}

func formTreeDockerCompose(serviceName string, finalParsedObj DockerComposeParsedStruct) Node {
	var (
		ServiceTable = map[string]*Node{}
		root      Node
	)
	addAllServicesToTree(serviceName, "0", finalParsedObj, ServiceTable, &root)
	show(&root)

	return root
}
func getServiceResponse(serviceName string)  string{

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_API_REQUEST_RESPONSE_APIS)

	var myobj APIResponseObj
	errPod := collection.Find(bson.M{"hostname": serviceName}).One(&myobj)
	if errPod != nil {
		log.Info("error ", errPod)
	}
	jsonParsedBytes, err := json.Marshal(myobj.Response)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(jsonParsedBytes)

}
func getServiceAPIEndPoint(serviceName string)  string{
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_API_REQUEST_RESPONSE_APIS)

	var myobj APIResponseObj
	errPod := collection.Find(bson.M{"hostname": serviceName}).One(&myobj)
	if errPod != nil {
		log.Info("error ", errPod)
	}

	return string(myobj.Api_endpoint)
}

func addDummyServiceToServices(root *Node, finalParsedObj DockerComposeParsedStruct ) {

	finalShortenedJson:=formYaml(root, finalParsedObj)
	for serviceName, serviceInfo :=range finalShortenedJson.Services{
		if searchForToRelevantDepndentServices(serviceName)==false {
			var individualServiceJson DockerComposeParsedStruct
			individualServiceJson.Version = finalShortenedJson.Version
			individualServiceJson.Services = make(map[string]ServicesStruct, len(serviceInfo.Depends_on)+1)
			if(len(serviceInfo.Depends_on) > 0) {

				for i := 0; i < len(serviceInfo.Depends_on); i++ {
					if (strings.Contains(serviceInfo.Depends_on[i], "dummy") == false) {
						individualServiceJson.Services[serviceInfo.Depends_on[i]] = finalShortenedJson.Services[serviceInfo.Depends_on[i]]
					} else {
						names := strings.Split(serviceInfo.Depends_on[i], "dummy_")
						serviceInfo.Depends_on[i] = names[1]
						if len(serviceInfo.Links) > 0{
							serviceInfo.Links[i] = names[1]
						}
						/*dummy := createDummyResponseApp(serviceInfo.Depends_on[i], getServicePort(names[1], finalParsedObj),
						getServiceResponse(names[1]), getServiceAPIEndPoint(names[1]))*/
						dummy := createDummyResponseApp(names[1], getServicePort(names[1], finalParsedObj), // need to keep the name here same as old so that request can be completed
							getServiceResponse(names[1]), getServiceAPIEndPoint(names[1]))

						individualServiceJson.Services[names[1]] = dummy // need to keep the name here same as old so that request can be completed
					}
				}
			}else if len(serviceInfo.Links) > 0 {

				for i := 0; i < len(serviceInfo.Links); i++ {
					if (strings.Contains(serviceInfo.Links[i], "dummy") == false) {
						individualServiceJson.Services[serviceInfo.Links[i]] = finalShortenedJson.Services[serviceInfo.Links[i]]
					} else {
						names := strings.Split(serviceInfo.Links[i], "dummy_")
						serviceInfo.Links[i] = names[1]
						/*dummy := createDummyResponseApp(serviceInfo.Depends_on[i], getServicePort(names[1], finalParsedObj),
						getServiceResponse(names[1]), getServiceAPIEndPoint(names[1]))*/
						dummy := createDummyResponseApp(names[1], getServicePort(names[1], finalParsedObj), // need to keep the name here same as old so that request can be completed
							getServiceResponse(names[1]), getServiceAPIEndPoint(names[1]))

						individualServiceJson.Services[names[1]] = dummy // need to keep the name here same as old so that request can be completed
					}
				}
			}
			individualServiceJson.Services[serviceName] = serviceInfo
			storeNewServicesMongoDb(individualServiceJson,serviceName)
			log.Info("Stored", serviceName)

		}
	}
}

func storeNewServicesMongoDb(individualServiceJson DockerComposeParsedStruct, serviceName string)  {
	AllData := NewDockerComposeService{
			DockerComposeFilePath: "",
			MainServiceName:       serviceName,
			EndpointAPI:           getServiceAPIEndPoint(serviceName),
			ModifiedDockerCompose: individualServiceJson,
			ExternalPort:          getServicePort(serviceName, individualServiceJson),
			Original:false,
		}
		mongoSession := GetMongoSession()
		collection := mongoSession.DB(Database).C(Collection_Individual_Modified_Services)

		if err := collection.Insert(AllData); err != nil {
			log.Info("error ", err)
		} else {
			log.Info("#inserted into ", Collection_Individual_Modified_Services)
		}
		defer mongoSession.Close()
}