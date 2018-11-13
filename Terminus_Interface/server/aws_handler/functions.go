package TERMINUS

import (
	"net/http"
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
	"github.com/gorilla/mux"
	"strings"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"math"
)

// InitUserConfig godoc
// @Summary Initialize User KOPS Configuration
// @Description Initialize User KOPS Configuration
// @Tags internalUse
// @Accept json
// @Produce json
// @Param body body TERMINUS.AWSConfigStruct true "..."
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
	AWSConfig.AwsAccessKeyId = r.Form.Get("AwsAccessKeyId")
	AWSConfig.AwsSecretAccessKey = r.Form.Get("AwsSecretAccessKey")
	AWSConfig.Region = r.Form.Get("Region")
	AWSConfig.KeyPairName = r.Form.Get("KeyPairName")
	AWSConfig.SubnetId = r.Form.Get("SubnetId")
	AWSConfig.SecurityGroup = r.Form.Get("SecurityGroup")
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
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
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

// CollectDataNoLimits godoc
// @Summary Collect the non limits data for an application on a VM
// @Description Collect the non limits data for an application on a VM /collectDataNoLimits/1/1/1/t2.nano/3000/web/webacapp/_api_web
// @Tags START_TEST
// @Accept text/html
// @Produce json
// @Param numsvs query string true "number of services in microservice"
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param instancetype query string true " host instance type "
// @Param apiendpoint query string true " api end point specify like _api_test for /api/test"
// @Param maxRPS query string true " max RPS "
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /collectDataNoLimits/{numsvs}/{nodecount}/{replicas}/{instancetype}/{mastertype}/{testvmtype}/{maxRPS}/{apptype}/{appname}/{mainservicename}/{apiendpoint} [get]
func CollectDataNoLimits(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Info("numsvs=", vars["numsvs"])
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("nodecount=", vars["nodecount"])
	log.Info("replicas=", vars["replicas"])
	log.Info("instancetype=", vars["instancetype"])
	log.Info("mastertype=", vars["mastertype"])

	vars["apiendpoint"] = strings.Replace(vars["apiendpoint"], "_", "/", -1)
	log.Info("apiendpoint=", vars["apiendpoint"])
	log.Info("maxRPS=", vars["maxRPS"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("testVMType=", vars["testvmtype"])
	go LaunchVMandDeployInterfaceNoLimits(vars["numsvs"], vars["appname"], vars["apptype"],
		vars["nodecount"], vars["replicas"], vars["instancetype"], vars["apiendpoint"], vars["mastertype"],
		vars["maxRPS"],vars["testvmtype"], vars["mainservicename"])
	w.Write([]byte("Started the Launch Procedure"))
}

// CollectDataLimits godoc
// @Summary Collect the non limits data for an application on a VM
// @Description Collect the non limits data for an application on a VM /collectDataLimits/1/1/1/t2.medium/350/t2.nano/dbaccess/movieapp/_api_movies
// @Tags START_TEST
// @Accept text/html
// @Produce json
// @Param numsvs query string true "number of services in microservice"
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param instancetype query string true " host instance type "
// @Param limitinginstancetype query string true " limitinginstancetype instance type on main instance"
// @Param apiendpoint query string true " api end point specify like _api_test for /api/test"
// @Param maxRPS query string true " max RPS"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /collectDataLimits/{numsvs}/{nodecount}/{replicas}/{instancetype}/{mastertype}/{testvmtype}/{maxRPS}/{limitinginstancetype}/{apptype}/{appname}/{mainservicename}/{apiendpoint}/[get]
func CollectDataLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("numsvs=", vars["numsvs"])
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("nodecount=", vars["nodecount"])
	log.Info("replicas=", vars["replicas"])
	log.Info("instancetype=", vars["instancetype"])
	log.Info("mastertype=", vars["mastertype"])
	log.Info("limitinginstancetype=", vars["limitinginstancetype"])
	vars["apiendpoint"] = strings.Replace(vars["apiendpoint"], "_", "/", -1)
	log.Info("apiendpoint=", vars["apiendpoint"])
	log.Info("maxRPS=", vars["maxRPS"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("testVMType=", vars["testvmtype"])
	go LaunchVMandDeployInterfaceLimits(vars["numsvs"], vars["appname"], vars["apptype"], vars["nodecount"],
		vars["replicas"], vars["instancetype"], vars["apiendpoint"],
		vars["limitinginstancetype"],vars["mastertype"], vars["maxRPS"],vars["testvmtype"], vars["mainservicename"])
	w.Write([]byte("Started the Launch Procedure"))
}
// ConductTestToCalculatePodBootTime godoc
// @Summary Collect the data for pod botting time
// @Description Collect the data for pod botting time /conductTestToCalculatePodBootTime/1/1/1/t2.xlarge/40/t2.large/compute/primeapp/_api_prime
// @Tags START_TEST
// @Accept text/html
// @Produce json
// @Param numsvs query string true "number of services in microservice"
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param instancetype query string true " host instance type "
// @Param limitinginstancetype query string true " limitinginstancetype instance type on main instance"
// @Param apiendpoint query string true " api end point specify like _api_test for /api/test"
// @Param maxRPS query string true " max RPS"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /conductTestToCalculatePodBootTime/{numsvs}/{nodecount}/{replicas}/{instancetype}/{mastertype}/{testvmtype}/{maxRPS}/{limitinginstancetype}/{apptype}/{appname}/{mainservicename}/{apiendpoint}/ [get]
func ConductTestToCalculatePodBootTime(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Info("numsvs=", vars["numsvs"])
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("nodecount=", vars["nodecount"])
	log.Info("replicas=", vars["replicas"])
	log.Info("instancetype=", vars["instancetype"])
	log.Info("mastertype=", vars["mastertype"])
	log.Info("limitinginstancetype=", vars["limitinginstancetype"])
	vars["apiendpoint"] = strings.Replace(vars["apiendpoint"], "_", "/", -1)
	log.Info("apiendpoint=", vars["apiendpoint"])
	log.Info("maxRPS=", vars["maxRPS"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("testVMType=", vars["testvmtype"])
	go conductTestToCalculatePodBootTime(vars["numsvs"], vars["appname"], vars["apptype"], vars["nodecount"],
		vars["replicas"], vars["instancetype"], vars["apiendpoint"],
		vars["limitinginstancetype"],vars["mastertype"], vars["maxRPS"],vars["testvmtype"], vars["mainservicename"])
	w.Write([]byte("Started the Launch Procedure"))
}
// GetRelevantInstancesDataBasedOnLimits godoc
// @Summary Get the resources data for that particular instance from all host instances on which it can be possible deployed
// @Description Get the resources data for that particular instance from all host instances on which it it can be possible deployed
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.NodeValues ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstancesDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas} [get]
func GetRelevantInstancesDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("comparinginstance=", vars["comparinginstance"])

	res := get_relevant_Instances_Data_Based_On_Limits(vars["comparinginstance"], vars["replicas"], vars["apptype"], vars["appname"],vars["mainservicename"])
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write([]byte(b))
}
// GetRelevantInstanceDataBasedOnLimits godoc
// @Summary Get the resources data for that particular instance from the given host instance on which it is deployed
// @Description Get the resources data for that particular instance from the given host instance on which it is deployed
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param hostinstance query string true " hostinstance on which the comparing instance is deployed "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.NodeValues ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstanceDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{hostinstance}/{comparinginstance}/{replicas} [get]
func GetRelevantInstanceDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("hostInstance=", vars["hostinstance"])
	log.Info("comparinginstance=", vars["comparinginstance"])

	res := get_relevant_Instance_Data_Based_On_Limits(vars["hostinstance"], vars["comparinginstance"],
		vars["replicas"], vars["apptype"], vars["appname"],vars["mainservicename"])
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write([]byte(b))
}
// GetRelevantInstancesActualInstanceDataBasedOnLimits godoc
// @Summary Get the resources data for that particular instance from all host instances on which it can be possible deployed and also the actual VM performance
// @Description Get the resources data for that particular instance from all host instances on which it can be possible deployed and also the actual VM performance
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.NodeValues ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstancesActualInstanceDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas} [get]
func GetRelevantInstancesActualInstanceDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("comparinginstance=", vars["comparinginstance"])
	res := get_relevant_Instances_actual_instance_Data_Based_On_Limits(vars["comparinginstance"], vars["replicas"], vars["apptype"], vars["appname"],vars["mainservicename"])
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write([]byte(b))
}
// GetRelevantInstanceActualInstanceDataBasedOnLimits godoc
// @Summary Get the resources data for that particular instance from the given host instance on which it is deployed along with the actual Instance data
// @Description Get the resources data for that particular instance from the given host instance on which it is deployed along with the actual Instance data
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param hostinstance query string true " hostinstance on which the comparing instance is deployed "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.NodeValues ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstanceActualInstanceDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{hostinstance}/{comparinginstance}/{replicas} [get]
func GetRelevantInstanceActualInstanceDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("hostInstance=", vars["hostinstance"])
	log.Info("comparinginstance=", vars["comparinginstance"])
	res := get_relevant_Instance_actual_instance_Data_Based_On_Limits(vars["hostinstance"], vars["comparinginstance"],
		vars["replicas"], vars["apptype"], vars["appname"],vars["mainservicename"])
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write([]byte(b))
}
// GetRelevantInstancesRequestsDataBasedOnLimits godoc
// @Summary Get the Request data for that particular instance from all host instances on which it can be possible deployed
// @Description Get the Requests data for that particular instance from all host instances on which it it can be possible deployed
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.RequestsData ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstancesRequestsDataBasedOnLimits/{apptype}/{appname}/{comparinginstance}/{replicas} [get]
func GetRelevantInstancesRequestsDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("comparinginstance=", vars["comparinginstance"])

	res := get_relevant_instances_requests_data_Based_on_limits(vars["comparinginstance"], vars["replicas"], vars["apptype"], vars["appname"],vars["mainservicename"])

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write([]byte(b))

}
// GetRelevantInstanceRequestsDataBasedOnLimits godoc
// @Summary Get the Requests data for that particular instance from the given host instance on which it is deployed
// @Description Get the Requests data for that particular instance from the given host instance on which it is deployed
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param hostinstance query string true " hostinstance on which the comparing instance is deployed "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.RequestsData ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstanceRequestsDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{hostinstance}/{comparinginstance}/{replicas} [get]
func GetRelevantInstanceRequestsDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("comparinginstance=", vars["comparinginstance"])
	log.Info("hostInstance=", vars["hostinstance"])
	res := get_relevant_instance_requests_data_Based_on_limits(vars["hostinstance"], vars["comparinginstance"], vars["replicas"], vars["apptype"], vars["appname"],vars["mainservicename"])

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write([]byte(b))

}
// GetRelevantInstancesActualInstanceRequestsDataBasedOnLimits godoc
// @Summary Get the requests data for that particular instance from all host instances on which it can be possible deployed and also the actual VM performance
// @Description Get the requests data for that particular instance from all host instances on which it can be possible deployed and also the actual VM performance
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.RequestsData ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstancesActualInstanceRequestsDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas} [get]
func GetRelevantInstancesActualInstanceRequestsDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("comparinginstance=", vars["comparinginstance"])
	res := get_relevant_instances_actual_instance_requests_data_based_on_limits(vars["comparinginstance"], vars["replicas"], vars["apptype"], vars["appname"],vars["mainservicename"])

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write([]byte(b))
}
// GetRelevantInstanceActualInstanceRequestsDataBasedOnLimits godoc
// @Summary Get the requests data for that particular instance from the given host instance on which it is deployed along with the actual Instance data
// @Description Get the requests data for that particular instance from the given host instance on which it is deployed along with the actual Instance data
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param hostinstance query string true " hostinstance on which the comparing instance is deployed "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.RequestsData ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRelevantInstanceActualInstanceRequestsDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{hostinstance}/{comparinginstance}/{replicas} [get]
func GetRelevantInstanceActualInstanceRequestsDataBasedOnLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("comparinginstance=", vars["comparinginstance"])
	log.Info("hostInstance=", vars["hostinstance"])
	res := get_relevant_instance_actual_instance_requests_data_based_on_limits(vars["hostinstance"], vars["comparinginstance"], vars["replicas"], vars["apptype"], vars["appname"], vars["mainservicename"])

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write([]byte(b))

}
// GetInstancesPerfData godoc
// @Summary Get the resources data for that particular instance from all host instances on which it can be possible deployed and also the actual VM performance
// @Description Get the resources data for that particular instance from all host instances on which it can be possible deployed and also the actual VM performance
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.NodeValues ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getInstancesPerfData/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}/{hostinstance} [get]
func GetInstancesPerfData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("hostinstance=", vars["hostinstance"])
	log.Info("comparinginstance=", vars["comparinginstance"])
	instances := strings.Split(vars["comparinginstance"], ",")
	replicas := strings.Split(vars["replicas"], ",")
	res := get_instances_perf_data(instances, replicas, vars["apptype"], vars["appname"], vars["mainservicename"], vars["hostinstance"])

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write([]byte(b))

}
// GetInstancesRequestsData godoc
// @Summary Get the requests data for that particular instance from the given host instance on which it is deployed along with the actual Instance data
// @Description Get the requests data for that particular instance from the given host instance on which it is deployed along with the actual Instance data
// @Tags CONDUCTED_EXPERIMENTS_DATA_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "number of replicas of service "
// @Param hostinstance query string true " hostinstance on which the comparing instance is deployed "
// @Param comparinginstance query string true " Instance to compare "
// @Success 200 {array} TERMINUS.RequestsData ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getInstancesRequestsData/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}/{hostinstance} [get]
func GetInstancesRequestsData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	log.Info("replicas=", vars["replicas"])
	log.Info("hostinstance=", vars["hostinstance"])
	log.Info("comparinginstance=", vars["comparinginstance"])
	instances := strings.Split(vars["comparinginstance"], ",")
	replicas := strings.Split(vars["replicas"], ",")
	res := get_instances_requests_data(instances, replicas, vars["apptype"], vars["appname"], vars["mainservicename"],  vars["hostinstance"])

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write([]byte(b))

}
// DoBruteForceTraining godoc
// @Summary Do brute force training for a particular app
// @Description Do brute force training for a particular app
// @Tags TRAINING_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {string} string "Started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /doBruteForceTraining/{apptype}/{appname}/{mainservicename}[get]
func DoBruteForceTraining(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	go func(appName, appType, mainServiceName string) {
		//#http://localhost:9002/pretrain?appname=primeapp&apptype=compute&instancefamily=t2&colname=ALL_BRUTE_FORCE_CONDUCTED_TEST_NAMES&dbname=TERMINUS
		instanceFamily := "t2"
		url := "http://regression:9002/pretrain?appname=" + appName + "&apptype=" + appType + "&instancefamily=" +
			instanceFamily + "&colname=" + Collection_All_Test_Names + "&dbname=" + Database+
			"&mainServiceName="+ mainServiceName
		log.Info(url)
		response, err := http.Get(url)
		if err != nil {
			log.Error("%s", err)
			http.Error(w, "cannot query to regression service", http.StatusBadRequest)
			return
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Error("%s", err)
				http.Error(w, "contents returned from regression are not readable", http.StatusBadRequest)
				return
			}
			log.Info(string(contents))
			//w.Write([]byte(string(contents)));
			mongoSession := GetMongoSession()
			collection := mongoSession.DB(Database).C(Collection_All_TrainingsError)

			var foundTrainingObject TrainingObject
			errIns := collection.Find(bson.M{"servicename": appName,"mainservicename":mainServiceName  }).One(&foundTrainingObject)
			if errIns != nil {
				log.Info("error ", errIns)
				AllData := TrainingObject{
					ServiceType:    appType,
					ServiceName:    appName,
					MainServiceName: mainServiceName,
					InstanceFamily: instanceFamily,
					RMSErrorPerMin: StringToFloat(string(contents)),
				}

				if err := collection.Insert(AllData); err != nil {
					log.Info("error ", err)
				} else {
					log.Info("#inserted into ", Collection_All_TrainingsError)
				}
			} else {
				errMonFin := collection.Update(bson.M{"servicename": appName,"mainservicename":mainServiceName}, bson.M{"$set": bson.M{"rmserrorpermin": StringToFloat(string(contents))}})
				if errMonFin != nil {
					log.Info("Error::%s", errMonFin)
				}
			}
			defer mongoSession.Close()
		}
	}(vars["appname"],vars["apptype"],vars["mainservicename"] )

	w.Write([]byte("Started"));

}
// DoSmartTestAllTraining godoc
// @Summary Do smart training for a particular app
// @Description Do smart training for a particular app
// @Tags TRAINING_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {string} string "Started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /doSmartTestAllTraining/{apptype}/{appname}/{mainservicename}[get]
func DoSmartTestAllTraining(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])

	go func(appName, appType,mainservicename string) {

		var resultDbNameInfo []InfluDbTabaseNamesInfo
		mongoSession := GetMongoSession()
		collection := mongoSession.DB(Database).C(Collection_All_Test_Names)
		err := collection.Find(bson.M{
			"servicetype": appType,
			"servicename": appName}).All(&resultDbNameInfo)
		if err != nil {
			log.Info("Db Error : ", err)
		}
		defer mongoSession.Close()
		for _, dbInfo := range resultDbNameInfo {

			folderName := dbInfo.FolderName
			instanceFamily := "t2"
			//http://localhost:9002/smartTestTrain?appname=primeapp&apptype=compute&instancefamily=t2&containerName=s1t1rc1nc1t2xlargecomputeprimeappt2nanob8j

			url := "http://regression:9002/smartTestTrain?appname=" + appName + "&apptype=" + appType + "&instancefamily=" +
				instanceFamily + "&containerName=" + folderName + "&mainServiceName="+mainservicename
			log.Info(url)
			response, err := http.Get(url)
			if err != nil {
				log.Error("%s", err)
				http.Error(w, "cannot query to regression service", http.StatusBadRequest)
				return
			} else {
				defer response.Body.Close()
				contents, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Error("%s", err)
					http.Error(w, "contents returned from regression are not readable", http.StatusBadRequest)
					return
				}
				log.Info(string(contents))
				//w.Write([]byte(string(contents)));
				mongoSession := GetMongoSession()
				collection := mongoSession.DB(Database).C(Collection_All_Smart_Trainings_Error)

				var foundTrainingObject SmartTrainingErrorObject
				errIns := collection.Find(bson.M{"servicename": appName, "foldername": folderName, "mainservicename":mainservicename }).One(&foundTrainingObject)
				if errIns != nil {
					log.Info("error ", errIns)
					AllData := SmartTrainingErrorObject{
						ServiceType:    	appType,
						ServiceName:    	appName,
						MainServiceName:    mainservicename,
						InstanceFamily: 	instanceFamily,
						RMSErrorPerMin: 	StringToFloat(string(contents)),
						FolderName:     	folderName,
					}

					if err := collection.Insert(AllData); err != nil {
						log.Info("error ", err)
					} else {
						log.Info("#inserted into ", Collection_All_Smart_Trainings_Error)
					}
				} else {
					errMonFin := collection.Update(bson.M{"servicename": appName, "foldername": folderName, "mainservicename":mainservicename}, bson.M{"$set": bson.M{"rmserrorpermin": StringToFloat(string(contents))}})
					if errMonFin != nil {
						log.Info("Error::%s", errMonFin)
					}
				}
				defer mongoSession.Close()
			}
		}
	}(vars["appname"],vars["apptype"],vars["mainservicename"] )

	w.Write([]byte("Started"));
}
// DoBruteForceTrainingReplicas godoc
// @Summary Do brute force training for replicas prediction for a particular app
// @Description Do brute force training for replicas prediction for a particular app
// @Tags TRAINING_MSC
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {string} string "Started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /doBruteForceTrainingReplicas/{apptype}/{appname}/{mainservicename}[get]
func DoBruteForceTrainingReplicas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])

	go func(appName, appType,mainservicename string) {

		//#http://localhost:9002/pretrain?appname=primeapp&apptype=compute&instancefamily=t2&colname=ALL_BRUTE_FORCE_CONDUCTED_TEST_NAMES&dbname=TERMINUS
		instanceFamily := "t2"
		url := "http://regression:9002/trainedreplicas?appname=" + appName + "&apptype=" + appType + "&instancefamily=" +
			instanceFamily + "&colname=" + Collection_All_Test_Names + "&dbname=" + Database + "&mainServiceName="+mainservicename
		log.Info(url)
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "cannot query to regression service", http.StatusBadRequest)
			return
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Error("%s", err)
				http.Error(w, "contents returned from regression are not readable", http.StatusBadRequest)
				return
			}
			log.Info(string(contents))
			//w.Write([]byte(string(contents)));
			mongoSession := GetMongoSession()
			collection := mongoSession.DB(Database).C(Collection_All_Replicas_TrainingsError)

			var foundTrainingObject TrainingObject
			errIns := collection.Find(bson.M{"servicename": appName, "mainservicename":mainservicename}).One(&foundTrainingObject)
			if errIns != nil {
				log.Info("error ", errIns)
				AllData := TrainingObject{
					ServiceType:    appType,
					ServiceName:    appName,
					MainServiceName:    mainservicename,
					InstanceFamily: instanceFamily,
					RMSErrorPerMin: StringToFloat(string(contents)),
				}

				if err := collection.Insert(AllData); err != nil {
					log.Info("error ", err)
				} else {
					log.Info("#inserted into ", Collection_All_Replicas_TrainingsError)
				}
			} else {
				errMonFin := collection.Update(bson.M{"servicename":appName,"mainservicename":mainservicename }, bson.M{"$set": bson.M{"rmserrorpermin": StringToFloat(string(contents))}})
				if errMonFin != nil {
					log.Info("Error::%s", errMonFin)
				}
			}
			defer mongoSession.Close()
		}
	}(vars["appname"],vars["apptype"],vars["mainservicename"] )

	w.Write([]byte("Started"));
}
// GetPredictedRegressionTRN godoc
// @Summary Get predicted MSC value based on the inputs (in seconds)
// @Description Get predicted MSC value based on the inputs
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param replicas query string true "replicas for application"
// @Param numcoresutil query string true "num of cores utilization"
// @Param numcoreslimit query string true "num of cores limit"
// @Param nummemlimit query string true "num of mem limit"
// @Success 200 {object} TERMINUS.MSCInfo "Predicted MSC replicas"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPredictedRegressionTRN/{apptype}/{appname}/{mainservicename}{replicas}/{numcoresutil}/{numcoreslimit}/{nummemlimit} [get]
func GetPredictedRegressionTRN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("replicas=", vars["replicas"])
	log.Info("numcoresutil=", vars["numcoresutil"])
	log.Info("numcoreslimit=", vars["numcoreslimit"])
	log.Info("nummemlimit=", vars["nummemlimit"])
	log.Info("mainservicename=", vars["mainservicename"])
	instancefamily:="t2"
	requestDuration := "1000"
	//# http://localhost:9002/getPredictionPreTrained?appname=primeapp&apptype=compute&replicas=2&numcoresutil=0.1
	//# &numcoreslimit=0.1&nummemlimit=0.1&instancefamily=t2&requestduration=1000

	url := "http://regression:9002/getPredictionPreTrained?appname=" + vars["appname"] + "&apptype=" + vars["apptype"] +
		"&replicas=" + vars["replicas"]+
		"&numcoresutil=" + vars["numcoresutil"]+
		"&numcoreslimit=" + vars["numcoreslimit"]+
		"&nummemlimit=" + vars["nummemlimit"]+
		"&instancefamily=" + instancefamily+
		"&requestduration=" + requestDuration+
		"&mainServiceName="+ vars["mainservicename"]
	log.Info(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "cannot query to regression service", http.StatusBadRequest)
		return
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "contents returned from regression are not readable", http.StatusBadRequest)
			return
		}
		log.Info(string(contents))
		mongoSessionPodBooting := GetMongoSession()
		collectionPodBooting := mongoSessionPodBooting.DB(Database).C(Collection_Pod_Bootingtime)

		var podBootingTimeObject PodBootingTimeStruct
		replicaFlot64:=StringToFloat(vars["replicas"])
		errPod := collectionPodBooting.Find(bson.M{"replicas": replicaFlot64, "appname": vars["appname"], "apptype": vars["apptype"], "mainservicename": vars["mainservicename"]}).One(&podBootingTimeObject)
		if errPod != nil {
			log.Info("error ", errPod)
		}
		defer mongoSessionPodBooting.Close()

		mscInformation:= MSCInfo{Replicas: replicaFlot64,
			Maximum_service_capacity_per_min: MSCDetails{RegBruteForce: StringToFloat(string(contents))},
			Maximum_service_capacity_per_sec: MSCDetails{RegBruteForce: StringToFloat(string(contents))/60},
			Pod_boot_time_ms: podBootingTimeObject.MeanBootingTime, Sd_Pod_boot_time_ms:podBootingTimeObject.SdBootingTime,}

		all, err := json.Marshal(mscInformation)

		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(all)

	}
}

