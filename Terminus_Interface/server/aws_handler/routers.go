package TERMINUS

import (
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	// Declare the static file directory and point it to the directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	router.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"approachVMBootingTime",
		"GET",
		"/approachVMBootingTime",
		ApproachVMBootingTime,
	},
	Route{
		"logsPage",
		"GET",
		"/logsPage",
		LogsPage,
	},
	Route{
		"CompareLimits",
		"GET",
		"/compareLimits",
		CompareLimits,
	},
	Route{
		"ConductTestMSCBruteForce",
		"GET",
		"/conductTestMSCBruteForce",
		ConductTestMSCBruteForce,
	},
	Route{
		"InterInstancePerfCompare",
		"GET",
		"/interInstancePerfCompare",
		InterInstancePerfCompare,
	},
	Route{
		"InterReplicasPerfCompare",
		"GET",
		"/interReplicasPerfCompare",
		InterReplicasPerfCompare,
	},
	Route{
		"MSCAnalyzed",
		"GET",
		"/mSCAnalyzed",
		MSCAnalyzed,
	},
	Route{
		"MSCAnalyzedAll",
		"GET",
		"/mSCAnalyzedAll",
		MSCAnalyzedAll,
	},
	Route{
		"MSCPredictionAll",
		"GET",
		"/mSCPredictionAll",
		MSCPredictionAll,
	},
	Route{
		"MSCAllTests",
		"GET",
		"/mSCAllTests",
		MSCAllTests,
	},
	Route{
		"PodBootingTimeResults",
		"GET",
		"/podBootingTimeResults",
		PodBootingTimeResults,
	},
	Route{
		"Conduct_pod_booting_timetest",
		"GET",
		"/conduct_pod_booting_timetest",
		Conduct_pod_booting_timetest,
	},
	Route{
		"VMBootingTimeResults",
		"GET",
		"/vMBootingTimeResults",
		VMBootingTimeResults,
	},
	Route{
		"Sandboxing",
		"GET",
		"/sandboxing",
		Sandboxing,
	},
	Route{
		"Sandboxing",
		"GET",
		"/sandboxing",
		Sandboxing,
	},
	Route{
		"SandboxingTree",
		"GET",
		"/sandboxingTree",
		SandboxingTree,
	},

	Route{
		"SandboxingYaml",
		"GET",
		"/sandboxedYaml",
		SandboxedYaml,
	},

	Route{
		"SandboxingTests",
		"GET",
		"/sandboxingTests",
		SandboxingTests,
	},
	Route{
		"Init",
		strings.ToUpper("Post"),
		"/init",
		InitUserConfig,
	},

	Route{
		"ListAllInstances",
		strings.ToUpper("Get"),
		"/instances",
		ListAllInstances,
	},

	//collectDataNoLimits/1/1/1/t2.medium/t2.large/t2.xlarge/400/compute/primeapp/_api_prime
	//collectDataNoLimits/1/1/1/t2.large/t2.large/t2.xlarge/400/compute/primeapp/_api_prime
	//collectDataNoLimits/1/1/1/t2.small/t2.large/t2.xlarge/350/compute/primeapp/_api_prime
	//collectDataNoLimits/1/1/1/t2.micro/t2.large/t2.xlarge/250/compute/primeapp/_api_prime
	//collectDataNoLimits/1/1/1/t2.nano/t2.large/t2.xlarge/150/compute/primeapp/_api_prime

	//collectDataNoLimits/1/1/1/t2.medium/t2.large/t2.xlarge/500/dbaccess/movieapp/_api_movies
	//collectDataNoLimits/1/1/1/t2.large/t2.large/t2.xlarge/500/dbaccess/movieapp/_api_movies
	//collectDataNoLimits/1/1/1/t2.small/t2.large/t2.xlarge/350/dbaccess/movieapp/_api_movies
	//collectDataNoLimits/1/1/1/t2.micro/t2.large/t2.xlarge/350/dbaccess/movieapp/_api_movies
	//collectDataNoLimits/1/1/1/t2.nano/t2.large/t2.xlarge/350/dbaccess/movieapp/_api_movies


	//collectDataNoLimits/1/1/1/t2.medium/t2.large/t2.xlarge/3500/web/webacapp/_api_web
	//collectDataNoLimits/1/1/1/t2.large/t2.large/t2.xlarge/4000/web/webacapp/_api_web //
	//collectDataNoLimits/1/1/1/t2.small/t2.large/t2.xlarge/3000/web/webacapp/_api_web
	//collectDataNoLimits/1/1/1/t2.micro/t2.large/t2.xlarge/3000/web/webacapp/_api_web
	//collectDataNoLimits/1/1/1/t2.nano/t2.large/t2.xlarge/3000/web/webacapp/_api_web


	Route{
		"CollectDataNoLimits",
		strings.ToUpper("Get"),
		"/collectDataNoLimits/{numsvs}/{nodecount}/{replicas}/{instancetype}/{mastertype}/{testvmtype}/{maxRPS}/{apptype}/{appname}/{mainservicename}/{apiendpoint}",
		CollectDataNoLimits,
	},

	//collectDataLimits/1/1/1/t2.micro/80/100/compute/primeapp/_api_prime

	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/150/t2.nano/compute/primeapp/_api_prime
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/150/t2.nano/compute/primeapp/_api_prime
	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/250/t2.micro/compute/primeapp/_api_prime
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/250/t2.micro/compute/primeapp/_api_prime
	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/350/t2.small/compute/primeapp/_api_prime
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/350/t2.small/compute/primeapp/_api_prime

	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/150/t2.nano/compute/primeapp/_api_prime //
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/250/t2.micro/compute/primeapp/_api_prime //
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/350/t2.small/compute/primeapp/_api_prime //
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/400/t2.medium/compute/primeapp/_api_prime //
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/400/t2.large/compute/primeapp/_api_prime //


	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/350/t2.nano/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/350/t2.nano/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/350/t2.micro/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/350/t2.micro/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/350/t2.small/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/350/t2.small/dbaccess/movieapp/_api_movies

	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/350/t2.nano/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/350/t2.micro/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/350/t2.small/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/500/t2.medium/dbaccess/movieapp/_api_movies
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/500/t2.large/dbaccess/movieapp/_api_movies  //


	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/3000/t2.nano/web/webacapp/_api_web
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/300/t2.nano/web/webacapp/_api_web //
	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/3000/t2.micro/web/webacapp/_api_web
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/3000/t2.micro/web/webacapp/_api_web  //
	//collectDataLimits/1/1/1/t2.medium/t2.large/t2.xlarge/3000/t2.small/web/webacapp/_api_web
	//collectDataLimits/1/1/1/t2.large/t2.large/t2.xlarge/3000/t2.small/web/webacapp/_api_web

	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/3000/t2.nano/web/webacapp/_api_web
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/3000/t2.micro/web/webacapp/_api_web
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/3000/t2.small/web/webacapp/_api_web
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/3500/t2.medium/web/webacapp/_api_web
	//collectDataLimits/1/1/1/t2.xlarge/t2.large/t2.xlarge/4000/t2.large/web/webacapp/_api_web




	//collectDataLimits/4/1/1/t2.xlarge/t2.large/t2.xlarge/150/100/mix/mixallapp/serveapp/_test/
	Route{
		"CollectDataLimits",
		strings.ToUpper("Get"),
		"/collectDataLimits/{numsvs}/{nodecount}/{replicas}/{instancetype}/{mastertype}/{testvmtype}/{maxRPS}/{limitinginstancetype}/{apptype}/{appname}/{mainservicename}/{apiendpoint}",
		CollectDataLimits,
	},
	//conductTestToCalculatePodBootTime/1/1/1/t2.xlarge/t2.large/t2.xlarge/40/t2.large/compute/primeapp/primeapp/_api_prime
	Route{
		"ConductTestToCalculatePodBootTime",
		strings.ToUpper("Get"),
		"/conductTestToCalculatePodBootTime/{numsvs}/{nodecount}/{replicas}/{instancetype}/{mastertype}/{testvmtype}/{maxRPS}/{limitinginstancetype}/{apptype}/{appname}/{mainservicename}/{apiendpoint}",
		ConductTestToCalculatePodBootTime,
	},
	// get all instances data wheerever it is possible to host that instance
	//getRelevantInstancesDataBasedOnLimits/compute/primeapp/primeapp/t2.small/1
	Route{
		"GetRelevantInstancesDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstancesDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}",
		GetRelevantInstancesDataBasedOnLimits,
	},
	// to tell where specificlly the comparing instance is hosted
	//getRelevantInstanceDataBasedOnLimits/compute/primeapp/primeapp/t2.large/t2.small/1
	//getRelevantInstanceDataBasedOnLimits/compute/primeapp/primeapp/t2.large/100/1
	Route{
		"GetRelevantInstanceDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstanceDataBasedOnLimits/{apptype}/{appname}/{hostinstance}/{mainservicename}/{comparinginstance}/{replicas}",
		GetRelevantInstanceDataBasedOnLimits,
	},
	// get all instances data wheerever it is possible to host that instance along with actual VM data
	//getRelevantInstancesActualInstanceDataBasedOnLimits/compute/primeapp/primeapp/t2.small/1
	Route{
		"GetRelevantInstancesActualInstanceDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstancesActualInstanceDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}",
		GetRelevantInstancesActualInstanceDataBasedOnLimits,
	},
	// to tell where specificlly the comparing instance is hosted along with actual VM data of comparing instance
	//getRelevantInstanceActualInstanceDataBasedOnLimits/compute/primeapp/primeapp/t2.large/t2.small/1
	Route{
		"GetRelevantInstanceActualInstanceDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstanceActualInstanceDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{hostinstance}/{comparinginstance}/{replicas}",
		GetRelevantInstanceActualInstanceDataBasedOnLimits,
	},
	// get all instances requests data wheerever it is possible to host that instance
	//getRelevantInstancesRequestsDataBasedOnLimits/compute/primeapp/primeapp/t2.small/1
	Route{
		"GetRelevantInstancesRequestsDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstancesRequestsDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}",
		GetRelevantInstancesRequestsDataBasedOnLimits,
	},
	// to tell where specificlly the comparing instance is hosted
	//getRelevantInstanceRequestsDataBasedOnLimits/compute/primeapp/primeapp/t2.small/100/1
	Route{
		"GetRelevantInstanceRequestsDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstanceRequestsDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{hostinstance}/{comparinginstance}/{replicas}",
		GetRelevantInstanceRequestsDataBasedOnLimits,
	},
	// get all instances data wheerever it is possible to host that instance along with actual VM data
	//getRelevantInstancesActualInstanceRequestsDataBasedOnLimits/compute/primeapp/primeapp/t2.small/1
	Route{
		"GetRelevantInstancesActualInstanceRequestsDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstancesActualInstanceRequestsDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}",
		GetRelevantInstancesActualInstanceRequestsDataBasedOnLimits,
	},
	// to tell where specificlly the comparing instance is hosted along with actual VM data of comparing instance
	//getRelevantInstanceActualInstanceRequestsDataBasedOnLimits/compute/primeapp/primeapp/t2.large/t2.small/1
	Route{
		"GetRelevantInstanceActualInstanceRequestsDataBasedOnLimits",
		strings.ToUpper("Get"),
		"/getRelevantInstanceActualInstanceRequestsDataBasedOnLimits/{apptype}/{appname}/{mainservicename}/{hostinstance}/{comparinginstance}/{replicas}",
		GetRelevantInstanceActualInstanceRequestsDataBasedOnLimits,
	},
	// get actual vm performance data
	//getInstancesPerfData/compute/primeapp/primeapp/t2.small/1/none
	Route{
		"GetInstancesPerfData",
		strings.ToUpper("Get"),
		"/getInstancesPerfData/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}/{hostinstance}",
		GetInstancesPerfData,
	},
	// get actual vm requests data
	Route{
		"GetInstancesRequestsData",
		strings.ToUpper("Get"),
		"/getInstancesRequestsData/{apptype}/{appname}/{mainservicename}/{comparinginstance}/{replicas}/{hostinstance}",
		GetInstancesRequestsData,
	},


	Route{
		"DoBruteForceTraining",
		strings.ToUpper("Get"),
		"/doBruteForceTraining/{apptype}/{appname}/{mainservicename}",
		DoBruteForceTraining,
	},
	Route{
		"DoSmartTestAllTraining",
		strings.ToUpper("Get"),
		"/doSmartTestAllTraining/{apptype}/{appname}/{mainservicename}",
		DoSmartTestAllTraining,
	},
	Route{
		"DoBruteForceTrainingReplicas",
		strings.ToUpper("Get"),
		"/doBruteForceTrainingReplicas/{apptype}/{appname}/{mainservicename}",
		DoBruteForceTrainingReplicas,
	},


	Route{
		"GetPredictedRegressionTRN",
		strings.ToUpper("Get"),
		"/getPredictedRegressionTRN/{apptype}/{appname}/{mainservicename}/{replicas}/{numcoresutil}/{numcoreslimit}/{nummemlimit}",
		GetPredictedRegressionTRN,
	},
	Route{
		"GetPredictedRegressionReplicas",
		strings.ToUpper("Get"),
		"/getPredictedRegressionReplicas/{apptype}/{appname}/{mainservicename}/{msc}/{numcoresutil}/{numcoreslimit}/{nummemlimit}",
		GetPredictedRegressionReplicas,
	},
	Route{
		"GetPredictedRegressionTRNSmart",
		strings.ToUpper("Get"),
		"/getPredictedRegressionTRNSmart/{apptype}/{appname}/{mainservicename}/{numcoresutil}{nummemutil}/{foldername}",
		GetPredictedRegressionTRNSmart,
	},
	Route{
		"StoreAllExperimentalTRN",
		strings.ToUpper("Get"),
		"/storeAllExperimentalTRN",
		StoreAllExperimentalTRN,
	},
	Route{
		"StoreAllTRNRegressionIntoMongo",
		strings.ToUpper("Get"),
		"/storeAllTRNRegressionIntoMongo",
		StoreAllTRNRegressionIntoMongo,
	},
	Route{
		"GetAllExperimentalTRN",
		strings.ToUpper("Get"),
		"/getAllExperimentalTRN",
		GetAllExperimentalTRN,
	},
	Route{
		"GetExperimentalTRNsMongoDB",
		strings.ToUpper("Get"),
		"/getExperimentalTRNsMongoDB/{apptype}/{appname}/{mainservicename}/{hostinstance}/",
		GetExperimentalTRNsMongoDB,
	},
	Route{
		"GetExperimentalTRNsMongoDBAll",
		strings.ToUpper("Get"),
		"/getExperimentalTRNsMongoDBAll/{apptype}/{appname}/{mainservicename}",
		GetExperimentalTRNsMongoDBAll,
	},
	Route{
		"GetRegressionTRNsMongoDBAll",
		strings.ToUpper("Get"),
		"/getRegressionTRNsMongoDBAll/{apptype}/{appname}/{mainservicename}",
		GetRegressionTRNsMongoDBAll,
	},
	Route{
		"GetRMSErrorMongoDB",
		strings.ToUpper("Get"),
		"/getRMSErrorMongoDB/{apptype}/{appname}/{mainservicename}",
		GetRMSErrorMongoDB,
	},
