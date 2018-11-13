package TERMINUS

import (
	"net/http"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
)

// here need to remove this once INSTANCE is converted into package
type InstancesData struct {
	Pending      int64 `json:"Pending"`
	Running      int64 `json:"Running"`
	SshLogin     int64 `json:"SshLogin"`
	ShuttingDown int64 `json:"ShuttingDown"`
	Stopped      int64 `json:"Stopped"`
	Terminated   int64 `json:"Terminated"`
	Other        int64 `json:"Other"`
}

type InstanceBootTime struct {
	Avg float64 `json:"Avg"`
	Min int64 `json:"Min"`
	Max int64 `json:"Max"`
}

type InstanceTime struct {
	Avg float64 `json:"Avg"`
	Min int64 `json:"Min"`
	Max int64 `json:"Max"`
}
type InstanceShutDownTime struct {
	Avg float64 `json:"Avg"`
	Min int64 `json:"Min"`
	Max int64 `json:"Max"`
}

type InstanceBootShutdownRate struct {
	NumInstances   int `json:"NumInstances"`
	NumExperiments int `json:"NumExperiments"`
	BootTime       InstanceTime `json:"BootTime"`
	ShutdownTime   InstanceTime `json:"ShutdownTime"`
}

type ExperimentSetting struct {
	ExperimentNum              string `json:"ExperimentNum"`
	NumInstances               int `json:"NumInstances"`
	Instances                  []InstancesData `json:"Instances"`
	InstancesBootTime          InstanceTime `json:"InstancesBootTime"`
	InstancesShutDownTime      InstanceTime `json:"InstancesShutDownTime"`
	TotalInstancesBootTime     int64 `json:"TotalInstancesBootTime"`
	TotalInstancesShutDownTime int64 `json:"TotalInstancesShutDownTime"`
}
type ExperimentsLoop struct {
	Experiments []ExperimentSetting `json:"Experiments"`
}

type VmTemplateData struct {
	InstanceType     string `json:"InstanceType"`
	Region           string `json:"Region"`
	AvailabilityZone string `json:"AvailabilityZone"`
	ImageId          string `json:"ImageId"`
	CoreCount        int64 `json:"CoreCount"`
	ExperimentLoop   []ExperimentsLoop `json:"ExperimentLoop"`
	BootShutdownRate []InstanceBootShutdownRate `json:"BootShutdownRate"`
}
type InstanceValue struct {
	NumInstances int `json:"NumInstances"`
	BootTime string `json:"BootTime"`
	ShutDownTime string `json:"ShutDownTime"`
}

type InstanceRegression struct {
	InstanceType string `json:"InstanceType"`
	Region string `json:"Region"`
	InstanceValues []InstanceValue `json:"InstanceValues"`
}

type VMBootShutDownRatePerInstanceTypeAll struct {
	InstanceValues []InstanceValue `json:"InstanceValues"`
}

type VMBootShutDownRatePerInstanceTypeOne struct {
	BootTime string `json:"BootTime"`
	ShutDownTime string `json:"ShutDownTime"`
}

