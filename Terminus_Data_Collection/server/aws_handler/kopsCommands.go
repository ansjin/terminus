package TERMINUS

import (
	"path/filepath"
	"fmt"
	"encoding/json"
	"time"
	"net/http"
	"bytes"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)
func createAllScriptFiles(params KopsConfig){

	var initClusterScript 		=	makeStartScript(params.AwsAccessKeyId, params.AwsSecretAccessKey, params.Region,
													params.KubeClusterName, params.S3BucketName)
	var createClusterScript 	=	makeCreateClusterScript(params.NodeCount, params.NodeSize, params.MasterSize,params.Zone)
	var updateClusterScript 	=	makeUpdateClusterScript()
	var validateClusterScript 	=	makeValidateClusterScript()
	var getPasswdScript 		=	makeGetPasswdScriptt()
	var getTokenScript 			=	makeGetTokenScript()
	var deleteClusterScript 	=	makeDeleteClusterScript()

	CreateScriptFile("/data/"+params.ContainerName+"/"+AllScriptsNames.InitClusterScript, []byte(initClusterScript))
	CreateScriptFile("/data/"+params.ContainerName+"/"+AllScriptsNames.CreateClusterScript, []byte(createClusterScript))
	CreateScriptFile("/data/"+params.ContainerName+"/"+AllScriptsNames.UpdateClusterScript, []byte(updateClusterScript))
	CreateScriptFile("/data/"+params.ContainerName+"/"+AllScriptsNames.ValidateClusterScript, []byte(validateClusterScript))
	CreateScriptFile("/data/"+params.ContainerName+"/"+AllScriptsNames.GetPasswdScript, []byte(getPasswdScript))
	CreateScriptFile("/data/"+params.ContainerName+"/"+AllScriptsNames.GetTokenScript, []byte(getTokenScript))
	CreateScriptFile("/data/"+params.ContainerName+"/"+AllScriptsNames.DeleteClusterScript, []byte(deleteClusterScript))

}
func initialKopsConfiguration() {
	// step1: Remove the container directory if it exists
	RemoveDirectory(UserConfig.ContainerName)
	// step2: make the container directory
	MakeDirectory(UserConfig.ContainerName)

	LogVar.Info("Directory Created", UserConfig.ContainerName)
	// step3: create scripts for different operationd
	createAllScriptFiles(UserConfig)
	LogVar.Info("All Scripts Created")
	// step4: Run the init script
	runScriptWithoutReturningOutput(UserConfig.ContainerName, AllScriptsNames.InitClusterScript)

	LogVar.Info("Initial Settings Done")
}
func testKopsContainer(){
	LogVar.Info("Run the Kops help command to check the kops working")
	// step1: Run the Kops help command to check the kops working
	runKopsHelpCommand()
	return
}

func createKubeCluster(){

	LogVar.Info("Create Cluster script to be run")
	runScriptWithoutReturningOutput(UserConfig.ContainerName,AllScriptsNames.CreateClusterScript)
	LogVar.Info("Cluster is created")
	return
}
func updateKopsCluster() {
	LogVar.Info("Update Cluster script to be run")
	runScriptWithoutReturningOutput(UserConfig.ContainerName,AllScriptsNames.UpdateClusterScript)
	LogVar.Info("Cluster is updated")
	return
}
func validateKopsCluster() {
	LogVar.Info("Validate Cluster script to be run")
	runScriptWithoutReturningOutput(UserConfig.ContainerName,AllScriptsNames.ValidateClusterScript)
	LogVar.Info("Cluster is validated")
	return
}

func createDashboard(){
	LogVar.Info("Dashboard yaml will be run")
	runKubectlApplyCommand("https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml")
	LogVar.Info("Dashboard is deployed")
	return
}
func deleteDashboard(){
	LogVar.Info("Dashboard deployment will be removed")
	runKubectlDeleteCommand("https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml")
	LogVar.Info("Dashboard is removed")
	return
}

