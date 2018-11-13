package TERMINUS

import (
	"net/http"
	"strings"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"time"
	"strconv"
	"github.com/gorilla/mux"
)

// InitUserConfig godoc
// @Summary Initialize User KOPS Configuration
// @Description Initialize User KOPS Configuration
// @Tags internalUse
// @Accept json
// @Produce json
// @Param body body TERMINUS.KopsConfig true "..."
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /initUserConfig [POST]
func InitUserConfig(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	// In case of any error, we respond with an error to the user
	if err != nil {
		log.Error(err)
		//fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	UserConfig.AwsAccessKeyId = r.Form.Get("AwsAccessKeyId")
	UserConfig.AwsSecretAccessKey = r.Form.Get("AwsSecretAccessKey")
	UserConfig.Region = r.Form.Get("Region")
	UserConfig.Zone = r.Form.Get("Zone")
	UserConfig.ContainerName = r.Form.Get("ContainerName")
	UserConfig.S3BucketName = strings.ToLower(UserConfig.ContainerName) + "-terminus-kops-state-store"
	UserConfig.KubeClusterName = strings.ToLower(UserConfig.ContainerName) + "terminus.k8s.local"
	UserConfig.NodeCount = r.Form.Get("NodeCount")
	UserConfig.NodeSize = r.Form.Get("NodeSize")
	UserConfig.MasterSize = r.Form.Get("MasterSize")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// ListAllInstances godoc
// @Summary List all instances on AWS
// @Description List all instances on AWS
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {array} TERMINUS.Ec2Instances ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /listAllInstances [get]
func ListAllInstances(w http.ResponseWriter, r *http.Request) {
	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(UserConfig.AwsAccessKeyId, UserConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(UserConfig.Region),
	}))

	svc := ec2.New(session)

	var allInstances []Ec2Instances

	/*input := ec2.DescribeInstancesInput{InstanceIds: []*string{
		aws.String("i-01bdeaedd35d3d55a"),
	},}*/
	input := ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(&input)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to perform descibe Instances", http.StatusBadRequest)
		return
	}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {

			oneInstance := Ec2Instances{InstanceId: ValueAssignString(instance.InstanceId, ""),
				ImageId: ValueAssignString(instance.ImageId, ""),
				InstanceType: ValueAssignString(instance.InstanceType, ""),
				LaunchTime: *instance.LaunchTime,
				InstanceState: ValueAssignString(instance.State.Name, ""),
				AvailabilityZone: ValueAssignString(instance.Placement.AvailabilityZone, ""),
				CoreCount: ValueAssignInt64(instance.CpuOptions.CoreCount, 0),
				PublicIpAddress: ValueAssignString(instance.PublicIpAddress, "")}

			allInstances = append(allInstances, oneInstance)
		}
	}

	all, err := json.Marshal(allInstances)

	if err != nil {
		log.Error(err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
		return
	}
	w.Write(all)
}

// InitKubeCluster godoc
// @Summary Initialize the KOPS container with all the environment variables
// @Description Initialize the KOPS container with all the environment variables
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /initKubeCluster [get]
func InitKubeCluster(w http.ResponseWriter, r *http.Request){

	LogVar.WithFields(log.Fields{"functionName": "InitKubeCluster"}).Info("Inside InitKubeCluster function")
	//StopAndRemoveKopsContainer(UserConfig.ContainerName,w)
	initialKopsConfiguration()

	testKopsContainer()

	w.WriteHeader(http.StatusOK)
	return
}

