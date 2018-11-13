package sandbox_microservices

import (
	"github.com/sirupsen/logrus"
	"time"
	"net/http"
)

type ServicesStruct struct {
	Container_name       string   `json:"container_name,omitempty"`
	Build       string   `json:"build,omitempty"`
	Image       string   `json:"image,omitempty"`
	Ports       []string `json:"ports,omitempty"`
	Environment []string `json:"environment,omitempty"`
	Depends_on   []string `json:"depends_on,omitempty"`
	Links       []string `json:"links,omitempty"`
	Volumes		 []string `json:"volumes,omitempty"`
	Restart 	 string `json:"restart,omitempty"`
	Command		 string `json:"command,omitempty"`
	Ulimits 	map[string]interface{} `json:"ulimits,omitempty"`
}
type DockerComposeParsedStructTemp struct {
	Version  string `json:"version"`
	Services map[string]interface{} `json:"-"` // Rest of the fields should go here.
}
type DockerComposeParsedStruct struct {
	Version  string `json:"version"`
	Services map[string]ServicesStruct `json:"services"` // Rest of the fields should go here.
}
type AWSConfigStruct struct {
	AwsAccessKeyId   		string `json:"AwsAccessKeyId"`
	AwsSecretAccessKey   	string `json:"AwsSecretAccessKey"`
	Region 					string `json:"Region"`
	KeyPairName					string `json:"KeyName"`
	SubnetId				string `json:"SubnetId"`
	SecurityGroup			string `json:"SecurityGroup"`
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

type APIResponseObj struct {

	Headers       				http.Header 		`json:"headers"`
	Protocol     				string 		`json:"protocol"`
	Method      				string 		`json:"method"`
	Hostname      				string 		`json:"hostname"`
	Port 						string 		`json:"port"`
	Api_endpoint 				string 		`json:"api_endpoint"`
	Response				interface{} 		`json:"response"`
}
var B 									*Broker
var LogVar *logrus.Logger
var AWSConfig 							AWSConfigStruct
var AllInstanceTypes 					[]string
var DefaultRegion 						[]string
var DefaultZone							[]string
var DefaultAMI 							[]string