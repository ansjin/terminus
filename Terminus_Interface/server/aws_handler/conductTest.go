package TERMINUS

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
	"time"
	"strconv"
	"strings"
	log "github.com/sirupsen/logrus"
	b64 "encoding/base64"
	"gopkg.in/mgo.v2/bson"
	"crypto/tls"
	"math"
)

func GetVMStartScript(clusterName, containerName, s3BucketName string)string{
	var VMStartScript = "#!bin/sh \n"+
		"echo \"setup\"  \n"+
		"sudo apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual  \n"+
		"sudo apt-get update  \n"+
		"sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common  \n"+
		"#Microservice Run  \n"+
		"git clone https://github.com/ansjin/terminus.git \n"+
		"cd terminus/Terminus_Data_Collection/data \n"+
		"git clone https://github.com/ansjin/apps.git \n"+
		"cd ../../../ \n"+
		"FILE_TERMINUS=\"terminus/Terminus_Data_Collection/.env\"  \n"+
		"/bin/cat <<EOM >$FILE_TERMINUS  \n"+
		"MONGODB_HOST=mongodb  \n"+
		"MONGODB_PORT=27017  \n"+
		"INFLUXDB_USER=root  \n"+
		"INFLUXDB_PASS=root  \n"+
		"INFLUXDB_HOST=influxdb  \n"+
		"INFLUXDB_PORT=8086  \n"+
		"AWS_KEY="+AWSConfig.AwsAccessKeyId+" \n"+
		"AWS_SECRET="+AWSConfig.AwsSecretAccessKey+" \n"+
		"AWS_DEFAULT_REGION="+AWSConfig.Region+"  \n"+
		"AWS_DEFAULT_ZONE="+GetZoneFromRegion(AWSConfig.Region)+" \n"+
		"KUBE_CLUSTER_NAME="+clusterName+" \n"+
		"KOPS_CONTAINER_NAME="+containerName+" \n"+
		"KOPS_S3_BUCKET_NAME="+s3BucketName+" \n"+
		"EOM\n"+
		"cd terminus/Terminus_Data_Collection/scripts  \n"+
		"sudo sh deploy_app.sh  \n"

	encodedString:=b64.StdEncoding.EncodeToString([]byte(VMStartScript))

	return encodedString
}

func startLoadTest(publicAddress string, TobeSentObj LoadTestQueryObj) bool{

	urlTest := "http://"+publicAddress+":8081/deployAndStartLoadTesting"
	log.Info("URL  :>", urlTest)
	log.Info("Values Sent  :>", TobeSentObj)
	all, err := json.Marshal(TobeSentObj)

	if err != nil {
		log.Error(err)
		return false
	}
	req, err := http.NewRequest("GET", urlTest, bytes.NewBuffer(all))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return false
	}
	if resp != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		log.Info("response Status:", resp.Status)
		log.Info("response Headers:", resp.Header)
		log.Info("response Body:", string(body))

		return true
	}
	return false
}

func deleteClusterFunction(publicAddress string){

	urlDeleteCluster := "http://"+publicAddress+":8081/deleteKubeCluster"
	log.Info("URL  :>", urlDeleteCluster)

	req, err := http.NewRequest("GET", urlDeleteCluster, http.NoBody)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		log.Info("response Status:", resp.Status)
		log.Info("response Headers:", resp.Header)
		log.Info("response Body:", string(body))

		return
	}

	return
}

func terminateTestVM(instanceId string) {

	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
	}))

	svc2 := ec2.New(session)

	var allInstances []string

	input2 := ec2.TerminateInstancesInput{InstanceIds: []*string{
		aws.String(instanceId),
	},}
	result2, er2r := svc2.TerminateInstances(&input2)
	if er2r != nil {
		log.Error(er2r)
		return
	}
	for _, instance := range result2.TerminatingInstances {
		allInstances = append(allInstances, ValueAssignString(instance.InstanceId, ""))
	}

	log.Info("Terminate Instances with id: ", allInstances)

}