// StartKubeClusterKops godoc
// @Summary Start the Kube cluster
// @Description Start the Kube cluster
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Param nodeCount query string true "number of slaves"
// @Param nodeSize query string true "type of Instance for slave nodes"
// @Param masterSize query string true "type of Instance for master node default is m3.large"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /initKubeCluster [get]
func StartKubeClusterKops(w http.ResponseWriter, r *http.Request){

	LogVar.WithFields(log.Fields{"functionName": "StartKubeClusterKops"}).Info("Entry to function")

	nodeCount := r.URL.Query().Get("nodeCount")
	nodeSize := r.URL.Query().Get("nodeSize")
	masterSize := r.URL.Query().Get("masterSize")

	if nodeCount == "" && nodeSize == "" && masterSize ==""{
		LogVar.WithFields(log.Fields{"nodeCount": UserConfig.NodeCount, "nodeSize": UserConfig.NodeSize,
									"masterSize": UserConfig.MasterSize}).Info(
										"No query parameters specified so using the default params")
		LogVar.Info("Starting Kube cluster with master size ",UserConfig.MasterSize)
		go createKubeCluster()
		w.WriteHeader(http.StatusOK)
		return
	} else if nodeCount!="" && nodeSize != "" && masterSize !="" && StringInSlice(nodeSize, AllInstanceTypes) && StringInSlice(masterSize, AllInstanceTypes) {
		UserConfig.NodeCount = nodeCount
		UserConfig.NodeSize = nodeSize
		UserConfig.MasterSize = masterSize

		LogVar.WithFields(log.Fields{"nodeCount": UserConfig.NodeCount, "nodeSize": UserConfig.NodeSize,
			"masterSize": UserConfig.MasterSize}).Info(
			"Using specified query parameters ")
		log.Info("Starting Kube cluster with master size ",UserConfig.MasterSize)
		go createKubeCluster()
		w.WriteHeader(http.StatusOK)
		return
	} else if nodeCount!="" && StringInSlice(nodeSize, AllInstanceTypes) && masterSize == "" {
		UserConfig.NodeCount = nodeCount
		UserConfig.NodeSize = nodeSize
		LogVar.WithFields(log.Fields{"nodeCount": UserConfig.NodeCount, "nodeSize": UserConfig.NodeSize,
			"masterSize": UserConfig.MasterSize}).Info(
			"Using specified query parameters ")
		log.Info("Starting Kube cluster with master size ",UserConfig.MasterSize)
		go createKubeCluster()
		w.WriteHeader(http.StatusOK)
		return
	} else {
		LogVar.WithFields(log.Fields{"nodeCount": UserConfig.NodeCount, "nodeSize": UserConfig.NodeSize,
			"masterSize": UserConfig.MasterSize}).Info(
			"No query parameters specified so using the default params")
		log.Info("Starting Kube cluster with master size ",UserConfig.MasterSize)
		go createKubeCluster()
		w.WriteHeader(http.StatusOK)
	}
}
// UpdateKubeCluster godoc
// @Summary This api confirms the cluster creation
// @Description This api confirms the cluster creation
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /updateKubeCluster [get]
func UpdateKubeCluster(w http.ResponseWriter, r *http.Request)  {

	LogVar.WithFields(log.Fields{"functionName": "UpdateKubeCluster"}).Info("Entry to function")
	go updateKopsCluster()
	w.WriteHeader(http.StatusOK)
}
// ValidateKubeCluster godoc
// @Summary This api validates the cluster and wait till its fully validated
// @Description This api validates the cluster and wait till its fully validated
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /validateKubeCluster [get]
func ValidateKubeCluster(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "ValidateKubeCluster"}).Info("Entry to function")
	go validateKopsCluster()
	w.WriteHeader(http.StatusOK)
}

