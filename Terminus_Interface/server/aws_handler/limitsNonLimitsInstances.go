package TERMINUS

import (
	"sync"
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
)

type LimitingInstanceInfo struct {
	InstanceType string
	Replicas string
}
var RelevantInstancesNames 		= make(map[LimitingInstanceInfo][]string)
/*
func InitLimitsDatabaseNames()  {
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.nano", "100",
		"1"}] = "s1t1rc1nc1t2nanocomputeprimeapp1004n9"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.micro", "100",
		"1"}] = "s1t1rc1nc1t2microcomputeprimeapp100512"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.small", "100",
		"1"}] = "s1t1rc1nc1t2smallcomputeprimeapp100vz7"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "100",
		"1"}] = "s1t1rc1nc1t2mediumcomputeprimeapp100k5z"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "100",
		"1"}] = "s1t1rc1nc1t2largecomputeprimeapp100aki"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "100",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeapp100m49"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.nano", "100",
		"2"}] = "s1t1rc2nc1t2nanocomputeprimeapp100zip"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.micro", "100",
		"2"}] = "s1t1rc2nc1t2microcomputeprimeapp100nn2"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.small", "100",
		"2"}] = "s1t1rc2nc1t2smallcomputeprimeapp10019m"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "100",
		"2"}] = "s1t1rc2nc1t2mediumcomputeprimeapp1003ka"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "100",
		"2"}] = "s1t1rc2nc1t2largecomputeprimeapp1009xe"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "100",
		"2"}] = "s1t1rc2nc1t2xlargecomputeprimeapp100pxv"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.nano", "100",
		"3"}] = "s1t1rc3nc1t2nanocomputeprimeapp100eyd"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.micro", "100",
		"3"}] = "s1t1rc3nc1t2microcomputeprimeapp100998"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.small", "100",
		"3"}] = ""
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "100",
		"3"}] = ""
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "100",
		"3"}] = ""
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "100",
		"3"}] = ""

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.nano", "200",
		"1"}] = "s1t1rc1nc1t2nanocomputeprimeapp2007ie"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.micro", "200",
		"1"}] = "s1t1rc1nc1t2microcomputeprimeapp200hcb"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.small", "200",
		"1"}] = "s1t1rc1nc1t2smallcomputeprimeapp200h2y"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "200",
		"1"}] = "s1t1rc1nc1t2mediumcomputeprimeapp200nbc"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "200",
		"1"}] = "s1t1rc1nc1t2largecomputeprimeapp200ulp"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "200",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeapp200mqr"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.nano", "200",
		"2"}] = "s1t1rc2nc1t2nanocomputeprimeapp200qtk"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.micro", "200",
		"2"}] = "s1t1rc2nc1t2microcomputeprimeapp200uk8"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.small", "200",
		"2"}] = "s1t1rc2nc1t2smallcomputeprimeapp200n5w"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "200",
		"2"}] = "s1t1rc2nc1t2mediumcomputeprimeapp2002np"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "200",
		"2"}] = "s1t1rc2nc1t2largecomputeprimeapp2009i9"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "200",
		"2"}] = "s1t1rc2nc1t2xlargecomputeprimeapp200cgq"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.micro", "500",
		"1"}] = "s1t1rc1nc1t2microcomputeprimeapp500eh1"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.small", "500",
		"1"}] = "s1t1rc1nc1t2smallcomputeprimeapp500ue4"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "500",
		"1"}] = "s1t1rc1nc1t2mediumcomputeprimeapp500q5q"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "500",
		"1"}] = "s1t1rc1nc1t2largecomputeprimeapp5004bz"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "500",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeapp500yzo"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "500",
		"2"}] = "s1t1rc2nc1t2mediumcomputeprimeapp5008c1"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "500",
		"2"}] = "s1t1rc2nc1t2largecomputeprimeapp5008ad"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "500",
		"2"}] = "s1t1rc2nc1t2xlargecomputeprimeapp500tzd"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "500",
		"3"}] = "s1t1rc3nc1t2mediumcomputeprimeapp50053s"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "500",
		"3"}] = "s1t1rc3nc1t2largecomputeprimeapp500gzb"


	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "t2.nano",
		"1"}] = "s1t1rc1nc1t2largecomputeprimeappt2nanobnr"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "t2.micro",
		"1"}] = "s1t1rc1nc1t2largecomputeprimeappt2microbby"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.large", "t2.small",
		"1"}] = "s1t1rc1nc1t2largecomputeprimeappt2smallh1i"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "t2.nano",
		"1"}] = "s1t1rc1nc1t2mediumcomputeprimeappt2nanoaba"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "t2.micro",
		"1"}] = "s1t1rc1nc1t2mediumcomputeprimeappt2micro8hd"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.medium", "t2.small",
		"1"}] = "s1t1rc1nc1t2mediumcomputeprimeappt2smallfyu"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.nano",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeappt2nanob8j"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.micro",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeappt2micro9vq"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.small",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeappt2smallm6o"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.medium",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeappt2mediumpdi"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.large",
		"1"}] = "s1t1rc1nc1t2xlargecomputeprimeappt2largescn"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.nano",
		"2"}] = "s1t1rc2nc1t2xlargecomputeprimeappt2nano51t"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.micro",
		"2"}] = "s1t1rc2nc1t2xlargecomputeprimeappt2microl27"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.small",
		"2"}] = ""
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.medium",
		"2"}] = ""
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"compute", "primeapp", "t2.xlarge", "t2.large",
		"2"}] = ""


	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.large", "t2.nano",
		"1"}] = "s1t1rc1nc1t2largedbaccessmovieappt2nanofch"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.large", "t2.micro",
		"1"}] = "s1t1rc1nc1t2largedbaccessmovieappt2micro8uf"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.large", "t2.small",
		"1"}] = "s1t1rc1nc1t2largedbaccessmovieappt2smallhs9"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.medium", "t2.nano",
		"1"}] = "s1t1rc1nc1t2mediumdbaccessmovieappt2nanomtp"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.medium", "t2.micro",
		"1"}] = "s1t1rc1nc1t2mediumdbaccessmovieappt2microy64"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.medium", "t2.small",
		"1"}] = "s1t1rc1nc1t2mediumdbaccessmovieappt2small3rz"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.xlarge", "t2.nano",
		"1"}] = "s1t1rc1nc1t2xlargedbaccessmovieappt2nanor1n"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.xlarge", "t2.micro",
		"1"}] = "s1t1rc1nc1t2xlargedbaccessmovieappt2micro8zu"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.xlarge", "t2.small",
		"1"}] = "s1t1rc1nc1t2xlargedbaccessmovieappt2smallpn8"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.xlarge", "t2.medium",
		"1"}] = "s1t1rc1nc1t2xlargedbaccessmovieappt2mediumdvw"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"dbaccess", "movieapp", "t2.xlarge", "t2.large",
		"1"}] = "s1t1rc1nc1t2xlargedbaccessmovieappt2largejn8"


	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.medium", "t2.nano",
		"1"}] = "s1t1rc1nc1t2mediumwebwebacappt2nanozz6"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.medium", "t2.micro",
		"1"}] = "s1t1rc1nc1t2mediumwebwebacappt2micro3ni"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.medium", "t2.small",
		"1"}] = "s1t1rc1nc1t2mediumwebwebacappt2smallsnd"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.large", "t2.nano",
		"1"}] = "s1t1rc1nc1t2largewebwebacappt2nano49g"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.large", "t2.micro",
		"1"}] = "s1t1rc1nc1t2largewebwebacappt2micro2w4"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.large", "t2.small",
		"1"}] = "s1t1rc1nc1t2largewebwebacappt2smalloj1"

	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.xlarge", "t2.nano",
		"1"}] = "s1t1rc1nc1t2xlargewebwebacappt2nanogxk"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.xlarge", "t2.micro",
		"1"}] = "s1t1rc1nc1t2xlargewebwebacappt2microx7x"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.xlarge", "t2.small",
		"1"}] = "s1t1rc1nc1t2xlargewebwebacappt2small5u7"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.xlarge", "t2.medium",
		"1"}] = "s1t1rc1nc1t2xlargewebwebacappt2mediumjxf"
	LimitsInstancesDatabaseNames[LimitsInstancesDatabase{"web", "webacapp", "t2.xlarge", "t2.large",
		"1"}] = "s1t1rc1nc1t2xlargewebwebacappt2large4i9"
}
func InitNonLimitsDatabaseNames()  {

	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"compute", "primeapp", "t2.nano",
		"1"}] = "sv1t1rc1nc1t2nanocomputeprimeappfalse82vj"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"compute", "primeapp", "t2.micro",
		"1"}] = "sv1t1rc1nc1t2microcomputeprimeappfalsexybp"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"compute", "primeapp", "t2.small",
		"1"}] = "sv1t1rc1nc1t2smallcomputeprimeappfalsegfqh"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"compute", "primeapp", "t2.medium",
		"1"}] = "sv1t1rc1nc1t2mediumcomputeprimeappfalse7mpr"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"compute", "primeapp", "t2.large",
		"1"}] = "sv1t1rc1nc1t2largecomputeprimeappfalse85hp"

	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"dbaccess", "movieapp", "t2.nano",
		"1"}] = "sv1t1rc1nc1t2nanodbaccessmovieappfalsezoss"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"dbaccess", "movieapp", "t2.micro",
		"1"}] = "sv1t1rc1nc1t2microdbaccessmovieappfalseloya"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"dbaccess", "movieapp", "t2.small",
		"1"}] = "sv1t1rc1nc1t2smalldbaccessmovieappfalsec2tp"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"dbaccess", "movieapp", "t2.medium",
		"1"}] = "sv1t1rc1nc1t2mediumdbaccessmovieappfalsedv3acs"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"dbaccess", "movieapp", "t2.large",
		"1"}] = "sv1t1rc1nc1t2largedbaccessmovieappfalsee1hics"

	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"web", "webacapp", "t2.nano",
		"1"}] = "sv1t1rc1nc1t2nanowebwebacappfalsej3za"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"web", "webacapp", "t2.micro",
		"1"}] = "sv1t1rc1nc1t2microwebwebacappfalseoibv"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"web", "webacapp", "t2.small",
		"1"}] = "sv1t1rc1nc1t2smallwebwebacappfalseehv6"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"web", "webacapp", "t2.medium",
		"1"}] = "sv1t1rc1nc1t2mediumwebwebacappfalsemeci"
	NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{"web", "webacapp", "t2.large",
		"1"}] = "sv1t1rc1nc1t2largewebwebacappfalseyowvcs"
}
*/
func InitRelevantInstanceTypesForLimits(){

	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"100", Replicas:"1"}] = []string{"t2.nano", "t2.micro",
		"t2.small", "t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"100", Replicas:"2"}] = []string{"t2.nano", "t2.micro",
		"t2.small", "t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"100", Replicas:"3"}] = []string{"t2.nano", "t2.micro",
		"t2.small", "t2.medium", "t2.large", "t2.xlarge" }

	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"200", Replicas:"1"}] = []string{"t2.nano", "t2.micro",
		"t2.small", "t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"200", Replicas:"2"}] = []string{"t2.nano", "t2.micro",
		"t2.small", "t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"200", Replicas:"3"}] = []string{"t2.micro",
		"t2.small", "t2.medium", "t2.large", "t2.xlarge" }

	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"500", Replicas:"1"}] = []string{"t2.micro",
		"t2.small", "t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"500", Replicas:"2"}] = []string{"t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"500", Replicas:"3"}] = []string{"t2.medium", "t2.large", "t2.xlarge" }

	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.nano", Replicas:"1"}] = []string{"t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.nano", Replicas:"2"}] = []string{"t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.nano", Replicas:"3"}] = []string{"t2.xlarge" }

	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.micro", Replicas:"1"}] = []string{"t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.micro", Replicas:"2"}] = []string{"t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.micro", Replicas:"3"}] = []string{"t2.xlarge" }

	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.small", Replicas:"1"}] = []string{"t2.medium", "t2.large", "t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.small", Replicas:"2"}] = []string{"t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.small", Replicas:"3"}] = []string{"t2.xlarge" }

	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.medium", Replicas:"1"}] = []string{"t2.xlarge" }
	RelevantInstancesNames[LimitingInstanceInfo{InstanceType:"t2.large", Replicas:"1"}] = []string{"t2.xlarge" }
}

