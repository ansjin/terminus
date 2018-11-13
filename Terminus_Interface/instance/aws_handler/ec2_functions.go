package INSTANCE

import (
	"net/http"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"encoding/json"
	"strconv"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"time"
	"math/rand"
	"strings"
	"sync"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
)
type InitAWSFormData struct {
	Access_key   string `json:"access_key"`
	Secret   string `json:"secret"`
	Region string `json:"region"`
}
// InitUserConfig godoc
// @Summary Initialize User Configuration
// @Description Initialize User Configuration
// @Tags internalUse
// @Accept json
// @Produce json
// @Param body body INSTANCE.InitAWSFormData true "{'access_key': '', 'secret': '', 'region':''}"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /initUserConfig [POST]
func InitUserConfig(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	// In case of any error, we respond with an error to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	UserConfig.AwsAccessKeyId = r.Form.Get("access_key")
	UserConfig.AwsSecretAccessKey = r.Form.Get("secret")
	UserConfig.Region =aws.String(r.Form.Get("region"))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// ListAllInstances godoc
// @Summary List all instances on AWS
// @Description List all instances on AWS
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {array} INSTANCE.Ec2Instances ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /listAllInstances [get]
func ListAllInstances(w http.ResponseWriter, r *http.Request) {
	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(UserConfig.AwsAccessKeyId, UserConfig.AwsSecretAccessKey, ""),
		Region:      UserConfig.Region,
	}))

	svc := ec2.New(session)

	var allInstances []Ec2Instances

	/*input := ec2.DescribeInstancesInput{InstanceIds: []*string{
		aws.String("i-01bdeaedd35d3d55a"),
	},}*/
	input := ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(&input)
	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Unable to perform descibe Instances", http.StatusBadRequest)
		return
	}
	for _, reservation := range result.Reservations {
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

	all, err := json.Marshal(allInstances)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Error converting json", http.StatusBadRequest)
		return
	}
	w.Write(all)
}
func getInstancesPeriodicInformation(collectionName string, InstanceIds []*string, region string) {

	mongoSession := GetMongoSession()

	collection := mongoSession.DB(Database).C(collectionName)

	var timestamp = time.Now().Unix()
	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(UserConfig.AwsAccessKeyId, UserConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(region),
	}))

	svc := ec2.New(session)

	var allInstances []Ec2Instances

	/*input := ec2.DescribeInstancesInput{InstanceIds: []*string{
		aws.String("i-01bdeaedd35d3d55a"),
	},}*/
	input := ec2.DescribeInstancesInput{
		InstanceIds: InstanceIds,
	}
	result, err := svc.DescribeInstances(&input)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	for _, reservation := range result.Reservations {
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
	AllData := Ec2InstancesTime{
		Timestamp: timestamp,
		Instances: allInstances,
	}
	// Insert
	if err := collection.Insert(AllData); err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Println("#inserted into " + collectionName)
	defer mongoSession.Close()
}
func storeData(collectionName string, InstanceIds []*string, region string, wg *sync.WaitGroup) {
	stop := Schedule(func() {
		getInstancesPeriodicInformation(collectionName, InstanceIds, region)
	}, 1*time.Second)

	time.Sleep(3 * time.Minute)

	stop <- true

	fmt.Println("############ Executing Terminate Instances #############")

	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(UserConfig.AwsAccessKeyId, UserConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(region),
	}))

	svc := ec2.New(session)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: InstanceIds,
	}

	result, err := svc.TerminateInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

	stop1 := Schedule(func() {
		getInstancesPeriodicInformation(collectionName, InstanceIds, region)
	}, 1*time.Second)

	time.Sleep(3 * time.Minute)

	stop1 <- true

	defer wg.Done()
}

func executeTestBootShutDown(totalExperiments int, instanceType string, numInstances int64, region string) {
	var wg sync.WaitGroup
	wg.Add(totalExperiments)
	for experimentCount := 1; experimentCount <= totalExperiments; experimentCount++ {
		session := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(UserConfig.AwsAccessKeyId, UserConfig.AwsSecretAccessKey, ""),
			Region:      aws.String(region),
		}))

		svc := ec2.New(session)

		input := &ec2.RunInstancesInput{
			/*BlockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName: aws.String("/dev/sdh"),
					Ebs: &ec2.EbsBlockDevice{
						VolumeSize: aws.Int64(100),
					},
				},
			},*/
			ImageId:      aws.String("ami-5e8bb23b"),
			InstanceType: aws.String(instanceType),
			//KeyName:      aws.String("my-key-pair"),
			MaxCount: aws.Int64(numInstances),
			MinCount: aws.Int64(numInstances),
			//SecurityGroupIds: []*string{
			//	aws.String("sg-1a2b3c4d"),
			//},
			//SubnetId: aws.String("subnet-6e7f829e"),
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
		}
		var InstanceIdsPointer []*string
		result, err := svc.RunInstances(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		for _, instance := range result.Instances {
			InstanceIdsPointer = append(InstanceIdsPointer, aws.String(ValueAssignString(instance.InstanceId, "")))
		}

		collectionName := "num_" + instanceType + "_" + strconv.Itoa(int(numInstances)) + "_" + region + "_" + strconv.Itoa(experimentCount)

		fmt.Println(collectionName)
		mongoSession := GetMongoSession()
		db := mongoSession.DB(Database) // Get db, use db name if not given in connection url

		collection := mongoSession.DB(Database).C(collectionName)
		names, err := db.CollectionNames()
		if err != nil {
			// Handle error
			panic(err)
			return
		}

		// Simply search in the names slice, e.g.
		for _, name := range names {
			if name == collectionName {
				errCollection := collection.DropCollection()
				if errCollection != nil {
					panic(errCollection)
				}
				break
			}
		}
		defer mongoSession.Close()
		go storeData(collectionName, InstanceIdsPointer, region, &wg)
	}
	wg.Wait()
}