// CreateDashboardKops godoc
// @Summary This api deploys the dashboard to kube cluster
// @Description This api deploys the dashboard to kube cluster
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /createDashboardKops [get]
func CreateDashboardKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "CreateDashboardKops"}).Info("Entry to function")
	go createDashboard()
	w.WriteHeader(http.StatusOK)
}
// DeleteDashboardKops godoc
// @Summary This api undeploys the dashboard from kube cluster
// @Description This api undeploys the dashboard from kube cluster
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /deleteDashboardKops [get]
func DeleteDashboardKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "DeleteDashboardKops"}).Info("Entry to function")
	go deleteDashboard()
	w.WriteHeader(http.StatusOK)
}
// CreateHeapsterKops godoc
// @Summary This api deploys the Heapster to kube cluster
// @Description This api deploys the Heapster to kube cluster
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /createHeapsterKops [get]
func CreateHeapsterKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "CreateHeapsterKops"}).Info("Entry to function")
	go createHeapster()
	w.WriteHeader(http.StatusOK)
}
// DeleteHeapsterKops godoc
// @Summary This api undeploys the heapster from kube cluster
// @Description This api undeploys the heapster from kube cluster
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /deleteHeapsterKops [get]
func DeleteHeapsterKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "DeleteHeapsterKops"}).Info("Entry to function")
	go deleteHeapster()
	w.WriteHeader(http.StatusOK)
}
// CreateMetricsConfigurationKops godoc
// @Summary This api deploys the metric-configuration to kube cluster
// @Description This api deploys the metric-configuration  to kube cluster
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /createMetricsConfigurationKops [get]
func CreateMetricsConfigurationKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "CreateMetricsConfigurationKops"}).Info("Entry to function")
	go createMetricsConfiguration()
	w.WriteHeader(http.StatusOK)
}
// DeleteMetricsConfigurationKops godoc
// @Summary This api deletes the Metrics Configuration
// @Description This api deletes the Metrics Configuration
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /deleteMetricsConfigurationKops [get]
func DeleteMetricsConfigurationKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "DeleteMetricsConfigurationKops"}).Info("Entry to function")
	go deleteMetricsConfiguration()
	w.WriteHeader(http.StatusOK)
}
// CreateApplicationKops godoc
// @Summary This api deploys the seleted application with provided configuration
// @Description This api deploys the seleted application with provided configuration
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Param numMicroservice query string true "number of microservices inside application"
// @Param typeService query string true "type of microservices eg compute"
// @Param appName query string true "name of the application"
// @Param limits query bool true "whether to use limits or not"
// @Param numReplicas query string true "number of replicas inside microservice"
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /createApplicationKops [get]
func CreateApplicationKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "CreateApplicationKops"}).Info("Entry to function")

	numMicroservice := r.URL.Query().Get("numMicroservice")
	typeService := r.URL.Query().Get("typeService")
	appName := r.URL.Query().Get("appName")
	limits := r.URL.Query().Get("limits")
	numReplicas := r.URL.Query().Get("numReplicas")
	instanceType := r.URL.Query().Get("instanceType")
	limitingResourcesInstanceType := r.URL.Query().Get("limitingResourcesInstanceType")

	if numMicroservice == "" && typeService == "" && appName =="" && limits == "" && numReplicas == ""{
		go createApplication("1", "compute", DefaultAppNames[0], false, "", "1", "t2.medium")
	}else if StringInSlice(numMicroservice, SupportedNumMicroservicesApp) && StringInSlice(typeService, SupportedTypesOfMicroservices) &&
		StringInSlice(appName, DefaultAppNames) && StringInSlice(numReplicas, SupportedNumMicroserviceReplicas){

			if limits == "true" && limitingResourcesInstanceType !=""{
				go createApplication(numMicroservice, typeService, appName, true,limitingResourcesInstanceType,  numReplicas, instanceType)
			} else if (limits=="false"){
					go createApplication(numMicroservice, typeService, appName, false,"",  numReplicas, instanceType)
				}else{
				http.Error(w, "params not specified properly", http.StatusNotFound)
				return
			}
		} else{
		http.Error(w, "params not specified properly", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
// DeleteApplicationKops godoc
// @Summary This api deletes the deployed application
// @Description This api deletes the deployed application
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Param numMicroservice query string true "number of microservices inside application"
// @Param typeService query string true "type of microservices eg compute"
// @Param appName query string true "name of the application"
// @Param limits query bool true "whether to use limits or not"
// @Param numReplicas query string true "number of replicas inside microservice"
// @Param limitingResourcesInstanceType query string true "limiting instance type"
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /deleteApplicationKops [get]
func DeleteApplicationKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "DeleteApplicationKops"}).Info("Entry to function")

	numMicroservice := r.URL.Query().Get("numMicroservice")
	typeService := r.URL.Query().Get("typeService")
	appName := r.URL.Query().Get("appName")
	limits := r.URL.Query().Get("limits")
	numReplicas := r.URL.Query().Get("numReplicas")
	limitingResourcesInstanceType := r.URL.Query().Get("limitingResourcesInstanceType")

	if numMicroservice == "" && typeService == "" && appName =="" && limits == "" && numReplicas == ""{
		go deleteApplication("1", "compute", DefaultAppNames[0], false, "", "1")
	}else if StringInSlice(numMicroservice, SupportedNumMicroservicesApp) && StringInSlice(typeService, SupportedTypesOfMicroservices) &&
		StringInSlice(appName, DefaultAppNames) && StringInSlice(numReplicas, SupportedNumMicroserviceReplicas){

		if limits == "true" && limitingResourcesInstanceType!=""{
			go deleteApplication(numMicroservice, typeService, appName, true,limitingResourcesInstanceType, numReplicas)
		} else if limits=="false" {
			go deleteApplication(numMicroservice, typeService, appName, false,"", numReplicas)
		}else{
			http.Error(w, "params not specified properly", http.StatusNotFound)
			return
		}
	} else{
		http.Error(w, "params not specified properly", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetClusterInfoKops godoc
// @Summary This api gets the cluster information
// @Description This api gets the cluster information
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getClusterInfoKops [get]
func GetClusterInfoKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "GetClusterInfoKops"}).Info("Entry to function")
	outClusterInfo:=getClusterInfo()
	s := strings.Split(outClusterInfo, "\n")

	log.Info("Sent cluster information")
	log.Info(s)
	log.Info(s[0])
	w.Write([]byte(outClusterInfo))
}

