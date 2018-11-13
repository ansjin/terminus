package TERMINUS

import (
	"os"
	"net/http"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"math"
)
func InitUserConfigurationFirst(){

	AllInstanceTypes = []string{ "100","200","500","700", "t2.micro", "t2.nano", "t2.small", "t2.medium", "t2.large","t2.xlarge"}

	DefaultRegion = []string{ "us-east-2", "us-east-1", "us-west-1", "eu-central-1", "ap-south-1", "eu-west-1", "ap-southeast-2"}

	DefaultZone = []string{ "us-east-2a","us-east-1a","us-west-1a","eu-central-1a","ap-south-1a","eu-west-1a", "ap-southeast-2a"}

	DefaultAMI = []string{ "ami-5e8bb23b", "ami-759bc50a", "ami-4aa04129", "ami-de8fb135", "ami-188fba77","ami-2a7d75c0", "ami-47c21a25"}

	B = &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}
	AWSConfig =  AWSConfigStruct{AwsAccessKeyId: os.Getenv("AWS_KEY"), AwsSecretAccessKey: os.Getenv("AWS_SECRET"),
		Region:  os.Getenv("AWS_DEFAULT_REGION"), KeyPairName:os.Getenv("AWS_KEY_PAIR_NAME"), SubnetId:os.Getenv("AWS_SUBNET_ID"),
		SecurityGroup:os.Getenv("AWS_SECURITY_GROUP")}

	SupportedNumMicroserviceReplicas = []string{"1", "2", "3", "4", "5"}
	TestNumMicroserviceReplicas = []string{"1", "2", "3", "4", "5"}
	SupportedNumMicroservicesApp = []string{"1", "4"}
	DefaultAppNames = []string{"primeapp", "movieapp","webacapp", "mixalapp", "serveapp"}
	SupportedTypesOfMicroservices = []string{"compute", "dbaccess", "web", "mix", "sandbox"}
	AppsAPIEndPoint = []string{"/api/prime", "api/movies", "/api/web", "/test", "/test"}

	computeAppNames = []string{"primeapp"}
	dbaccessAppNames = []string{"movieapp"}
	webAppNames = []string{"webacapp"}

	mixAppNames = []string{"mixalapp"}
	sandboxAppNames = []string{"serveapp"}
	mainServiceNamesMix = []string{"primeapp","movieapp", "webacapp", "serveapp"}

	MaxTestIterations = 1
	//InitLimitsDatabaseNames()
	//InitNonLimitsDatabaseNames()
	InitRelevantInstanceTypesForLimits()
}


func getLimits(instanceType string) LimitsInfo{

	switch instanceType {
	case "100":
		return LimitsInfo{0.100, 0.100, 0}
	case "200":
		return LimitsInfo{0.200, 0.200, 0}
	case "500":
		return LimitsInfo{0.500, 0.500, 0}
	case "700":
		return LimitsInfo{0.700, 0.700, 0}
	case "t2.nano":
		return LimitsInfo{1, 0.500, 0}
	case "t2.micro":
		return LimitsInfo{1, 1, 0}
	case "t2.small":
		return LimitsInfo{1, 2, 0}
	case "t2.medium":
		return LimitsInfo{2, 4, 0}
	case "t2.large":
		return LimitsInfo{2, 8, 0}
	case "t2.xlarge":
		return LimitsInfo{4, 16, 0}
	case "t2.2xlarge":
		return LimitsInfo{8, 32, 0}
	case "m5.xlarge":
		return LimitsInfo{4, 16, 0}
	default:
		return LimitsInfo{}
	}

}
/*
func storeAllNamesInMongo() {
	for idx, appName := range DefaultAppNames {
		appType := SupportedTypesOfMicroservices[idx]
		appAPIEndPoint := AppsAPIEndPoint[idx]
		for _, hostInstance := range AllInstanceTypes {
			for _, replica := range SupportedNumMicroserviceReplicas {
			for _, limitingInstance := range AllInstanceTypes {
					containerName := LimitsInstancesDatabaseNames[LimitsInstancesDatabase{
						appType, appName, hostInstance, limitingInstance,
						replica}]

					if containerName != "" {
						AllData := InfluDbTabaseNamesInfo{HostInstanceType: hostInstance,
							ServiceType: appType,
							ServiceName: appName,
							Replicas: replica,
							LimitingInstanceType: limitingInstance,
							Limits: getLimits(limitingInstance),
							TestAPI: appAPIEndPoint,
							FolderName: containerName,
							InfluxDbk6Name: containerName + "_TestK6",
							InfluxdbK8sName: containerName + "_k8s",
							NodeCount: "1",
							NumMicroSvs: "1",
							TestIter: "1",
						}
						mongoSession := GetMongoSession()
						collection := mongoSession.DB(Database).C(Collection_All_Test_Names)
						if err := collection.Insert(AllData); err != nil {
							log.Info("error ", err)
						} else {
							log.Info("#inserted into ", Collection_All_Test_Names)
						}
						defer mongoSession.Close()

					}
				}
				containerName := NonLimitsInstancesDatabaseNames[NonLimitsInstancesDatabase{
					appType, appName, hostInstance,
					replica}]

				if containerName != "" {
					AllData := InfluDbTabaseNamesInfo{HostInstanceType: hostInstance,
						ServiceType: appType,
						ServiceName: appName,
						Replicas: replica,
						LimitingInstanceType: "",
						Limits: LimitsInfo{},
						TestAPI: appAPIEndPoint,
						FolderName: containerName,
						InfluxDbk6Name: containerName + "_TestK6",
						InfluxdbK8sName: containerName + "_k8s",
						NodeCount: "1",
						NumMicroSvs: "1",
						TestIter: "1",
					}
					mongoSession := GetMongoSession()
					collection := mongoSession.DB(Database).C(Collection_All_Test_Names)
					if err := collection.Insert(AllData); err != nil {
						log.Info("error ", err)
					} else {
						log.Info("#inserted into ", Collection_All_Test_Names)
					}
					defer mongoSession.Close()
				}
			}
		}
	}
}
*/