func TestInstanceBootingShuttingTime(w http.ResponseWriter, r *http.Request) {

	numInstances, errconv := strconv.Atoi(r.URL.Query().Get("numInstances"))
	if errconv != nil {
		fmt.Println(errconv)
	}

	totalExperiments, errconv2 := strconv.Atoi(r.URL.Query().Get("totalExperiments"))
	if errconv2 != nil {
		fmt.Println(errconv2)
	}

	instanceType := r.URL.Query().Get("instanceType")

	region := r.URL.Query().Get("region")

	if (numInstances > 0) {

		if StringInSlice(instanceType, InstanceTypes) {

			go executeTestBootShutDown(totalExperiments, instanceType, int64(numInstances), region)
			w.Write([]byte("started"))
		} else {
			w.Write([]byte("Please Write in proper format"))
		}

	} else {
		w.Write([]byte("Please Write in proper format"))
	}

}

func randomExecuteTestBootShutDown(totalExperiments int, instanceType string) {

	for loopCount := 1; loopCount <= TotalExperimentsLoop; loopCount++ {
		var wg sync.WaitGroup
		rand.Seed(time.Now().Unix())
		wg.Add(totalExperiments)
		for experimentCount := 1; experimentCount <= totalExperiments; experimentCount++ {
			randNum := rand.Intn(len(DefaultRegion))

			region := DefaultRegion[randNum]
			ami := DefaultAMI[randNum]
			numInstances := DefaultNumInstances[rand.Intn(len(DefaultNumInstances))]
			instanceType := instanceType
			session := session.Must(session.NewSession(&aws.Config{
				Credentials: credentials.NewStaticCredentials(UserConfig.AwsAccessKeyId, UserConfig.AwsSecretAccessKey, ""),
				Region:      aws.String(region),
			}))

			svc := ec2.New(session)

			input := &ec2.RunInstancesInput{
				ImageId:      aws.String(ami),
				InstanceType: aws.String(instanceType),
				//KeyName:      aws.String("my-key-pair"),
				MaxCount: aws.Int64(numInstances),
				MinCount: aws.Int64(numInstances),
				//SecurityGroupIds: []*string{
				//	aws.String("sg-1a2b3c4d"),
				//},
				//SubnetId: aws.String("subnet-6e7f829e"),
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
			}
			var InstanceIdsPointer []*string
			result, err := svc.RunInstances(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					default:
						fmt.Println(aerr.Error())
					}
				} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
				}
				return
			}
			for _, instance := range result.Instances {
				InstanceIdsPointer = append(InstanceIdsPointer, aws.String(ValueAssignString(instance.InstanceId, "")))
			}

			collectionName := "num_" + instanceType + "_" + strconv.Itoa(int(numInstances)) + "_" + region + "_" + strconv.Itoa(experimentCount) + "_" + strconv.Itoa(loopCount)

			fmt.Println(collectionName)


			mongoSession := GetMongoSession()
			db := mongoSession.DB(Database) // Get db, use db name if not given in connection url
			collection := mongoSession.DB(Database).C(collectionName)
			names, err := db.CollectionNames()
			if err != nil {
				// Handle error
				fmt.Printf("%s", err)
				return
			}
			// Simply search in the names slice, e.g.
			for _, name := range names {
				if name == collectionName {
					errCollection := collection.DropCollection()
					if errCollection != nil {
						fmt.Printf("%s", errCollection)
						return
					}
					break
				}
			}
			defer mongoSession.Close()
			go storeData(collectionName, InstanceIdsPointer, region, &wg)
		}
		wg.Wait()
	}
}
// RandomTestInstanceBootingShuttingTime godoc
// @Summary Starts the RMIT experiment
// @Description Starts the RMIT experiment
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Param instanceType query string true "instance type"
// @Param totalExperiments query string true "totalExperiments (batch) to perform each ecperiment will be performed 5 times"
// @Success 200 {string} string "started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /randomTestInstanceBootingShuttingTime [get]