/*
	Route{
		"StoreAllNamesToMongo",
		strings.ToUpper("Get"),
		"/storeAllNamesToMongo",
		StoreAllNamesToMongo,
	},
*/
	Route{
		"GetAppTypes",
		strings.ToUpper("Get"),
		"/getAppTypes",
		GetAppTypes,
	},
	Route{
		"GetAppNames",
		strings.ToUpper("Get"),
		"/getAppNames/{apptype}",
		GetAppNames,
	},
	Route{
		"GetMetrics",
		strings.ToUpper("Get"),
		"/getMetrics",
		GetMetricsAPI,
	},
	Route{
		"GetAllTestsInfo",
		strings.ToUpper("Get"),
		"/getAllTestsInfo",
		GetAllTestsInfo,
	},
	Route{
		"AnalyzeKubeEvents",
		strings.ToUpper("Get"),
		"/analyzeKubeEvents/{apptype}/{appname}/{mainservicename}",
		AnalyzeKubeEvents,
	},
	Route{
		"GetAnalyzedKubeEvents",
		strings.ToUpper("Get"),
		"/getAnalyzedKubeEvents/{apptype}/{appname}/{mainservicename}",
		GetAnalyzedKubeEvents,
	},
	Route{
		"GetPodBootingTime",
		strings.ToUpper("Get"),
		"/getPodBootingTime/{apptype}/{appname}/{mainservicename}",
		GetPodBootingTime,
	},

	Route{
		"TestInstanceBootingShuttingTime",
		strings.ToUpper("Get"),
		"/testbootshutdowntime",
		TestInstanceBootingShuttingTime,
	},

	Route{
		"RandomTestInstanceBootingShuttingTime",
		strings.ToUpper("Get"),
		"/randomtestbootshutdowntime",
		RandomTestInstanceBootingShuttingTime,
	},

	Route{
		"GetAllVMTypesBootShutDownDataAvg",
		strings.ToUpper("Get"),
		"/getAllVMTypesBootShutDownDataAvg",
		GetAllVMTypesBootShutDownDataAvg,
	},

	Route{
		"TrainDataSetRegression",
		strings.ToUpper("Get"),
		"/trainDataSetRegression",
		TrainDataSetRegression,
	},

	Route{
		"GetandStoreRegressionValues",
		strings.ToUpper("Get"),
		"/getandStoreRegressionValues",
		GetandStoreRegressionValues,
	},

	Route{
		"GetAllVMTypesBootShutDownDataRegression",
		strings.ToUpper("Get"),
		"/getAllVMTypesBootShutDownDataRegression",
		GetAllVMTypesBootShutDownDataRegression,
	},

	Route{
		"GetPerVMTypeAllBootShutDownData",
		strings.ToUpper("Get"),
		"/getPerVMTypeAllBootShutDownData",
		GetPerVMTypeAllBootShutDownData,
	},

	Route{
		"GetPerVMTypeOneBootShutDownData",
		strings.ToUpper("Get"),
		"/getPerVMTypeOneBootShutDownData",
		GetPerVMTypeOneBootShutDownData,
	},

	// sandboxing intf routes
	Route{
		"FormTreeFromDockerCompose",
		strings.ToUpper("Get"),
		"/formTreeFromDockerCompose/{filename}/{mainservicename}",
		FormTreeFromDockerCompose,
	},

	Route{
		"GetModifiedDockerComposeYaml",
		strings.ToUpper("GET"),
		"/getModifiedDockerComposeYaml/{serviceName}/{original}/",
		GetModifiedDockerComposeYaml,
	},

	Route{
		"GetAvailableApps",
		strings.ToUpper("GET"),
		"/getAvailableApps",
		GetAvailableApps,
	},

	Route{
		"getModifiedServiceEndpoint",
		strings.ToUpper("GET"),
		"/getModifiedServiceEndpoint/{serviceName}/{original}/",
		GetModifiedServiceEndpoint,
	},
	Route{
		"formTreeFromDockerComposeMongo",
		strings.ToUpper("GET"),
		"/formTreeFromDockerComposeMongo/{serviceName}/{original}/",
		FormTreeFromDockerComposeMongo,
	},
	Route{
		"ParseYamlToJson",
		strings.ToUpper("GET"),
		"/parseyamltojson/{filename}",
		ParseYamlToJson,
	},
	Route{
		"ParseYamlToYaml",
		strings.ToUpper("GET"),
		"/parseYamlToYaml/{filename}",
		ParseYamlToYaml,
	},
	Route{
		"DeployserviceToRecordAPIAndResponse",
		strings.ToUpper("GET"),
		"/deployserviceToRecordAPIAndResponse/{filename}/{endpointapi}/{mainservicename}",
		DeployserviceToRecordAPIAndResponse,
	},
	Route{
		"GetAllDockerComposeFileNames",
		strings.ToUpper("GET"),
		"/getAllDockerComposeFileNames/",
		GetAllDockerComposeFileNames,
	},
	Route{
		"GetAllDockerComposeServiceNames",
		strings.ToUpper("GET"),
		"/getAllDockerComposeServiceNames/{filename}",
		GetAllDockerComposeServiceNames,
	},

	Route{
		"GetAllTreesFromDockerComposeMongo",
		strings.ToUpper("GET"),
		"/getAllTreesFromDockerComposeMongo/{filename}/{servicename}",
		GetAllTreesFromDockerComposeMongo,
	},

	Route{
		"GetAllTestsInformationSandbox",
		strings.ToUpper("GET"),
		"/getAllTestsInformationSandbox/",
		GetAllTestsInformationSandbox,
	},


	Route{
		"Swaggeryaml",
		strings.ToUpper("GET"),
		"/swaggeryaml",
		Swaggeryaml,
	},



}
