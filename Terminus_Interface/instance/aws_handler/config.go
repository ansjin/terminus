package INSTANCE

import (
	"net/http"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"time"
	"strconv"
)


type AwsUserConfig struct {
	AwsAccessKeyId   string `json:"AwsAccessKeyId"`
	AwsSecretAccessKey   string `json:"AwsSecretAccessKey"`
	Region *string `json:"Region"`
}


type Ec2Instances struct {
	InstanceId       		string 		`json:"InstanceId"`
	InstanceState     		string 		`json:"InstanceState"`
	AvailabilityZone      	string 		`json:"AvailabilityZone"`
	PublicIpAddress      	string 		`json:"PublicIpAddress"`
	InstanceType 			string 		`json:"InstanceType"`
	ImageId 				string 		`json:"ImageId"`
	CoreCount				int64 		`json:"CoreCount"`
	LaunchTime				time.Time 	`json:"LaunchTime"`
}

type Ec2InstancesTime struct {
	Timestamp      			int64  			`json:"timestamp"`
	Instances 				[]Ec2Instances 	`json:"instances"`
}
type InstancesData struct {
	Pending      int64 `json:"Pending"`
	Running      int64 `json:"Running"`
	SshLogin     int64 `json:"SshLogin"`
	ShuttingDown int64 `json:"ShuttingDown"`
	Stopped      int64 `json:"Stopped"`
	Terminated   int64 `json:"Terminated"`
	Other        int64 `json:"Other"`
}

type InstanceBootTime struct {
	Avg float64 `json:"Avg"`
	Min int64 `json:"Min"`
	Max int64 `json:"Max"`
}

type InstanceTime struct {
	Avg float64 `json:"Avg"`
	Min int64 `json:"Min"`
	Max int64 `json:"Max"`
}
type InstanceShutDownTime struct {
	Avg float64 `json:"Avg"`
	Min int64 `json:"Min"`
	Max int64 `json:"Max"`
}

type InstanceBootShutdownRate struct {
	NumInstances   int `json:"NumInstances"`
	NumExperiments int `json:"NumExperiments"`
	BootTime       InstanceTime `json:"BootTime"`
	ShutdownTime   InstanceTime `json:"ShutdownTime"`
}

type ExperimentSetting struct {
	ExperimentNum              string `json:"ExperimentNum"`
	NumInstances               int `json:"NumInstances"`
	Instances                  []InstancesData `json:"Instances"`
	InstancesBootTime          InstanceTime `json:"InstancesBootTime"`
	InstancesShutDownTime      InstanceTime `json:"InstancesShutDownTime"`
	TotalInstancesBootTime     int64 `json:"TotalInstancesBootTime"`
	TotalInstancesShutDownTime int64 `json:"TotalInstancesShutDownTime"`
}
type ExperimentsLoop struct {
	Experiments []ExperimentSetting `json:"Experiments"`
}

type VmTemplateData struct {
	InstanceType     string `json:"InstanceType"`
	Region           string `json:"Region"`
	AvailabilityZone string `json:"AvailabilityZone"`
	ImageId          string `json:"ImageId"`
	CoreCount        int64 `json:"CoreCount"`
	ExperimentLoop   []ExperimentsLoop `json:"ExperimentLoop"`
	BootShutdownRate []InstanceBootShutdownRate `json:"BootShutdownRate"`
}
type PerExperiment struct {
	Instances          []InstancesData `json:"Instances"`
	InstanceType       string `json:"InstanceType"`
	NumInstances       int `json:"NumInstances"`
	InstanceBootRate   InstanceTime `json:"InstanceBootRate"`
	ShutDownRate       InstanceTime `json:"ShutDownRate"`
	ExperimentNum      string `json:"ExperimentNum"`
	ExperimenLoopCount int `json:"ExperimenLoopCount"`
	Region             string `json:"Region"`
	AvailabilityZone   string `json:"AvailabilityZone"`
	ImageId            string `json:"ImageId"`
	CoreCount          int64 `json:"CoreCount"`
	TotalStartTime     int64 `json:"TotalStartTime"`
	TotalShutDownTime  int64 `json:"TotalShutDownTime"`
}

type InstanceValue struct {
	NumInstances int `json:"NumInstances"`
	BootTime float64 `json:"BootTime"`
	ShutDownTime float64 `json:"ShutDownTime"`
}

type InstanceRegression struct {
	InstanceType string `json:"InstanceType"`
	Region string `json:"Region"`
	InstanceValues []InstanceValue `json:"InstanceValues"`
}

type VMBootShutDownRatePerInstanceTypeAll struct {
	InstanceValues []InstanceValue `json:"InstanceValues"`
}

type VMBootShutDownRatePerInstanceTypeOne struct {
	BootTime float64 `json:"BootTime"`
	ShutDownTime float64 `json:"ShutDownTime"`
}



var TotalExperimentsLoop = 5
var InstanceTypes = []string{ "t2.micro"}
var AllInstanceTypes = []string{ "t2.micro", "t2.nano", "t2.small", "t2.medium", "t2.large"}

var DefaultNumInstances = []int64{1,2, 3, 4, 5}

//var DefaultRegion = []string{ "us-east-2", "us-east-1", "us-west-1", "eu-central-1", "ap-south-1", "eu-west-1"}
var DefaultRegion = []string{ "us-east-2"}

//var DefaultAMI = []string{ "ami-5e8bb23b", "ami-759bc50a", "ami-4aa04129", "ami-de8fb135", "ami-188fba77","ami-2a7d75c0"}
var DefaultAMI = []string{ "ami-5e8bb23b"}


func Schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
func StringToFloat(stringVal string) float64 {
	// to convert a float number to a string
	if s, err := strconv.ParseFloat(stringVal, 64); err == nil {
		return s// 3.14159265
	}
	return 0
}
func ValueAssignString(value *string, fallback string) string{

	if value!=nil {
		return *value
	} else {
		return fallback
	}
}
func ValueAssignInt64(value *int64, fallback int64) int64{

	if value!=nil {
		return *value
	} else {
		return fallback
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}


func InitVMStartScriptBootRate(dummyURL string, templateName string) string {

	var vmStartScript =	"#!/bin/bash\n"+
		"mkdir test_VM\n"+
		"myip =\"$(hostname --ip-address)\"\n"+
		"echo \"My WAN/Public IP address: $myip\" >test_VM/test\n"+
		"curl \""+dummyURL+"/?vmip=$myip&templateName="+templateName+"\""

	return vmStartScript

}
var UserConfig =  AwsUserConfig{AwsAccessKeyId: "AKIAJHTP7RGHGKF2DVPA", AwsSecretAccessKey: "KK8FHFtzRv5qfyE7mbxI22WiLIvelkrO/7oVlGnQ",
								Region:  aws.String("us-east-2")}


func InitUserConfigEnvVars(w http.ResponseWriter, r *http.Request){
	UserConfig.AwsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	UserConfig.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	UserConfig.Region = aws.String(os.Getenv("REGION"))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}