func RandomTestInstanceBootingShuttingTime(w http.ResponseWriter, r *http.Request) {

	instanceType := r.URL.Query().Get("instanceType")
	totalExperimentsStr := r.URL.Query().Get("totalExperiments")

	if(instanceType==""|| totalExperimentsStr == ""){
		http.Error(w, "Query params not specified properly", http.StatusBadRequest)
		return
	}else{
		totalExperiments, errconv := strconv.Atoi(totalExperimentsStr)
		if errconv != nil {
			fmt.Println(errconv)
			http.Error(w, "Error converting string to integer", http.StatusBadRequest)
			return
		}
		if StringInSlice(instanceType, AllInstanceTypes) {
			go randomExecuteTestBootShutDown(totalExperiments, instanceType)
			w.Write([]byte("started"))
		}else {
			http.Error(w, "Instance type wrongly specified", http.StatusBadRequest)
			return
		}
	}
}

func formatDataforVMBootShutRate(w http.ResponseWriter) []VmTemplateData{
	mongoSession := GetMongoSession()
	db := mongoSession.DB(Database) // Get db, use db name if not given in connection url
	collectionNames, err := db.CollectionNames()
	defer mongoSession.Close()
	if err != nil {
		// Handle error
		fmt.Printf("%s", err)
		http.Error(w, "db connection error", http.StatusBadRequest)
		return nil
	}
	var AllVMTypesData []PerExperiment

	var AllProfiles []VmTemplateData

	for iter_collection := 0; iter_collection < len(collectionNames); iter_collection++ {
		// Find All
		substrParts := strings.Split(string(collectionNames[iter_collection]), "_")

		if len(substrParts) > 2 {
			expLoop, errconv := strconv.Atoi(substrParts[5])
			if errconv != nil {
				fmt.Printf("%s", errconv)
				http.Error(w, "string to int conversion error", http.StatusBadRequest)
				return nil
			}

			numInstancesS, errconv1 := strconv.Atoi(substrParts[2])
			if errconv1 != nil {
				fmt.Printf("%s", errconv1)
				http.Error(w, "string to int conversion error", http.StatusBadRequest)
				return nil
			}

			var allData []Ec2InstancesTime
			mongoSession := GetMongoSession()
			collection := mongoSession.DB(Database).C(collectionNames[iter_collection])
			err = collection.Find(nil).All(&allData)
			if err != nil {
				fmt.Printf("%s", err)
				http.Error(w, "Db connection error", http.StatusBadRequest)
				return nil
			}
			defer mongoSession.Close()
			var instances []InstancesData
			for instance := 0; instance < len(allData[0].Instances); instance++ {
				instances = append(instances, InstancesData{Pending: 0, Running: 0, SshLogin: 0, ShuttingDown: 0, Stopped: 0})
			}

			for iter := 0; iter < len(allData); iter++ {

				for instance := 0; instance < len(allData[iter].Instances); instance++ {

					switch allData[iter].Instances[instance].InstanceState {
					case "pending":
						if instances[instance].Pending == 0 {
							instances[instance].Pending = allData[iter].Timestamp
						}
					case "running":
						if instances[instance].Running == 0 {
							instances[instance].Running = allData[iter].Timestamp
						}
					case "app_running":
						if instances[instance].SshLogin == 0 {
							instances[instance].SshLogin = allData[iter].Timestamp
						}
					case "shutting-down":
						if instances[instance].ShuttingDown == 0 {
							instances[instance].ShuttingDown = allData[iter].Timestamp
						}
					case "terminated":
						if instances[instance].Terminated == 0 {
							instances[instance].Terminated = allData[iter].Timestamp
						}
					default:
						if instances[instance].Other == 0 {
							instances[instance].Other = allData[iter].Timestamp
						}
					}
				}
			}

			var maxStartTime, maxShutDownTime int64
			var minStartTime = int64(99999999)
			var minShutDownTime = int64(99999999)
			var totalStartTimeForAllInstances int64
			var totalShutDownTimeForAllInstances int64
			var leastPendingStateForInstance = int64(9999999999999)
			var maxRunningStateForInstance int64

			var leastShuttingDownStateForInstance = int64(9999999999999)
			var maxTerminatedStateForInstance int64

			for iter := 0; iter < len(instances); iter++ {

				if instances[iter].SshLogin == 0 {
					totalStartTimeForAllInstances += instances[iter].Running - instances[iter].Pending

					if maxStartTime < instances[iter].Running-instances[iter].Pending {
						maxStartTime = instances[iter].Running - instances[iter].Pending
					}

					if minStartTime > instances[iter].Running-instances[iter].Pending {
						minStartTime = instances[iter].Running - instances[iter].Pending
					}

					if (maxRunningStateForInstance < instances[iter].Running) {
						maxRunningStateForInstance = instances[iter].Running
					}

				} else {
					totalStartTimeForAllInstances += instances[iter].SshLogin - instances[iter].Pending

					if maxStartTime < instances[iter].SshLogin-instances[iter].Pending {
						maxStartTime = instances[iter].SshLogin - instances[iter].Pending
					}

					if minStartTime > instances[iter].SshLogin-instances[iter].Pending {
						minStartTime = instances[iter].SshLogin - instances[iter].Pending
					}

					if (maxRunningStateForInstance < instances[iter].SshLogin) {
						maxRunningStateForInstance = instances[iter].SshLogin
					}

				}
				totalShutDownTimeForAllInstances += instances[iter].Terminated - instances[iter].ShuttingDown
				if maxShutDownTime < instances[iter].Terminated-instances[iter].ShuttingDown {
					maxShutDownTime = instances[iter].Terminated - instances[iter].ShuttingDown
				}

				if minShutDownTime > instances[iter].Terminated-instances[iter].ShuttingDown {
					minShutDownTime = instances[iter].Terminated - instances[iter].ShuttingDown
				}
				if (leastPendingStateForInstance > instances[iter].Pending) {
					leastPendingStateForInstance = instances[iter].Pending
				}

				if (leastShuttingDownStateForInstance > instances[iter].ShuttingDown) {
					leastShuttingDownStateForInstance = instances[iter].ShuttingDown
				}

				if (maxTerminatedStateForInstance < instances[iter].Terminated) {
					maxTerminatedStateForInstance = instances[iter].Terminated
				}
			}
			var totalStartTime = maxRunningStateForInstance - leastPendingStateForInstance
			var totalShutdownTime = maxTerminatedStateForInstance - leastShuttingDownStateForInstance

			VMData := PerExperiment{InstanceType: substrParts[1], NumInstances: numInstancesS, Instances: instances, InstanceBootRate: InstanceTime{
				Avg: float64(totalStartTimeForAllInstances / int64(len(instances))),
				Max: maxStartTime,
				Min: minStartTime,
			}, ShutDownRate: InstanceTime{
				Avg: float64(totalShutDownTimeForAllInstances / int64(len(instances))),
				Max: maxShutDownTime,
				Min: minShutDownTime,
			},
				Region: substrParts[3],
				ExperimentNum: substrParts[4],
				ExperimenLoopCount: expLoop,
				AvailabilityZone: allData[0].Instances[0].AvailabilityZone,
				CoreCount: allData[0].Instances[0].CoreCount,
				ImageId: allData[0].Instances[0].ImageId,
				TotalStartTime: totalStartTime,
				TotalShutDownTime: totalShutdownTime,
			}
			AllVMTypesData = append(AllVMTypesData, VMData)
		}
	}

	for iter := 0; iter < len(AllVMTypesData); iter++ {

		found := 0;
		for iter2 := 0; iter2 < len(AllProfiles); iter2++ {

			if (AllProfiles[iter2].CoreCount == AllVMTypesData[iter].CoreCount) &&
				(AllProfiles[iter2].InstanceType == AllVMTypesData[iter].InstanceType) &&
				(AllProfiles[iter2].Region == AllVMTypesData[iter].Region) &&
				(AllProfiles[iter2].ImageId == AllVMTypesData[iter].ImageId) {

				AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].NumInstances = AllVMTypesData[iter].NumInstances
				AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].NumExperiments++
				AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].BootTime.Avg += float64(AllVMTypesData[iter].TotalStartTime)
				AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].ShutdownTime.Avg += float64(AllVMTypesData[iter].TotalShutDownTime)

				if (AllVMTypesData[iter].TotalStartTime < AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].BootTime.Min) {

					AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].BootTime.Min = AllVMTypesData[iter].TotalStartTime
				}

				if (AllVMTypesData[iter].TotalShutDownTime < AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].ShutdownTime.Min) {

					AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].ShutdownTime.Min = AllVMTypesData[iter].TotalShutDownTime
				}
				if (AllVMTypesData[iter].TotalStartTime > AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].BootTime.Max) {

					AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].BootTime.Max = AllVMTypesData[iter].TotalStartTime
				}

				if (AllVMTypesData[iter].TotalShutDownTime > AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].ShutdownTime.Max) {

					AllProfiles[iter2].BootShutdownRate[AllVMTypesData[iter].NumInstances-1].ShutdownTime.Max = AllVMTypesData[iter].TotalShutDownTime
				}

				AllProfiles[iter2].ExperimentLoop[AllVMTypesData[iter].ExperimenLoopCount-1].Experiments = append(AllProfiles[iter2].ExperimentLoop[AllVMTypesData[iter].ExperimenLoopCount-1].Experiments, ExperimentSetting{AllVMTypesData[iter].ExperimentNum, AllVMTypesData[iter].NumInstances,
					AllVMTypesData[iter].Instances, AllVMTypesData[iter].InstanceBootRate,
					AllVMTypesData[iter].ShutDownRate, AllVMTypesData[iter].TotalStartTime,
					AllVMTypesData[iter].TotalShutDownTime,})
				found = 1
				break;
			}
		}
		if (found == 0) {
			tempVMDTypeData := VmTemplateData{
				InstanceType:     AllVMTypesData[iter].InstanceType,
				ImageId:          AllVMTypesData[iter].ImageId,
				Region:           AllVMTypesData[iter].Region,
				AvailabilityZone: AllVMTypesData[iter].AvailabilityZone,
				CoreCount:        AllVMTypesData[iter].CoreCount,
			}

			for test := 0; test < TotalExperimentsLoop; test++ {

				var tempExp2 ExperimentsLoop
				tempVMDTypeData.ExperimentLoop = append(tempVMDTypeData.ExperimentLoop, tempExp2)
			}

			for test := 0; test < int(DefaultNumInstances[len(DefaultNumInstances)-1]); test++ {

				var tempExp2 = InstanceBootShutdownRate{NumExperiments: 0, NumInstances: 0,
					BootTime: InstanceTime{Avg: float64(0), Min: 999999999, Max: 0},
					ShutdownTime: InstanceTime{Avg: float64(0), Min: 999999999, Max: 0},}

				tempVMDTypeData.BootShutdownRate = append(tempVMDTypeData.BootShutdownRate, tempExp2)
			}

			tempVMDTypeData.BootShutdownRate[AllVMTypesData[iter].NumInstances-1] = InstanceBootShutdownRate{NumExperiments: 1, NumInstances: AllVMTypesData[iter].NumInstances,
				BootTime: InstanceTime{Avg: float64(AllVMTypesData[iter].TotalStartTime), Min: 999999999, Max: AllVMTypesData[iter].TotalStartTime},
				ShutdownTime: InstanceTime{Avg: float64(AllVMTypesData[iter].TotalShutDownTime), Min: 999999999, Max: AllVMTypesData[iter].TotalShutDownTime},}

			var tempExp ExperimentsLoop
			tempExp.Experiments = append(tempExp.Experiments, ExperimentSetting{AllVMTypesData[iter].ExperimentNum, AllVMTypesData[iter].NumInstances,
				AllVMTypesData[iter].Instances, AllVMTypesData[iter].InstanceBootRate,
				AllVMTypesData[iter].ShutDownRate, AllVMTypesData[iter].TotalStartTime,
				AllVMTypesData[iter].TotalShutDownTime,}, )

			tempVMDTypeData.ExperimentLoop[AllVMTypesData[iter].ExperimenLoopCount-1] = tempExp
			AllProfiles = append(AllProfiles, tempVMDTypeData)
		}

	}

	for iter2 := 0; iter2 < len(AllProfiles); iter2++ {

		for test := 0; test < len(AllProfiles[iter2].BootShutdownRate); test++ {
			if AllProfiles[iter2].BootShutdownRate[test].NumExperiments != 0 {
				AllProfiles[iter2].BootShutdownRate[test].BootTime.Avg = AllProfiles[iter2].BootShutdownRate[test].BootTime.Avg / float64(AllProfiles[iter2].BootShutdownRate[test].NumExperiments)
				AllProfiles[iter2].BootShutdownRate[test].ShutdownTime.Avg = AllProfiles[iter2].BootShutdownRate[test].ShutdownTime.Avg / float64(AllProfiles[iter2].BootShutdownRate[test].NumExperiments)
			}
		}
	}
	return AllProfiles
}
// GetAllVMTypesBootShutDownDataAvg godoc
// @Summary Gives all the Booting and shutting down data following avg approach
// @Description Gives all the Booting and shutting down data following avg approach
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Param takeNewValues query string false "specify it yes if to take walues from the db and update the resulting average"
// @Success 200 {array} INSTANCE.VmTemplateData ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAllVMTypesBootShutDownDataAvg [get]
func GetAllVMTypesBootShutDownDataAvg(w http.ResponseWriter, r *http.Request) {

	takeNewValues := r.URL.Query().Get("takeNewValues")


	mongoSession := GetMongoSession()
	db := mongoSession.DB(Database_Final) // Get db, use db name if not given in connection url
	collectionNames, err := db.CollectionNames()
	defer mongoSession.Close()
	if err != nil {
		// Handle error
		fmt.Printf("%s", err)
		http.Error(w, "db connection error", http.StatusBadRequest)
		return
	}
	var AllProfiles []VmTemplateData
	var found = 0

	if(takeNewValues == "yes"){
		AllProfiles=formatDataforVMBootShutRate(w)
		mongoSession := GetMongoSession()
		collection := mongoSession.DB(Database_Final).C("data")
		// Insert
		for i:=0;i<len(AllProfiles); i++{
			if err := collection.Insert(AllProfiles[i]); err != nil {
				fmt.Printf("%s", err)
				http.Error(w, "Db insertion error", http.StatusBadRequest)
				return
			}
			fmt.Println("#inserted into " + "data")
		}

		defer mongoSession.Close()
	}else{
		for iter_collection := 0; iter_collection < len(collectionNames); iter_collection++ {

			if string(collectionNames[iter_collection]) == "data" {
				found =1
				mongoSession := GetMongoSession()
				collection := mongoSession.DB(Database_Final).C(collectionNames[iter_collection])
				err = collection.Find(nil).All(&AllProfiles)
				if err != nil {
					fmt.Printf("%s", err)
					http.Error(w, "db query error", http.StatusBadRequest)
					return
				}
				defer mongoSession.Close()
				break;
			}
		}
		if found==0{
			AllProfiles=formatDataforVMBootShutRate(w)
			mongoSession := GetMongoSession()
			collection := mongoSession.DB(Database_Final).C("data")
			// Insert
			for i:=0;i<len(AllProfiles); i++{
				if err := collection.Insert(AllProfiles[i]); err != nil {
					fmt.Printf("%s", err)
					http.Error(w, "Db insertion error", http.StatusBadRequest)
					return
				}
				fmt.Println("#inserted into " + "data")
			}

			defer mongoSession.Close()
		}
	}

	all, err := json.Marshal(AllProfiles)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	}
	w.Write(all)
}