// GetPodsRunningKops godoc
// @Summary This api gets all the pods information
// @Description This api gets all the pods information
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPodsRunningKops [get]
func GetPodsRunningKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "GetPodsRunningKops"}).Info("Entry to function")
	outPodsInfo:=getPodsCluster()
	log.Info("Sent Pods information")
	w.Write([]byte(outPodsInfo))
}
// GetServicesRunningKops godoc
// @Summary This api gets all the services information
// @Description This api gets all the services information
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getServicesRunningKops [get]
func GetServicesRunningKops(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "GetServicesRunningKops"}).Info("Entry to function")
	outServicesInfo:=getServicesCluster()
	log.Info("Sent Services information")
	w.Write([]byte(outServicesInfo))
}
// GetPasswordDashboard godoc
// @Summary This api gets the password to connection to Kube cluster with default username as admin
// @Description This api gets the password to connection to Kube cluster with default username as admin
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Password"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPasswordDashboard [get]
func GetPasswordDashboard(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "GetPasswordDashboard"}).Info("Entry to function")
	outPassword:=getDashboardPassword()
	log.Info("sent password ")
	w.Write([]byte(outPassword))
}
// GetTokenDashboard godoc
// @Summary This api gets the token to login into Kubedashboard
// @Description This api gets the token to login into Kubedashboard
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Token"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getTokenDashboard [get]
func GetTokenDashboard(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "GetTokenDashboard"}).Info("Entry to function")
	outToken:=getDashboardToken()
	log.Info("Sent Token")
	w.Write([]byte(outToken))
}
// DeleteKubeCluster godoc
// @Summary This api deletes the Kube cluster and undeploys everything
// @Description This api deletes the Kube cluster and undeploys everything
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /deleteKubeCluster [get]
func DeleteKubeCluster(w http.ResponseWriter, r *http.Request){

	LogVar.WithFields(log.Fields{"functionName": "DeleteKubeCluster"}).Info("Entry to function")
	go deleteKopsCluster()
	w.WriteHeader(http.StatusOK)
}
// RunLoadGen godoc
// @Summary This api deletes the Kube cluster and undeploys everything
// @Description This api deletes the Kube cluster and undeploys everything
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /runLoadGen [get]
func RunLoadGen(w http.ResponseWriter, r *http.Request){

	LogVar.WithFields(log.Fields{"functionName": "RunLoadGen"}).Info("Entry to function")
	go generateLoad("resultsdb", "www.google.com", "3", "3")
	w.WriteHeader(http.StatusOK)
}
// DeployAndStartLoadTesting godoc
// @Summary This api deploys the cluster and do load testing
// @Description This api deploys the cluster and do load testing
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Param body body TERMINUS.LoadTestQueryObj true "..."
// @Success 200 {string} string "Status"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /deployAndStartLoadTesting [get]
func DeployAndStartLoadTesting(w http.ResponseWriter, r *http.Request)  {
	LogVar.WithFields(log.Fields{"functionName": "DeployAndStartLoadTesting"}).Info("Entry to function")

	log.Info(r.Body)
	decoder := json.NewDecoder(r.Body)

	var data LoadTestQueryObj
	err := decoder.Decode(&data)
	if err != nil {
		LogVar.Error(err)
	}
	numMicroservice := data.NumMicroservice
	numReplicas := data.IterReplicas
	testIter := data.TestIter
	nodeCount := data.NodeCount
	InstanceType := data.InstanceType
	MasterType := data.MasterType
	limits := data.Limits
	appType := data.AppType
	appname := data.AppName
	testAPI := data.TestApi
	limitingResourcesInstanceType := data.LimitingResourcesInstanceType
	maxRPS := data.MaxRPS
	mainServiceName := data.MainServiceName

	log.Info("numsvs=", numMicroservice)
	log.Info("appname=", appname)
	log.Info("apptype=", appType)
	log.Info("nodecount=", nodeCount)
	log.Info("replicas=", numReplicas)
	log.Info("instancetype=", InstanceType)
	log.Info("master instancetype=", MasterType)
	log.Info("testiter=", testIter)
	log.Info("limits=", limits)
	log.Info("testAPI=", testAPI)
	log.Info("limitinginstancetype=", limitingResourcesInstanceType)
	log.Info("maxRPS=", maxRPS)
	log.Info("mainServiceName=", mainServiceName)


	numMicroserviceInt, err := strconv.Atoi(numMicroservice)
	if err != nil {
		LogVar.Error("Error : ", err)
		return
	}
	log.Info("numMicroserviceInt=", numMicroserviceInt)
	numMaxRPS, errRPS := strconv.Atoi(maxRPS)
	if errRPS != nil {
		LogVar.Error("Error : ", errRPS)
		return
	}
	if limits == "true" && limitingResourcesInstanceType!=""{
		go startLoadTest(numMicroserviceInt,numReplicas,testIter,nodeCount,InstanceType,MasterType, appType, appname,testAPI, limitingResourcesInstanceType,true ,numMaxRPS, mainServiceName )
	}else if limits =="false"{
		go startLoadTest(numMicroserviceInt,numReplicas,testIter,nodeCount,InstanceType,MasterType, appType, appname,testAPI,"", false ,numMaxRPS, mainServiceName )
	}else{
		http.Error(w, "params not specified properly", http.StatusNotFound)
		return
	}


	w.WriteHeader(http.StatusOK)
}