func startTestVM( clusterName, containerName, s3BucketName, testVMType string)  string {

	sessionAWS := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
	}))

	svc := ec2.New(sessionAWS)
	var allInstancesStarted []Ec2Instances

	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(20),
				},
			},
		},
		ImageId:      aws.String(GetImageFromRegion(AWSConfig.Region)),
		InstanceType: aws.String(testVMType),
		KeyName:      aws.String(AWSConfig.KeyPairName),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String(AWSConfig.SecurityGroup),
		},
		SubnetId: aws.String(AWSConfig.SubnetId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Purpose"),
						Value: aws.String("test"),
					},
				},
			},
		},
		UserData: aws.String(GetVMStartScript(clusterName, containerName,
			s3BucketName)),
	}

	result, err := svc.RunInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Info(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Info(err.Error())
		}
		return ""
	}

	for _, instance := range result.Instances {

		oneInstance := Ec2Instances{InstanceId: ValueAssignString(instance.InstanceId, ""),
			ImageId: ValueAssignString(instance.ImageId, ""),
			InstanceType: ValueAssignString(instance.InstanceType, ""),
			LaunchTime: *instance.LaunchTime,
			InstanceState: ValueAssignString(instance.State.Name, ""),
			AvailabilityZone: ValueAssignString(instance.Placement.AvailabilityZone, ""),
			CoreCount: ValueAssignInt64(instance.CpuOptions.CoreCount, 0),
			PublicIpAddress: ValueAssignString(instance.PublicIpAddress, "")}

		allInstancesStarted = append(allInstancesStarted, oneInstance)
	}
	log.Info(allInstancesStarted)
	return allInstancesStarted[0].InstanceId
}

func getVMPublicIP(startedInstanceId string)  string{
	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
	}))

	svc2 := ec2.New(session)

	var allInstances []Ec2Instances

	input2 := ec2.DescribeInstancesInput{InstanceIds: []*string{
		aws.String(startedInstanceId),
	},}
	result2, er2r := svc2.DescribeInstances(&input2)
	if er2r != nil {
		log.Error(er2r)
		return ""
	}
	for _, reservation := range result2.Reservations {
		for _, instance := range reservation.Instances {

			oneInstance := Ec2Instances{InstanceId: ValueAssignString(instance.InstanceId, ""),
				ImageId: ValueAssignString(instance.ImageId, ""),
				InstanceType: ValueAssignString(instance.InstanceType, ""),
				LaunchTime: *instance.LaunchTime,
				InstanceState: ValueAssignString(instance.State.Name, ""),
				AvailabilityZone: ValueAssignString(instance.Placement.AvailabilityZone, ""),
				CoreCount: ValueAssignInt64(instance.CpuOptions.CoreCount, 0),
				PublicIpAddress: ValueAssignString(instance.PublicIpAddress, "")}

			allInstances = append(allInstances, oneInstance)
		}
	}
	log.Info(allInstances[0].PublicIpAddress)
	return allInstances[0].PublicIpAddress
}
func saveInfluxDbDump(publicAddress, backUpDirectoryName string)  {

	commandmkdir := "mkdir "+backUpDirectoryName
	exe_cmd(commandmkdir)
	hostAddress:=publicAddress+":8088"
	commandGetDump := "influxd backup -portable -host "+hostAddress + " "+backUpDirectoryName
	log.Info(commandGetDump)
	exe_cmd(commandGetDump)
	log.Info("Stored the dump")
}
/*
func restoreInfluxDbDump(backUpDirectoryName,containerName, restoreAddress string)  {

	restoreAddressPort:=restoreAddress+":8088"
	commandRestoreDumpK8s := "influxd restore -portable -host "+restoreAddressPort + " -db k8s -newdb "+containerName+"_k8s"+ " " +
						backUpDirectoryName
	commandRestoreDumpK6 := "influxd restore -portable -host "+restoreAddressPort + " -db TestK6 -newdb "+containerName+"_TestK6"+ " " +
		backUpDirectoryName

	exe_cmd(commandRestoreDumpK8s)
	exe_cmd(commandRestoreDumpK6)
	log.Info("Restored the dump")
}
*/


