package TERMINUS

import (
	"time"
)
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

type DBTimeStampValues struct {
	Timestamps []int64 `json:"Timestamps"`
	Values []float64 `json:"Values"`
}
type PodValues struct {
	Timestamps []int64 `json:"Timestamps"`
	PodName string `json:"PodName"`
	CpuUtil []float64 `json:"CpuUtil"`
	CpuLimit []float64 `json:"CpuLimit"`
	CpuRequest []float64 `json:"CpuRequest"`
	MemUtil []float64 `json:"MemUtil"`
	MemLimit []float64 `json:"MemLimit"`
	MemRequest []float64 `json:"MemRequest"`
}

type NodeValues struct {
	Timestamps []int64 `json:"Timestamps"`
	InstanceType string `json:"InstanceType"`
	CpuUtil []float64 `json:"CpuUtil"`
	MemUtil []float64 `json:"MemUtil"`
	Cores []float64 `json:"Cores"`
	Mem []float64 `json:"Mem"`
	Pod []PodValues `json:"PodValues"`
	Requests []float64 `json:"Requests"`
}

type RequestsData struct {
	Timestamps []int64 `json:"Timestamps"`
	InstanceType string `json:"InstanceType"`
	Vus []float64 `json:"Vus"`
	Requests []float64 `json:"Requests"`
	ReqDurationPercentile95 []float64 `json:"ReqDurationPercentile95"`
	ReqDurationPercentile90 []float64 `json:"ReqDurationPercentile90"`
	ReqDurationMax []float64 `json:"ReqDurationMax"`
	ReqDurationMin []float64 `json:"ReqDurationMin"`
	ReqDurationMean []float64 `json:"ReqDurationMean"`
	ReqDurationMedian []float64 `json:"ReqDurationMedian"`
}


type MSCDetails struct {

	Experimental float64 `json:"Experimental"`
	RegBruteForce float64 `json:"RegBruteForce"`
	RegSmart  float64 `json:"RegSmart"`

}

type MSCInfo struct {
	Replicas           					float64 `json:"Replicas"`
	Maximum_service_capacity_per_min	MSCDetails `json:"Maximum_service_capacity_per_min"`
	Maximum_service_capacity_per_sec	MSCDetails `json:"Maximum_service_capacity_per_sec"`
	Pod_boot_time_ms					float64 `json:"Pod_boot_time_ms"`
	Sd_Pod_boot_time_ms					float64 `json:"Sd_Pod_boot_time_ms"`
	PredictedReplicas 					float64 `json:"PredictedReplicas"`
}
type LimitsInfo struct{
	Cpu_cores float64 `json:"Cpu_cores"`
	Mem_gb float64 `json:"Mem_gb"`
	Request_per_sercond float64 `json:"Request_per_sercond"`
}

type MSCProfiles struct {
	Limits LimitsInfo `json:"Limits"`
	MSCs		[]MSCInfo `json:"MSCs"`
}
type MSCValueObject struct {
	HostInstanceType   			string `json:"HostInstanceType"`
	ServiceType            		string `json:"ServiceType"`
	ServiceName           		string `json:"ServiceName"`
	MainServiceName           		string `json:"MainServiceName"`
	TestAPI            			string `json:"TestAPI"`
	Profiles 					[]MSCProfiles `json:"Profiles"`
}

type TrainingObject struct {
	ServiceType            		string `json:"ServiceType"`
	ServiceName           		string `json:"ServiceName"`
	MainServiceName           		string `json:"MainServiceName"`
	InstanceFamily            	string `json:"InstanceFamily"`
	RMSErrorPerMin				float64 `json:"RMSErrorPerMin"`

}
type SmartTrainingErrorObject struct {
	ServiceType            		string `json:"ServiceType"`
	ServiceName           		string `json:"ServiceName"`
	MainServiceName           		string `json:"MainServiceName"`
	InstanceFamily            	string `json:"InstanceFamily"`
	RMSErrorPerMin				float64 `json:"RMSErrorPerMin"`
	FolderName					string `json:"FolderName"`

}

type InfluDbTabaseNamesInfo struct {
	HostInstanceType   			string `json:"HostInstanceType"`
	ServiceType            		string `json:"ServiceType"`
	ServiceName           		string `json:"ServiceName"`
	TestAPI            			string `json:"TestAPI"`
	Replicas           			string `json:"Replicas"`
	LimitingInstanceType		string `json:"LimitingInstanceType"`
	Limits 						LimitsInfo `json:"Limits"`
	FolderName 					string `json:"FolderName"`
	InfluxDbk6Name 				string `json:"InfluxDbk6Name"`
	InfluxdbK8sName 			string `json:"InfluxdbK8sName"`
	NodeCount 					string `json:"NodeCount"`
	NumMicroSvs 				string `json:"NumMicroSvs"`
	TestIter 					string `json:"TestIter"`
}
/*
type CompareLimitsInstancesQueryStruct struct {
	AppType, AppName, HostInstance1,HostInstance2, LimitingInstance, Replicas string
}

type LimitsInstancesDatabase struct {
	AppType, AppName, HostInstance, LimitingInstance, Replicas string
}
type NonLimitsInstancesDatabase struct {
	AppType, AppName, HostInstance, Replicas string
}
//var LimitsInstancesDatabaseNames 		= make(map[LimitsInstancesDatabase]string)

//var NonLimitsInstancesDatabaseNames 		= make(map[NonLimitsInstancesDatabase]string)
*/

type EventObjectDescribe struct {
	AppType, AppName, Replicas, MainServiceName, Name string
}
type EventObjectItems struct {
	Kind string
	Reason string
	Message string
	FirstTimestamp time.Time
	LastTimestamp time.Time
}
type KubeEventsMongoStoreObject struct {
	Replicas string `json:"Replicas"`
	AppName string `json:"AppName"`
	AppType string `json:"AppType"`
	MainServiceName string  `json:"MainServiceName"`
	Name string  `json:"Name"`
	TotalTime time.Duration `json:"TotalTime"`
	KubeEvents []EventObjectItems `json:"KubeEvents"`
}

var KubeEventsAnalyzed 		= make(map[EventObjectDescribe][]EventObjectItems)
var PodBootTime				= make(map[string] []time.Duration)
var PodBootTimeTests			= make(map[string] int64)

type PodBootingTimeStruct struct {
	Replicas float64 `json:"Replicas"`
	AppName string `json:"AppName"`
	AppType string `json:"AppType"`
	MainServiceName string  `json:"MainServiceName"`
	MeanBootingTime float64 `json:"MeanBootingTime"`
	SdBootingTime float64 `json:"SdBootingTime"`
	NumTests int64 `json:"NumTests"`

}

var AWSConfig 							AWSConfigStruct
var AllInstanceTypes 					[]string
var DefaultRegion 						[]string
var DefaultZone							[]string
var DefaultAMI 							[]string
var B 									*Broker
var SupportedNumMicroservicesApp		[]string
var DefaultAppNames 					[]string
var SupportedTypesOfMicroservices 		[]string
var AppsAPIEndPoint						[]string
var SupportedNumMicroserviceReplicas	[]string
var TestNumMicroserviceReplicas			[]string
var MaxTestIterations					int
var computeAppNames 					[]string
var dbaccessAppNames 					[]string
var webAppNames 						[]string
var mixAppNames 						[]string
var sandboxAppNames 						[]string
var mainServiceNamesMix 						[]string