func GetExternalServiceIp(w http.ResponseWriter, r *http.Request){

	LogVar.WithFields(log.Fields{"functionName": "RunLoadGen"}).Info("Entry to function")
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	outExternalIp := getExternalIpSvc(vars["appname"])
	log.Info(" External IP address: ", outExternalIp[9])
	w.Write([]byte(outExternalIp[9]))
}
func GetMetricsAPI(w http.ResponseWriter, r *http.Request){

	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(UserConfig.AwsAccessKeyId, UserConfig.AwsSecretAccessKey, ""),
		Region:      aws.String("us-east-1"),
	}))

	svc := cloudwatch.New(session)

	now:=time.Now()

	input := &cloudwatch.GetMetricDataInput{
		StartTime: aws.Time(now.Add(-120*time.Minute)),
		EndTime: aws.Time(now),
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id:   aws.String("ec2Data"),
				MetricStat: &cloudwatch.MetricStat{

					Metric: &cloudwatch.Metric{
						Dimensions : []*cloudwatch.Dimension{
							{
								Name:   aws.String("InstanceId"),
								Value: aws.String("i-050c36c809e41e491"),
							},
						},
						MetricName: aws.String("CPUUtilization"),
						Namespace: aws.String("AWS/EC2"),

					},
					Period: aws.Int64(1),
					Stat:aws.String("Average"),
					Unit:aws.String("Percent"),
				},

			},
			{
				Id:   aws.String("ec2Data2"),
				MetricStat: &cloudwatch.MetricStat{

					Metric: &cloudwatch.Metric{
						Dimensions : []*cloudwatch.Dimension{
							{
								Name:   aws.String("InstanceId"),
								Value: aws.String("i-08cb2efb755b52573"),
							},
						},
						MetricName: aws.String("CPUUtilization"),
						Namespace: aws.String("AWS/EC2"),

					},
					Period: aws.Int64(1),
					Stat:aws.String("Average"),
					Unit:aws.String("Percent"),
				},

			},
		},
	}

	result, err := svc.GetMetricData(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

	b, err := json.Marshal(result)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write([]byte(b))
}