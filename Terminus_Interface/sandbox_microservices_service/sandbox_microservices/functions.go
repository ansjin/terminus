package sandbox_microservices

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"encoding/json"
	"github.com/ghodss/yaml"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"strings"
)
func ParseYamlToJson(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])

	path:="/data/docker_compose_files/"+vars["filename"]+".yml"
	finalParsedObj:=parseyamlToJson(path)
	jsonParsedBytes, err := json.Marshal(finalParsedObj)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
		return
	}
	w.Write(jsonParsedBytes)
}
func GetServicesFromDockerCompose(w http.ResponseWriter, r *http.Request)  {

	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])

	path:="/data/docker_compose_files/"+vars["filename"]+".yml"
	finalParsedObj:=parseyamlToJson(path)

	var services []string
	for serviceName, _ :=range finalParsedObj.Services{

		services = append(services, serviceName)
	}
	jsonParsedBytes, err := json.Marshal(services)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
		return
	}
	w.Write(jsonParsedBytes)

}
func ParseYamlToYaml(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])
	path:="/data/docker_compose_files/"+vars["filename"]+".yml"
	finalParsedObj:=parseyamlToJson(path)
	jsonParsedBytes, err := json.Marshal(finalParsedObj)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
		return
	}
	yamlParsedBytes, err := yaml.JSONToYAML(jsonParsedBytes)
	if err != nil {
		log.Info( err)
		http.Error(w, "Error converting yaml", http.StatusBadRequest)
		return
	}
	w.Write(yamlParsedBytes)
}
func DeployserviceToRecordAPIAndResponse(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])
	log.Info("endpointapi=", vars["endpointapi"])
	log.Info("mainservicename=", vars["mainservicename"])
	go func(endpointAPI, mainserviceName, dockerComposeFileName string) {

		result:= deployserviceToRecordAPIAndResponse(endpointAPI,mainserviceName,
			"/data/docker_compose_files/"+dockerComposeFileName+".yml")

		log.Info("Result of Deployment : ", result)

		if(result==true){

			path:="/data/docker_compose_files/"+vars["filename"]+".yml"
			finalParsedObj:=parseyamlToJson(path)
			root:=formTreeDockerCompose(vars["mainservicename"],finalParsedObj )
			addDummyServiceToServices(&root, finalParsedObj)
		}

	}(vars["endpointapi"],vars["mainservicename"],vars["filename"] )

	w.Write([]byte("Started Will update Status in Logs"))
}
func FormTreeFromDockerCompose(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])
	log.Info("mainservicename=", vars["mainservicename"])

	path:="/data/docker_compose_files/"+vars["filename"]+".yml"
	finalParsedObj:=parseyamlToJson(path)

	root:=formTreeDockerCompose(vars["mainservicename"],finalParsedObj )
	jsonParsedBytes, err := json.Marshal(root)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
	}
	w.Write(jsonParsedBytes)
}

func Test12(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	log.Info("filename=", vars["filename"])
	log.Info("mainservicename=", vars["serviceName"])

	path:="/data/docker_compose_files/"+vars["filename"]+".yml"
	finalParsedObj:=parseyamlToJson(path)
	root:=formTreeDockerCompose(vars["mainservicename"],finalParsedObj )
	addDummyServiceToServices(&root, finalParsedObj)
}
func GetModifiedDockerComposeYaml(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	log.Info("serviceName=", vars["serviceName"])
	log.Info("original=", vars["original"])

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_Individual_Modified_Services)

	var obj NewDockerComposeService
	if(vars["original"]=="true"){
		errPod := collection.Find(bson.M{"mainservicename": vars["serviceName"], "original": true}).One(&obj)
		if errPod != nil {
			log.Info("error ", errPod)
			http.Error(w, "Service Not found", http.StatusBadRequest)
			return
		}
	}else{
		errPod := collection.Find(bson.M{"mainservicename": vars["serviceName"], "original": false}).One(&obj)
		if errPod != nil {
			log.Info("error ", errPod)
			http.Error(w, "Service Not found", http.StatusBadRequest)
			return
		}
	}

	bsonBytes, _ := bson.Marshal(obj.ModifiedDockerCompose)

	var parsedDockerCmp DockerComposeParsedStruct
	bson.Unmarshal(bsonBytes, &parsedDockerCmp)

	jsonParsedBytes, err := json.Marshal(parsedDockerCmp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
		return
	}

	yamlParsedBytes, err := yaml.JSONToYAML(jsonParsedBytes)
	if err != nil {
		log.Info( err)
		http.Error(w, "Error converting yaml", http.StatusBadRequest)
		return
	}
	w.Write(yamlParsedBytes)

}
func FormTreeFromDockerComposeMongo(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	log.Info("serviceName=", vars["serviceName"])
	log.Info("original=", vars["original"])

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_Individual_Modified_Services)

	var obj NewDockerComposeService
	if(vars["original"]=="true"){
		errPod := collection.Find(bson.M{"mainservicename": vars["serviceName"], "original": true}).One(&obj)
		if errPod != nil {
			log.Info("error ", errPod)
			http.Error(w, "Service Not found", http.StatusBadRequest)
			return
		}
	}else{
		errPod := collection.Find(bson.M{"mainservicename": vars["serviceName"], "original": false}).One(&obj)
		if errPod != nil {
			log.Info("error ", errPod)
			http.Error(w, "Service Not found", http.StatusBadRequest)
			return
		}
	}

	bsonBytes, _ := bson.Marshal(obj.ModifiedDockerCompose)

	var parsedDockerCmp DockerComposeParsedStruct
	bson.Unmarshal(bsonBytes, &parsedDockerCmp)

	root:=formTreeDockerCompose(vars["serviceName"],parsedDockerCmp )
	jsonParsedBytes, err := json.Marshal(root)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
	}
	w.Write(jsonParsedBytes)
}