func createHeapster(){
	LogVar.Info("Heapster yaml will be run")
	runKubectlApplyCommand("/data/etc/heapster.yaml")
	LogVar.Info("Heapster is deployed")
	return
}
func deleteHeapster(){
	LogVar.Info("Heapster deployment will be removed")
	runKubectlDeleteCommand("/data/etc/heapster.yaml")
	LogVar.Info("Heapster is removed")
	return
}
func createMetricsConfiguration(){
	LogVar.Info("Metrics Configuration yaml will be run")
	runKubectlApplyCommand("/data/etc/metrics_config.yaml")
	LogVar.Info("Metrics Configuration is deployed")
	return
}
func deleteMetricsConfiguration(){
	LogVar.Info("Metrics Configuration deployment will be removed")
	runKubectlDeleteCommand("/data/etc/metrics_config.yaml")
	LogVar.Info("Metrics Configuration is removed")
	return
}
func createApplication(numMicroservice string, typeService string, appName string, limits bool,limitingResourcesInstanceType,  numReplicas, instanceType string){

	LogVar.Info("App yaml will be run")
	var path string
	if limits {

		path = "/data/apps/"+ numMicroservice+"_microservice/"+typeService+"/"+appName+"/kops/limits/"+limitingResourcesInstanceType+"/"+
				numReplicas+"_replicas/web-deployment.yaml"

	} else{
		path = "/data/apps/"+ numMicroservice+"_microservice/"+typeService+"/"+appName+"/kops/nolimits/"+
			numReplicas+"_replicas/web-deployment.yaml"
	}
	LogVar.Info("path: ", path)
	runKubectlApplyCommand(path)
	LogVar.Info("App is deployed")
	return
}
func deleteApplication(numMicroservice string, typeService string, appName string, limits bool,limitingResourcesInstanceType, numReplicas string){

	LogVar.Info("App yaml will be undeployed")
	var path string
	if limits {

		path = "/data/apps/"+ numMicroservice+"_microservice/"+typeService+"/"+appName+"/kops/limits/"+limitingResourcesInstanceType+
			numReplicas+"_replicas/web-deployment.yaml"


	} else{
		path = "/data/apps/"+ numMicroservice+"_microservice/"+typeService+"/"+appName+"/kops/nolimits/"+
			numReplicas+"_replicas/web-deployment.yaml"
	}
	LogVar.Info("path: ", path)
	runKubectlDeleteCommand(path)
	LogVar.Info("App is undeployed")
	return
}

func getClusterInfo() string{
	output:= runKubectlClusterInfoCommand()
	return output
}

func getPodsCluster() string{
	output:= runKubectlGetPodsCommand()
	return output
}

func getServicesCluster() string{
	output:= runKubectlGetServicesCommand()
	return output
}

func getDashboardPassword() string{

	LogVar.Info("Get password from Cluster script to be run")
	output:=runScriptWithOutput(UserConfig.ContainerName,AllScriptsNames.GetPasswdScript)
	return output
}
func getDashboardToken() string{

	LogVar.Info("Get dashboard token from Cluster script to be run")
	output:=runScriptWithOutput(UserConfig.ContainerName,AllScriptsNames.GetTokenScript)
	return output
}

func getExternalIpSvc(appName string) []string{

	LogVar.Info("Geting External URL for the service")
	output:= runKubectlGetExternalIpSvcCommand(appName)
	LogVar.Info(output)
	return output
}

func deleteKopsCluster() {

	LogVar.Info("Delete Cluster script to be run")
	runScriptWithoutReturningOutput(UserConfig.ContainerName,AllScriptsNames.DeleteClusterScript)
	LogVar.Info("Cluster is Deleted")
	return
}