func launchVMandDeploy(TobeSentObj LoadTestQueryObj, clusterName, containerName, s3BucketName , backUpDirectoryName, testVMType string ){

	log.Info("Starting a test VM ", testVMType, " and deploying data collector")

	startedInstanceId :=startTestVM(clusterName, containerName, s3BucketName, testVMType)
	if( startedInstanceId==""){
		log.Info("Cannot start test VM, terminating test start again latter")
		return
	}
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database).C(Collection_MSC)

	AllData := TestInformation{
		ContainerName: 		containerName,
		StartTimestamp:      time.Now().Unix(),
		InfluxDbK6Name:     containerName+"_TestK6",
		InfluxDbK8sName:    containerName+"_k8s",
		NodeCount:          TobeSentObj.NodeCount,
		NodeInstanceType:   TobeSentObj.InstanceType,
		MasterInstanceType: TobeSentObj.MasterType,
		NumMicroSvs:        "1",
		TypeApp:            TobeSentObj.AppType,
		AppName:            TobeSentObj.AppName,
		Limits:             TobeSentObj.Limits,
		Replicas:           TobeSentObj.IterReplicas,
		TestIter:           TobeSentObj.TestIter,
		TestAPI:            TobeSentObj.TestApi,
		Phase:              "Deployment",
		LimitingInstanceType: TobeSentObj.LimitingResourcesInstanceType,
	}
	if err := collection.Insert(AllData); err != nil {
		log.Info("error ", err)
	} else {
		log.Info("#inserted into ", Collection_MSC)
	}

	time.Sleep(20 * time.Second)
	stopDataCollectionDeployment := Schedule(func() {
		log.Info("waiting for the deployment of the dataCollection to get finish.....")
		getVMPublicIP(startedInstanceId)
	}, 30*time.Second)
	time.Sleep(12 * time.Minute)
	// assuming that it might be finished need to add some check conditions here
	stopDataCollectionDeployment <- true
	publicAddress:= getVMPublicIP(startedInstanceId)
	log.Info("Public Ip Address : ",publicAddress )
	log.Info("Data Collection is deployed")
	log.Info("Starting the test")
	errMongoU := collection.Update(bson.M{"containername": containerName}, bson.M{"$set": bson.M{"testvmip": publicAddress,
		"grafana": "http://" + publicAddress + ":3000", "kibana": "http://" + publicAddress + ":5601", "logs": "http://" + publicAddress + ":8081",
		"phase": "Deployed"}})
	if errMongoU != nil {
		log.Info("Error : %s", errMongoU)
	}


	machineStarted:=startLoadTest(publicAddress, TobeSentObj)

	if machineStarted == false{
		log.Info("Cannot send request to Start Load test, terminating test start again afterwards")

		log.Info(" Deleting Kube Cluster")
		deleteClusterFunction(publicAddress)
		time.Sleep(5 * time.Minute)
		log.Info(" Terminating Test VM")
		terminateTestVM(startedInstanceId)

		errMonFin := collection.Update(bson.M{"containername": containerName}, bson.M{"$set": bson.M{"endtimestamp": time.Now().Unix(),
			"phase": "Errored"}})
		if errMonFin != nil {
			log.Info("Error::%s", errMonFin)
		}
		defer mongoSession.Close()
		return
	}


	timeNow := time.Now().Local();
	var rateIncr = 1
	if TobeSentObj.AppType == SupportedTypesOfMicroservices[1]{

		rateIncr = 6
	}else if TobeSentObj.AppType == SupportedTypesOfMicroservices[2]{

		rateIncr = 50
	} else if TobeSentObj.AppType == SupportedTypesOfMicroservices[0] {
		rateIncr = 4
	} else {
		rateIncr = 20
	}

	// here to sleep for around 3 hours for a test
	numMaxRPS, errRPS := strconv.Atoi(TobeSentObj.MaxRPS)
	if errRPS != nil {
		log.Error("Error : ", errRPS)
		return
	}
	errMonT := collection.Update(bson.M{"containername": containerName}, bson.M{"$set": bson.M{"phase": "Ongoing"}})
	if errMonT != nil {
		log.Info("Error: %s", errMongoU)
	}

	stopDTest:= Schedule(func() {
		log.Info("Completion at: ", timeNow.Add(time.Duration(numMaxRPS/rateIncr) * time.Minute))
	}, 30*time.Second)

	time.Sleep(time.Duration(numMaxRPS/rateIncr) * time.Minute)
	time.Sleep(time.Duration(numMaxRPS/rateIncr) * 15 * time.Second)
	time.Sleep(15 * time.Minute)

	stopDTest <-true
	log.Info("Test Completed .")
	errMonCom := collection.Update(bson.M{"containername": containerName}, bson.M{"$set": bson.M{"phase": "Backingup"}})
	if errMonCom != nil {
		log.Info("Error: %s", errMonCom)
	}

	log.Info(" assuming that it might be finished here we need to import the dump of influxdb")
	saveInfluxDbDump(publicAddress, backUpDirectoryName)
	time.Sleep(16 * time.Minute)

	//restoreInfluxDbDump(backUpDirectoryName,containerName, restoreAddress)

	AllDataDBName := InfluDbTabaseNamesInfo{HostInstanceType: TobeSentObj.InstanceType,
		ServiceType: TobeSentObj.AppType,
		ServiceName: TobeSentObj.AppName,
		Replicas: TobeSentObj.IterReplicas,
		LimitingInstanceType: TobeSentObj.LimitingResourcesInstanceType,
		Limits: getLimits(TobeSentObj.LimitingResourcesInstanceType),
		TestAPI: TobeSentObj.TestApi,
		FolderName: containerName,
		InfluxDbk6Name: containerName + "_TestK6",
		InfluxdbK8sName: containerName + "_k8s",
		NodeCount: TobeSentObj.NodeCount,
		NumMicroSvs: TobeSentObj.NumMicroservice,
		TestIter: TobeSentObj.TestIter,
	}
	mongoSessionDBName := GetMongoSession()
	collectionDBName := mongoSessionDBName.DB(Database).C(Collection_All_Test_Names)
	if errDBName := collectionDBName.Insert(AllDataDBName); errDBName != nil {
		log.Info("error ", errDBName)
	} else {
		log.Info("#inserted into ", Collection_All_Test_Names)
	}
	defer mongoSessionDBName.Close()


	log.Info(" Deleting Kube Cluster")
	deleteClusterFunction(publicAddress)
	time.Sleep(5 * time.Minute)
	log.Info(" Terminating Test VM")
	terminateTestVM(startedInstanceId)

	errMonFin := collection.Update(bson.M{"containername": containerName}, bson.M{"$set": bson.M{"endtimestamp": time.Now().Unix(),
																									"phase": "Completed"}})
	if errMonFin != nil {
		log.Info("Error::%s", errMonFin)
	}

	defer mongoSession.Close()

}
func deployAndTestApp(TobeSentObj LoadTestQueryObj, testVmType string)  {
	containerName:="sv"+TobeSentObj.NumMicroservice+"t"+TobeSentObj.TestIter+"rc"+TobeSentObj.IterReplicas+"nc"+
		TobeSentObj.NodeCount+TobeSentObj.InstanceType+TobeSentObj.AppType+TobeSentObj.AppName+TobeSentObj.Limits+randSeq(4)
	containerName = strings.Replace(containerName, ".", "", -1)
	clusterName := containerName+"cs.k8s.local"
	s3bucketName:=containerName+"-kops"

	backupDirName:="/data/"+containerName+"/"
	log.Info("ContainerName:", containerName)
	log.Info("ClusterName:", clusterName)
	log.Info("S3bucketName:", s3bucketName)
	log.Info("BackupDirName:", backupDirName)
	launchVMandDeploy (TobeSentObj,clusterName,containerName, s3bucketName,  backupDirName, testVmType )
}
func deployAndTestAppLimits(TobeSentObj LoadTestQueryObj, testVmType string)  {
	containerName:="s"+TobeSentObj.NumMicroservice+"t"+TobeSentObj.TestIter+"rc"+TobeSentObj.IterReplicas+"nc"+
		TobeSentObj.NodeCount+TobeSentObj.InstanceType+TobeSentObj.AppType+TobeSentObj.AppName+TobeSentObj.LimitingResourcesInstanceType+randSeq(3)
	containerName = strings.Replace(containerName, ".", "", -1)
	clusterName := containerName+".k8s.local"
	s3bucketName:=containerName+"-kops"

	backupDirName:="/data/"+containerName+"/"
	log.Info("ContainerName:", containerName)
	log.Info("ClusterName:", clusterName)
	log.Info("S3bucketName:", s3bucketName)
	log.Info("BackupDirName:", backupDirName)
	launchVMandDeploy (TobeSentObj,clusterName,containerName, s3bucketName,  backupDirName, testVmType   )
}
func LaunchVMandDeployInterfaceNoLimits(numMicroSvs, appName, appType,nodecount, replicas, instanceType, testAPI,masterType, MaxRPS, testVmType , mainServiceName string){


	for iterTest:=1; iterTest<=MaxTestIterations;iterTest++{
		TobeSentObj:=LoadTestQueryObj{numMicroSvs,replicas, strconv.Itoa(iterTest), nodecount, instanceType,masterType,
			"false", appName, appType, testAPI, "", MaxRPS, mainServiceName}
		log.Info("Params: ",TobeSentObj)
		go deployAndTestApp(TobeSentObj, testVmType)
	}


}
func LaunchVMandDeployInterfaceLimits(numMicroSvs, appName, appType,nodecount, replicas, instanceType, testAPI, limitingInstanceType,masterType,MaxRPS, testVmType, mainServiceName string){

	for iterTest:=1; iterTest<=MaxTestIterations;iterTest++{
		TobeSentObj:=LoadTestQueryObj{numMicroSvs,replicas, strconv.Itoa(iterTest), nodecount, instanceType,masterType,
			"true", appName, appType, testAPI, limitingInstanceType,MaxRPS, mainServiceName}
		log.Info("Params: ",TobeSentObj)
		go deployAndTestAppLimits(TobeSentObj, testVmType)
	}


}

