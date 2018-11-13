package TERMINUS

import (
	"github.com/influxdata/influxdb/client/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"fmt"
	"encoding/json"
	"sync"
)

func influxDBClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://"+os.Getenv("INFLUXDB_HOST")+":"+os.Getenv("INFLUXDB_PORT"),
		Username: os.Getenv("INFLUXDB_USER"),
		Password: os.Getenv("INFLUXDB_PASS"),
	})
	if err != nil {
		log.Info("Error: ", err)
	}
	return c
}
func getInfluxDbQueryResult(query, databaseName string) [][]interface{}{
	c := influxDBClient()
	log.Info("Query: ", query)
	log.Info("databaseName: ", databaseName)
	q := client.Query{
		Command:  query,
		Database: databaseName,
	}

	resp, err1 := c.Query(q)
	if err1 != nil {
		log.Info("Error: ", err1)
	}
	if resp.Error() != nil {
		log.Info("Error: ", resp.Error())
	}

	if(len(resp.Results) > 0){

		if(len(resp.Results[0].Series)>0){
			res := resp.Results[0].Series[0].Values
			return res
		}
	}
	return nil
}
/*
	Node Queries
 */

func getNodeNamesQuery() string{
	return fmt.Sprintf("SHOW TAG VALUES FROM uptime WITH KEY=nodename")
}
func getNodeCPUUtilizationQuery(nodeName string) string{
	return fmt.Sprintf("SELECT value FROM \"cpu/node_utilization\" where nodename = '%s' AND type='node'",  nodeName)
}
func getNodeMemUtilizationQuery(nodeName string) string{
	return fmt.Sprintf("SELECT value FROM \"memory/node_utilization\" where nodename = '%s' AND type='node'",  nodeName)
}
func getNodeCoresQuery(nodeName string) string{
	return fmt.Sprintf("SELECT value FROM \"cpu/node_capacity\" where nodename = '%s' AND type = 'node'",  nodeName)
}
func getNodeMemQuery(nodeName string) string{
	return fmt.Sprintf("SELECT value FROM \"memory/node_capacity\" where nodename = '%s' AND type = 'node'",  nodeName)
}
/*
	Pod Queries
 */
func getPodNamesQuery(nsName,  nodeName string) string{
	return fmt.Sprintf("SHOW TAG VALUES FROM uptime WITH KEY = pod_name WHERE namespace_name = '%s' AND nodename = '%s'", nsName, nodeName)
}
func getPodCPUUtilizationQuery(nodeName, podName, nsName string) string{
	return fmt.Sprintf("SELECT value FROM \"cpu/usage_rate\" where nodename = '%s' AND pod_name ='%s' AND namespace_name = '%s'  AND type='pod'",  nodeName, podName, nsName)
}
func getPodCPULimitQuery(nodeName, podName, nsName string) string{
	return fmt.Sprintf("SELECT value FROM \"cpu/limit\" where nodename = '%s' AND pod_name ='%s' AND namespace_name = '%s'  AND type='pod'",
		nodeName, podName, nsName)
}
func getPodCPURequestsQuery(nodeName, podName, nsName string) string{
	return fmt.Sprintf("SELECT value FROM \"cpu/request\" where nodename = '%s' AND pod_name ='%s' AND namespace_name = '%s'  AND type='pod'",
		nodeName, podName, nsName)
}
func getPodMemUtilizationQuery(nodeName, podName, nsName string) string{
	return fmt.Sprintf("SELECT value FROM \"memory/usage\" where nodename = '%s' AND pod_name ='%s' AND namespace_name = '%s'  AND type='pod'",  nodeName, podName, nsName)
}
func getPodMemLimitQuery(nodeName, podName, nsName string) string{
	return fmt.Sprintf("SELECT value FROM \"memory/limit\" where nodename = '%s' AND pod_name ='%s' AND namespace_name = '%s'  AND type='pod'",
		nodeName, podName, nsName)
}
func getPodMemRequestsQuery(nodeName, podName, nsName string) string{
	return fmt.Sprintf("SELECT value FROM \"memory/request\" where nodename = '%s' AND pod_name ='%s' AND namespace_name = '%s'  AND type='pod'",
		nodeName, podName, nsName)
}
/*
	Request Data Queries
 */