func generateLoad(dbName, testUrl, vus, rps string) {
	LogVar.Info("Entered into generateLoad function")
	url := "http://k6:9002/execute"
	LogVar.Info("URL:>", url)

	TobeSentObj:=K6QueryObj{"test", dbName, testUrl, vus,rps}

	all, err := json.Marshal(TobeSentObj)

	if err != nil {
		LogVar.Error(err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(all))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		LogVar.Error(err)
	}
	defer resp.Body.Close()

	LogVar.Info("response Status:", resp.Status)
	LogVar.Info("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	LogVar.Info("response Body:", string(body))
	LogVar.Info("Completed")
	return
}

func startLoadTest(NumMicroserviceApp int,iterReplicas, testIter , nodeCount, InstanceType, MasterType, appType,
	appName, testApi,limitingResourcesInstanceType  string, limits bool, MaxRPS int, mainServiceName string ) {

	UserConfig.NodeCount =  nodeCount
	UserConfig.NodeSize = InstanceType
	UserConfig.MasterSize = MasterType

	initialKopsConfiguration()
	testKopsContainer()

	LogVar.WithFields(log.Fields{"nodeCount": UserConfig.NodeCount, "nodeSize": UserConfig.NodeSize,
		"masterSize": UserConfig.MasterSize}).Info(
		"Using specified query parameters ")
	LogVar.Info("Starting Kube cluster with master size ", UserConfig.MasterSize)

	LogVar.Info("Entering Cluster creation Step")
	createKubeCluster()
	LogVar.Info("Cluster created")
	LogVar.Info("Entering Updation Step")
	updateKopsCluster()
	LogVar.Info("Updation step Complete")
	time.Sleep(1 * time.Minute)
	LogVar.Info("Entering Validation Step")
	validateKopsCluster()
	LogVar.Info("waiting for validation ")

	stopVal := Schedule(func() {
		validateKopsCluster()
		LogVar.Info("waiting for validation ")
	}, 30*time.Second)

	//check for is ready string in the last line
	time.Sleep(6 * time.Minute) // some buffer time given so that cluster gets ready

	stopVal <- true

	createDashboard()
	createHeapster()
	createMetricsConfiguration()
	createApplication(strconv.Itoa(NumMicroserviceApp), appType,appName, limits,limitingResourcesInstanceType, iterReplicas, InstanceType)

	//time.Sleep(2 * time.Minute)

	stop := Schedule(func() {
		outputPod := getPodsCluster()
		outputSvc := getServicesCluster()
		LogVar.Info("pods Info : ", outputPod)
		LogVar.Info("pods Info : ", outputSvc)
		log.Info("pods Info : ", outputPod)
		log.Info("pods Info : ", outputSvc)
		outClusterInfo := getClusterInfo()
		outPass := getDashboardPassword()
		outToken := getDashboardToken()

		LogVar.Info(" token: ", outClusterInfo)

		LogVar.Info("Password for user admin : ", outPass)
		LogVar.Info(" token: ", outToken)

		log.Info(" token: ", outClusterInfo)

		log.Info("Password for user admin : ", outPass)

		log.Info(" token: ", outToken)
	}, 30*time.Second)

	time.Sleep(5 * time.Minute) // some buffer time given so that cluster gets ready

	stop <- true

	outExternalIp := getExternalIpSvc(mainServiceName)

	// get the port from it

	LogVar.Info(" outExternalIp : ", outExternalIp)
	LogVar.Info(" port : ", outExternalIp[12])

	portsArr := strings.Split(outExternalIp[12], ":")

	LogVar.Info(" portsArr : ", portsArr)

	// get the test URL here from the services external endpoint

	testURL := "http://" + outExternalIp[9] + ":"+ portsArr[0] + testApi

	LogVar.Info(" External IP address: ", testURL)
	log.Info(" External IP address: ", testURL)

	var database = "NodeCount_" + UserConfig.NodeCount + "_InstanceType_" + UserConfig.NodeSize + "_MasterType_" + UserConfig.MasterSize
	database = strings.Replace(database, ".", "_", -1)
	var collectionName = "TestK6"
	LogVar.Info("collectionName : ", collectionName)
	log.Info("collectionName : ", collectionName)
	LogVar.Info("dbName : ", database)
	log.Info("dbName : ", database)
	mongoSession := GetMongoSession()

	collection := mongoSession.DB(database).C(collectionName)

	var timestamp = time.Now().Unix()

	AllData := TestInformation{
		Timestamp:          timestamp,
		DbName:             collectionName,
		NodeCount:          UserConfig.NodeCount,
		NodeInstanceType:   UserConfig.NodeSize,
		MasterInstanceType: UserConfig.MasterSize,
		NumMicroSvs:        string(NumMicroserviceApp),
		TypeApp:            appType,
		AppName:            appName,
		Limits:             limits,
		Replicas:           iterReplicas,
		TestIter:           testIter,
		TestURL:            testURL,
		Phase:              "Start",
	}
	// Insert
	if err := collection.Insert(AllData); err != nil {
		LogVar.Info("error ", err)
	} else {
		LogVar.Info("#inserted into ", collectionName)
	}
	defer mongoSession.Close()

	var rateIncr = 1
	if appType == SupportedTypesOfMicroservices[1]{

		rateIncr = 6
	}else if appType == SupportedTypesOfMicroservices[2]{

		rateIncr = 50
	} else if appType == SupportedTypesOfMicroservices[0]{
		rateIncr = 4
	}	else {
		rateIncr = 20
	}
	for iterRPS := 1; iterRPS < MaxRPS; iterRPS += rateIncr {

		LogVar.Info("#RPS ", iterRPS)
		var vus int
		if(iterRPS>500) {
			vus = 500
		} else{
			vus = iterRPS
		}

		generateLoad(collectionName, testURL, strconv.Itoa(vus), strconv.Itoa(iterRPS))

		time.Sleep(75 * time.Second)
	}
	// after it is deployed need to get the ip address of the application
	// start load testing
	// save the information of test running

	time.Sleep(5 * time.Minute) //wait for some time to collect after test data
	mongoSessionEnd := GetMongoSession()
	collectionEnd := mongoSessionEnd.DB(database).C(collectionName)
	var timestampEnd = time.Now().Unix()
	AllDataEnd := TestInformation{
		Timestamp:          timestampEnd,
		DbName:             collectionName,
		NodeCount:          UserConfig.NodeCount,
		NodeInstanceType:   UserConfig.NodeSize,
		MasterInstanceType: UserConfig.MasterSize,
		NumMicroSvs:        string(NumMicroserviceApp),
		TypeApp:            appType,
		AppName:            appName,
		Limits:             limits,
		Replicas:           iterReplicas,
		TestIter:           testIter,
		TestURL:            testURL,
		Phase:              "End",
	}
	// Insert
	if err := collectionEnd.Insert(AllDataEnd); err != nil {
		fmt.Printf("%s", err)
		return
	}
	LogVar.Info("#inserted into " + collectionName)
	defer mongoSessionEnd.Close()

	deleteKopsCluster()
	time.Sleep(2 * time.Minute) //again wait for some time to end cleanly

	//copy of data need to be performed
	//var dirName = database+"_NumMicroSvs_"+SupportedNumMicroservicesApp[NumMicroserviceApp -1 ]+
	//	"_" + appType + "_" +appName + "_nolimits_replicas_" +	 strconv.Itoa(iterReplicas) + "_testIter" + strconv.Itoa(testIter)
	//LogVar.Info("dirName : ", dirName)
	//CopyDataToDirectory(dirName)

}

func GetMetrics(containerName string) {

	dir, err := filepath.Abs(filepath.Dir(containerName))
	if err != nil {
		LogVar.Fatal(err)
	}
	fmt.Println(dir)

	type Dims struct {
		Name string
		Value string
	}

	type Metrics struct {
		Namespace string
		MetricName string
		Dimensions []Dims
	}

	type MetricStats struct {

		Metric Metrics
		Period int
		Stat string
		Unit string

	}

	type MetricQuery struct {
		ID     string
		MetricStat MetricStats
		Expression string
		Label string
		ReturnData bool

	}
	var allMetricQueries []MetricQuery

	var Dim = Dims{"InstanceId", "i-050c36c809e41e491"}
	var mDims []Dims
	mDims=append(mDims, Dim)
	testQuery:=MetricQuery{"1", MetricStats{Metrics{"AWS/EC2", "CPUUtilization", mDims},600, "Average","Percent" },
							"","", true }

	allMetricQueries = append(allMetricQueries, testQuery)
	b, err := json.Marshal(allMetricQueries)
	if err != nil {
		fmt.Println("error:", err)
	}

	var getMetrics =	"#!/bin/bash\n"+
		"source ~/.bashrc\n"+
		"aws cloudwatch get-metric-data --metric-data-queries "+string(b) + " --start-time "+string(time.Now().Unix() - 30000) +  " --end-time "+string(time.Now().Unix())

	CreateScriptFile(dir+"/"+containerName+"/getMetrics.sh", []byte(getMetrics))

	command_run := "sh /usr/personal/getMetrics.sh"

	exe_cmd(command_run)

	return

}