func getPasswordKopsCluster(publicAddress string) string{
	url := "http://"+publicAddress+":8081/getPasswordDashboard"
	log.Info(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Info(err)
		}
		return string(contents)
	}
	return ""
}
func getKopsClusterAddress(publicAddress string)  string{
	url := "http://"+publicAddress+":8081/getClusterInfoKops"
	var clusterAddress string
	log.Info(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Info(err)
		}

		s := strings.Split(string(contents), "\n")
		if(s[0]!=""){
			s_sub := strings.Split(s[0], "Kubernetes master")
			if(len(s_sub) > 1){
				s_sub2 := strings.Split(s_sub[1], "is running at")
				if(len(s_sub2)>1){
					log.Info(s_sub2[1])
					s_sub2 := strings.Split(s_sub2[1], "https://")
					if(len(s_sub2)>1) {
						clusterAddress = s_sub2[1]
					}else{
						clusterAddress = s_sub2[0]
					}
				}
			}
		}
	}
	return clusterAddress
}
func getKubernetesEvents(clusterAddress, passwd string)  KubeEventsObject{

	urlGetEvents := "https://"+string(clusterAddress)+"/api/v1/events"
	urlGetEvents = strings.Replace(urlGetEvents, "\x1b", "", -1)
	urlGetEvents = strings.Replace(urlGetEvents, "[", "", -1)
	urlGetEvents = strings.Replace(urlGetEvents, "0m", "", -1)
	log.Info(urlGetEvents)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}

	passwd = strings.Replace(passwd, "\n", "", -1)
	req, err := http.NewRequest("GET", urlGetEvents, nil)
	req.SetBasicAuth("admin", string(passwd))

	clientEvents := &http.Client{Transport: tr}
	resp, err := clientEvents.Do(req)
	if err != nil {
		log.Error(err)
	}else{
		defer resp.Body.Close()
		events, _ := ioutil.ReadAll(resp.Body)
		var kubeEvents KubeEventsObject
		if err := json.Unmarshal([]byte(events), &kubeEvents); err != nil {
			log.Error(err)
		}
		return kubeEvents
	}
	return KubeEventsObject{}
}
func conductTestToCalculatePodBootTime(numMicroSvs, appName, appType,nodecount, replicas, instanceType, testAPI, limitingInstanceType,masterType,MaxRPS,testVMType, mainservicename string){

	for iterTest:=1; iterTest<=1;iterTest++ {
		TobeSentObj := LoadTestQueryObj{numMicroSvs, replicas, strconv.Itoa(iterTest), nodecount, instanceType, masterType,
			"true", appName, appType, testAPI, limitingInstanceType, MaxRPS,mainservicename }
		log.Info("Params: ", TobeSentObj)
		containerName := "s" + TobeSentObj.NumMicroservice + "t" + TobeSentObj.TestIter + "rc" + TobeSentObj.IterReplicas + "nc" +
			TobeSentObj.NodeCount + TobeSentObj.InstanceType + TobeSentObj.AppType + TobeSentObj.AppName + TobeSentObj.LimitingResourcesInstanceType + randSeq(3)
		containerName = strings.Replace(containerName, ".", "", -1)
		clusterName := containerName + ".k8s.local"
		s3bucketName := containerName + "-kops"

		backupDirName := "/data/" + containerName + "/"
		log.Info("ContainerName:", containerName)
		log.Info("ClusterName:", clusterName)
		log.Info("S3bucketName:", s3bucketName)
		log.Info("BackupDirName:", backupDirName)

		log.Info("Starting a test VM t2.xlarge and deploying data collector")
		startedInstanceId := startTestVM(clusterName, containerName, s3bucketName, testVMType)
		if( startedInstanceId==""){
			log.Info("Cannot start test VM, terminating test start again latter")
			return
		}
		time.Sleep(20 * time.Second)
		stopDataCollectionDeployment := Schedule(func() {
			log.Info("waiting for the deployment of the dataCollection to get finish.....")
			getVMPublicIP(startedInstanceId)
		}, 30*time.Second)
		time.Sleep(10 * time.Minute)
		// assuming that it might be finished need to add some check conditions here
		stopDataCollectionDeployment <- true
		publicAddress := getVMPublicIP(startedInstanceId)
		log.Info("Public Ip Address : ", publicAddress)
		log.Info("Data Collection is deployed")
		log.Info("Starting the test")


		machineStarted:=startLoadTest(publicAddress, TobeSentObj)

		if machineStarted == false{
			log.Info("Cannot send request to Start Load test, terminating test start again afterwards")

			log.Info(" Deleting Kube Cluster")
			deleteClusterFunction(publicAddress)
			time.Sleep(5 * time.Minute)
			log.Info(" Terminating Test VM")
			terminateTestVM(startedInstanceId)
			return
		}

		stopDTest := Schedule(func() {
			log.Info("Deploying the kube cluster: ")
		}, 30*time.Second)

		time.Sleep(15 * time.Minute)

		stopDTest <- true
		kopsClusteraddr := getKopsClusterAddress(publicAddress)
		kopsPasswd := getPasswordKopsCluster(publicAddress)

		kubeEvents := getKubernetesEvents(kopsClusteraddr, kopsPasswd)

		mongoStoreObj := KubeMongoEventsObject{Replicas:replicas, AppName: appName, AppType: appType, KubeEvents: kubeEvents}
		mongoSession := GetMongoSession()
		collectionDBName := mongoSession.DB(Database).C(Collection_Pod_Events)
		if errDBName := collectionDBName.Insert(mongoStoreObj); errDBName != nil {
			log.Info("error ", errDBName)
		} else {
			log.Info("#inserted into ", Collection_Pod_Events)
		}
		defer mongoSession.Close()

		time.Sleep(5 * time.Minute)
		log.Info("Test Completed .")

		log.Info(" Deleting Kube Cluster")
		deleteClusterFunction(publicAddress)
		time.Sleep(5 * time.Minute)
		log.Info(" Terminating Test VM")
		terminateTestVM(startedInstanceId)
	}
}

