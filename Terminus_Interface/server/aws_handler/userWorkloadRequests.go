package TERMINUS

import "sync"

func getInstanceRequestsData(hostInstance, comparingInstance, numReplicas, appType, appName string, limits bool)RequestsData{
	
	k6sDb:=GetTestFolderName(hostInstance,comparingInstance, numReplicas, appType, appName, limits)+"_TestK6"
	requestsData:=getRequestsData(k6sDb,hostInstance)
	return requestsData
}

func get_relevant_instances_requests_data_Based_on_limits(comparingInstance, numReplicas, appType, appName,mainServiceName string) []RequestsData{

	var allInstancesRequestsData []RequestsData
	allRelevantInstances :=  RelevantInstancesNames[LimitingInstanceInfo{comparingInstance, numReplicas}]
	var wg sync.WaitGroup
	wg.Add(len(allRelevantInstances))
	for _, relevantInstance := range allRelevantInstances {
		go func(relevantInstance string) {
			defer wg.Done()
			instanceRequestData:= getInstanceRequestsData(relevantInstance,comparingInstance, numReplicas, appType, appName, true)
			allInstancesRequestsData = append(allInstancesRequestsData, instanceRequestData)
		}(relevantInstance)

	}
	wg.Wait()
	return allInstancesRequestsData
}
func get_relevant_instance_requests_data_Based_on_limits(relevantInstance, comparingInstance, numReplicas, appType, appName,mainServiceName string) []RequestsData{

	var allInstancesRequestsData []RequestsData
	instanceRequestData:= getInstanceRequestsData(relevantInstance,comparingInstance, numReplicas, appType, appName, true)
	allInstancesRequestsData = append(allInstancesRequestsData, instanceRequestData)
	return allInstancesRequestsData
}

func get_relevant_instances_actual_instance_requests_data_based_on_limits(comparingInstance, numReplicas, appType, appName , mainServiceName string) []RequestsData{

	var allInstancesRequestsData []RequestsData
	allRelevantInstances :=  RelevantInstancesNames[LimitingInstanceInfo{comparingInstance, numReplicas}]
	var wg sync.WaitGroup
	wg.Add(len(allRelevantInstances))
	for _, relevantInstance := range allRelevantInstances {
		go func(relevantInstance string) {
			defer wg.Done()
			instanceRequestData:= getInstanceRequestsData(relevantInstance,comparingInstance, numReplicas, appType, appName, true)
			allInstancesRequestsData = append(allInstancesRequestsData, instanceRequestData)
		}(relevantInstance)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		vmData := getInstanceRequestsData(comparingInstance, "", numReplicas, appType, appName, false)
		allInstancesRequestsData = append(allInstancesRequestsData, vmData)
	}()
	wg.Wait()
	return allInstancesRequestsData
}
func get_relevant_instance_actual_instance_requests_data_based_on_limits(relevantInstance,comparingInstance, numReplicas, appType, appName, mainServiceName string) []RequestsData{

	var allInstancesRequestsData []RequestsData
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		instanceRequestData := getInstanceRequestsData(relevantInstance, comparingInstance, numReplicas, appType, appName, true)
		allInstancesRequestsData = append(allInstancesRequestsData, instanceRequestData)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		vmData := getInstanceRequestsData(comparingInstance, "", numReplicas, appType, appName, false)
		allInstancesRequestsData = append(allInstancesRequestsData, vmData)
	}()
	wg.Wait()
	return allInstancesRequestsData
}

func get_instances_requests_data(comparingInstances , replicasSel []string, appType, appName,mainServiceName, hostInstance string) []RequestsData {
	var allInstancesRequestsData []RequestsData
	var wg sync.WaitGroup
	wg.Add(len(comparingInstances) * len(replicasSel))
	for _, comparingInstance := range comparingInstances {
		for _, numReplicas := range replicasSel {
			go func(comparingInstance, numReplicas string) {
				defer wg.Done()
				var vmData RequestsData
				if hostInstance=="none"{
					vmData = getInstanceRequestsData(comparingInstance, "", numReplicas, appType, appName, false)
					vmData.InstanceType = vmData.InstanceType + "-Replicas-"+numReplicas
				}else {
					vmData = getInstanceRequestsData(hostInstance, comparingInstance, numReplicas, appType, appName, true)
					vmData.InstanceType = comparingInstance + "-Replicas-"+numReplicas
				}
				allInstancesRequestsData = append(allInstancesRequestsData, vmData)
			}(comparingInstance, numReplicas)
		}
	}
	wg.Wait()
	return allInstancesRequestsData
}
