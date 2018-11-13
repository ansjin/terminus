package TERMINUS

import (
	"time"
	"github.com/sirupsen/logrus"
)
type KopsConfig struct {
	AwsAccessKeyId   		string `json:"AwsAccessKeyId"`
	AwsSecretAccessKey   	string `json:"AwsSecretAccessKey"`
	Region 					string `json:"Region"`
	ContainerName 			string `json:"ContainerName"`
	KubeClusterName 		string `json:"KubeClusterName"`
	S3BucketName 			string `json:"S3BucketName"`
	Zone 					string `json:"Zone"`
	NodeCount 				string `json:"NodeCount"`
	NodeSize 				string `json:"NodeSize"`
	MasterSize 				string `json:"MasterSize"`
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
type ScriptTypeNames struct {
	InitClusterScript       string 		`json:"InitClusterScript"`
	CreateClusterScript     string 		`json:"CreateClusterScript"`
	UpdateClusterScript     string 		`json:"UpdateClusterScript"`
	ValidateClusterScript   string 		`json:"ValidateClusterScript"`
	GetPasswdScript 		string 		`json:"GetPasswdScript"`
	GetTokenScript 			string 		`json:"GetTokenScript"`
	DeleteClusterScript		string 		`json:"DeleteClusterScript"`
}
type K6QueryObj struct {
	Scriptname 				string 		`json:"scriptname"`
	Dbname 					string		`json:"dbname"`
	Testurl 				string		`json:"testurl"`
	Vus 					string		`json:"vus"`
	Rps 					string		`json:"rps"`

}
type TestInformation struct {
	Timestamp          		int64
	DbName             		string
	Phase              		string
	NodeCount          		string
	NodeInstanceType   		string
	MasterInstanceType 		string
	NumMicroSvs        		string
	TypeApp            		string
	AppName            		string
	Limits             		bool
	Replicas           		string
	TestIter           		string
	TestURL            		string
}
type LoadTestQueryObj struct {
	NumMicroservice string `json:"NumMicroservice"`
	IterReplicas string `json:"IterReplicas"`
	TestIter string `json:"TestIter"`
	NodeCount string `json:"NodeCount"`
	InstanceType string `json:"InstanceType"`
	MasterType string `json:"MasterType"`
	Limits string `json:"Limits"`
	AppName string `json:"AppName"`
	AppType string `json:"AppType"`
	TestApi string `json:"TestApi"`
	LimitingResourcesInstanceType string `json:"LimitingResourcesInstanceType"`
	MaxRPS string  `json:"MaxRPS"`
	MainServiceName string `json:"MainServiceName"`

}
var UserConfig 				KopsConfig
var AllScriptsNames 		ScriptTypeNames
var AllInstanceTypes 					[]string
var DefaultRegion 						[]string
var DefaultZone							[]string
var DefaultAMI 							[]string
var B 									*Broker
var SupportedNumMicroservicesApp		[]string
var DefaultAppNames 					[]string
var SupportedTypesOfMicroservices 		[]string
var SupportedNumMicroserviceReplicas	[]string
var LogVar *logrus.Logger