func analyzeKubernetesEvents(appType, appName, mainserviceName string){

	var kubeEvents  []KubeMongoEventsObject
	for key, _ := range KubeEventsAnalyzed {
		delete(KubeEventsAnalyzed, key)
	}

	mongoSession := GetMongoSession()
	collectionDBName := mongoSession.DB(Database).C(Collection_Pod_Events)
	err := collectionDBName.Find(bson.M{"appname": appName,"apptype": appType}).All(&kubeEvents)
	if err != nil {
		log.Info("%s", err)
	}else{

		for _, event := range kubeEvents{

			items:=event.KubeEvents.Items

			for _,item := range items {

				if (strings.Contains(item.InvolvedObject.Name, mainserviceName) && item.InvolvedObject.Kind=="Pod" ) {

					firstTimeStamp, err := time.Parse(time.RFC3339, item.FirstTimestamp)
					if err != nil {
						log.Info(err)
					}
					secondTimeStamp, err := time.Parse(time.RFC3339, item.LastTimestamp)
					if err != nil {
						log.Info(err)
					}

					mainObj := EventObjectDescribe{event.AppType, event.AppName, event.Replicas,
						mainserviceName, item.InvolvedObject.Name}

					itemObj := EventObjectItems{item.InvolvedObject.Kind,
						item.Reason, item.Message, firstTimeStamp,
						secondTimeStamp}
					KubeEventsAnalyzed[mainObj] = append(KubeEventsAnalyzed[mainObj], itemObj)
				}
			}
			}

		for key, value := range KubeEventsAnalyzed {
			var mongoStoreObj KubeEventsMongoStoreObject
			if(len(value)>1){
				t1 := value[len(value) - 1].LastTimestamp
				t2 := value[0].FirstTimestamp
				mongoStoreObj = KubeEventsMongoStoreObject{key.Replicas, key.AppName, key.AppType,key.MainServiceName, key.Name,
					 t1.Sub(t2)/1000000, value}

				PodBootTime[key.Replicas] = append(PodBootTime[key.Replicas], t1.Sub(t2)/1000000)
				PodBootTimeTests[key.Replicas]++
				mongoSession := GetMongoSession()
				collectionDBName := mongoSession.DB(Database).C(Collection_Pod_Events_Analyzed)

				var mongoFind KubeEventsMongoStoreObject
				err := collectionDBName.Find(bson.M{"appname": key.AppName,"apptype": key.AppType, "mainservicename": key.MainServiceName, "name": key.Name }).One(&mongoFind)
				if err != nil {
					log.Info("%s", err)
					if errDBName := collectionDBName.Insert(mongoStoreObj); errDBName != nil {
						log.Info("error ", errDBName)
					} else {
						log.Info("#inserted into ", Collection_Pod_Events_Analyzed)
					}
				}else{
					log.Info("Analyze Kube Events entry Alraeady Exists ")

				}

				defer mongoSession.Close()
			}
		}
		for key, value := range PodBootTime {

			sum := int64(0)
			sd := float64(0)
			for _,podBootTime := range  value{

				sum+=int64(podBootTime)
			}
			mean:=float64(sum/PodBootTimeTests[key])


			for _,podBootTime := range  value{

				sd += math.Pow(float64(podBootTime) - mean, 2)
			}
			sd = math.Sqrt(sd/float64(PodBootTimeTests[key]))

			mongoSession := GetMongoSession()
			collectionDBName := mongoSession.DB(Database).C(Collection_Pod_Bootingtime)

			var mongoBootTimeFind PodBootingTimeStruct
			err := collectionDBName.Find(bson.M{"appname": appName,"apptype": appType, "mainservicename": mainserviceName, "replicas": StringToFloat(key) }).One(&mongoBootTimeFind)
			if err != nil {
				log.Info("%s", err)
				obj := PodBootingTimeStruct{AppType: appType, AppName: appName, Replicas: StringToFloat(key), NumTests: PodBootTimeTests[key], MainServiceName: mainserviceName,
					MeanBootingTime: mean, SdBootingTime: sd}
				if errDBName := collectionDBName.Insert(obj); errDBName != nil {
					log.Info("error ", errDBName)
				} else {
					log.Info("#inserted into ", Collection_Pod_Bootingtime)
				}
			}else{
				errMongoU := collectionDBName.Update(bson.M{"appname": appName,"apptype": appType, "mainservicename": mainserviceName, "replicas": StringToFloat(key) },
							bson.M{"$set": bson.M{"meanbootingtime": mean,
					"sdbootingtime": sd}})
				if errMongoU != nil {
					log.Info("Error : %s", errMongoU)
				}
			}
			defer mongoSession.Close()
		}

	}
}
func getAnalyzedKubernetesEvents(appName string, mainServiceName string) []KubeEventsMongoStoreObject{

	var kubeEvents []KubeEventsMongoStoreObject

	mongoSession := GetMongoSession()
	collectionDBName := mongoSession.DB(Database).C(Collection_Pod_Events_Analyzed)
	err := collectionDBName.Find(bson.M{"appname": appName,"mainservicename": mainServiceName}).All(&kubeEvents)
	if err != nil {
		log.Info("%s", err)
	}else{

	}
	defer mongoSession.Close()

	return kubeEvents
}

func getPodBootingTime(appName, mainServiceName string) []PodBootingTimeStruct{

	var podBootingtime []PodBootingTimeStruct

	mongoSession := GetMongoSession()
	collectionDBName := mongoSession.DB(Database).C(Collection_Pod_Bootingtime)
	err := collectionDBName.Find(bson.M{"appname": appName,"mainservicename": mainServiceName}).All(&podBootingtime)
	if err != nil {
		log.Info("%s", err)
	}else{

	}
	defer mongoSession.Close()

	return podBootingtime
}