func GetandStorePredictedValuesInternal(w http.ResponseWriter){
	var csp="aws"
	mongoSession := GetMongoSession()
	collection := mongoSession.DB(Database_predicted_boot).C("data")
	db := mongoSession.DB(Database_predicted_boot) // Get db, use db name if not given in connection url
	names, err := db.CollectionNames()
	if err != nil {
		// Handle error
		fmt.Printf("%s", err)
		http.Error(w, "Db connection error", http.StatusBadRequest)
		return
	}
	// Simply search in the names slice, e.g.
	for _, name := range names {
		if name == "data" {
			errCollection := collection.DropCollection()
			if errCollection != nil {
				fmt.Printf("%s", errCollection)
				http.Error(w, "Db connection error", http.StatusBadRequest)
				return
			}
			break
		}
	}
	defer mongoSession.Close()

	for iterInst:=0; iterInst<int(len(AllInstanceTypes));iterInst++{

		var InstanceAllValues []InstanceValue
		InstanceType:=AllInstanceTypes[iterInst];
		for iter:=1; iter<=int(DefaultNumInstances[len(DefaultNumInstances) - 1]); iter++{

			url :="http://regression_vm_boot:9001/getPrediction?csp="+csp+"&instanceType="+InstanceType+"&numInstances="+strconv.Itoa(iter)+"&type=boot"
			fmt.Println(url)
			response, err := http.Get(url)
			if err != nil {
				fmt.Printf("%s", err)
				http.Error(w, "cannot query to regression_vm_boot service", http.StatusBadRequest)
				return
			} else {
				defer response.Body.Close()
				contentsboot, err := ioutil.ReadAll(response.Body)
				if err != nil {
					http.Error(w, "contents returned from regression_vm_boot are not readable", http.StatusBadRequest)
					return
				}
				fmt.Printf("%s\n", string(contentsboot))

				numInstances:= iter
				url :="http://regression_vm_boot:9001/getPrediction?csp="+csp+"&instanceType="+InstanceType+"&numInstances="+strconv.Itoa(iter)+"&type=shut"
				fmt.Println(url)
				response, err := http.Get(url)
				if err != nil {
					fmt.Printf("%s", err)
					http.Error(w, "cannot query to regression_vm_boot service", http.StatusBadRequest)
					return
				} else {
					defer response.Body.Close()
					contentsshut, err := ioutil.ReadAll(response.Body)
					if err != nil {
						fmt.Printf("%s", err)
						http.Error(w, "contents returned from regression_vm_boot are not readable", http.StatusBadRequest)
						return
					}
					fmt.Printf("%s\n", string(contentsshut))

					tempInstanceValue:=InstanceValue{numInstances, StringToFloat(string(contentsboot)), StringToFloat(string(contentsshut))}

					InstanceAllValues = append(InstanceAllValues, tempInstanceValue)
				}
			}
		}
		tempInstObj:=InstanceRegression{InstanceType:InstanceType,InstanceValues: InstanceAllValues, Region: DefaultRegion[0]}
		collection := mongoSession.DB(Database_predicted_boot).C("data")
		// Insert
		if err := collection.Insert(tempInstObj); err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "Db connection error", http.StatusBadRequest)
			return
		}
		fmt.Println("#inserted into " + "data")
	}

	defer mongoSession.Close()
}
// TrainDataSetRegression godoc
// @Summary Train the dataset
// @Description train the dataset
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /trainDataSetRegression [get]
func TrainDataSetRegression(w http.ResponseWriter, r *http.Request){

	url :="http://regression_vm_boot:9001/train"
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "cannot query to regression_vm_boot service", http.StatusBadRequest)
		return
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "contents returned from regression_vm_boot are not readable", http.StatusBadRequest)
			return
		}
		fmt.Printf("%s\n", string(contents))
		w.Write([]byte(string(contents)));
	}
}