func GetAllTreesFromDockerComposeMongo(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	log.Info("mainservicename=", vars["serviceName"])
	log.Info("filename=", vars["filename"])

	path:="/data/docker_compose_files/"+vars["filename"]+".yml"
	finalParsedObj:=parseyamlToJson(path)

	var composeServices []Node
	var obj NewDockerComposeService
	var parsedDockerCmp DockerComposeParsedStruct

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_Individual_Modified_Services)

	errPod := collection.Find(bson.M{"mainservicename": vars["serviceName"], "original": true}).One(&obj)
	if errPod != nil {
		log.Info("error ", errPod)
		http.Error(w, "Service Not found", http.StatusBadRequest)
		return
	}
	bsonBytes, _ := bson.Marshal(obj.ModifiedDockerCompose)
	bson.Unmarshal(bsonBytes, &parsedDockerCmp)
	root:=formTreeDockerCompose(vars["serviceName"],parsedDockerCmp )
	composeServices = append(composeServices, root)

	log.Info("Now will search for all services")
	for serviceName, _ :=range finalParsedObj.Services{

		var obj NewDockerComposeService
		var parsedDockerCmp DockerComposeParsedStruct
		log.Info(serviceName)
		errPod := collection.Find(bson.M{"mainservicename": serviceName, "original": false}).One(&obj)
		if errPod != nil {
			log.Info("error ", errPod)
		}
		bsonBytes, _ := bson.Marshal(obj.ModifiedDockerCompose)
		bson.Unmarshal(bsonBytes, &parsedDockerCmp)
		root:=formTreeDockerCompose(serviceName,parsedDockerCmp )

		composeServices = append(composeServices, root)

	}

	jsonParsedBytes, err := json.Marshal(composeServices)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
	}
	w.Write(jsonParsedBytes)
}
func GetAvailableApps(w http.ResponseWriter, r *http.Request)  {

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_Individual_Modified_Services)

	var obj []NewDockerComposeService
	errPod := collection.Find(nil).All(&obj)
	if errPod != nil {
		log.Info("error ", errPod)
		http.Error(w, "Services Not found", http.StatusBadRequest)
		return
	}

	type ServiceInfo struct {
		ServiceName string
		Original bool
	}
	var servicesInfo []ServiceInfo
	for _, service := range obj{

		service:=ServiceInfo{service.MainServiceName, service.Original}
		servicesInfo = append(servicesInfo, service)
	}
	jsonParsedBytes, err := json.Marshal(servicesInfo)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
		return
	}
	w.Write(jsonParsedBytes)
}

func GetModifiedServiceEndpoint(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	log.Info("serviceName=", vars["serviceName"])
	log.Info("original=", vars["original"])

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_Individual_Modified_Services)

	var obj NewDockerComposeService
	if(vars["original"]=="true"){
		errPod := collection.Find(bson.M{"mainservicename": vars["serviceName"], "original": true}).One(&obj)
		if errPod != nil {
			log.Info("error ", errPod)
			http.Error(w, "Service Not found", http.StatusBadRequest)
			return
		}
	}else{
		errPod := collection.Find(bson.M{"mainservicename": vars["serviceName"], "original": false}).One(&obj)
		if errPod != nil {
			log.Info("error ", errPod)
			http.Error(w, "Service Not found", http.StatusBadRequest)
			return
		}
	}
	w.Write([]byte(obj.EndpointAPI))

}
func GetAllDockerComposeFileNames (w http.ResponseWriter, r *http.Request)  {

	var foundFiles []string
	files, err := ioutil.ReadDir("/data/docker_compose_files/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		log.Info(f.Name())
		if(strings.Contains(f.Name(), ".yml")){

			namesExtension:=strings.Split(f.Name(), ".")

			foundFiles = append(foundFiles, namesExtension[0])
		}

	}
	all, err := json.Marshal(foundFiles)

	if err != nil {
		log.Error("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
}

func GetAllTestsInformation (w http.ResponseWriter, r *http.Request)  {

	var AllData []TestInformation

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_REQUEST_RESPONSE_TESTS)
	err := collection.Find(nil).All(&AllData)
	if err != nil {
		log.Info("%s", err)
		http.Error(w, "Db error", http.StatusBadRequest)
		return
	}
	defer mongoSession.Close()

	all, err := json.Marshal(AllData)

	if err != nil {
		log.Info("%s", err)
		http.Error(w, "JSON conversion error", http.StatusBadRequest)
		return
	} else {
		w.Write(all)
	}
}