func getPredictionMsc(replicas, appName, appType, numCoresUtil, numCoresLimit,nummemlimit, instanceFamily, requestDuration,mainservicename string) float64 {

	url := "http://regression:9002/getPredictionPreTrained?appname=" + appName + "&apptype=" + appType +
		"&replicas=" + replicas +
		"&numcoresutil=" + numCoresUtil +
		"&numcoreslimit=" + numCoresLimit +
		"&nummemlimit=" + nummemlimit +
		"&instancefamily=" + instanceFamily +
		"&requestduration=" + requestDuration+
		"&mainServiceName="+ mainservicename
	log.Info(url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		return 0
	} else {
		defer response.Body.Close()
		predictedMSC, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return 0
		}
		log.Info(string(predictedMSC))
		predictedMSCFloat := StringToFloat(string(predictedMSC)) / 60

		return predictedMSCFloat
	}
}
// GetPredictedRegressionReplicas godoc
// @Summary Get Number of replicas predicted
// @Description Get Number of replicas predicted
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param msc query string true "msc expected per seconds (RPS) "
// @Param numcoresutil query string true "num of cores utilization"
// @Param numcoreslimit query string true "num of cores limit"
// @Param nummemlimit query string true "num of mem limit"
// @Success 200 {object} TERMINUS.MSCInfo "Predicted MSC replicas"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPredictedRegressionReplicas/{apptype}/{appname}/{mainservicename}{msc}/{numcoresutil}/{numcoreslimit}/{nummemlimit} [get]
func GetPredictedRegressionReplicas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("msc=", vars["msc"]) // its in seconds
	log.Info("numcoresutil=", vars["numcoresutil"])
	log.Info("numcoreslimit=", vars["numcoreslimit"])
	log.Info("nummemlimit=", vars["nummemlimit"])
	log.Info("mainservicename=", vars["mainservicename"])
	instancefamily := "t2"
	requestDuration := "1000"

	//# http://localhost:9002/getPredictionPreTrained?appname=primeapp&apptype=compute&replicas=2&numcoresutil=0.1
	//# &numcoreslimit=0.1&nummemlimit=0.1&instancefamily=t2&requestduration=1000

	providedMSCFloat := StringToFloat(string(vars["msc"]))
	urlReplicasPrediction := "http://regression:9002/getPredictionReplicas?appname=" + vars["appname"] + "&apptype=" + vars["apptype"] +
		"&msc=" + FloatToString(StringToFloat(string(vars["msc"]))*60) +
		"&numcoresutil=" + vars["numcoresutil"] +
		"&numcoreslimit=" + vars["numcoreslimit"] +
		"&nummemlimit=" + vars["nummemlimit"] +
		"&instancefamily=" + instancefamily +
		"&requestduration=" + requestDuration+
		"&mainServiceName="+ vars["mainservicename"]
	response, err := http.Get(urlReplicasPrediction)
	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "cannot query to regression service", http.StatusBadRequest)
		return
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "contents returned from regression are not readable", http.StatusBadRequest)
			return
		}
		result := StringToFloat(string(contents))
		result = math.Ceil(result)
		log.Info(result)
		predictedReplicas := FloatToString(result)
		predictedReplicasFloat := result
		mongoSessionPodBooting := GetMongoSession()
		collectionPodBooting := mongoSessionPodBooting.DB(Database).C(Collection_Pod_Bootingtime)
		var podBootingTimeObject PodBootingTimeStruct
		replicaFlot64 := result
		errPod := collectionPodBooting.Find(bson.M{"replicas": replicaFlot64, "appname": vars["appname"], "apptype": vars["apptype"], "mainservicename": vars["mainservicename"]}).One(&podBootingTimeObject)
		if errPod != nil {
			log.Info("error ", errPod)
		}
		defer mongoSessionPodBooting.Close()
		var predictedMSCFloatPerSec float64
		predictedMSCFloatPerSec = getPredictionMsc(predictedReplicas, vars["appname"], vars["apptype"], vars["numcoresutil"], vars["numcoreslimit"],
			vars["nummemlimit"], instancefamily, requestDuration,  vars["mainservicename"])
		for i := 0; i < 10; i++ { // takes 10 tests
			if predictedMSCFloatPerSec < providedMSCFloat {
				predictedReplicasFloat = predictedReplicasFloat + 1
				predictedMSCFloatPerSec = getPredictionMsc(FloatToString(predictedReplicasFloat), vars["appname"], vars["apptype"], vars["numcoresutil"], vars["numcoreslimit"],
					vars["nummemlimit"], instancefamily, requestDuration, vars["mainservicename"])
			} else {
				// found
				break
			}
		}
		log.Info(FloatToString(predictedMSCFloatPerSec))
		mscInformation := MSCInfo{Replicas: predictedReplicasFloat,
			Maximum_service_capacity_per_min: MSCDetails{RegBruteForce: predictedMSCFloatPerSec * 60},
			Maximum_service_capacity_per_sec: MSCDetails{RegBruteForce: predictedMSCFloatPerSec},
			Pod_boot_time_ms: podBootingTimeObject.MeanBootingTime, Sd_Pod_boot_time_ms: podBootingTimeObject.SdBootingTime,}
		all, err := json.Marshal(mscInformation)

		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(all)
		return
	}
	all, err := json.Marshal(MSCInfo{})

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
	return
}
// GetPredictedRegressionTRNSmart godoc
// @Summary Get MSC prediction based on smart test for a particular test
// @Description Get Number of replicas predicted for a particular test
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param testname query string true "name of the test"
// @Param numcoresutil query string true "num of cores utilization"
// @Param nummemUtil query string true "num of mem limit"
// @Success 200 {string} string "Predicted MSC"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPredictedRegressionTRNSmart/{apptype}/{appname}/{mainservicename}{numcoresutil}{nummemutil}{testname} [get]
func GetPredictedRegressionTRNSmart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("numcoresutil=", vars["numcoresutil"])
	log.Info("nummemutil=", vars["nummemutil"])
	log.Info("containerName=", vars["testname"])
	log.Info("mainservicename=", vars["mainservicename"])
	instancefamily:="t2"
	requestDuration := "1000"
	//#http://localhost:9002/smartTestGetResult?appname=primeapp&apptype=compute&numcoresutil=0.1
	//# &nummemutil=0.1&instancefamily=t2&requestduration=1000&containerName=s1t1rc1nc1t2xlargecomputeprimeappt2nanob8j

	url := "http://regression:9002/smartTestGetResult?appname=" + vars["appname"] + "&apptype=" + vars["apptype"] +
		"&numcoresutil=" + vars["numcoresutil"]+
		"&nummemutil=" + vars["nummemutil"]+
		"&instancefamily=" + instancefamily+
		"&requestduration=" + requestDuration+
		"&containerName=" + vars["testname"]+
		"&mainServiceName="+ vars["mainservicename"]
	log.Info(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "cannot query to regression service", http.StatusBadRequest)
		return
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "contents returned from regression are not readable", http.StatusBadRequest)
			return
		}
		log.Info(string(contents))
		w.Write([]byte(string(contents)));
	}
}
// StoreAllExperimentalTRN godoc
// @Summary Get and store all MSC experimental, predicted and smart along with replicas predicted in DB
// @Description Get and store all MSC experimental, predicted and smart along with replicas predicted in DB
// @Tags STORE_MSCs
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Completed"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /storeAllExperimentalTRN [get]
func StoreAllExperimentalTRN(w http.ResponseWriter, r *http.Request) {
	storeAllTRNIntoMongo()
	w.Write([]byte("Completed"));
}
// StoreAllTRNRegressionIntoMongo godoc
// @Summary Get and store all MSC predicted  in DB
// @Description Get and store all MSC predicted  in DB
// @Tags STORE_MSCs
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Completed"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /storeAllTRNRegressionIntoMongo [get]
func StoreAllTRNRegressionIntoMongo(w http.ResponseWriter, r *http.Request) {
	go storeAllTRNRegressionIntoMongo()
	w.Write([]byte("Completed"));
}
// GetAllExperimentalTRN godoc
// @Summary Get all MSC experimental, predicted and smart along with replicas predicted for all apps from DB
// @Description Get all MSC experimental, predicted and smart along with replicas predicted for all apps from DB
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Success 200 {array} TERMINUS.MSCValueObject ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAllExperimentalTRN [get]
func GetAllExperimentalTRN(w http.ResponseWriter, r *http.Request) {
	result:=getAllTRNsMongo()
	all, err := json.Marshal(result)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
}
// GetAllExperimentalTRN godoc
// @Summary Get all MSC experimental, predicted and smart along with replicas predicted for an app from DB
// @Description Get all MSC experimental, predicted and smart along with replicas predicted for an app from DB
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {array} TERMINUS.MSCValueObject ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getExperimentalTRNsMongoDBAll/{apptype}/{appname}/{mainservicename} [get]
func GetExperimentalTRNsMongoDBAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	result:=getTRNsMongoAll( vars["appname"],vars["apptype"], vars["mainservicename"])
	all, err := json.Marshal(result)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
}
// GetAllExperimentalTRN godoc
// @Summary Get all MSC experimental, predicted and smart along with replicas predicted for an app on a particular instance type from DB
// @Description Get all MSC experimental, predicted and smart along with replicas predicted for an app on a particular instance type from DB
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Param hostinstance query string true "type of host instance"
// @Success 200 {array} TERMINUS.MSCValueObject ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getExperimentalTRNsMongoDB/{apptype}/{appname}/{mainservicename}{hostinstance} [get]
func GetExperimentalTRNsMongoDB(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("hostInstanceType=", vars["hostinstance"])
	log.Info("mainservicename=", vars["mainservicename"])

	result:=getTRNsMongo( vars["appname"],vars["apptype"], vars["hostinstance"], vars["mainservicename"])
	all, err := json.Marshal(result)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
}
// GetAllExperimentalTRN godoc
// @Summary Get  MSCs for different configurations predicted for an app  using regression
// @Description Get  MSCs for different configurations predicted for an app  using regression
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {object} TERMINUS.MSCValueObject ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRegressionTRNsMongoDBAll/{apptype}/{appname}/{mainservicename} [get]
func GetRegressionTRNsMongoDBAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	result:=getRegressionTRNsMongoAll( vars["appname"],vars["apptype"], vars["mainservicename"])
	all, err := json.Marshal(result)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
}
// GetRMSErrorMongoDB godoc
// @Summary Get  RMS error trained using the brute force
// @Description Get  RMS error trained using the brute force
// @Tags GET_ANALYZED_MSCs
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {object} TERMINUS.MSCValueObject ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getRMSErrorMongoDB/{apptype}/{appname}/{mainservicename} [get]
func GetRMSErrorMongoDB(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	result:=getRMSTrainingError( vars["appname"],vars["apptype"], vars["mainservicename"])
	all, err := json.Marshal(result)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
}
// AnalyzeKubeEvents godoc
// @Summary Analyze the kube events stored in db from conducting pod booting time test and store them back to db
// @Description Analyze the kube events stored in db from conducting pod booting time test and store them back to db
// @Tags POD
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {string} string "completed"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /analyzeKubeEvents/{apptype}/{appname}/{mainservicename} [get]
func AnalyzeKubeEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	analyzeKubernetesEvents( vars["apptype"], vars["appname"],vars["mainservicename"] )
	w.Write([]byte("Completed"))

}
// GetAnalyzedKubeEvents godoc
// @Summary Get Analyzed Kube events from db
// @Description Get Analyzed Kube events from db
// @Tags POD
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {array} TERMINUS.KubeEventsMongoStoreObject ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAnalyzedKubeEvents/{apptype}/{appname}/{mainservicename} [get]
func GetAnalyzedKubeEvents(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("mainservicename=", vars["mainservicename"])

	events:=getAnalyzedKubernetesEvents(vars["appname"], vars["mainservicename"])

	b, err := json.Marshal(events)
	if err != nil {
		log.Info("error:", err)
	}

	w.Write([]byte(b))
}
// GetPodBootingTime godoc
// @Summary Get Pod Booting time from db
// @Description Get Analyzed Kube events from db
// @Tags POD
// @Accept text/html
// @Produce json
// @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
// @Param apptype query string true "application type[compute, dbaccess, web]"
// @Success 200 {array} TERMINUS.PodBootingTimeStruct ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPodBootingTime/{apptype}/{appname}/{mainservicename} [get]
func GetPodBootingTime(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Info("appname=", vars["appname"])
	log.Info("apptype=", vars["apptype"])
	log.Info("mainservicename=", vars["mainservicename"])
	bootingTime:=getPodBootingTime(vars["appname"], vars["mainservicename"])

	b, err := json.Marshal(bootingTime)
	if err != nil {
		log.Info("error:", err)
	}

	w.Write([]byte(b))
}
/*
func StoreAllNamesToMongo(w http.ResponseWriter, r *http.Request) {
	storeAllNamesInMongo()
	w.Write([]byte("Completed"));
}
*/
// GetAppTypes godoc
// @Summary Get App Types
// @Description Get App Types
// @Tags GENERAL
// @Accept text/html
// @Produce json
// @Success 200 {array} string ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAppTypes [get]
func GetAppTypes(w http.ResponseWriter, r *http.Request) {

	b, err := json.Marshal(SupportedTypesOfMicroservices)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write([]byte(b))
}
// GetAppNames  godoc
// @Summary Get App Names based on app type
// @Description Get App Names based on app type
// @Tags GENERAL
// @Accept text/html
// @Produce json
// @Success 200 {array} string ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAppNames/{apptype} [get]
func GetAppNames(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Info("apptype=", vars["apptype"])

	if (vars["apptype"] == "compute") {
		b, err := json.Marshal(computeAppNames)
		if err != nil {
			fmt.Println("error:", err)
		}

		w.Write([]byte(b))
	} else if (vars["apptype"] == "dbaccess") {
		b, err := json.Marshal(dbaccessAppNames)
		if err != nil {
			fmt.Println("error:", err)
		}

		w.Write([]byte(b))
	} else if (vars["apptype"] == "web") {
		b, err := json.Marshal(webAppNames)
		if err != nil {
			fmt.Println("error:", err)
		}

		w.Write([]byte(b))
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
// GetAllTestsInfo  godoc
// @Summary Get all ongoing test information
// @Description Get all ongoing test information
// @Tags GENERAL
// @Accept text/html
// @Produce json
// @Success 200 {array} TERMINUS.TestInformation ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAllTestsInfo [get]
func GetAllTestsInfo(w http.ResponseWriter, r *http.Request) {

	var AllData []TestInformation

	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_MSC)
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
func GetMetricsAPI(w http.ResponseWriter, r *http.Request) {

	sessionAWS := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
	}))

	svc := cloudwatch.New(sessionAWS)

	now := time.Now()

	input := &cloudwatch.GetMetricDataInput{
		StartTime: aws.Time(now.Add(-120 * time.Minute)),
		EndTime:   aws.Time(now),
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("ec2Data"),
				MetricStat: &cloudwatch.MetricStat{

					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String("i-050c36c809e41e491"),
							},
						},
						MetricName: aws.String("CPUUtilization"),
						Namespace:  aws.String("AWS/EC2"),
					},
					Period: aws.Int64(1),
					Stat:   aws.String("Average"),
					Unit:   aws.String("Percent"),
				},
			},
			{
				Id: aws.String("ec2Data2"),
				MetricStat: &cloudwatch.MetricStat{

					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String("i-08cb2efb755b52573"),
							},
						},
						MetricName: aws.String("CPUUtilization"),
						Namespace:  aws.String("AWS/EC2"),
					},
					Period: aws.Int64(1),
					Stat:   aws.String("Average"),
					Unit:   aws.String("Percent"),
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