func getInstancesData(hostInstance,comparingInstance, numReplicas, appType, appName, mainServiceName string, limits bool) [] NodeValues{

	var allRelevantNodes = [] NodeValues{}
	k8sDb:=GetTestFolderName(hostInstance,comparingInstance, numReplicas, appType, appName, limits)+"_k8s"
	if(len(k8sDb)>5) {
		nodeNames := getInfluxDbQueryResult(getNodeNamesQuery(), k8sDb)
		log.Info("Data: ", nodeNames)

		var wg sync.WaitGroup
		wg.Add(len(nodeNames))
		for _, nodeName := range nodeNames {
			go func(nodeName []interface{}, appName,mainServiceName, hostInstance string) {
				nodeValues := getNodeUtilization(k8sDb, fmt.Sprint(nodeName[1]), hostInstance)
				podNames := getInfluxDbQueryResult(getPodNamesQuery("default", fmt.Sprint(nodeName[1])), k8sDb)
				log.Info("podNames: ", podNames)
				if (podNames != nil) {
					for idxPod := 0; idxPod < len(podNames); idxPod++ {
						if strings.Contains(fmt.Sprint(podNames[idxPod][1]), mainServiceName) {

							podValues := getPodUtilization(k8sDb, fmt.Sprint(nodeName[1]), fmt.Sprint(podNames[idxPod][1]),
								"default", hostInstance)
							podValues.PodName = fmt.Sprint(podNames[idxPod][1])
							nodeValues.Pod = append(nodeValues.Pod, podValues)
							k6sDb:=GetTestFolderName(hostInstance,comparingInstance, numReplicas, appType, appName, limits)+"_TestK6"
							lastResultRequestsData:=getInfluxDbQueryResult(getLastResultFromRequestsDataQuery(), k6sDb)
							var maxTimeStamp string
							if(len(lastResultRequestsData) > 0) {
								maxTimeStamp = lastResultRequestsData[0][0].(string)
							}
							numRequestsIntf:=getInfluxDbQueryResult(getNumRequestsQuery(maxTimeStamp),k6sDb )
							numRequests:=convertResults(numRequestsIntf, len(numRequestsIntf))
							nodeValues.Requests = numRequests.Values
						}
					}
				}
				if (len(nodeValues.Pod) > 0) {
					allRelevantNodes = append(allRelevantNodes,  nodeValues)
				}
				wg.Done()
			}(nodeName, appName,mainServiceName, hostInstance)
		}
		wg.Wait()
	}
	return allRelevantNodes
}
func get_relevant_Instances_Data_Based_On_Limits(comparingInstance, numReplicas, appType, appName, mainServiceName string) [] NodeValues{
	var allrelevantInstancesNodesData [] NodeValues
	allRelevantInstances :=  RelevantInstancesNames[LimitingInstanceInfo{comparingInstance, numReplicas}]

	var wg sync.WaitGroup
	wg.Add(len(allRelevantInstances))
	for _, relevantInstance := range allRelevantInstances {
		go func(relevantInstance,comparingInstance,numReplicas, appType, appName string) {
			defer wg.Done()
			nodesValues:= getInstancesData(relevantInstance,comparingInstance, numReplicas, appType, appName, mainServiceName, true)
			for _, nodeValue := range nodesValues {
				allrelevantInstancesNodesData = append(allrelevantInstancesNodesData, nodeValue)
			}
		}(relevantInstance,comparingInstance,numReplicas, appType, appName)
	}
	wg.Wait()
	return allrelevantInstancesNodesData
}
func get_relevant_Instance_Data_Based_On_Limits(relevantInstance, comparingInstance, numReplicas, appType, appName, mainServiceName string) [] NodeValues{
	allRelevantInstances :=  RelevantInstancesNames[LimitingInstanceInfo{comparingInstance, numReplicas}]
	var nodesValues [] NodeValues
	if StringInSlice(relevantInstance,allRelevantInstances ){
		nodesValues= getInstancesData(relevantInstance,comparingInstance, numReplicas, appType, appName, mainServiceName, true)
	}
	return nodesValues
}