// GetandStoreRegressionValues godoc
// @Summary Query python regression_vm_boot to get the values for each instance type and then store in db
// @Description Query python regression_vm_boot to get the values for each instance type and then store in db
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Done"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getandStoreRegressionValues [get]
func GetandStoreRegressionValues(w http.ResponseWriter, r *http.Request){

	GetandStorePredictedValuesInternal(w)
	w.Write([]byte("Done"));

}

// GetAllVMTypesBootShutDownDataRegression godoc
// @Summary Gives all the Booting and shutting down data
// @Description Gives all the Booting and shutting down data
// @Tags internalUse
// @Accept text/html
// @Produce json
// @Success 200 {array} INSTANCE.InstanceRegression ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAllVMTypesBootShutDownDataRegression [get]
func GetAllVMTypesBootShutDownDataRegression(w http.ResponseWriter, r *http.Request) {
	mongoSession := GetMongoSession()
	db := mongoSession.DB(Database_predicted_boot) // Get db, use db name if not given in connection url
	collectionNames, err := db.CollectionNames()
	defer mongoSession.Close()
	if err != nil {
		// Handle error
		fmt.Printf("%s", err)
		http.Error(w, "Db connection error", http.StatusBadRequest)
		return
	}
	var AllData []InstanceRegression
	var found = 0
	for iter_collection := 0; iter_collection < len(collectionNames); iter_collection++ {

		if string(collectionNames[iter_collection]) == "data" {
			found =1
			mongoSession := GetMongoSession()
			collection := mongoSession.DB(Database_predicted_boot).C(collectionNames[iter_collection])
			err = collection.Find(nil).All(&AllData)
			if err != nil {
				fmt.Printf("%s", err)
				http.Error(w, "Db error", http.StatusBadRequest)
				return
			}
			defer mongoSession.Close()
			break;
		}
	}

	if found==0{
		GetandStorePredictedValuesInternal(w)
		mongoSession := GetMongoSession()
		collection := mongoSession.DB(Database_predicted_boot).C("data")
		err = collection.Find(nil).All(&AllData)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "Db error", http.StatusBadRequest)
			return
		}
		defer mongoSession.Close()
	}

	all, err := json.Marshal(AllData)

	if err != nil {
		fmt.Printf("%s", err)
		http.Error(w, "JSON conversion error", http.StatusBadRequest)
		return
	}else {
		w.Write(all)
	}

}
// GetPerVMTypeAllBootShutDownData godoc
// @Summary Gives all the Booting and shutting down data for a VM type
// @Description Gives all the Booting and shutting down data
// @Tags internal_External_Use
// @Accept text/html
// @Produce json
// @Param instanceType query string true "instance type"
// @Param region query string true "aws region "
// @Param appraoch query string true "appraoch avg or regression_vm_boot, by default it is avg"
// @Param csp query string true "cloud service provider..current;y it is aws only"
// @Success 200 {object} INSTANCE.VMBootShutDownRatePerInstanceTypeAll ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPerVMTypeAllBootShutDownData [get]
func GetPerVMTypeAllBootShutDownData(w http.ResponseWriter, r *http.Request) {

	/*numInstances, errconv := strconv.Atoi(r.URL.Query().Get("numInstances"))
	if errconv != nil {
		fmt.Println(errconv)
	}
	csp, errconv := strconv.Atoi(r.URL.Query().Get("csp"))
	if errconv != nil {
		fmt.Println(errconv)
	}*/

	instanceType := r.URL.Query().Get("instanceType")

	region := r.URL.Query().Get("region")

	approach := r.URL.Query().Get("approach")

	if (instanceType == "" && region =="" ){
		http.Error(w, "check input format", http.StatusBadRequest)
		return
	}
	 if(approach =="regression") {

		mongoSession := GetMongoSession()

		var foundData InstanceRegression

		collection := mongoSession.DB(Database_predicted_boot).C("data")
		err :=  collection.Find(bson.M{"instancetype": instanceType, "region": region}).One(&foundData)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "Db Error", http.StatusBadRequest)
			return
		}
		defer mongoSession.Close()

		 tempObj:= VMBootShutDownRatePerInstanceTypeAll{InstanceValues:foundData.InstanceValues}
		all, err := json.Marshal(tempObj)

		if err != nil {
			http.Error(w, "JSON conversion error", http.StatusBadRequest)
			return
		}else {
			w.Write(all)
		}

	}else {

		 mongoSession := GetMongoSession()

		 var foundData VmTemplateData

		 collection := mongoSession.DB(Database_Final).C("data")
		 err :=  collection.Find(bson.M{"instancetype": instanceType, "region": region}).One(&foundData)
		 defer mongoSession.Close()
		 if err != nil {
			 fmt.Printf("%s", err)
			 http.Error(w, "Db Error", http.StatusBadRequest)
			 return
		 } else {
				 var dataTObeReturn InstanceRegression

				 dataTObeReturn.InstanceType = instanceType
				 dataTObeReturn.Region = region
				 for iter:=0;iter<len(foundData.BootShutdownRate); iter++{

					 tempObj:=InstanceValue{foundData.BootShutdownRate[iter].NumInstances, foundData.BootShutdownRate[iter].BootTime.Avg,
						 foundData.BootShutdownRate[iter].ShutdownTime.Avg}

					 dataTObeReturn.InstanceValues  = append(dataTObeReturn.InstanceValues, tempObj)
				 }
			 		tempObj:= VMBootShutDownRatePerInstanceTypeAll{InstanceValues:dataTObeReturn.InstanceValues}
				 all, err := json.Marshal(tempObj)

				 if err != nil {
					 fmt.Printf("%s", err)
					 http.Error(w, "JSON conversion error", http.StatusBadRequest)
					 return
				 }else {
					 w.Write(all)
				 }

			 }
		 }
 }