func getLastResultFromRequestsDataQuery() string{
	return fmt.Sprintf("SELECT * FROM \"vus\" GROUP BY * ORDER BY DESC LIMIT 1")
}
func getNumVusQuery(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT mean(\"value\") FROM \"vus\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func getNumRequestsQuery(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT sum(\"value\") FROM \"http_reqs\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func getRequestsDurationPercentile95Query(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT percentile(\"value\", 95) FROM \"http_req_duration\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func getRequestsDurationPercentile90Query(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT percentile(\"value\", 90) FROM \"http_req_duration\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func getRequestsDurationMaxQuery(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT max(\"value\") FROM \"http_req_duration\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func getRequestsDurationMinQuery(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT min(\"value\") FROM \"http_req_duration\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func getRequestsDurationMeanQuery(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT mean(\"value\") FROM \"http_req_duration\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func getRequestsDurationMedianQuery(maxTimeStamp string) string{
	return fmt.Sprintf("SELECT median(\"value\") FROM \"http_req_duration\" where time < '%s' group by time(1m)", maxTimeStamp)
}
func convertResults(res4 [][]interface{}, length int) DBTimeStampValues{

	if(res4==nil){
		return DBTimeStampValues{}
	}
	values := make([]float64, length)
	timestamps := make([]int64, length)

	for idx2:=0;idx2<length;idx2++ {
		if(res4[idx2][1]!=nil) {
			res, err := res4[idx2][1].(json.Number).Float64()
			if err != nil {
				log.Fatalln("Error: ", err)
			}

			values[idx2] = res //provided it's indeed int. you can add a check here
			/*t, e := time.Parse(
				time.RFC3339,
				res4[idx2][0].(string))
			if e != nil {
				log.Fatal(e)
			}*/
			timestamps[idx2] = int64(idx2)
		} else{
			values[idx2] =0 //provided it's indeed int. you can add a check here
			timestamps[idx2] = int64(idx2)
		}
	}
	result:=DBTimeStampValues{Timestamps:timestamps, Values:values}
	return result

}
func getNodeUtilization(databaseName, nodeName,hostInstance string) NodeValues{

	var result  = NodeValues{}
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		nodeCpuUtilIntf:=getInfluxDbQueryResult(getNodeCPUUtilizationQuery(fmt.Sprint(nodeName)),databaseName)
		cpuUtil:=convertResults(nodeCpuUtilIntf, len(nodeCpuUtilIntf))
		result.Timestamps = cpuUtil.Timestamps
		result.CpuUtil = cpuUtil.Values
	}()
	go func() {
		defer wg.Done()
		nodeMemUtilIntf:=getInfluxDbQueryResult(getNodeMemUtilizationQuery(fmt.Sprint(nodeName)),databaseName)
		memUtil:=convertResults(nodeMemUtilIntf,len(nodeMemUtilIntf))
		result.MemUtil = memUtil.Values
	}()

	go func() {
		defer wg.Done()
		nodeCoresIntf:=getInfluxDbQueryResult(getNodeCoresQuery(fmt.Sprint(nodeName)),databaseName)
		cores:=convertResults(nodeCoresIntf,len(nodeCoresIntf))
		result.Cores = cores.Values
	}()

	go func() {
		defer wg.Done()
		nodeMemIntf:=getInfluxDbQueryResult(getNodeMemQuery(fmt.Sprint(nodeName)),databaseName)
		mem:=convertResults(nodeMemIntf,len(nodeMemIntf))
		result.Mem = mem.Values
	}()

	result.InstanceType = hostInstance
	wg.Wait()
	return result
}

