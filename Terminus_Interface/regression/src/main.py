from http.server import BaseHTTPRequestHandler, HTTPServer
import socketserver
import pickle
import urllib.request
import json
from pprint import pprint
from pandas.io.json import json_normalize
import pandas as pd
from sklearn import preprocessing
from sklearn.preprocessing import PolynomialFeatures
from sklearn import datasets, linear_model
from sklearn.metrics import mean_squared_error, r2_score
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import scale
from sklearn.preprocessing import PolynomialFeatures
import numpy as np
from sklearn.preprocessing import StandardScaler
from sklearn.pipeline import Pipeline
from sklearn.metrics import mean_squared_error
from sklearn.linear_model import Ridge
from math import sqrt
import os
import errno
from pymongo import MongoClient
import urllib.parse as urlparse
from influxdb import InfluxDBClient
from pymongo import MongoClient
import pandas as pd
from pandas.io.json import json_normalize
from sklearn.linear_model import Ridge
from sklearn.preprocessing import PolynomialFeatures
from sklearn.pipeline import make_pipeline
from sklearn.linear_model import TheilSenRegressor
from sklearn.datasets import make_regression

class Terminus(BaseHTTPRequestHandler):
    def getAllNodeNames(self,client):
        queryResult = client.query("SHOW TAG VALUES FROM uptime WITH KEY=nodename;")
        nodeNames_temp = list(queryResult.get_points())
        dfnodeNames = pd.DataFrame(nodeNames_temp)
        allNodeNames = dfnodeNames[:]["value"]
        return allNodeNames
    def getNamespaceNames(self,client,node):
        nsQuery = client.query("SHOW TAG VALUES FROM uptime WITH KEY=namespace_name WHERE nodename = '"+node+"';")
        nsQuery_temp = list(nsQuery.get_points())
        dfnsNames = pd.DataFrame(nsQuery_temp)
        allnsNames = dfnsNames[:]["value"]
        return allnsNames
    def getAllPodNames(self,client,node,ns_name):
        queryResult = client.query("SHOW TAG VALUES FROM uptime WITH KEY = pod_name WHERE namespace_name = '"+ns_name+"' AND nodename = '"+node+"';")
        podNames_temp = list(queryResult.get_points())
        dfpodNames = pd.DataFrame(podNames_temp)
        if dfpodNames.empty:
            return dfpodNames
        else:
            allpodNames = dfpodNames[:]["value"]
            return allpodNames
    def getCPUUtilizationNode(self,client, node):
        queryResult = client.query('SELECT * FROM "cpu/node_utilization" where nodename = \''+node+'\' AND type=\'node\';')
        dfcpuUtilization = pd.DataFrame(queryResult['cpu/node_utilization'])
        return dfcpuUtilization
    def getCPUUtilizationPod(self,client, node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "cpu/usage_rate" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfcpuUtilization = pd.DataFrame(queryResult['cpu/usage_rate'])
        return dfcpuUtilization
    def getCPUUtilizationPodContainer(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "cpu/usage_rate" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\' AND type=\'pod_container\';')
        dfcpuUtilization = pd.DataFrame(queryResult['cpu/usage_rate'])
        return dfcpuUtilization
    def prepareCpuUtilization(self,client,node,ns_name, pod_name):
        cpuUtilization = self.getCPUUtilizationNode(client,node)
        podCpuUtilization = self.getCPUUtilizationPod(client,node,ns_name, pod_name)
        containercpuUtilization = self.getCPUUtilizationPodContainer(client,node,ns_name, pod_name)
        plt.plot(cpuUtilization.index, cpuUtilization['value'] *1000, 'r', label="node") # plotting t, a separately
        plt.plot(podCpuUtilization.index, podCpuUtilization['value'], 'b', label="pod") # plotting t, b separately
        plt.plot(containercpuUtilization.index, containercpuUtilization['value'], 'g', label="container") # plotting t, c separately
        plt.legend(loc='upper left')
        plt.show()
    def getMemoryUtilizationNode(self,client,node):
        queryResult = client.query('SELECT * FROM "memory/node_utilization" where nodename = \''+node+'\' AND type=\'node\';')
        dfmemUtilization = pd.DataFrame(queryResult['memory/node_utilization'])
        return dfmemUtilization
    def getMemoryUtilizationPod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "memory/usage" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['memory/usage'])
        return dfmemUtilization
    def getMemoryUtilizationPodContainer(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "memory/usage" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\' AND type=\'pod_container\';')
        dfmemUtilization = pd.DataFrame(queryResult['memory/usage'])
        return dfmemUtilization
    def prepareMemoryUtilization(self,client,node,ns_name, pod_name):
        memoryUtilization = self.getMemoryUtilizationNode(client,node)
        podMemoryUtilization = self.getMemoryUtilizationPod(client,node,ns_name, pod_name)
        containerMemoryUtilization = self.getMemoryUtilizationPodContainer(client,node,ns_name, pod_name)
        plt.plot(memoryUtilization.index, memoryUtilization['value'], 'r', label="node") # plotting t, a separately
        plt.plot(podMemoryUtilization.index, podMemoryUtilization['value'], 'b', label="pod") # plotting t, b separately
        plt.plot(containerMemoryUtilization.index, containerMemoryUtilization['value'], 'g', label="container") # plotting t, c separately
        plt.legend(loc='upper left')
        plt.show()
    def getNetworkTxRatePod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/tx_rate" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/tx_rate'])
        return dfmemUtilization
    def getNetworkTxPod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/tx" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/tx'])
        return dfmemUtilization
    def getNetworkTxErrorsPod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/tx_errors" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/tx_errors'])
        return dfmemUtilization
    def getNetworkTxErrorsRatePod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/tx_errors_rate" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/tx_errors_rate'])
        return dfmemUtilization
    def prepareNetworkTxRateUtilization(self,client,node,ns_name, pod_name):
        podNetworTxRate = self.getNetworkTxRatePod(client,node,ns_name, pod_name)
        podNetworTx = self.getNetworkTxPod(client,node,ns_name, pod_name)
        podNetworkError = self.getNetworkTxErrorsPod(client,node,ns_name, pod_name)
        podNetworkErrorRate = self.getNetworkTxErrorsRatePod(client,node,ns_name, pod_name)
        plt.plot(podNetworTxRate.index, podNetworTxRate['value'], 'b') # plotting t, b separately
        #plt.plot(podNetworTx.index, podNetworTx['value'], 'g') # plotting t, b separately
        #plt.plot(podNetworkError.index, podNetworkError['value'], 'y') # plotting t, b separately
        plt.plot(podNetworkErrorRate.index, podNetworkErrorRate['value'], 'r') # plotting t, b separately
        plt.show()
    def getNetworkRxRatePod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/rx_rate" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/rx_rate'])
        return dfmemUtilization
    def getNetworkRxPod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/rx" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/rx'])
        return dfmemUtilization

    def getNetworkRxErrorsPod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/rx_errors" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/rx_errors'])
        return dfmemUtilization
    def getNetworkRxErrorsRatePod(self,client,node,ns_name, pod_name):
        queryResult = client.query('SELECT * FROM "network/rx_errors_rate" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\';')
        dfmemUtilization = pd.DataFrame(queryResult['network/rx_errors_rate'])
        return dfmemUtilization
    def prepareNetworkRxRateUtilization(self,client,node,ns_name, pod_name):
        podNetworRxRate = self.getNetworkRxRatePod(client,node,ns_name, pod_name)
        podNetworRx = self.getNetworkRxPod(client,node,ns_name, pod_name)
        podNetworkError = self.getNetworkRxErrorsPod(client,node,ns_name, pod_name)
        podNetworkErrorRate = self.getNetworkRxErrorsRatePod(client,node,ns_name, pod_name)
        plt.plot(podNetworRxRate.index, podNetworRxRate['value'], 'b') # plotting t, b separately
        #plt.plot(podNetworRx.index, podNetworRx['value'], 'g') # plotting t, b separately
        #plt.plot(podNetworkError.index, podNetworkError['value'], 'y') # plotting t, b separately
        plt.plot(podNetworkErrorRate.index, podNetworkErrorRate['value'], 'r') # plotting t, b separately
        plt.show()
    def getRelevantNodeName(self,client,ns_name):
        allNodeNames  = self.getAllNodeNames(client)
        #nsNames = getNamespaceNames(allNodeNames[0])
        relevantNodes = []
        for node in allNodeNames:
            allPodNamesNode = self.getAllPodNames(client,node,'default')
            if(not allPodNamesNode.empty):
                relevantNodes.append(node)
        return relevantNodes
    def getNodeResourceUtilizationDf(self,client, nodeName):
        Result_node_CPU = client.query("SELECT value from \"cpu/node_utilization\" where nodename = '"+nodeName+"' AND type = 'node' ")
        Result_node_MEM = client.query("SELECT value from \"memory/node_utilization\" where nodename = '"+nodeName+"' AND type = 'node' ")

        Result_node_CPU_Cores = client.query("SELECT mean(\"value\") FROM \"cpu/node_capacity\" where nodename = '"+nodeName+
                                             "' AND type = 'node' GROUP BY time(1m)")
        Result_node_mem_node = client.query("SELECT mean(\"value\")FROM \"memory/node_capacity\" where nodename = '"+
                                            nodeName+"' AND type = 'node' GROUP BY time(1m)")

        cpu_points = pd.DataFrame(Result_node_CPU.get_points())
        cpu_points['time'] = pd.to_datetime(cpu_points['time'])
        cpu_points = cpu_points.set_index('time')
        cpu_points.columns = ['node_cpu_util']
        mem_points = pd.DataFrame(Result_node_MEM.get_points())
        mem_points['time'] = pd.to_datetime(mem_points['time'])
        mem_points = mem_points.set_index('time')
        mem_points.columns = ['node_mem_util']

        cores_points = pd.DataFrame(Result_node_CPU_Cores.get_points())
        cores_points['time'] = pd.to_datetime(cores_points['time'])
        cores_points = cores_points.set_index('time')
        cores_points.columns = ['node_cores']

        mem_node_points = pd.DataFrame(Result_node_mem_node.get_points())
        mem_node_points['time'] = pd.to_datetime(mem_node_points['time'])
        mem_node_points = mem_node_points.set_index('time')
        mem_node_points.columns = ['node_mem']

        df_node =pd.concat([cpu_points, mem_points,cores_points,mem_node_points], axis=1)
        return df_node
    def getPodResourceUtilizationDf(self,client, node, ns_name, pod_name):
        Result_Pod_CPU_usage = client.query('SELECT value FROM "cpu/usage_rate" where nodename = \''+node+
                                            '\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+
                                            '\'  AND type=\'pod\';')
        Result_Pod_MEM_usage = client.query('SELECT value from \"memory/usage\" where nodename = \''+node+
                                            '\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+
                                            '\'  AND type=\'pod\';')

        Result_Pod_CPU_limit = client.query('SELECT mean(\"value\") FROM "cpu/limit" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\' group by time(1m);')
        Result_Pod_MEM_limit = client.query('SELECT mean(\"value\") from \"memory/limit\" where nodename = \''+node+
                                            '\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+
                                            '\'  AND type=\'pod\' group by time(1m);')

        Result_Pod_CPU_requests = client.query('SELECT mean(\"value\") FROM "cpu/request" where nodename = \''+node+'\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+'\'  AND type=\'pod\' group by time(1m);')
        Result_Pod_MEM_requests = client.query('SELECT mean(\"value\") from \"memory/request\" where nodename = \''+node+
                                               '\' AND pod_name = \''+pod_name+'\' AND namespace_name = \''+ns_name+
                                               '\'  AND type=\'pod\' group by time(1m);')


        cpu_points_usage = pd.DataFrame(Result_Pod_CPU_usage.get_points())
        cpu_points_usage['time'] = pd.to_datetime(cpu_points_usage['time'])
        cpu_points_usage = cpu_points_usage.set_index('time')
        cpu_points_usage.columns = ['pod_cpu_usage']


        mem_points_usage = pd.DataFrame(Result_Pod_MEM_usage.get_points())
        mem_points_usage['time'] = pd.to_datetime(mem_points_usage['time'])
        mem_points_usage = mem_points_usage.set_index('time')
        mem_points_usage.columns = ['pod_mem_usage']


        cpu_points_limits = pd.DataFrame(Result_Pod_CPU_limit.get_points())
        cpu_points_limits['time'] = pd.to_datetime(cpu_points_limits['time'])
        cpu_points_limits = cpu_points_limits.set_index('time')
        cpu_points_limits.columns = ['pod_cpu_limit']


        mem_points_limits = pd.DataFrame(Result_Pod_MEM_limit.get_points())
        mem_points_limits['time'] = pd.to_datetime(mem_points_limits['time'])
        mem_points_limits = mem_points_limits.set_index('time')
        mem_points_limits.columns = ['pod_mem_limit']


        cpu_points_request = pd.DataFrame(Result_Pod_CPU_requests.get_points())
        cpu_points_request['time'] = pd.to_datetime(cpu_points_request['time'])
        cpu_points_request = cpu_points_request.set_index('time')
        cpu_points_request.columns = ['pod_cpu_request']


        mem_points_request = pd.DataFrame(Result_Pod_MEM_requests.get_points())
        mem_points_request['time'] = pd.to_datetime(mem_points_request['time'])
        mem_points_request = mem_points_request.set_index('time')
        mem_points_request.columns = ['pod_mem_request']

        df_pod =pd.concat([cpu_points_usage, mem_points_usage,cpu_points_limits,mem_points_limits,cpu_points_request,mem_points_request ], axis=1)

        return df_pod
    def getRequestsDf(self,clientK6):
        queryResult = clientK6.query('SELECT sum("value") FROM "vus" group by time(1m);')
        vus = pd.DataFrame(queryResult['vus'])
        vus.columns = ['vus','time']
        vus = vus.set_index('time')


        queryResultReqs = clientK6.query('SELECT sum("value") FROM "http_reqs" group by time(1m);')
        reqs = pd.DataFrame(queryResultReqs['http_reqs'])
        reqs.columns = ['requests','time']
        reqs = reqs.set_index('time')
        queryResultReqsDuration95 = clientK6.query('SELECT percentile("value", 95) FROM "http_req_duration" group by time(1m) ;')
        reqs_duration95 = pd.DataFrame(queryResultReqsDuration95['http_req_duration'])
        reqs_duration95.columns = [ 'requests_duration_percentile_95','time']
        reqs_duration95 = reqs_duration95.set_index('time')
        queryResultReqsDuration90 = clientK6.query('SELECT percentile("value", 90) FROM "http_req_duration" group by time(1m) ;')
        reqs_duration90 = pd.DataFrame(queryResultReqsDuration90['http_req_duration'])
        reqs_duration90.columns = ['requests_duration_percentile_90','time']
        reqs_duration90 = reqs_duration90.set_index('time')

        queryResultMaxDuration = clientK6.query('SELECT max("value") FROM "http_req_duration" group by time(1m);')
        reqs_duration_max = pd.DataFrame(queryResultMaxDuration['http_req_duration'])
        reqs_duration_max.columns = ['requests_duration_max','time']
        reqs_duration_max = reqs_duration_max.set_index('time')

        queryResultMinDuration = clientK6.query('SELECT min("value") FROM "http_req_duration" group by time(1m);')
        reqs_duration_min = pd.DataFrame(queryResultMinDuration['http_req_duration'])
        reqs_duration_min.columns = ['requests_duration_min','time']
        reqs_duration_min = reqs_duration_min.set_index('time')

        queryResultMeanDuration = clientK6.query('SELECT mean("value") FROM "http_req_duration" group by time(1m);')
        reqs_duration_mean = pd.DataFrame(queryResultMeanDuration['http_req_duration'])
        reqs_duration_mean.columns = ['requests_duration_mean','time']
        reqs_duration_mean = reqs_duration_mean.set_index('time')

        queryResultMedianDuration = clientK6.query('SELECT median("value") FROM "http_req_duration" group by time(1m);')
        reqs_duration_median = pd.DataFrame(queryResultMedianDuration['http_req_duration'])
        reqs_duration_median.columns = ['requests_duration_median','time']
        reqs_duration_median = reqs_duration_median.set_index('time')

        finalDF = pd.merge(vus, reqs, left_index=True, right_index=True)
        finalDF = pd.merge(finalDF, reqs_duration95, left_index=True, right_index=True)
        finalDF = pd.merge(finalDF, reqs_duration90, left_index=True, right_index=True)
        finalDF = pd.merge(finalDF,reqs_duration_max, left_index=True, right_index=True)
        finalDF = pd.merge(finalDF,reqs_duration_min, left_index=True, right_index=True)
        finalDF = pd.merge(finalDF,reqs_duration_mean, left_index=True, right_index=True)
        finalDF = pd.merge(finalDF,reqs_duration_median, left_index=True, right_index=True)
        finalDF.index = pd.to_datetime(finalDF.index)

        return finalDF

    def getPodsNodesRequestsDf(self,appNames, client,  clientK6):
        default_ns_name =  "default"
        df_pods_node = []
        relevantNodeNames = self.getRelevantNodeName(client,default_ns_name)
        for relevantNodeName in relevantNodeNames:
            if relevantNodeName is not None:
                podNames = self.getAllPodNames(client,relevantNodeName, default_ns_name)
                df_node = self.getNodeResourceUtilizationDf(client,relevantNodeName)

                for podName in podNames:
                    if appNames[0] in podName:
                        df_pod = self.getPodResourceUtilizationDf(client,relevantNodeName, default_ns_name, podName)
                        finalDF = pd.merge(df_node,df_pod, left_index=True, right_index=True)
                        requestsDF = self.getRequestsDf(clientK6)
                        finalDF = pd.merge(finalDF,requestsDF, left_index=True, right_index=True)
                        if(finalDF['pod_cpu_limit'].values[0]==0):
                            finalDF['pod_cpu_usage'] = finalDF['pod_cpu_usage']/(finalDF['node_cores'])
                            finalDF['pod_cpu_limit'] = finalDF['node_cores']/1000
                            finalDF['pod_cpu_request'] = finalDF['node_cores']/1000
                        else:
                            finalDF['pod_cpu_usage'] = finalDF['pod_cpu_usage']/(finalDF['pod_cpu_limit'])
                            finalDF['pod_cpu_limit'] = finalDF['pod_cpu_limit']/1000
                            finalDF['pod_cpu_request'] = finalDF['pod_cpu_request']/1000

                        if(finalDF['pod_mem_limit'].values[0]==0):
                            finalDF['pod_mem_usage'] = finalDF['pod_mem_usage']/(finalDF['node_mem'])
                            finalDF['pod_mem_limit'] = finalDF['node_mem']/(1073741824)
                            finalDF['pod_mem_request'] = finalDF['node_mem']/(1073741824)
                        else:
                            finalDF['pod_mem_usage'] = finalDF['pod_mem_usage']/(finalDF['pod_mem_limit'])
                            finalDF['pod_mem_limit'] = finalDF['pod_mem_limit']/(1073741824)
                            finalDF['pod_mem_request'] = finalDF['pod_mem_request']/(1073741824)
                        finalDF['node_cores'] = finalDF['node_cores']/1000
                        finalDF['node_mem'] = finalDF['node_mem']/(1073741824)

                        finalDF = finalDF.fillna(0)
                        finalDF = finalDF[(finalDF.T != 0).any()]
                        df_pods_node.append(finalDF)
        return df_pods_node

    def getAndCombineAllDbs(self, host, port, username, password,appNames, folderNames):
        allFinalDFs = []
        print("FolderNames len = ", len(folderNames))
        for folderName in folderNames:
            client = InfluxDBClient(host, port,username , password, folderName+'_k8s')
            clientK6 = InfluxDBClient(host, port, username, password, folderName+'_TestK6')
            df_pods_node = self.getPodsNodesRequestsDf(appNames, client, clientK6)
            print(folderName)
            if(len(df_pods_node)>0):
                finalDF = pd.DataFrame()
                finalDF['pod_util_cpu_sum'] = 0
                finalDF['pod_util_mem_sum'] = 0
                first = 1
                for i in range(len(df_pods_node)):
                    df_pods_node[i] = df_pods_node[i].reset_index(drop=True)
                    if(first==1):
                        finalDF['pod_util_cpu_sum'] = df_pods_node[i]['pod_cpu_usage']
                        finalDF['pod_util_mem_sum'] = df_pods_node[i]['pod_mem_usage']
                        first=0
                    else:
                        finalDF['pod_util_cpu_sum'] = finalDF['pod_util_cpu_sum'] +  df_pods_node[i]['pod_cpu_usage']
                        finalDF['pod_util_mem_sum'] = finalDF['pod_util_mem_sum'] +  df_pods_node[i]['pod_mem_usage']

                finalDF['num_pods'] = int(len(df_pods_node))
                finalDF['pod_util_cpu_avg'] = finalDF['pod_util_cpu_sum']/finalDF['num_pods']
                finalDF['pod_util_mem_avg'] = finalDF['pod_util_mem_sum']/finalDF['num_pods']

                finalDF = pd.concat([finalDF, df_pods_node[0][['node_cores', 'node_mem','node_cpu_util','node_mem_util', 'pod_cpu_limit', 'pod_cpu_request','pod_mem_limit',
                                                               'pod_mem_request','vus','requests','requests_duration_percentile_95',
                                                               'requests_duration_percentile_90','requests_duration_max', 'requests_duration_min',
                                                               'requests_duration_mean', 'requests_duration_median'
                                                               ]]], axis=1)

                allFinalDFs.append(finalDF)
        df = pd.DataFrame()
        print("All Dfs len = ", len(allFinalDFs))
        for idx in range(len(allFinalDFs)):
            df = df.append(allFinalDFs[idx])

        final_df  = df[['requests','requests_duration_mean','num_pods','pod_cpu_limit','node_cores', 'node_mem','pod_mem_limit','pod_util_cpu_avg','pod_util_mem_avg',
                        ]]
        final_df['pod_util_cpu_avg'] = final_df['pod_util_cpu_avg']*final_df['pod_cpu_limit']
        final_df['pod_util_mem_avg'] = final_df['pod_util_mem_avg']*final_df['pod_mem_limit']
        final_df = final_df.sort_values(['requests'])
        final_df = final_df[(final_df[['pod_util_cpu_avg','pod_util_mem_avg','requests_duration_mean']] != 0).all(axis=1)]
        final_df = final_df[np.isfinite(final_df['requests'])]
        final_df = final_df[np.isfinite(final_df['requests_duration_mean'])]
        final_df = final_df[np.isfinite(final_df['pod_util_cpu_avg'])]
        final_df = final_df[np.isfinite(final_df['pod_util_mem_avg'])]
        final_df = final_df[final_df.requests_duration_mean < 2500]
        final_df = final_df.reset_index(drop=True)
        return final_df

    def train_and_return_model(self,host, port, username, password,appType, appNames, folderNames ):
        df = self.getAndCombineAllDbs(host, port, username, password,appNames, folderNames)
        df['total_cpu_util'] = df['pod_util_cpu_avg']*df['num_pods']
        df['total_mem_util'] = df['pod_util_mem_avg']*df['num_pods']
        df_X = df[['total_cpu_util']].values
        df_Y = df[['requests']].values
        X_train, X_test, y_train, y_test = train_test_split(df_X, df_Y, test_size=0.33, random_state=42)
        X, y = make_regression(n_samples=df_X.shape[0], n_features=1, noise=4.0, random_state=0)
        regr = TheilSenRegressor(random_state=0).fit(X_train, y_train)
        regr.score(X, y)
        y_pred = regr.predict(X_test)

        rms = sqrt(mean_squared_error(y_test, y_pred))
        print('RMs score: %.2f' % rms)
        return regr, rms
    def train_and_return_model_replicas(self,host, port, username, password,appType, appNames, folderNames ):
        df = self.getAndCombineAllDbs(host, port, username, password,appNames, folderNames)
        df['total_cpu_util'] = df['pod_util_cpu_avg']*df['num_pods']
        df['total_mem_util'] = df['pod_util_mem_avg']*df['num_pods']
        df_X = df[['requests']].values
        df_Y = df[['total_cpu_util']].values
        X_train, X_test, y_train, y_test = train_test_split(df_X, df_Y, test_size=0.33, random_state=42)
        X, y = make_regression(n_samples=df_X.shape[0], n_features=1, noise=4.0, random_state=0)
        regr = TheilSenRegressor(random_state=0).fit(X_train, y_train)
        regr.score(X, y)
        y_pred = regr.predict(X_test)

        rms = sqrt(mean_squared_error(y_test, y_pred))
        print('RMs score: %.2f' % rms)
        return regr, rms

    def train_and_return_model_smart(self,host, port, username, password,appType, appNames, folderNames ):
        df = self.getAndCombineAllDbs(host, port, username, password,appNames, folderNames)
        df['total_cpu_util'] = df['pod_util_cpu_avg']*df['num_pods']
        df['total_mem_util'] = df['pod_util_mem_avg']*df['num_pods']

        df = df.head(46)
        df_X = df[['total_cpu_util']].values
        df_Y = df[['requests']].values
        if(df.shape[0] < 15):
            testSize = 0
            X_train, X_test, y_train, y_test = train_test_split(df_X, df_Y, test_size=testSize, random_state=42)
            X_test = X_train
            y_test = y_train

        else:
            testSize = 0.33
            X_train, X_test, y_train, y_test = train_test_split(df_X, df_Y, test_size=testSize, random_state=42)
        lin_reg = linear_model.LinearRegression()
        regr = Pipeline([
            ("lin_reg", lin_reg),
        ])
        regr.fit(X_train, y_train)

        # Make predictions using the testing set
        y_pred = regr.predict(X_test)

        # The coefficients

        # The mean squared error
        print("Mean squared error: %.2f"
              % mean_squared_error(y_test, y_pred))
        # Explained variance score: 1 is perfect prediction
        print('Variance score: %.2f' % r2_score(y_test, y_pred))

        #print ('Test score %.2f', regr.score(X_test, y_test) )
        print("Train Mean squared error: %.2f"
              % mean_squared_error(y_train, regr.predict(X_train)))
        rms = sqrt(mean_squared_error(y_test, y_pred))
        print('RMs score: %.2f' % rms)
        return regr, rms
    def do_training(self, host, port, username, password, instanceFamily, appNames, appType, folderNames, filename):
        model, rms = self.train_and_return_model(host, port, username, password,appType, appNames,folderNames)
        print(filename)
        if not os.path.exists(os.path.dirname(filename)):
            try:
                os.makedirs(os.path.dirname(filename))
            except OSError as exc: # Guard against race condition
                if exc.errno != errno.EEXIST:
                    raise
        pickle.dump(model, open(filename, 'wb'))
        return rms
    def do_training_replicas(self, host, port, username, password, instanceFamily, appNames, appType, folderNames, filename):
        model, rms = self.train_and_return_model_replicas(host, port, username, password,appType, appNames,folderNames)
        print(filename)
        if not os.path.exists(os.path.dirname(filename)):
            try:
                os.makedirs(os.path.dirname(filename))
            except OSError as exc: # Guard against race condition
                if exc.errno != errno.EEXIST:
                    raise
        pickle.dump(model, open(filename, 'wb'))
        return rms

    def do_training_smart(self, host, port, username, password, instanceFamily, appNames, appType, folderNames, filename):
        model, rms = self.train_and_return_model_smart(host, port, username, password,appType, appNames,folderNames)
        print(filename)
        if not os.path.exists(os.path.dirname(filename)):
            try:
                os.makedirs(os.path.dirname(filename))
            except OSError as exc: # Guard against race condition
                if exc.errno != errno.EEXIST:
                    raise
        pickle.dump(model, open(filename, 'wb'))
        return rms
    def get_folder_names(self, appName, colName, dbName):
        mongoclient = MongoClient(mongo_host, 27017, username=mongo_username, password=mongo_password)
        db = mongoclient[dbName]
        col = db[colName]
        datapoints = list(col.find({"servicename": appName}))
        dfMongo = json_normalize(datapoints)
        return dfMongo.foldername


    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def do_GET(self):
        self._set_headers()
        output=''
        #http://localhost:9002/pretrain?appname=primeapp&apptype=compute&instancefamily=t2&colname=ALL_BRUTE_FORCE_CONDUCTED_TEST_NAMES&dbname=TERMINUS
        if '/pretrain' in self.path:
            parsed = urlparse.urlparse(self.path)

            appName = urlparse.parse_qs(parsed.query)['appname'][0]
            appType = urlparse.parse_qs(parsed.query)['apptype'][0]
            instanceFamily = urlparse.parse_qs(parsed.query)['instancefamily'][0]
            colName = urlparse.parse_qs(parsed.query)['colname'][0]
            dbName = urlparse.parse_qs(parsed.query)['dbname'][0]
            mainServiceName = urlparse.parse_qs(parsed.query)['mainServiceName'][0]

            folderNames = self.get_folder_names(appName, colName, dbName)
            filename = "/app/training/preTrained/"+appType+"/"+appName+"/"+mainServiceName+"/"+instanceFamily+"/"+"trained.sav"
            appNames = [mainServiceName]
            rms = self.do_training(host, port, username, password, instanceFamily, appNames, appType, folderNames, filename)
            output=""+str(rms)
            self.wfile.write(output.encode())
        #http://localhost:9002/trainedreplicas?appname=primeapp&apptype=compute&instancefamily=t2&colname=ALL_BRUTE_FORCE_CONDUCTED_TEST_NAMES&dbname=TERMINUS
        elif '/trainedreplicas' in self.path:
            parsed = urlparse.urlparse(self.path)
            appName = urlparse.parse_qs(parsed.query)['appname'][0]
            appType = urlparse.parse_qs(parsed.query)['apptype'][0]
            instanceFamily = urlparse.parse_qs(parsed.query)['instancefamily'][0]
            colName = urlparse.parse_qs(parsed.query)['colname'][0]
            dbName = urlparse.parse_qs(parsed.query)['dbname'][0]
            mainServiceName = urlparse.parse_qs(parsed.query)['mainServiceName'][0]


            folderNames = self.get_folder_names(appName, colName, dbName)
            filename = "/app/training/preTrainedReplicas/"+appType+"/"+appName+"/"+mainServiceName+"/"+instanceFamily+"/"+"trained.sav"
            appNames = [mainServiceName]
            rms = self.do_training_replicas(host, port, username, password, instanceFamily, appNames, appType, folderNames, filename)
            output=""+str(rms)
            self.wfile.write(output.encode())

        # http://localhost:9002/getPredictionPreTrained?appname=primeapp&apptype=compute&replicas=2&numcoresutil=0.1
        # &numcoreslimit=0.1&nummemlimit=0.1&instancefamily=t2&requestduration=1000
        elif '/getPredictionPreTrained' in self.path:
            parsed = urlparse.urlparse(self.path)
            appName = urlparse.parse_qs(parsed.query)['appname'][0]
            appType = urlparse.parse_qs(parsed.query)['apptype'][0]
            replicas = urlparse.parse_qs(parsed.query)['replicas'][0]
            numcoresUtil = urlparse.parse_qs(parsed.query)['numcoresutil'][0]
            numcoresLimit = urlparse.parse_qs(parsed.query)['numcoreslimit'][0]
            nummemLimit = urlparse.parse_qs(parsed.query)['nummemlimit'][0]
            requestDuration = urlparse.parse_qs(parsed.query)['requestduration'][0]
            instanceFamily = urlparse.parse_qs(parsed.query)['instancefamily'][0]
            mainServiceName = urlparse.parse_qs(parsed.query)['mainServiceName'][0]

            filename = "/app/training/preTrained/"+appType+"/"+appName+"/"+mainServiceName+"/"+instanceFamily+"/"+"trained.sav"
            print(filename)
            loaded_model = pickle.load(open(filename, 'rb'))
            val = [[float(numcoresUtil)*float(replicas)]]
            predict = loaded_model.predict(val)
            output=""+str(predict[0])
            print (predict)
            print (output)
            self.wfile.write(output.encode())
        # http://localhost:9002/getPredictionPreTrained?appname=primeapp&apptype=compute&msc=1000&numcoresutil=0.1
        # &numcoreslimit=0.1&nummemlimit=0.1&instancefamily=t2&requestduration=1000
        elif '/getPredictionReplicas' in self.path:
            parsed = urlparse.urlparse(self.path)
            appName = urlparse.parse_qs(parsed.query)['appname'][0]
            appType = urlparse.parse_qs(parsed.query)['apptype'][0]
            msc = urlparse.parse_qs(parsed.query)['msc'][0]
            numcoresUtil = urlparse.parse_qs(parsed.query)['numcoresutil'][0]
            numcoresLimit = urlparse.parse_qs(parsed.query)['numcoreslimit'][0]
            nummemLimit = urlparse.parse_qs(parsed.query)['nummemlimit'][0]
            requestDuration = urlparse.parse_qs(parsed.query)['requestduration'][0]
            instanceFamily = urlparse.parse_qs(parsed.query)['instancefamily'][0]
            mainServiceName = urlparse.parse_qs(parsed.query)['mainServiceName'][0]

            filename = "/app/training/preTrainedReplicas/"+appType+"/"+appName+"/"+mainServiceName+"/"+instanceFamily+"/"+"trained.sav"
            print(filename)
            loaded_model = pickle.load(open(filename, 'rb'))
            val = [[float(msc)]]
            predict = loaded_model.predict(val)
            output=""+str(predict[0]/float(numcoresLimit))
            print (predict)
            print (output)
            self.wfile.write(output.encode())
        #http://localhost:9002/smartTestTrain?appname=primeapp&apptype=compute&instancefamily=t2&containerName=s1t1rc1nc1t2xlargecomputeprimeappt2nanob8j
        elif '/smartTestTrain' in self.path:
            parsed = urlparse.urlparse(self.path)
            appName = urlparse.parse_qs(parsed.query)['appname'][0]
            appType = urlparse.parse_qs(parsed.query)['apptype'][0]
            folderName = urlparse.parse_qs(parsed.query)['containerName'][0]
            instanceFamily = urlparse.parse_qs(parsed.query)['instancefamily'][0]
            mainServiceName = urlparse.parse_qs(parsed.query)['mainServiceName'][0]
            appNames = [mainServiceName]
            folderNames = [folderName]
            filename = "/app/training/smartTest/"+appType+"/"+appName+"/"+mainServiceName+"/"+instanceFamily+"/"+folderName+".sav"
            rms = self.do_training_smart(host, port, username, password, instanceFamily, appNames, appType, folderNames, filename)
            output=""+str(rms)
            self.wfile.write(output.encode())
        #http://localhost:9002/smartTestGetResult?appname=primeapp&apptype=compute&numcoresutil=0.1
        # &nummemutil=0.1&instancefamily=t2&requestduration=1000&containerName=s1t1rc1nc1t2xlargecomputeprimeappt2nanob8j

        elif '/smartTestGetResult' in self.path:
            parsed = urlparse.urlparse(self.path)
            folderName = urlparse.parse_qs(parsed.query)['containerName'][0]
            appName = urlparse.parse_qs(parsed.query)['appname'][0]
            appType = urlparse.parse_qs(parsed.query)['apptype'][0]
            mainServiceName = urlparse.parse_qs(parsed.query)['mainServiceName'][0]
            numcoresUtil = urlparse.parse_qs(parsed.query)['numcoresutil'][0]

            requestDuration = urlparse.parse_qs(parsed.query)['requestduration'][0]
            instanceFamily = urlparse.parse_qs(parsed.query)['instancefamily'][0]
            nummemUtil = urlparse.parse_qs(parsed.query)['nummemutil'][0]

            filename = "/app/training/smartTest/"+appType+"/"+appName+"/"+mainServiceName+"/"+instanceFamily+"/"+folderName+".sav"
            loaded_model = pickle.load(open(filename, 'rb'))
            val = [[]]
            if(appType=="compute"):
                val = [[float(numcoresUtil)]]
            else:
                val = [[float(numcoresUtil)]]

            predict = loaded_model.predict(val)
            output=""+str(predict[0][0])
            print (predict)
            print (output)
            self.wfile.write(output.encode())
            # /getActualTRN?appname=primeapp&containerName=s1t1rc1nc1t2xlargecomputeprimeappt2nanob8j&requestduration=1000
        elif '/getActualTRN' in self.path:
            parsed = urlparse.urlparse(self.path)
            appName = urlparse.parse_qs(parsed.query)['appname'][0]
            folderName = urlparse.parse_qs(parsed.query)['containerName'][0]
            requestDuration = urlparse.parse_qs(parsed.query)['requestduration'][0]
            mainServiceName = urlparse.parse_qs(parsed.query)['mainServiceName'][0]
            appNames = [mainServiceName]
            folderNames = [folderName]
            df = self.getAndCombineAllDbs(host, port, username, password,appNames, folderNames)
            hit=0
            finaldf = df
            idxt = 0
            for idxt, valt in enumerate(df['requests_duration_mean']):
                if(valt > float(requestDuration)):
                    hit+=1
                    print(valt)
                if(hit >=10):
                    break
            if(hit==0):
                for idxt, valt in enumerate(df['pod_util_cpu_avg']):
                    threshhold = df['pod_cpu_limit'][idxt] - 0.3*df['pod_cpu_limit'][idxt]
                    if(valt > threshhold):
                        hit+=1
                        print (valt)
                    if(hit >=10):
                        break
            finaldf = df.head(idxt)

            finaldf= finaldf.sort_values(['requests'])
            val = float(finaldf.tail(1).requests)
            output=""+str(val)
            print (val)
            print (output)
            self.wfile.write(output.encode())
        else:
            self.wfile.write("Path not found".encode())

    def do_HEAD(self):
        self._set_headers()

host = os.environ['INFLUXDB_HOST']
port = os.environ['INFLUXDB_PORT']
username= os.environ['INFLUXDB_USER']
password = os.environ['INFLUXDB_PASS']

mongo_host = os.environ['MONGODB_HOST']
mongo_port = os.environ['MONGODB_PORT']
mongo_username= os.environ['MONGODB_USER']
mongo_password = os.environ['MONGODB_PASS']

def run(server_class=HTTPServer, handler_class=Terminus, port=9002):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print ('Starting httpd...')
    httpd.serve_forever()

if __name__ == "__main__":
    from sys import argv

    if len(argv) == 2:
        run(port=int(argv[1]))
    else:
        run()
        