// GetPerVMTypeOneBootShutDownData godoc
// @Summary Gives all the Booting and shutting down data for a VM type
// @Description Gives all the Booting and shutting down data
// @Tags internal_External_Use
// @Accept text/html
// @Produce json
// @Param numInstances query string true "number of instances"
// @Param instanceType query string true "instance type"
// @Param region query string true "aws region "
// @Param approach query string true "approach avg or regression, by default it is avg"
// @Param csp query string true "cloud service provider..current;y it is aws only"
// @Success 200 {object} INSTANCE.VMBootShutDownRatePerInstanceTypeOne ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPerVMTypeOneBootShutDownData [get]
func GetPerVMTypeOneBootShutDownData(w http.ResponseWriter, r *http.Request) {

	numInstances, errconv := strconv.Atoi(r.URL.Query().Get("numInstances"))
	if errconv != nil {
		fmt.Println(errconv)
		http.Error(w, "Number of instances not provided", http.StatusBadRequest)
		return
	}
	/*csp, errconv := strconv.Atoi(r.URL.Query().Get("csp"))
	if errconv != nil {
		fmt.Println(errconv)
	}*/

	instanceType := r.URL.Query().Get("instanceType")

	region := r.URL.Query().Get("region")

	approach := r.URL.Query().Get("approach")

	if (instanceType == "" && region =="" ){
		http.Error(w, "check input format", http.StatusBadRequest)
		return
	}
	if(approach =="regression") {

		var csp="aws"

		mongoSession := GetMongoSession()

		var foundData InstanceRegression

		collection := mongoSession.DB(Database_predicted_boot).C("data")
		err :=  collection.Find(bson.M{"instancetype": instanceType, "region": region}).One(&foundData)
		if err != nil {
			http.Error(w, "Db Error", http.StatusBadRequest)
			return
		}
		defer mongoSession.Close()

		url :="http://regression_vm_boot:9001/getPrediction?csp="+csp+"&instanceType="+instanceType+"&numInstances="+strconv.Itoa(numInstances)+"&type=boot"
		fmt.Println(url)
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "cannot query to regression_vm_boot service", http.StatusBadRequest)
			return
		} else {
			defer response.Body.Close()
			contentsboot, err := ioutil.ReadAll(response.Body)
			if err != nil {
				http.Error(w, "contents returned from regression_vm_boot are not readable", http.StatusBadRequest)
				return
			}
			fmt.Printf("%s\n", string(contentsboot))

			url := "http://regression_vm_boot:9001/getPrediction?csp=" + csp + "&instanceType=" + instanceType + "&numInstances=" + strconv.Itoa(numInstances) + "&type=shut"
			fmt.Println(url)
			response, err := http.Get(url)
			if err != nil {
				fmt.Printf("%s", err)
				http.Error(w, "cannot query to regression_vm_boot service", http.StatusBadRequest)
				return
			} else {
				defer response.Body.Close()
				contentsshut, err := ioutil.ReadAll(response.Body)
				if err != nil {
					fmt.Printf("%s", err)
					http.Error(w, "contents returned from regression_vm_boot are not readable", http.StatusBadRequest)
					return
				}
				fmt.Printf("%s\n", string(contentsshut))

				tempObj:=VMBootShutDownRatePerInstanceTypeOne{ StringToFloat(string(contentsboot)),
					StringToFloat(string(contentsshut))}
				all, err := json.Marshal(tempObj)
				if err != nil {
					http.Error(w, "JSON conversion error", http.StatusBadRequest)
					return
				}
				w.Write(all)

			}
		}
	}else {

		mongoSession := GetMongoSession()

		var foundData VmTemplateData

		collection := mongoSession.DB(Database_Final).C("data")
		err :=  collection.Find(bson.M{"instancetype": instanceType, "region": region}).One(&foundData)
		defer mongoSession.Close()
		if err != nil {
			fmt.Printf("%s", err)
			http.Error(w, "Db Error", http.StatusBadRequest)
			return
		} else {
			if(numInstances <=len(foundData.BootShutdownRate)){

				tempObj:=VMBootShutDownRatePerInstanceTypeOne{ foundData.BootShutdownRate[numInstances -1].BootTime.Avg,
					foundData.BootShutdownRate[numInstances -1].ShutdownTime.Avg}
				all, err := json.Marshal(tempObj)

				if err != nil {
					http.Error(w, "JSON conversion error", http.StatusBadRequest)
					return
				}
				w.Write(all)
			}else {
				http.Error(w, "Number of Instances specified more than what data has", http.StatusBadRequest)
				return
			}

		}
	}
}