func getPodUtilization(databaseName, nodeName, podName, nsName,hostInstance1 string) PodValues{

	var result = PodValues{}
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		podCpuUtilIntf := getInfluxDbQueryResult(getPodCPUUtilizationQuery(fmt.Sprint(nodeName), fmt.Sprint(podName), nsName), databaseName)
		cpuUtil:=convertResults(podCpuUtilIntf, len(podCpuUtilIntf))
		result.Timestamps = cpuUtil.Timestamps
		result.CpuUtil = cpuUtil.Values
	}()
	go func() {
		defer wg.Done()
		podCpuLimitIntf:=getInfluxDbQueryResult(getPodCPULimitQuery(fmt.Sprint(nodeName),fmt.Sprint(podName), nsName),databaseName)
		cpuLimit:=convertResults(podCpuLimitIntf, len(podCpuLimitIntf))
		result.CpuLimit = cpuLimit.Values
	}()
	go func() {
		defer wg.Done()
		podCpuRequestIntf:=getInfluxDbQueryResult(getPodCPURequestsQuery(fmt.Sprint(nodeName),fmt.Sprint(podName), nsName),databaseName)
		cpuRequest:=convertResults(podCpuRequestIntf, len(podCpuRequestIntf))
		result.CpuRequest = cpuRequest.Values
	}()

	go func() {
		defer wg.Done()
		podMemUtilIntf:=getInfluxDbQueryResult(getPodMemUtilizationQuery(fmt.Sprint(nodeName),fmt.Sprint(podName), nsName),databaseName)
		memUtil:=convertResults(podMemUtilIntf,len(podMemUtilIntf))
		result.MemUtil = memUtil.Values
	}()
	go func() {
		defer wg.Done()
		podMemLimitIntf:=getInfluxDbQueryResult(getPodMemLimitQuery(fmt.Sprint(nodeName),fmt.Sprint(podName), nsName),databaseName)
		memLimit:=convertResults(podMemLimitIntf, len(podMemLimitIntf))
		result.MemLimit = memLimit.Values
	}()
	go func() {
		defer wg.Done()
		podMemRequestIntf:=getInfluxDbQueryResult(getPodMemRequestsQuery(fmt.Sprint(nodeName),fmt.Sprint(podName), nsName),databaseName)
		memRequest:=convertResults(podMemRequestIntf, len(podMemRequestIntf))
		result.MemRequest = memRequest.Values
	}()

	wg.Wait()
	return result
}

func getRequestsData(databaseName,hostInstance string) RequestsData {

	var result = RequestsData{}
	lastResultRequestsData:=getInfluxDbQueryResult(getLastResultFromRequestsDataQuery(), databaseName)
	if(len(lastResultRequestsData) > 0){
		maxTimeStamp:=lastResultRequestsData[0][0].(string)
		log.Info(maxTimeStamp)
		var wg sync.WaitGroup
		wg.Add(4)
		go func() {
			defer wg.Done()
			numVUsIntf:=getInfluxDbQueryResult(getNumVusQuery(maxTimeStamp),databaseName )
			numVUs:=convertResults(numVUsIntf, len(numVUsIntf))
			result.Timestamps = numVUs.Timestamps
			result.Vus = numVUs.Values
		}()

		go func() {
			defer wg.Done()
			numRequestsIntf:=getInfluxDbQueryResult(getNumRequestsQuery(maxTimeStamp),databaseName )
			numRequests:=convertResults(numRequestsIntf, len(numRequestsIntf))
			result.Requests = numRequests.Values
		}()

		go func() {
			defer wg.Done()
			requestsDurationMaxIntf:=getInfluxDbQueryResult(getRequestsDurationMaxQuery(maxTimeStamp),databaseName )
			requestsDurationMax:=convertResults(requestsDurationMaxIntf, len(requestsDurationMaxIntf))
			result.ReqDurationMax = requestsDurationMax.Values
		}()
		go func() {
			defer wg.Done()
			requestsDurationMeanIntf:=getInfluxDbQueryResult(getRequestsDurationMeanQuery(maxTimeStamp),databaseName )
			requestsDurationMean:=convertResults(requestsDurationMeanIntf, len(requestsDurationMeanIntf))
			result.ReqDurationMean = requestsDurationMean.Values
		}()

		result.InstanceType = hostInstance

		wg.Wait()
	}
	return result
}