// TestInstanceBootingShuttingTime godoc
// @Summary Conducts the VM booting and shutting down experiment
// @Description Conducts the VM booting and shutting down experiment
// @Tags INSTANCE_VM_BOOTING_TIME
// @Accept text/html
// @Produce json
// @Param numInstances query string true "number of instances"
// @Param totalExperiments query string true "total number of experiments to conduct"
// @Success 200 {string} string "started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /testInstanceBootingShuttingTime [get]
func TestInstanceBootingShuttingTime(w http.ResponseWriter, r *http.Request) {

	numInstances := r.URL.Query().Get("numInstances")
	totalExperiments := r.URL.Query().Get("totalExperiments")

	urlExp := "http://instance:8080/testInstanceBootingShuttingTime?numInstances="+numInstances+"&totalExperiments="+totalExperiments
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
// RandomTestInstanceBootingShuttingTime godoc
// @Summary Starts the RMIT experiment
// @Description Starts the RMIT experiment
// @Tags INSTANCE_VM_BOOTING_TIME
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

	urlExp := "http://instance:8080/getAllVMTypesBootShutDownDataAvg?instanceType="+instanceType+"&totalExperiments="+totalExperimentsStr
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

// GetAllVMTypesBootShutDownDataAvg godoc
// @Summary Gives all the Booting and shutting down data following avg approach
// @Description Gives all the Booting and shutting down data following avg approach
// @Tags INSTANCE_VM_BOOTING_TIME
// @Accept text/html
// @Produce json
// @Param takeNewValues query string false "specify it yes if to take walues from the db and update the resulting average"
// @Success 200 {array} TERMINUS.VmTemplateData ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAllVMTypesBootShutDownDataAvg [get]
func GetAllVMTypesBootShutDownDataAvg(w http.ResponseWriter, r *http.Request) {

	takeNewValues := r.URL.Query().Get("takeNewValues")
	urlExp := "http://instance:8080/getAllVMTypesBootShutDownDataAvg?takeNewValues="+takeNewValues
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}


// TrainDataSetRegression godoc
// @Summary Train the dataset
// @Description train the dataset
// @Tags INSTANCE_VM_BOOTING_TIME
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Started"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /trainDataSetRegression [get]
func TrainDataSetRegression(w http.ResponseWriter, r *http.Request){
	urlExp := "http://instance:8080/trainDataSetRegression"
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

// GetandStoreRegressionValues godoc
// @Summary Query python regression_vm_boot to get the values for each instance type and then store in db
// @Description Query python regression_vm_boot to get the values for each instance type and then store in db
// @Tags INSTANCE_VM_BOOTING_TIME
// @Accept text/html
// @Produce json
// @Success 200 {string} string "Done"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getandStoreRegressionValues [get]
func GetandStoreRegressionValues(w http.ResponseWriter, r *http.Request){
	urlExp := "http://instance:8080/getandStoreRegressionValues"
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}

// GetAllVMTypesBootShutDownDataRegression godoc
// @Summary Gives all the Booting and shutting down data
// @Description Gives all the Booting and shutting down data
// @Tags INSTANCE_VM_BOOTING_TIME
// @Accept text/html
// @Produce json
// @Success 200 {array} TERMINUS.InstanceRegression ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getAllVMTypesBootShutDownDataRegression [get]
func GetAllVMTypesBootShutDownDataRegression(w http.ResponseWriter, r *http.Request) {

	urlExp := "http://instance:8080/getAllVMTypesBootShutDownDataRegression"
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
// GetPerVMTypeAllBootShutDownData godoc
// @Summary Gives all the Booting and shutting down data for a VM type
// @Description Gives all the Booting and shutting down data
// @Tags INSTANCE_VM_BOOTING_TIME_EXTERNAL_USE
// @Accept text/html
// @Produce json
// @Param instanceType query string true "instance type"
// @Param region query string true "aws region "
// @Param appraoch query string true "appraoch avg or regression_vm_boot, by default it is avg"
// @Param csp query string true "cloud service provider..current;y it is aws only"
// @Success 200 {object} TERMINUS.VMBootShutDownRatePerInstanceTypeAll ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPerVMTypeAllBootShutDownData [get]
func GetPerVMTypeAllBootShutDownData(w http.ResponseWriter, r *http.Request) {
	instanceType := r.URL.Query().Get("instanceType")

	region := r.URL.Query().Get("region")

	approach := r.URL.Query().Get("approach")

	urlExp := "http://instance:8080/getPerVMTypeAllBootShutDownData?instanceType=" + instanceType + "&region=" + region+ "&approach=" + approach
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}


// GetPerVMTypeOneBootShutDownData godoc
// @Summary Gives all the Booting and shutting down data for a VM type
// @Description Gives all the Booting and shutting down data
// @Tags INSTANCE_VM_BOOTING_TIME_EXTERNAL_USE
// @Accept text/html
// @Produce json
// @Param numInstances query string true "number of instances"
// @Param instanceType query string true "instance type"
// @Param region query string true "aws region "
// @Param appraoch query string true "appraoch avg or regression, by default it is avg"
// @Param csp query string true "cloud service provider..current;y it is aws only"
// @Success 200 {object} TERMINUS.VMBootShutDownRatePerInstanceTypeOne ""
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /getPerVMTypeOneBootShutDownData [get]
func GetPerVMTypeOneBootShutDownData(w http.ResponseWriter, r *http.Request) {

	numInstances := r.URL.Query().Get("numInstances")

	instanceType := r.URL.Query().Get("instanceType")

	region := r.URL.Query().Get("region")

	approach := r.URL.Query().Get("approach")

	urlExp := "http://instance:8080/getPerVMTypeOneBootShutDownData?numInstances=" + numInstances + "&instanceType=" + instanceType + "&region=" + region+ "&approach=" + approach
	log.Info(urlExp)
	responseExp, err := http.Get(urlExp)
	if err != nil {
		log.Error(err)
		http.Error(w, "Json Conversion error", http.StatusBadRequest)
		return
	} else {
		defer responseExp.Body.Close()
		contentsExp, err := ioutil.ReadAll(responseExp.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Json Conversion error", http.StatusBadRequest)
			return
		}
		w.Write(contentsExp)
	}
}