func storeMSCMongo(appType, appName,mainServiceName, instanceFamily, requestDuration, appAPIEndPoint string ){

	for _, hostInstance := range AllInstanceTypes {

		var AllProfiles= []MSCProfiles{}
		for _, limitingInstance := range AllInstanceTypes {

			var mscPerReplica []MSCInfo
			for _, replica := range SupportedNumMicroserviceReplicas {
				containerName := GetTestFolderName(hostInstance, limitingInstance, replica, appType, appName, true)
				if containerName != "" {

					urlExp := "http://regression:9002/getActualTRN?appname=" + appName + "&containerName=" + containerName + "&requestduration=" + requestDuration+
						"&mainServiceName="+ mainServiceName
					log.Info(urlExp)
					responseExp, err := http.Get(urlExp)
					if err != nil {
						log.Info(err)
					} else {
						defer responseExp.Body.Close()
						contentsExp, err := ioutil.ReadAll(responseExp.Body)
						if err != nil {
							log.Info(err)
						}
						log.Info(string(contentsExp))
						limits := getLimits(limitingInstance)
						urlReg := "http://regression:9002/getPredictionPreTrained?appname=" + appName + "&apptype=" + appType +
							"&replicas=" + replica +
							"&numcoresutil=" + FloatToString(limits.Cpu_cores) +
							"&numcoreslimit=" + FloatToString(limits.Cpu_cores) +
							"&nummemlimit=" + FloatToString(limits.Mem_gb) +
							"&instancefamily=" + instanceFamily +
							"&requestduration=" + requestDuration+
							"&mainServiceName="+ mainServiceName
						log.Info(urlReg)
						responseReg, err := http.Get(urlReg)
						if err != nil {
							log.Info(err)
						} else {
							defer responseReg.Body.Close()
							contentsReg, err := ioutil.ReadAll(responseReg.Body)
							if err != nil {
								log.Info(err)
							}
							log.Info(string(contentsReg))

							urlReplicasPrediction := "http://regression:9002/getPredictionReplicas?appname=" + appName + "&apptype=" + appType +
								"&msc=" + string(contentsExp) +
								"&numcoresutil=" + FloatToString(limits.Cpu_cores) +
								"&numcoreslimit=" + FloatToString(limits.Cpu_cores) +
								"&nummemlimit=" + FloatToString(limits.Mem_gb) +
								"&instancefamily=" + instanceFamily +
								"&requestduration=" + requestDuration+
								"&mainServiceName="+ mainServiceName
							log.Info(urlReplicasPrediction)
							responseReplicaPred, err := http.Get(urlReplicasPrediction)
							if err != nil {
								log.Info(err)
							} else {
								defer responseReplicaPred.Body.Close()
								contentsReplicaPred, err := ioutil.ReadAll(responseReplicaPred.Body)
								if err != nil {
									log.Info(err)
								}
								log.Info(string(contentsReplicaPred))
								resultReplicaPred := StringToFloat(string(contentsReplicaPred))
								resultReplicaPred = math.Ceil(resultReplicaPred)
								log.Info(resultReplicaPred)

								url := "http://regression:9002/smartTestGetResult?appname=" + appName + "&apptype=" + appType +
									"&numcoresutil=" + FloatToString(limits.Cpu_cores * StringToFloat(replica)) +
									"&nummemutil=" + FloatToString(limits.Mem_gb) +
									"&instancefamily=" + instanceFamily +
									"&requestduration=" + requestDuration +
									"&containerName=" + containerName+
									"&mainServiceName="+ mainServiceName
								log.Info(url)
								response, err := http.Get(url)
								if err != nil {
									log.Info(err)
								} else {
									defer response.Body.Close()
									contentsSmart, err := ioutil.ReadAll(response.Body)
									if err != nil {
										log.Info(err)
									}
									log.Info(string(contentsSmart))

									mongoSessionPodBooting := GetMongoSession()
									collectionPodBooting := mongoSessionPodBooting.DB(Database).C(Collection_Pod_Bootingtime)
									var podBootingTimeObject PodBootingTimeStruct
									replicaFlot64 := StringToFloat(replica)
									errPod := collectionPodBooting.Find(bson.M{"replicas": replicaFlot64, "appname": appName, "apptype": appType, "mainservicename": mainServiceName}).One(&podBootingTimeObject)
									if errPod != nil {
										log.Info("error ", errPod)
									}
									defer mongoSessionPodBooting.Close()

									msc := MSCInfo{Replicas: replicaFlot64,
										Maximum_service_capacity_per_min: MSCDetails{StringToFloat(string(contentsExp)),
											StringToFloat(string(contentsReg)), StringToFloat(string(contentsSmart))},
										Maximum_service_capacity_per_sec: MSCDetails{StringToFloat(string(contentsExp)) / 60,
											StringToFloat(string(contentsReg)) / 60, StringToFloat(string(contentsSmart)) / 60},
										Pod_boot_time_ms: podBootingTimeObject.MeanBootingTime, Sd_Pod_boot_time_ms: podBootingTimeObject.SdBootingTime,
										PredictedReplicas: resultReplicaPred}
									mscPerReplica = append(mscPerReplica, msc)
								}
							}

						}
					}
				}
			}
			if len(mscPerReplica) > 0 {
				profile := MSCProfiles{getLimits(limitingInstance), mscPerReplica}
				AllProfiles = append(AllProfiles, profile)
			}
		}
		log.Info("Length: ", len(AllProfiles))
		if (len(AllProfiles) > 0) {
			AllData := MSCValueObject{
				HostInstanceType: hostInstance,
				ServiceType:      appType,
				ServiceName:      appName,
				MainServiceName: mainServiceName,
				Profiles:         AllProfiles,
				TestAPI:          appAPIEndPoint,
			}
			mongoSession := GetMongoSession()
			collection := mongoSession.DB(Database).C(Collection_MSC_EXPERIMENTAL)
			if err := collection.Insert(AllData); err != nil {
				log.Info("error ", err)
			} else {
				log.Info("#inserted into ", Collection_MSC_EXPERIMENTAL)
			}
			defer mongoSession.Close()
		}
	}
}
func storeAllTRNIntoMongo()  {
	requestDuration := "1000"
	///getActualTRN?appname=primeapp&containerName=s1t1rc1nc1t2xlargecomputeprimeappt2nanob8j&requestduration=1000
	instanceFamily:="t2"
	for idx, appName := range DefaultAppNames {

		appType := SupportedTypesOfMicroservices[idx]
		appAPIEndPoint := AppsAPIEndPoint[idx]
		if(appType=="mix") {

			for _,mainSericeName:=range mainServiceNamesMix{
				storeMSCMongo(appType, appName, mainSericeName, instanceFamily, requestDuration, appAPIEndPoint)
			}
		}else{
			storeMSCMongo(appType, appName, appName, instanceFamily, requestDuration, appAPIEndPoint)
		}
	}
}
func storeMSCRegression(appType, appName,mainServiceName, instanceFamily, requestDuration, appAPIEndPoint string){
	var AllProfiles = []MSCProfiles{}
	for _, limitingInstance := range AllInstanceTypes {
		var mscPerReplica []MSCInfo
		for _, replica := range TestNumMicroserviceReplicas {
			limits:=getLimits(limitingInstance)
			url := "http://regression:9002/getPredictionPreTrained?appname=" + appName + "&apptype=" + appType +
				"&replicas=" + replica+
				"&numcoresutil=" + FloatToString(limits.Cpu_cores)+
				"&numcoreslimit=" + FloatToString(limits.Cpu_cores)+
				"&nummemlimit=" + FloatToString(limits.Mem_gb)+
				"&instancefamily=" + instanceFamily+
				"&requestduration=" + requestDuration+
				"&mainServiceName="+ mainServiceName
			log.Info(url)
			response, err := http.Get(url)
			if err != nil {
				log.Info(err)
			} else {
				defer response.Body.Close()
				contents, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Info(err)
				}
				log.Info(string(contents))
				mongoSessionPodBooting := GetMongoSession()
				collectionPodBooting := mongoSessionPodBooting.DB(Database).C(Collection_Pod_Bootingtime)
				var podBootingTimeObject PodBootingTimeStruct
				replicaFlot64:=StringToFloat(replica)
				errPod := collectionPodBooting.Find(bson.M{"replicas": replicaFlot64, "appname": appName, "apptype": appType, "mainservicename": mainServiceName}).One(&podBootingTimeObject)
				if errPod != nil {
					log.Info("error ", errPod)
				}
				defer mongoSessionPodBooting.Close()

				msc := MSCInfo{Replicas: replicaFlot64,
					Maximum_service_capacity_per_min: MSCDetails{RegBruteForce:StringToFloat(string(contents))},
					Maximum_service_capacity_per_sec: MSCDetails{RegBruteForce:StringToFloat(string(contents))/60},
					Pod_boot_time_ms: podBootingTimeObject.MeanBootingTime, Sd_Pod_boot_time_ms:podBootingTimeObject.SdBootingTime}
				mscPerReplica = append(mscPerReplica, msc)
			}
		}
		if len(mscPerReplica) > 0{
			profile := MSCProfiles{getLimits(limitingInstance), mscPerReplica}
			AllProfiles = append(AllProfiles, profile)
		}
	}
	log.Info("Length: ",len(AllProfiles))
	if (len(AllProfiles) > 0) {
		AllData := MSCValueObject{
			HostInstanceType: "",
			ServiceType:      appType,
			ServiceName:      appName,
			MainServiceName: mainServiceName,
			Profiles:         AllProfiles,
			TestAPI:          appAPIEndPoint,
		}
		mongoSession := GetMongoSession()
		collection := mongoSession.DB(Database).C(Collection_MSC_REGRESSION)
		if err := collection.Insert(AllData); err != nil {
			log.Info("error ", err)
		} else {
			log.Info("#inserted into ", Collection_MSC_REGRESSION)
		}
		defer mongoSession.Close()
	}
}
func storeAllTRNRegressionIntoMongo()  {
	requestDuration := "1000"
	instanceFamily:="t2"
	for idx, appName := range DefaultAppNames {
		appType := SupportedTypesOfMicroservices[idx]
		appAPIEndPoint := AppsAPIEndPoint[idx]
		if(appType=="mix") {

			for _,mainSericeName:=range mainServiceNamesMix{
				storeMSCRegression(appType, appName, mainSericeName, instanceFamily, requestDuration, appAPIEndPoint)
			}
		}else{
			storeMSCRegression(appType, appName, appName, instanceFamily, requestDuration, appAPIEndPoint)
		}

	}
}
func getAllTRNsMongo() []MSCValueObject{
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_MSC_EXPERIMENTAL)
	var allData  []MSCValueObject
	err := collection.Find(nil).All(&allData)
	if err != nil {
		log.Info("error ", err)
	}
	defer mongoSession.Close()

	return  allData
}
func getTRNsMongo(appName, appType, hostInstance, mainServiceName string) MSCValueObject{
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_MSC_EXPERIMENTAL)
	var allData  MSCValueObject
	err := collection.Find(bson.M{"hostinstancetype": hostInstance, "servicename": appName, "mainservicename": mainServiceName}).One(&allData)
	if err != nil {
		log.Info("error ", err)
	}
	defer mongoSession.Close()

	return  allData
}
func getTRNsMongoAll(appName, appType, mainServiceName string) []MSCValueObject{
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_MSC_EXPERIMENTAL)
	var allData  []MSCValueObject
	err := collection.Find(bson.M{"servicename": appName,"mainservicename": mainServiceName}).All(&allData)
	if err != nil {
		log.Info("error ", err)
	}
	defer mongoSession.Close()

	return  allData
}
func getRMSTrainingError(appName, appType, mainServiceName  string) []TrainingObject{
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_All_TrainingsError)
	var allData  []TrainingObject
	err := collection.Find(bson.M{"servicename": appName,"mainservicename": mainServiceName}).All(&allData)
	if err != nil {
		log.Info("error ", err)
	}
	defer mongoSession.Close()

	return  allData
}
func getRegressionTRNsMongoAll(appName, appType, mainServiceName  string) MSCValueObject{
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_MSC_REGRESSION)
	var allData  MSCValueObject
	err := collection.Find(bson.M{"servicename": appName,"mainservicename": mainServiceName}).One(&allData)
	if err != nil {
		log.Info("error ", err)
	}
	defer mongoSession.Close()

	return  allData
}