func get_relevant_Instances_actual_instance_Data_Based_On_Limits(comparingInstance, numReplicas, appType, appName, mainServiceName string) [] NodeValues{
	allrelevantInstancesNodesData:=get_relevant_Instances_Data_Based_On_Limits(comparingInstance, numReplicas, appType, appName, mainServiceName)
	actualVMsData:=getInstancesData(comparingInstance,"", numReplicas, appType, appName, mainServiceName, false)
	for  _, vmData := range actualVMsData {
		allrelevantInstancesNodesData = append(allrelevantInstancesNodesData, vmData)
	}
	return allrelevantInstancesNodesData
}
func get_relevant_Instance_actual_instance_Data_Based_On_Limits(relevantInstance, comparingInstance, numReplicas, appType, appName, mainServiceName string) [] NodeValues{

	allRelevantInstances :=  RelevantInstancesNames[LimitingInstanceInfo{comparingInstance, numReplicas}]
	var relevantInstanceNodesData []NodeValues
	if StringInSlice(relevantInstance,allRelevantInstances ){
		relevantInstanceNodesData=get_relevant_Instance_Data_Based_On_Limits(relevantInstance, comparingInstance, numReplicas, appType, appName, mainServiceName)
		actualVMsData:=getInstancesData(comparingInstance,"", numReplicas, appType, appName, mainServiceName, false)
		for  _, vmData := range actualVMsData {
			relevantInstanceNodesData = append(relevantInstanceNodesData,  vmData)
		}
	}
	return relevantInstanceNodesData
}
func get_instances_perf_data(comparingInstances , replicasSel []string, appType, appName, mainServiceName, hostInstance string) [] NodeValues{

	var allInstancesLimitsData [] NodeValues
	log.Info("Replicas:::", replicasSel)
	var wg sync.WaitGroup
	wg.Add(len(comparingInstances) * len(replicasSel))
	for _, comparingInstance := range comparingInstances {
		for _, numReplicas := range replicasSel{
			go func(comparingInstance,numReplicas, appType, appName,mainServiceName string ) {
				defer wg.Done()
				var actualVmsData []NodeValues
				log.Info("numReplicas:::", numReplicas)
				if hostInstance=="none"{
					actualVmsData = getInstancesData(comparingInstance,"", numReplicas, appType, appName, mainServiceName, false)
					for  _, vmData := range actualVmsData {
						vmData.InstanceType = vmData.InstanceType + "-Replicas-"+numReplicas
						allInstancesLimitsData = append(allInstancesLimitsData,  vmData)
					}
				}else {
					actualVmsData = getInstancesData(hostInstance,comparingInstance, numReplicas, appType, appName, mainServiceName, true)
					for  _, vmData := range actualVmsData {
						vmData.InstanceType = comparingInstance + "-Replicas-"+numReplicas
						allInstancesLimitsData = append(allInstancesLimitsData,  vmData)
					}
				}
				log.Info("allInstancesLimitsData:::", len(allInstancesLimitsData))

			}(comparingInstance, numReplicas, appType, appName,mainServiceName)
		}
	}
	wg.Wait()
	return allInstancesLimitsData
}